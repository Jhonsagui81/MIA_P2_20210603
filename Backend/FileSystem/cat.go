package filesystem

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"fmt"
	"strings"
)

func Cat(path string) {
	fmt.Println("===== START Cat =====")
	if global.LoginValidacion() {
		id := global.InfoDisk()

		driveletter := string(id[0])

		filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
		file, err := Utilities.OpenFile(filepath)
		if err != nil {
			fmt.Println("Error: El disco no existe -> login")
			return
		}

		//Recupera el MBR
		var TempMBR Structs.MRB
		// Read object from bin file
		if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
			return
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
						return
					}
					break
				}
			}
		}

		if index != -1 {
			fmt.Println("Se encontro la particion")
		} else {
			fmt.Println("Partition not found")
			return
		}

		//Recupera el super bloque
		var tempSuperblock Structs.Superblock
		// Read object from bin file
		if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
			return
		}

		indexInode := UtilitiesInodes.InitSearch(path, file, tempSuperblock)

		if indexInode != -1 {
			//Recupero el Inodo que indica indexInode
			var crrInode Structs.Inode
			if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
				return
			}

			//Validar los permisos
			uid := fmt.Sprintf("%d", crrInode.I_uid)
			gid := fmt.Sprintf("%d", crrInode.I_gid)
			perm_u := strings.TrimRight(string(crrInode.I_perm[0]), "\x00")
			perm_g := strings.TrimRight(string(crrInode.I_perm[1]), "\x00")
			perm_o := strings.TrimRight(string(crrInode.I_perm[2]), "\x00")

			if global.DeterminarPermisoEscritura(gid, uid, perm_u, perm_g, perm_o) {
				//Tiene permiso de escritura
				fmt.Println("SI tiene permiso para crear el archivo")
				data := UtilitiesInodes.GetInodeFileData(crrInode, file, tempSuperblock)

				fmt.Println("########## Contenido de archivo ##########")
				fmt.Println(data)
				fmt.Println("########## Finaliza archivo ##########")
			} else {
				fmt.Println("Error: El usuario actual no tiene permiso de escritura -MKDIR")
				fmt.Println("====== END MKDIR =====")
				return
			}
			// getInodeFileData -> Iterate the I_Block n concat the data
			//Falta crear funcion para iterar todos los bloques del inodo archivo y concatenar
			//Validar apuntadores indirectos
			//Repuero el bloque del archivo, para obtener la data

		}

	}
	fmt.Println("===== END Cat =====")
}
