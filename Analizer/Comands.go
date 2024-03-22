package analizer

import (
	reportes "PROYECTO1_MIA/Reportes"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"fmt"
	"os"
	"strings"
	"time"
)

func rep(name string, path string, id string) {
	//Open file
	fmt.Println("======Start REP======")
	//Transformacio a dot
	pathdot := path[:len(path)-3]
	pathdot += "dot"
	//Transformarlo a pdf
	pathPDF := path[:len(path)-3]
	pathPDF += "pdf"

	switch name {
	case "mbr":
		byte_id := []byte(id)
		valor_disk := string(byte_id[0])
		fmt.Println(valor_disk)
		text := "digraph G {\n\ta0 [shape=none label=<\n\t<TABLE cellspacing=\"10\" cellpadding=\"10\" style=\"rounded\" bgcolor=\"red\">\n\t"
		text += "<TR>\n\t"
		text += "<TD bgcolor=\"yellow\">REPORTE MBR</TD>\n\t"
		text += "</TR>\n\n\t"

		//Validar que exista disco
		disco := Utilities.BusquedaArchivo(valor_disk)
		if disco {
			//Si existe el disco
			text += reportes.Recorrido_MBR(valor_disk)
			//Cerrar la estructura
			text += "</TABLE>>];\n"
			text += "}"

			Utilities.CrearGrafo(text, pathdot, pathPDF)

			// fmt.Println(text)
		} else {
			//no existe
			fmt.Println("Error: No existe el disco. rep")
			return
		}
	case "disk":
		byte_id := []byte(id)
		valor_disk := string(byte_id[0])
		fmt.Println(valor_disk)

		text := "digraph D {\n\t"
		text += "subgraph cluster_0 {\n\t\t"
		text += "bgcolor=\"#68d9e2\"\n\t\t"
		text += "node [style=\"rounded\" style=filled];\n\t\t"
		text += "node_A [shape=record   label=\"MBR\\n"

		//Validar que exista disco
		disco := Utilities.BusquedaArchivo(valor_disk)
		if disco {
			//Si existe el disco
			text += reportes.Recorrido_disk(valor_disk)
			// pathdot := path[:len(path)-3]
			newText := text[:len(text)-1]
			newText += "\"];\n\t"
			newText += "}\n"
			newText += "}"

			Utilities.CrearGrafo(newText, pathdot, pathPDF)

			// fmt.Println(text)
		} else {
			//no existe
			fmt.Println("Error: No existe el disco. rep")
			return
		}
	default:
		fmt.Println("Reporte no reconocido")
	}

}

func Unmount(id string) {
	fmt.Println("===== Start Unmount =====")
	byte_id := []byte(id)
	valor_disk := byte_id[0]
	valor_corre := byte_id[1]

	// fmt.Println("Valor id:", id)
	// fmt.Println("Valor de disco:", string(valor_disk))
	// fmt.Println("Valor de correlarivo:", string(valor_corre))

	disco := Utilities.BusquedaArchivo(string(valor_disk))
	if disco {
		//Si existe disco
		desmonto := Utilities.Desmontar(string(valor_disk), string(valor_corre), id)
		if desmonto {
			//Desmonto
			fmt.Println("Se desmonto la particion con id:", id)
		} else {
			return
		}
	} else {
		fmt.Println("Error: No existe el disco. unmount")
		return
	}

	fmt.Println("===== End Unmount ======")
}

func mount(letter string, name string) {
	fmt.Println("=====Start Mount=====")
	disco := Utilities.BusquedaArchivo(letter)
	if disco {
		//Si existe el disco
		monto := Utilities.MontarParticion(letter, name)
		if monto {
			//monto particion
			fmt.Println("Se monto la particion:", name)
			//Mostrar lista de particiones montadas
			fmt.Println("Particiones montadas: ")
			Utilities.ParticionesMontadas(letter)
		} else {
			//paso algo
			fmt.Println("===== End Mount ======")
			return
		}
		//Tratar de buscar la partcion y trata de montar
	} else {
		//No existe el disco
		fmt.Println("Error: No existe el disco. mount")
		fmt.Println("===== End Mount ======")
		return
	}

	fmt.Println("===== End Mount ======")
}

