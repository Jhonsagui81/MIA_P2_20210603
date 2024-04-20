package reportes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"fmt"
	"os"
	"strings"
)

func Recorrido_MBR(disco string) string {
	grafo := "" //En caso se produce error, retornar cadena vacio
	text := ""  //contenido a retornar
	//existe el disco
	filepath := "./MIA/P1/" + disco + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return grafo
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return grafo
	}
	//Tamano MBR
	sizeMBR := fmt.Sprintf("%d", TempMBR.MbrSize)
	text += "<TR>\n\t"
	text += "<TD bgcolor=\"yellow\">mbr_tamano</TD>\n\t"
	text += "<TD bgcolor=\"yellow\">" + sizeMBR + "</TD>\n\t"
	text += "</TR>\n\t"

	//FechaCreacion
	fecha := strings.TrimRight(string(TempMBR.CreationDate[:]), "\x00")
	text += "<TR>\n\t"
	text += "<TD bgcolor=\"yellow\">mbr_fecha_creatopm</TD>\n\t"
	text += "<TD bgcolor=\"yellow\">" + string(fecha) + "</TD>\n\t"
	text += "</TR>\n\t"

	//signature (dudoso proceso de parseo)
	signature := fmt.Sprintf("%d", TempMBR.Signature)
	text += "<TR>\n\t"
	text += "<TD bgcolor=\"yellow\">mbr_disk_signature</TD>\n\t"
	text += "<TD bgcolor=\"yellow\">" + signature + "</TD>\n\t"
	text += "</TR>\n\t"

	//Iterar particiones
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			//Particion
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">PARTICION</TD>\n\t"
			text += "</TR>\n\n\t"
			//Estado
			status := strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_status</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(status) + "</TD>\n\t"
			text += "</TR>\n\t"
			//type
			tipo := string(TempMBR.Partitions[i].Type[:])
			type_ := strings.TrimRight(string(TempMBR.Partitions[i].Type[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_type</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(type_) + "</TD>\n\t"
			text += "</TR>\n\t"
			//fit
			fit := strings.TrimRight(string(TempMBR.Partitions[i].Fit[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_fit</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(fit) + "</TD>\n\t"
			text += "</TR>\n\t"
			//start
			startParticion := fmt.Sprintf("%d", TempMBR.Partitions[i].Start)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_start</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + startParticion + "</TD>\n\t"
			text += "</TR>\n\t"
			//size
			sizeParticion := fmt.Sprintf("%d", TempMBR.Partitions[i].Size)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_size</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + sizeParticion + "</TD>\n\t"
			text += "</TR>\n\t"
			//name
			name_ := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_name</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(name_) + "</TD>\n\t"
			text += "</TR>\n\t"

			if tipo == "e" {
				//Recorrer los EBR, cambia text
				text += Recorrido_logicas(file, int64(TempMBR.Partitions[i].Start))
			}
		}
	}

	return text
}

func Recorrido_logicas(file *os.File, SaltosEBR int64) string {
	continuar := true
	text := ""
	for continuar {
		//Recuperar el primer EBR
		var TempEBR Structs.EBR
		if err := Utilities.ReadObject(file, &TempEBR, SaltosEBR); err != nil {
			return ""
		}
		if TempEBR.Part_next == -1 {
			//Es el ultimo EBR si se grafica
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">PARTICION_LOGICA</TD>\n\t"
			text += "</TR>\n\n\t"
			//Next
			nextFinal := fmt.Sprintf("%d", TempEBR.Part_next)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_next</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + nextFinal + "</TD>\n\t"
			text += "</TR>\n\t"
			//FIT
			fit1 := strings.TrimRight(string(TempEBR.Part_fit[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_fit</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(fit1) + "</TD>\n\t"
			text += "</TR>\n\t"
			//start
			startFinal := fmt.Sprintf("%d", TempEBR.Part_start)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_start</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + startFinal + "</TD>\n\t"
			text += "</TR>\n\t"
			//size
			sizeFinal := fmt.Sprintf("%d", TempEBR.Part_s)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_size</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + sizeFinal + "</TD>\n\t"
			text += "</TR>\n\t"
			//name
			namee := strings.TrimRight(string(TempEBR.Part_name[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_name</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(namee) + "</TD>\n\t"
			text += "</TR>\n\t"

			continuar = false
		} else if TempEBR.Part_next > 0 {
			//Hay EBR conectados - Graficar
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">PARTICION_LOGICA</TD>\n\t"
			text += "</TR>\n\n\t"
			//Next
			nextFinal := fmt.Sprintf("%d", TempEBR.Part_next)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_next</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + nextFinal + "</TD>\n\t"
			text += "</TR>\n\t"
			//FIT
			fit1 := strings.TrimRight(string(TempEBR.Part_fit[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_fit</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(fit1) + "</TD>\n\t"
			text += "</TR>\n\t"
			//start
			startFinal := fmt.Sprintf("%d", TempEBR.Part_start)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_start</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + startFinal + "</TD>\n\t"
			text += "</TR>\n\t"
			//size
			sizeFinal := fmt.Sprintf("%d", TempEBR.Part_s)
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_size</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + sizeFinal + "</TD>\n\t"
			text += "</TR>\n\t"
			//name
			namee := strings.TrimRight(string(TempEBR.Part_name[:]), "\x00")
			text += "<TR>\n\t"
			text += "<TD bgcolor=\"yellow\">part_name</TD>\n\t"
			text += "<TD bgcolor=\"yellow\">" + string(namee) + "</TD>\n\t"
			text += "</TR>\n\t"
			SaltosEBR = int64(TempEBR.Part_next)
		}
	}
	return text
}

// func Parseo_int32(valor int32) {
// 	numero, err := strconv.Atoi(valor)
// 	if err != nil {
// 		fmt.Println("Error: Correlativo incorrecto")
// 		return desmonto
// 	}
// }
