package filesystem

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func Mkdir(path string, r bool) {
	fmt.Println("===== START Mkdir =====")
	fmt.Println("path: ", path)

	if global.LoginValidacion() { //Valida que haya un usuario logeado
		// partition := global.InfoPartition()
		driveletter := global.InfoDisk()
		id := global.InfoID()

		filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
		file, err := Utilities.OpenFile(filepath)
		if err != nil {
			fmt.Println("Error: El disco no existe -> login")
			return
		}
		defer file.Close()

		//Recupera el MBR
		var TempMBR Structs.MRB
		// Read object from bin file
		if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
			return
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
						return
					}
					break
				}
			}
		}

		if PartitionStart != -1 {
			fmt.Println("Se encontro la particion")
		} else {
			fmt.Println("Partition not found")
			return
		}

		//Recupera el super bloque
		var tempSuperblock Structs.Superblock
		// Read object from bin file
		if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
			return
		}

		TempStepsPath := strings.Split(path, "/")
		StepsPath := TempStepsPath[1:]

		//Valida Parametro R
		if r {
			//SI no existen carpetas crearlas

			if len(StepsPath) > 1 {
				//no es en la carpeta raiz
				indexInode := UtilitiesInodes.InitSearch(path, file, tempSuperblock)
				if indexInode != -1 {
					//Todas las carpetas existen
					fmt.Println("Las carpetas ya existen")
				} else {
					fmt.Println("Algunas carpetas no existen")
					indiceLastCoincidencia := int32(0)
					creadas := ""
					for _, folder := range StepsPath {
						var tempSuperblock Structs.Superblock
						// Read object from bin file
						if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
							return
						}
						fmt.Println("-------REEEERR")
						fmt.Println(folder)
						creadas += "/" + folder
						//Cada iteracion verificar si existe carpeta o archivo
						indexInode := UtilitiesInodes.InitSearch(creadas, file, tempSuperblock)
						if indexInode != int32(-1) {
							fmt.Println("Esta carpera si eEXISTE")
							indiceLastCoincidencia = indexInode
						} else {
							fmt.Println("Esta carpeta no existe CREARLA")
							//recupero inodo
							fmt.Println("Indice la reci ----------dad", indiceLastCoincidencia)

							var CRRRInode0 Structs.Inode
							// Read object from bin file
							if err := Utilities.ReadObject(file, &CRRRInode0, int64(tempSuperblock.S_inode_start+indiceLastCoincidencia*int32(binary.Size(Structs.Inode{})))); err != nil {
								return
							}
							//Validar los permisos
							uid := fmt.Sprintf("%d", CRRRInode0.I_uid)
							gid := fmt.Sprintf("%d", CRRRInode0.I_gid)
							perm_u := strings.TrimRight(string(CRRRInode0.I_perm[0]), "\x00")
							perm_g := strings.TrimRight(string(CRRRInode0.I_perm[1]), "\x00")
							perm_o := strings.TrimRight(string(CRRRInode0.I_perm[2]), "\x00")

							if global.DeterminarPermisoEscritura(gid, uid, perm_u, perm_g, perm_o) {
								//Tiene permiso de escritura
								CrearCarpeta(CRRRInode0, file, tempSuperblock, TempMBR, PartitionStart, folder, indiceLastCoincidencia)

							} else {
								fmt.Println("Error: El usuario actual no tiene permiso de escritura -MKDIR")
								fmt.Println("====== END MKDIR =====")
								return
							}
							indiceLastCoincidencia = UtilitiesInodes.InitSearch(creadas, file, tempSuperblock)
							fmt.Println("Indice la reciente creadad", indiceLastCoincidencia)
						}
					}
				}
			} else {
				//la nueva carpeta va en la raiz
				indexInode := UtilitiesInodes.InitSearch(path, file, tempSuperblock)
				if indexInode != -1 {
					fmt.Println("Error: La carpeta ya existe")
					fmt.Println("===== END MKDIR ======")
					return
				} else {
					fmt.Println("La direccion no existe")
					//Esta es la importante
				}
			}

		} else {
			//si no existe carpeta mostrar error
			if len(StepsPath) > 1 {
				//no es en la carpeta raiz
				indexInode := UtilitiesInodes.InitSearch(path, file, tempSuperblock)
				if indexInode != -1 {
					fmt.Println("FI xd")
				}
			} else {
				//la nueva carpeta va en la raiz
				indexInode := UtilitiesInodes.InitSearch(path, file, tempSuperblock)
				if indexInode != -1 {
					fmt.Println("Error: La carpeta ya existe")
					fmt.Println("===== END MKDIR ======")
					return
				} else {
					fmt.Println("--La direccion no existe")
					//Esta es la importante
					//Crear una funcion que me de el indice libre
					//una funcion que mande a crear esa carpeta en ese inodo
					var tempSuperblock Structs.Superblock
					// Read object from bin file
					if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
						return
					}
					//--> ya se sabe que va en la raiz
					var CRRRInode0 Structs.Inode
					// Read object from bin file
					if err := Utilities.ReadObject(file, &CRRRInode0, int64(tempSuperblock.S_inode_start)); err != nil {
						return
					}
					//Validar los permisos
					uid := fmt.Sprintf("%d", CRRRInode0.I_uid)
					gid := fmt.Sprintf("%d", CRRRInode0.I_gid)
					perm_u := strings.TrimRight(string(CRRRInode0.I_perm[0]), "\x00")
					perm_g := strings.TrimRight(string(CRRRInode0.I_perm[1]), "\x00")
					perm_o := strings.TrimRight(string(CRRRInode0.I_perm[2]), "\x00")

					if global.DeterminarPermisoEscritura(gid, uid, perm_u, perm_g, perm_o) {
						//Tiene permiso de escritura
						CrearCarpeta(CRRRInode0, file, tempSuperblock, TempMBR, PartitionStart, StepsPath[0], 0)

					} else {
						fmt.Println("Error: El usuario actual no tiene permiso de escritura -MKDIR")
						fmt.Println("====== END MKDIR =====")
						return
					}
				}
			}

		}

	} else {
		fmt.Println("Error: Para este comando debe haber un usuario logeado  -MKDIR")
		fmt.Println("====== END MKDIR =====")
		return
	}
	fmt.Println("===== END mkdir =====")
}

