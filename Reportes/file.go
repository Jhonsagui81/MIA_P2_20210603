package reportes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

func ReporteFile(disco string, id string, ruta string, path string) bool {
	encontrado := false
	filepath := "./MIA/P1/" + disco + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return encontrado
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return encontrado
	}

	PartitionStart := -1
	//Buscar la particion
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				if strings.Contains(string(TempMBR.Partitions[i].Status[:]), "1") {
					fmt.Println("Partition is mounted")
					PartitionStart = i

				} else {
					fmt.Println("Partition is not mounted")
					return encontrado
				}
				break
			}
		}
	}

	if PartitionStart != -1 {
		fmt.Println("Se encontro la particion")
	} else {
		fmt.Println("Partition not found")
		return encontrado
	}

	//Recupera el super bloque
	var tempSuperblock Structs.Superblock
	// Read object from bin file
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
		return encontrado
	}

	//Recuperar indice del archivo que busco
	indexInode := UtilitiesInodes.InitSearch(ruta, file, tempSuperblock)
	if indexInode != -1 {
		//el archivo fue encontrado
		encontrado = true
		var crrInode Structs.Inode
		if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
			return false
		}

		data := UtilitiesInodes.GetInodeFileData(crrInode, file, tempSuperblock)

		//Crear archivo
		err := Utilities.CreateFile(path)
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}

		//Abrir el archivo
		file, err := Utilities.OpenFile(path)
		if err != nil {
			fmt.Println("Error: No se pudo abrir el archivo -Rep file")
			return false
		}

		//Escribir archivo
		_, err = io.WriteString(file, data)
		if err != nil {
			fmt.Println("Error: No se pudo escribir el archivo -Rep file")
			return false
		}

		file.Close()
	} else {
		fmt.Println("Error: Archivo no existe  -REPFILE")
		return encontrado
	}

	return encontrado

}
