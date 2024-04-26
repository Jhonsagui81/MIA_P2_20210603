package main

import (
	analizer "PROYECTO1_MIA/Analizer"
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Command struct {
	Nombre string `json: "Nombre"`
	Id     int    `json:"ID"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	//Configuracion cors
	corsObj := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"*",
		},
	})

	handler := corsObj.Handler(router)

	// Endpoint para recibir el nombre
	router.HandleFunc("/input", GenerarArchivo).Methods("POST")
	router.HandleFunc("/comand", InputCommand).Methods("POST")
	router.HandleFunc("/tarea", getMensage).Methods("GET")
	router.HandleFunc("/discos", GetDiscos).Methods("GET")
	router.HandleFunc("/partitions", GetParticiones).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/sistema", ArchivosEspecificos).Methods("POST")
	router.HandleFunc("/reportes", GetReportes).Methods("GET")

	fmt.Println("Servidor levantado en puerto 3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func GetReportes(w http.ResponseWriter, r *http.Request) {
	var reportes []Structs.Contenido
	archivos, err := ioutil.ReadDir("/home/jhonatan/archivos/reports")
	if err != nil {
		fmt.Println("Error: No se puedo acceder a los discos ", err)
	}

	for _, archivo := range archivos {
		if strings.HasSuffix(archivo.Name(), ".dot") {
			nombre := archivo.Name()
			idd := len(reportes) + 1
			//Transformar el contenido
			file, err := ioutil.ReadFile("/home/jhonatan/archivos/reports/" + archivo.Name())
			if err != nil {
				fmt.Println("Error leer pdf")
			}

			// Convierte el contenido a base64
			pdfBase64 := string(file)

			var pdf Structs.Contenido
			pdf.Nombre = nombre
			pdf.Id = idd
			pdf.Content = pdfBase64
			pdf.Imagen = "https://img.freepik.com/vector-premium/icono-pdf-estilo-comic-ilustracion-dibujos-animados-vector-texto-documento-sobre-fondo-blanco-aislado-concepto-negocio-efecto-salpicadura-archivo_157943-16160.jpg"

			reportes = append(reportes, pdf)
		}

	}
	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	fmt.Println("envia a:", reportes)
	json.NewEncoder(w).Encode(reportes)
}

func GetDiscos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entra a get Discos")
	discos := []Command{}
	//Obtener los discos
	archivos, err := ioutil.ReadDir("./MIA/P1")
	if err != nil {
		fmt.Println("Error: No se puedo acceder a los discos ", err)
	}
	//rellenar la estructura
	for _, archivo := range archivos {
		nombre := archivo.Name()

		idd := len(discos) + 1
		nuevoDisco := Command{Nombre: nombre, Id: idd}

		discos = append(discos, nuevoDisco)
	}

	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	fmt.Println("envia a:", discos)
	json.NewEncoder(w).Encode(discos)

}

func GenerarArchivo(w http.ResponseWriter, r *http.Request) {
	var archivo Structs.Command
	requesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid data")
	}
	//asigna el dato recibido a la variable
	json.Unmarshal(requesBody, &archivo)

	errr := ioutil.WriteFile("archivo.txt", []byte(archivo.Nombre), 0644)
	if errr != nil {
		log.Fatal("Error al escribir el archivo:", err)
	} else {
		fmt.Println("Archivo creado o sobrescrito exitosamente")
	}

	analizer.AnalizerType("archivo.txt")

	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(archivo)

}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entra a login")
	var login Structs.Command
	requesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid data")
	}

	//asigna el dato recibido a la variable
	json.Unmarshal(requesBody, &login)

	analizer.Analize(login.Nombre)

	if global.LoginValidacion() {
		//Si se puedo iniciar sesion
		_, params := analizer.GetCommandAndParams(login.Nombre)
		//Devolver los doc de la carpeta raiz y algo para validar login
		var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
		user := ""
		pass := ""
		id := ""
		matches := re.FindAllStringSubmatch(params, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"") // Removing the double
			switch flagName {
			case "user":
				user = flagValue
			case "pass":
				pass = flagValue
			case "id":
				flagValue = strings.ToUpper(flagValue)
				id = flagValue
			default:
			}

		}
		//FUncion para obtener los archivos raiz
		archivos := BusquedaRuta(user, pass, id, "/")
		w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
		w.WriteHeader(http.StatusCreated)
		fmt.Println("envia a:", archivos)
		json.NewEncoder(w).Encode(archivos)

	} else {
		//No hay usuario logeado
		//fallo en el login -> no cambiar de ventana (mostrar error ingrese datos correctos)
		var archivo Structs.Contenido
		archivo.Nombre = ""
		w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
		w.WriteHeader(http.StatusCreated)
		fmt.Println("envia a:", archivo)
		json.NewEncoder(w).Encode(archivo)
	}

}
func GetParticiones(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entra a obtener particioes")
	var namePartition Structs.Command
	requesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid data")
	}

	//asigna el dato recibido a la variable
	json.Unmarshal(requesBody, &namePartition)

	// Open bin file
	filepath := "./MIA/P1/" + namePartition.Nombre
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	particiones := Utilities.ObtenerParticiones(file, TempMBR)
	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	fmt.Println("envia a:", particiones)
	json.NewEncoder(w).Encode(particiones)

}
func InputCommand(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Si llesdsdgo el comando")
	var comand Structs.Command
	requesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid data")
	}

	//asigna el dato recibido a la variable
	json.Unmarshal(requesBody, &comand)

	//Llamar analize
	var response Structs.RespuestaFron
	res := analizer.Analize(comand.Nombre)
	//Responder al cliente
	response.Respuesta = res
	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)                  //indicar que todo bien
	json.NewEncoder(w).Encode(response)

}

func ArchivosEspecificos(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Si llesdsdgo a la busqueda especifica")
	var comand Structs.Command
	requesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid data")
	}

	//asigna el dato recibido a la variable
	json.Unmarshal(requesBody, &comand)

	user, pass, id := global.DataLogin()

	archivos := BusquedaRuta(user, pass, id, comand.Nombre)

	if archivos == nil {
		var cont Structs.Contenido
		cont.Content = "Esta Carpeta esta vacia"
		archivos = append(archivos, cont)
	}

	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	fmt.Println("envia a:", archivos)
	json.NewEncoder(w).Encode(archivos)
}

func getMensage(w http.ResponseWriter, r *http.Request) {
	tasks := `{
		"mensaje" : "hoa"
	}`
	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tasks)
}

func BusquedaRuta(user string, pass string, id string, ruta string) []Structs.Contenido {
	driveletter := string(id[0])
	var archivos []Structs.Contenido

	// Open bin file
	filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: El disco no existe -> login")
		return archivos
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return archivos
	}

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				if strings.Contains(string(TempMBR.Partitions[i].Status[:]), "1") {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return archivos
				}
				break
			}
		}
	}

	if index != -1 {
		fmt.Println("Se encontro la particion")
	} else {
		fmt.Println("Partition not found")
		return archivos
	}

	var tempSuperblock Structs.Superblock
	// Read object from bin file
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		return archivos
	}

	if ruta == "/" {
		//La raiz del sistema -> inicia sesion por primera vez
		var Inode0 Structs.Inode
		// Read object from bin file
		if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
			return archivos
		}

		index = 0
		for _, block := range Inode0.I_block {
			if block != -1 {
				if index < 13 {
					//CASO DIRECTO

					var crrFolderBlock Structs.Folderblock
					// Read object from bin file
					if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
						return archivos
					}

					for _, folder := range crrFolderBlock.B_content {
						// fmt.Println("Folder found======")

						namefolder := strings.TrimRight(string(folder.B_name[:]), "\x00")
						if namefolder != "." && namefolder != ".." && folder.B_inodo != -1 {
							//Buscar el inodo para determinar si es carpeta o archivo
							var NextInode Structs.Inode
							// Read object from bin file
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return archivos
							}
							//En base a eso determinar la imagen y el tipo
							var cont Structs.Contenido
							cont.Nombre = strings.TrimRight(string(folder.B_name[:]), "\x00")
							cont.Content = ""
							count := len(archivos) + 1
							if string(NextInode.I_type[:]) == "0" {
								//carpeta
								fmt.Println("tipo carpeta -BAACK")
								cont.Imagen = "https://www.vhv.rs/dpng/d/552-5526475_carpeta-de-computadora-animada-hd-png-download.png"
							} else {
								fmt.Println("ESTEes archivo ---")
								cont.Imagen = "https://us.123rf.com/450wm/stockgiu/stockgiu1708/stockgiu170806199/84869128-documento-de-negocio-para-la-comercializaci%C3%B3n-estrategia-smm.jpg?ver=6"
							}
							cont.Id = count

							archivos = append(archivos, cont)
						}
					}

				} else {
					//CASO INDIRECTO
				}
			}
			index++
		}

	} else { //No se busca la raiz
		indexInode := UtilitiesInodes.InitSearch(ruta, file, tempSuperblock)
		//Inodo de la carpeta contenedora
		var CRRRInode0 Structs.Inode
		// Read object from bin file
		if err := Utilities.ReadObject(file, &CRRRInode0, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
			return archivos
		}

		if string(CRRRInode0.I_type[:]) == "0" {
			//Estoy en una carpeta
			fmt.Println("Entro a inodo tipo carpeta")

			index = 0
			for _, block := range CRRRInode0.I_block {
				if block != -1 {
					if index < 13 {
						//CASO DIRECTO

						var crrFolderBlock Structs.Folderblock
						// Read object from bin file
						if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
							return archivos
						}

						for _, folder := range crrFolderBlock.B_content {
							// fmt.Println("Folder found======")

							namefolder := strings.TrimRight(string(folder.B_name[:]), "\x00")
							if namefolder != "." && namefolder != ".." && folder.B_inodo != -1 {
								//Buscar el inodo para determinar si es carpeta o archivo
								var NextInode Structs.Inode
								// Read object from bin file
								if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
									return archivos
								}
								//En base a eso determinar la imagen y el tipo
								var cont Structs.Contenido
								cont.Nombre = strings.TrimRight(string(folder.B_name[:]), "\x00")
								cont.Content = ""
								count := len(archivos) + 1
								if string(NextInode.I_type[:]) == "0" {
									//carpeta
									fmt.Println("tipo carpeta -BAACK")
									cont.Imagen = "https://www.vhv.rs/dpng/d/552-5526475_carpeta-de-computadora-animada-hd-png-download.png"
								} else {
									fmt.Println("ESTEes archivo ---")
									cont.Imagen = "https://us.123rf.com/450wm/stockgiu/stockgiu1708/stockgiu170806199/84869128-documento-de-negocio-para-la-comercializaci%C3%B3n-estrategia-smm.jpg?ver=6"
								}
								cont.Id = count
								archivos = append(archivos, cont)
							}
						}

					} else {
						//CASO INDIRECTO
					}
				}
				index++
			}
		} else {
			//Estoy en un archivo
			data := UtilitiesInodes.GetInodeFileData(CRRRInode0, file, tempSuperblock)
			var cont Structs.Contenido
			cont.Content = data
			archivos = append(archivos, cont)
			return archivos

		}
	}

	return archivos
}
