package reportes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Recorrido_disk(disco string) string {
	text := ""
	filepath := "./MIA/P1/" + disco + ".dsk"

	ultima_parti := 3

	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return ""
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return ""
	}
	//Valores Disco
	tamano := binary.Size(TempMBR)                                        //size MBR
	TotalDisk := TempMBR.MbrSize                                          //size total disco
	porcenta_MBR := float32((float32(tamano) / float32(TotalDisk)) * 100) //porcentaje
	porMBR := fmt.Sprintf("%f", porcenta_MBR)                             //parseo
	text += string(porMBR) + "|"
	fmt.Println("totalsize:", porcenta_MBR)

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			//si particion existe se coloca el nombre
			tipo := string(TempMBR.Partitions[i].Type[:])
			if tipo == "e" {
				//Si es extendida es otro proceso
				text += "{Extendida|{"
				text += recorrido_extendida(file, int64(TempMBR.Partitions[i].Start), int32(TempMBR.Partitions[i].Size), TotalDisk)
				text += "|"
			} else {
				//particion primaria
				name_ := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
				text += name_ + "\\n"
				//agregar porcentaje
				porce_particion := float32(float32(TempMBR.Partitions[i].Size)/float32(TotalDisk)) * 100
				porParticion := fmt.Sprintf("%f", porce_particion)
				text += string(porParticion) + "|"

			}
		} else {
			ultima_parti = i - 1
			//No hay particion -> Libre
			text += "Libre\\n"
			//si llega aqui es porque ya no hay mas particiones por tanto es vacio
			sobrante := TotalDisk - (TempMBR.Partitions[i-1].Start + TempMBR.Partitions[i-1].Size)
			porcentaje := float32(float32(sobrante)/float32(TotalDisk)) * 100
			porVacio := fmt.Sprintf("%f", porcentaje)
			text += porVacio + "|"
			break
		}
	}

	//Si se utilizan las 4 particionse y no se lleno disco
	if ultima_parti == 3 {
		//No hay particion -> Libre
		text += "Libre\\n"
		//si llega aqui es porque ya no hay mas particiones por tanto es vacio
		sobrante := TotalDisk - (TempMBR.Partitions[ultima_parti].Start + TempMBR.Partitions[ultima_parti].Size)
		porcentaje := float32(float32(sobrante)/float32(TotalDisk)) * 100
		porVacio := fmt.Sprintf("%f", porcentaje)
		text += porVacio + "|"
	}

	return text
}

func recorrido_extendida(file *os.File, SaltosEBR int64, size_extendida int32, sizeTotal int32) string {
	continuar := true
	text := ""

	for continuar {
		//Recuperar el primer EBR
		var TempEBR Structs.EBR
		if err := Utilities.ReadObject(file, &TempEBR, SaltosEBR); err != nil {
			return ""
		}
		if TempEBR.Part_next == 0 {
			//No hay EBR por tanto toda la extendida esta vacia
			text += "Libre\\n"
			porcentaje := float32(float32(size_extendida)/float32(sizeTotal)) * 100
			porVacio := fmt.Sprintf("%f", porcentaje)
			text += porVacio + "}}"
			continuar = false
		} else if TempEBR.Part_next == -1 {
			//Es el ultimo EBR -> se debe calcular EBR, LOGICA y el sobrante para la libre
			text += "EBR\\n"
			sizeEBR := binary.Size(TempEBR)
			porcentaje_ := float32(float32(sizeEBR)/float32(sizeTotal)) * 100
			porVacio := fmt.Sprintf("%f", porcentaje_)
			text += porVacio + "|"
			//logica
			text += "Logica\\n"
			porcentajeLo := float32(float32(TempEBR.Part_s)/float32(sizeTotal)) * 100
			porVacio1 := fmt.Sprintf("%f", porcentajeLo)
			text += porVacio1
			//El posible espacio libre
			resto := (size_extendida) - (TempEBR.Part_start + TempEBR.Part_s)
			if resto > 0 {
				text += "|Libre\\n"
				porcenta := float32(float32(resto)/float32(sizeTotal)) * 100
				porVacio := fmt.Sprintf("%f", porcenta)
				text += porVacio
			}
			text += "}}"
			continuar = false
		} else if TempEBR.Part_next > 0 {
			//como estan enlazados tengo que hace EBR y Logica
			text += "EBR\\n"
			sizeEBR := binary.Size(TempEBR)
			porcentaje_ := float32(float32(sizeEBR)/float32(sizeTotal)) * 100
			porVacio := fmt.Sprintf("%f", porcentaje_)
			text += porVacio + "|"
			//logica
			text += "Logica\\n"
			porcentajeLo := float32(float32(TempEBR.Part_s)/float32(sizeTotal)) * 100
			porVacio1 := fmt.Sprintf("%f", porcentajeLo)
			text += porVacio1 + "|"
			SaltosEBR = int64(TempEBR.Part_next)
		}
	}

	return text
}
