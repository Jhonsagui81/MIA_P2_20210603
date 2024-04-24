package main

import (
	analizer "PROYECTO1_MIA/Analizer"
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
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
	router.HandleFunc("/comand", InputCommand).Methods("POST")
	router.HandleFunc("/tarea", getMensage).Methods("GET")
	router.HandleFunc("/discos", GetDiscos).Methods("GET")
	router.HandleFunc("/partitions", GetParticiones).Methods("POST")
	router.HandleFunc("/login").Methods("POST")

	fmt.Println("Servidor levantado en puerto 3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
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
		//Devolver los doc de la carpeta raiz y algo para validar login
		var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
		user := ""
		pass := ""
		id := ""
		flags := true
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
				flags = false
			}

		}

	} else {
		//No hay usuario logeado
		//fallo en el login -> no cambiar de ventana (mostrar error ingrese datos correctos)
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

func getMensage(w http.ResponseWriter, r *http.Request) {
	tasks := `{
		"mensaje" : "hoa"
	}`
	w.Header().Set("Content-Type", "application/json") //tipo de dato se enviara
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tasks)
}