func fdisk(size int, unit string, letter string, name string, type_ string, fit string, delete string, add int) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("LetterDrive:", letter)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)
	fmt.Println("delete:", delete)
	fmt.Println("add:", add)

	//Validar si las etiquetas delete y add estan vacias
	if delete == "full" {
		//Procede a trarar de eliminar la particion (solo primarias)
		//pasar type a null - para que no entre al siguiente switch
		type_ = "xd"
		add = 0
		size = 1
		//Buscar la particion - si no existe mostrar error - si existe espera confirmacion

		//Limpiar el espacio de la particion con 0

		//Limpiar la informacion en MBR
		eliminacion := Utilities.EliminarParticion(name, letter)
		if !eliminacion {
			//Eliminacion correcta
			fmt.Println("Se elimino particion correctamente")
			fmt.Println("===== END Fdisk =====")
			return
		} else {
			//Algo salio mal con la eliminacion
			return
		}
	} else if delete == "" {
		//no viene parametro delete
	} else {
		//parametro delete incorrecto
		fmt.Println("Error: parametro delete recibe parametro no valido")
		fmt.Println("===== END Fdisk =====")
	}

	if add != 0 {
		//Viene parametro add, unit, driveletter, name
		type_ = "xd"
		size = 1
		//Agregar o eliminar espacio
		agrego := Utilities.Add_Espacio(letter, name, unit, add)
		if agrego {
			//agrego o resto espacio
			fmt.Println("Se implemento Add")
			fmt.Println("===== END Fdisk =====")
			return
		} else {
			//Algo salio mal
			return
		}

	}

	// validate size > 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		fmt.Println("===== END Fdisk =====")
		return
	}

	//Validar si archivo existe
	existe := Utilities.BusquedaArchivo(letter)
	if !existe {
		//El Disco existe
		fmt.Println("Error: Disco no existe")
		fmt.Println("===== END Fdisk =====")
		return
	}

	// Open bin file
	filepath := "./MIA/P1/" + letter + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	//Valida nombre de particion no se repita
	se_repite := Utilities.RecorridoParticionesDisco(file, name, TempMBR)
	if se_repite {
		//fallo apertura|lectura|nombre de particion ya existe
		fmt.Println("Error: Nombre de la particion ya existe")
		fmt.Println("===== END Fdisk =====")
		return
	}

	// validate unit equals to b/k/m
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be b, k or m")
		fmt.Println("===== END Fdisk =====")
		return
	}

	// Set the size in bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// validate fit equals to b/w/f
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit must be b, w or f")
		fmt.Println("===== END Fdisk =====")
		return
	}

	// validate type equals to p/e/l
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Type must be p, e or l")
		fmt.Println("===== END Fdisk =====")
		return
	}

	switch type_ {
	case "p":
		//Validar que no haya mas de 4 particiones por disco(NO sirve)
		disco_lleno := Utilities.ConteoParticiones(TempMBR)
		if disco_lleno {
			fmt.Println("Error: El disco no permite mas particiones")
			fmt.Println("===== END Fdisk =====")
			return
		}
		//Inicia proceso de incersion
		aceptacion := Utilities.InsertarParticion(file, &TempMBR, size, unit, name, type_, fit)
		if !aceptacion {
			//Espacio insuficiente o mala sobre escritura del mbr
			fmt.Println("Error: La particion no cabe en el disco")
			fmt.Println("===== END Fdisk =====")
			return
		}
	case "e":
		//Validar que no haya 4 particiones por disco(NO sirve)
		disco_lleno := Utilities.ConteoParticiones(TempMBR)
		if disco_lleno {
			fmt.Println("Error: El disco no permite mas particiones")
			fmt.Println("===== END Fdisk =====")
			return
		}
		//Validar que no exista otra extendida
		existe_extend := Utilities.ParticionExtendida(TempMBR)
		if existe_extend {
			fmt.Println("Error: No se permite mas de UNA particion extendida")
			fmt.Println("===== END Fdisk =====")
			return
		} else {
			//Inicia proceso de insersion
			aceptacion := Utilities.InsertarParticion(file, &TempMBR, size, unit, name, type_, fit)
			if !aceptacion {
				//Espacio insuficiente o mala sobre escritura del mbr
				fmt.Println("Error: La particion no cabe en el disco")
				fmt.Println("===== END Fdisk =====")
				return
			}
		}
	case "l":
		//Validar que exista una extendida en el disco
		existe := Utilities.ParticionExtendida(TempMBR)
		if existe {
			//Como existe una extendida puede proceder
			Utilities.InsertaLogica(file, &TempMBR, size, name, fit)
		} else {
			fmt.Println("Error: No existe particion extendida")
			fmt.Println("===== END Fdisk =====")
			return
		}

	}

	// Close bin file
	defer file.Close()

	fmt.Println("======End FDISK======")
}

