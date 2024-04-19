package Utilities

import (
	"PROYECTO1_MIA/Structs"
	"fmt"
	"strings"
)

func ParticionesMontadas(disco string) {
	filepath := "./MIA/P1/" + disco + ".dsk"
	file, err := OpenFile(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			estado := strings.TrimRight(string(TempMBR.Partitions[i].Status[:]), "\x00")
			name_parti := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
			id_parti := strings.TrimRight(string(TempMBR.Partitions[i].Id[:]), "\x00")
			if estado == "1" {
				fmt.Println("Particion: " + name_parti + ", Id: " + id_parti + ", estado: " + estado)
			}

		}
	}
}

func MontarParticion(disco string, name string) bool {
	//abrir el archivo
	monto := false       //Esta de funcion, logro hacer tarea = true
	var posicion_par int //Posicion de la particion a montar
	//existe el disco
	filepath := "./MIA/P1/" + disco + ".dsk"
	file, err := OpenFile(filepath)
	if err != nil {
		return monto
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := ReadObject(file, &TempMBR, 0); err != nil {
		return monto
	}

	Se_repite := false

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			nombre_particion := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
			if nombre_particion == name {
				//Existe la particion a montar
				Se_repite = true
				posicion_par = i
			}
		}
	}

	if Se_repite {
		//Validar si la particion esta montada o no
		estado := strings.TrimRight(string(TempMBR.Partitions[posicion_par].Status[:]), "\x00")
		if estado == "0" {
			//Particion no montada
			Tipo := strings.TrimRight(string(TempMBR.Partitions[posicion_par].Type[:]), "\x00")

			if Tipo == "e" {
				fmt.Println("Error: No se puede montar particiones extendidas")
				return monto
			}
			//Si existe la particion, Procede a montar

			id := fmt.Sprintf("%s%d%s", disco, TempMBR.Partitions[posicion_par].Correlative, "03")
			copy(TempMBR.Partitions[posicion_par].Id[:], id)
			copy(TempMBR.Partitions[posicion_par].Status[:], "1")
			//Reescribir MBR
			if err := WriteObject(file, TempMBR, 0); err != nil {
				return monto
			}
			monto = true
		} else {
			//Particion ya montada
			fmt.Println("Error: La particion ya esta montada ")
			return monto
		}

	} else {
		//NO existe la particion a montar
		fmt.Println("Error: No existe la particion a montar.  mount")
		return monto
	}
	return monto
}