func CrearCarpeta(CRRRInode0 Structs.Inode, file *os.File, tempSuperblock Structs.Superblock, TempMBR Structs.MRB, PartitionStart int, StepsPath string, indiceO int32) {
	creoCarpeta := false
	for i, block := range CRRRInode0.I_block {
		if block != -1 {
			if i < 13 {
				// fmt.Println("entra??")
				var crrFolderBlock Structs.Folderblock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
					return
				}

				for i, folder := range crrFolderBlock.B_content {
					if folder.B_inodo == -1 {

						///Procedo a crear la referencias

						no_inodo := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count
						no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
						fmt.Println("Numero inodo:", no_inodo)
						fmt.Println("Numero bloque:", no_bloque)
						tempSuperblock.S_free_blocks_count -= 1
						tempSuperblock.S_free_inodes_count -= 1
						//Actualiza su bloque apuntador
						crrFolderBlock.B_content[i].B_inodo = no_inodo
						copy(crrFolderBlock.B_content[i].B_name[:], StepsPath)

						//Datos utilis para el nuevo inodo
						usuariID := global.InfoUsuario()
						grupoID := global.InfoGrupo()
						now_date := time.Now()
						formattedDate := now_date.Format("2006-01-02")
						///Crear el inodo
						var newInode Structs.Inode //Inode 0 -> carpeta raiz
						for i := int32(0); i < 15; i++ {
							newInode.I_block[i] = -1
						}

						newInode.I_block[0] = no_bloque

						newInode.I_uid = usuariID
						newInode.I_gid = grupoID
						newInode.I_size = 0
						copy(newInode.I_atime[:], formattedDate)
						copy(newInode.I_ctime[:], formattedDate)
						copy(newInode.I_mtime[:], formattedDate)
						copy(newInode.I_type[:], "0") //Carpetas
						copy(newInode.I_perm[:], "664")

						//Crear el bloque
						var NewFolderblock Structs.Folderblock
						//primero
						NewFolderblock.B_content[0].B_inodo = no_inodo
						copy(NewFolderblock.B_content[0].B_name[:], ".")
						//segundo
						NewFolderblock.B_content[1].B_inodo = 0
						copy(NewFolderblock.B_content[1].B_name[:], "..")
						//tercero
						NewFolderblock.B_content[2].B_inodo = -1
						//Cuarto
						NewFolderblock.B_content[3].B_inodo = -1

						//Escribir los datos en disco
						//Superbloque
						err := Utilities.WriteObject(file, tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start))

						//write Bitman Inodes
						err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_inode_start+no_inodo))
						//write bitman bloques
						err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_block_start+no_bloque))
						//write inodes
						err = Utilities.WriteObject(file, newInode, int64(tempSuperblock.S_inode_start+no_inodo*int32(binary.Size(Structs.Inode{}))))
						//Write bloques
						err = Utilities.WriteObject(file, NewFolderblock, int64(tempSuperblock.S_block_start+(no_bloque)*int32(binary.Size(Structs.Folderblock{}))))
						//Write bloque modificado
						err = Utilities.WriteObject(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{}))))

						if err != nil {
							fmt.Println("Error: ", err)
						}
						fmt.Println("Carpeta creada")
						creoCarpeta = true
						break
					}
				}
				if creoCarpeta {
					break
				}
			}
		} else {
			//la posicion es igual a -1
			fmt.Println("Entro en Caso No cabe en bloque")
			//apuntarlo a su nueva direccion

			no_bloque1 := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
			tempSuperblock.S_free_blocks_count -= 1
			CRRRInode0.I_block[i] = no_bloque1
			//CREAR UN nuevo bloque para su apuntador indirecto
			var NewFolderblock Structs.Folderblock
			for i := int32(0); i < 4; i++ {
				NewFolderblock.B_content[i].B_inodo = -1
			}

			//en una posicion libre apuntar al nuevo inodo que se va a crear

			for i, folder := range NewFolderblock.B_content {
				if folder.B_inodo == -1 {
					no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
					no_inodo := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count
					fmt.Println("Numero inodo:", no_inodo)
					fmt.Println("Numero bloque:", no_bloque)
					tempSuperblock.S_free_blocks_count -= 1
					tempSuperblock.S_free_inodes_count -= 1
					//Indicar a que inodo apunta el bloque
					NewFolderblock.B_content[i].B_inodo = no_inodo
					copy(NewFolderblock.B_content[i].B_name[:], StepsPath)

					//Datos utilis para el nuevo inodo
					usuariID := global.InfoUsuario()
					grupoID := global.InfoGrupo()
					now_date := time.Now()
					formattedDate := now_date.Format("2006-01-02")
					///Crear el inodo
					var newInode Structs.Inode //Inode 0 -> carpeta raiz
					for i := int32(0); i < 15; i++ {
						newInode.I_block[i] = -1
					}

					newInode.I_block[0] = no_bloque

					newInode.I_uid = usuariID
					newInode.I_gid = grupoID
					newInode.I_size = 0
					copy(newInode.I_atime[:], formattedDate)
					copy(newInode.I_ctime[:], formattedDate)
					copy(newInode.I_mtime[:], formattedDate)
					copy(newInode.I_type[:], "0") //Carpetas
					copy(newInode.I_perm[:], "664")
					//Crear el bloque del nuevo inodo y apuntarlos
					//Crear el bloque
					var BloqueNewInodo Structs.Folderblock
					//primero
					BloqueNewInodo.B_content[0].B_inodo = no_inodo
					copy(BloqueNewInodo.B_content[0].B_name[:], ".")
					//segundo
					BloqueNewInodo.B_content[1].B_inodo = 0
					copy(BloqueNewInodo.B_content[1].B_name[:], "..")
					//tercero
					BloqueNewInodo.B_content[2].B_inodo = -1
					//Cuarto
					BloqueNewInodo.B_content[3].B_inodo = -1
					//Escribir en el disco
					//Superbloque
					err := Utilities.WriteObject(file, tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start))

					//write Bitman Inodes
					err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_inode_start+no_inodo))
					//write bitman bloques
					err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_block_start+no_bloque))
					//write inodes
					err = Utilities.WriteObject(file, newInode, int64(tempSuperblock.S_inode_start+no_inodo*int32(binary.Size(Structs.Inode{}))))
					err = Utilities.WriteObject(file, CRRRInode0, int64(tempSuperblock.S_inode_start+indiceO*int32(binary.Size(Structs.Inode{}))))
					//Write bloques

					err = Utilities.WriteObject(file, NewFolderblock, int64(tempSuperblock.S_block_start+(no_bloque1)*int32(binary.Size(Structs.Folderblock{}))))
					//Write bloque modificado
					err = Utilities.WriteObject(file, BloqueNewInodo, int64(tempSuperblock.S_block_start+no_bloque*int32(binary.Size(Structs.Folderblock{}))))

					if err != nil {
						fmt.Println("Error: ", err)
					}
					fmt.Println("Carpeta creada")
					creoCarpeta = true
					break //se escribe en el nuevo inodo
				}
			}
			if creoCarpeta {
				break
			}

		}
	}
}