func mkdisk(size int, unit string, fit string) {
	fmt.Println("======Start MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Unit:", unit)
	fmt.Println("Unit:", fit)

	// validate size > 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate unit equals to k/m
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be k or m")
		return
	}

	// Validate fit equals to BF/FF/WF
	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: Fit must be bf, ff or wf")
		return
	}

	// Create file
	nombre_archivo := fmt.Sprintf("./MIA/P1/%c.dsk", Letra_Disco)

	err := Utilities.CreateFile(nombre_archivo)
	ContadorArchivos++
	Letra_Disco++
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Set the size in bytes
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	file, err := Utilities.OpenFile(nombre_archivo)
	if err != nil {
		return
	}

	// Write 0 binary data to the file
	kb := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		kb[i] = 0
	}
	for i := 0; i < size/1024; i++ {
		if _, err := file.Write(kb); err != nil {
			fmt.Println("Error: No puedo con los ceros", err)
			return
		}
	}

	// Create a new instance of MRB
	Aux_asignature := Utilities.NumeroRandom()
	now_date := time.Now()
	formattedDate := now_date.Format("2006-01-02")

	var newMRB Structs.MRB
	newMRB.MbrSize = int32(size)
	newMRB.Signature = int32(Aux_asignature) // random
	copy(newMRB.Fit[:], fit)
	copy(newMRB.CreationDate[:], formattedDate)

	// Write object in bin file
	if err := Utilities.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	fmt.Println("======End MKDISK======")
}

func rmdisk(letter string) {
	fmt.Printf("=====Start Rmdisk=====")
	fmt.Println("driveletter:", letter)

	//buscar los archivos contenidos en el directorio
	encontrado := Utilities.BusquedaArchivo(letter)

	if encontrado {
		var confirmacion string
		name_file := letter + ".dsk"
		fmt.Printf("->Disco '%s' encontrado\n", name_file)
		fmt.Println("->Desea Eliminar el disco " + name_file + " (S/N)? ")
		_, err := fmt.Scan(&confirmacion)

		//valida lectura de consola
		if err != nil {
			fmt.Println("Error: confirmacion rmdisk", err)
			fmt.Printf("=====END Rmdisk=====")
			return
		}

		//Valida si la respusta es valida
		aux_confir := strings.ToLower(confirmacion)
		if aux_confir != "s" && aux_confir != "n" {
			fmt.Println("Error: opcion incorrecta rmdisk ", err)
			fmt.Printf("=====END Rmdisk=====")
			return
		}

		if aux_confir == "s" {
			//eliminar el disco
			ruta_archivo := "./MIA/P1/" + name_file
			err := os.Remove(ruta_archivo)
			if err != nil {
				fmt.Println("Error: Delete disk  rmdisk", err)
				fmt.Printf("=====END Rmdisk=====")
				return
			}
			fmt.Printf("Disco %s Eliminado\n", name_file)

		} else {
			//No se eliminar
			fmt.Println("No se elimino el disco")
			fmt.Printf("=====END Rmdisk=====")
		}

		//peticion

	} else {
		fmt.Println("-> Disco buscado no existe")
		fmt.Printf("=====END Rmdisk=====")
	}
}
