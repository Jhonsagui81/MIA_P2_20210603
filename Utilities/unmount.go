package Utilities

import (
	"PROYECTO1_MIA/Structs"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Desmontar(disco string, correlativo string, id string) bool {
	desmonto := false
	var posicion_par int
	//existe el disco
	filepath := "./MIA/P1/" + disco + ".dsk"
	file, err := OpenFile(filepath)
	if err != nil {
		return desmonto
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := ReadObject(file, &TempMBR, 0); err != nil {
		return desmonto
	}

	//Existe el correlativo
	corre := false
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {

			numero, err := strconv.Atoi(correlativo)
			if err != nil {
				fmt.Println("Error: Correlativo incorrecto")
				return desmonto
			}

			if TempMBR.Partitions[i].Correlative == int32(numero) {
				//Existe la particion a montar
				corre = true
				posicion_par = i
				break
			}
		}
	}

	if corre {
		//Validar si la particion esta montada o no
		estado := strings.TrimRight(string(TempMBR.Partitions[posicion_par].Status[:]), "\x00")
		if estado == "0" {
			//Particion no montada
			fmt.Println("Error: La particion no esta montada. unmount")
			return desmonto

		} else {
			//Esta montada -> se debe desmontar
			reflect.ValueOf(&TempMBR.Partitions[posicion_par].Id).Elem().Set(reflect.Zero(reflect.TypeOf(TempMBR.Partitions[posicion_par].Id)))
			copy(TempMBR.Partitions[posicion_par].Status[:], "0")
			//Reescribir MBR
			if err := WriteObject(file, TempMBR, 0); err != nil {
				return desmonto
			}
			desmonto = true
		}

	} else {
		//NO existe la particion a montar
		fmt.Println("Error: No existe el correlativo a desmontar.  unmount")
		return desmonto
	}

	return desmonto
}
