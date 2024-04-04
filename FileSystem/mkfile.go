package filesystem

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func Mkfile(path string, r bool, size int, cont string) {
	fmt.Println("===== START MKFILE =====")
	fmt.Println("path:", path)
	fmt.Println("r:", r)
	fmt.Println("size:", size) //bytes
	fmt.Println("cont:", cont)

	//Validaciones
	//size
	if size < 0 {
		fmt.Println("Error: Size debe ser mayor a 0		MKFILE")
		fmt.Println("===== END MKFILE =====")
		return
	}
	//count
	data := ""
	if cont != "" { //Contengo una ruta
		//validar que archivo exista
		_, err := os.Stat(cont)
		if err == nil {
			//Archivo Existe
			contenido, err := ioutil.ReadFile(cont)
			if err != nil {
				fmt.Println("Error: No se puede leer archivo count		MKFILE")
				fmt.Println("===== END MKFILE =====")
				return
			} else {
				data = string(contenido)
			}
		} else if os.IsNotExist(err) {
			fmt.Println("Error: Ruta cont no existe		MKFILE")
			fmt.Println("===== END MKFILE =====")
			return
		} else {
			fmt.Println("Error: Ruta cont no existe		MKFILE")
			fmt.Println("===== END MKFILE =====")
			return
		}
		//abrirlo y extarer su contenido
	} else {
		data = ContenidoArchivo(size)
	}

	fmt.Println("data:", data)
	if global.LoginValidacion() {
		//Valida si hay usuario logeado
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
		//Ibtener directorio de carpetas y nombre del archivo
		TempStepsPath := strings.Split(path, "/")
		StepsPath := TempStepsPath[1:]
		//obtener el nombre del archivo
		lastIndex := len(StepsPath) - 1
		nameFile := StepsPath[lastIndex]
		fmt.Println("nombreArchivo por crear:", nameFile)
		//Eliminar el nombre del archivo del slide
		StepsPath = StepsPath[:lastIndex]

		//Crear direcctorio a buscar
		direccion := ""
		for _, step := range StepsPath {
			direccion += "/" + step
		}

		fmt.Println("Directorio a buscar:", direccion)
		if r {
			//Si no existen carpetas crearlas
			indexInode := UtilitiesInodes.InitSearch(direccion, file, tempSuperblock)
			if indexInode != -1 {
				fmt.Println("Las carpetas ya existe")
				//Si existe la direccion -> Crear el archivo
				var tempSuperblock Structs.Superblock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
					return
				}

				//Inodo de la carpeta contenedora
				var CRRRInode0 Structs.Inode
				// Read object from bin file
				if err := Utilities.ReadObject(file, &CRRRInode0, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
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
					fmt.Println("SI tiene permiso para crear el archivo")
					CrearArchivo(CRRRInode0, file, tempSuperblock, TempMBR, PartitionStart, nameFile, indexInode, int32(size), data)

				} else {
					fmt.Println("Error: El usuario actual no tiene permiso de escritura -MKDIR")
					fmt.Println("====== END MKDIR =====")
					return
				}

			} else {
				fmt.Println("Algunas carpetas no Existen")
				TempStepsPath := strings.Split(path, "/")
				StepsPath := TempStepsPath[1 : len(TempStepsPath)-1]

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
				//Ya se crearon las carpetas -> Escribir el archivo
				indexInode := UtilitiesInodes.InitSearch(direccion, file, tempSuperblock)
				if indexInode != -1 {
					fmt.Println("Las carpetas ya existe")
					//Si existe la direccion -> Crear el archivo
					var tempSuperblock Structs.Superblock
					// Read object from bin file
					if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
						return
					}

					//Inodo de la carpeta contenedora
					var CRRRInode0 Structs.Inode
					// Read object from bin file
					if err := Utilities.ReadObject(file, &CRRRInode0, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
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
						fmt.Println("SI tiene permiso para crear el archivo")
						CrearArchivo(CRRRInode0, file, tempSuperblock, TempMBR, PartitionStart, nameFile, indexInode, int32(size), data)

					} else {
						fmt.Println("Error: El usuario actual no tiene permiso de escritura -MKDIR")
						fmt.Println("====== END MKDIR =====")
						return
					}

				}

			}

		} else {
			//Si no existen carpetas mostrar error
			//Recortar path, para buscar la carpeta contenedor -> si existe crearle un archivo

			//Buscar la direccion del nuevo archivo
			indexInode := UtilitiesInodes.InitSearch(direccion, file, tempSuperblock)
			if indexInode != -1 {
				//Si existe la direccion -> Crear el archivo
				var tempSuperblock Structs.Superblock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
					return
				}

				//Inodo de la carpeta contenedora
				var CRRRInode0 Structs.Inode
				// Read object from bin file
				if err := Utilities.ReadObject(file, &CRRRInode0, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
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
					fmt.Println("SI tiene permiso para crear el archivo")
					CrearArchivo(CRRRInode0, file, tempSuperblock, TempMBR, PartitionStart, nameFile, indexInode, int32(size), data)

				} else {
					fmt.Println("Error: El usuario actual no tiene permiso de escritura -MKDIR")
					fmt.Println("====== END MKDIR =====")
					return
				}
			} else {
				fmt.Println("Error: La ruta para crear archivo no existe. 	MKFILE")
				fmt.Println("===== END MKFILE =====")
				return
			}

		}
	} else {
		fmt.Println("Error: NO se inicio sesion   MKFILE")
		fmt.Println("===== END MKFILE =====")
		return
	}

	fmt.Println("===== END MKFILE =====")
}

func ContenidoArchivo(Size int) string {
	content := "0123456789"
	repeticiones := Size / len(content)

	if Size%len(content) > 0 {
		repeticiones++
	}

	cadenaRellena := strings.Repeat(content, repeticiones)
	return cadenaRellena[:Size]
}

func CrearArchivo(CRRRInode0 Structs.Inode, file *os.File, tempSuperblock Structs.Superblock, TempMBR Structs.MRB, PartitionStart int, nameFile string, indiceO int32, size int32, data string) {
	creoArchivo := false
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
						fmt.Println("Numero inodo:", no_inodo)
						tempSuperblock.S_free_inodes_count -= 1
						//Actualiza su bloque apuntador
						crrFolderBlock.B_content[i].B_inodo = no_inodo
						copy(crrFolderBlock.B_content[i].B_name[:], nameFile)

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

						newInode.I_uid = usuariID
						newInode.I_gid = grupoID
						newInode.I_size = size
						copy(newInode.I_atime[:], formattedDate)
						copy(newInode.I_ctime[:], formattedDate)
						copy(newInode.I_mtime[:], formattedDate)
						copy(newInode.I_type[:], "1") //Carpetas
						copy(newInode.I_perm[:], "664")

						//Rellenar el archivo - de ser mayor a 64  crear mas bloques
						//Divide la cadena en trozos de 64 bytes
						trozos := make([]string, 0)
						inicio := 0

						for i := 0; i < len(data); i++ {
							if i-inicio >= 64 {
								trozos = append(trozos, data[inicio:i])
								inicio = i
							}
						}

						if len(data)-inicio > 0 {
							trozos = append(trozos, data[inicio:])
						}

						for i, trozo := range trozos {
							//Crear el bloque
							fmt.Println("no bloque siguiente:", i)

							var NewFileBlock Structs.Fileblock
							copy(NewFileBlock.B_content[:], trozo)

							//Crear apuntadores
							no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
							newInode.I_block[i] = no_bloque
							//Restar un bloque
							tempSuperblock.S_free_blocks_count -= 1

							//Escribir datos
							//Superbloque
							err := Utilities.WriteObject(file, tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start))
							//bitman inodes
							err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_inode_start+no_inodo))
							//write bitman bloques
							err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_block_start+no_bloque))
							//write inodes
							err = Utilities.WriteObject(file, newInode, int64(tempSuperblock.S_inode_start+no_inodo*int32(binary.Size(Structs.Inode{}))))
							//Write bloques
							err = Utilities.WriteObject(file, NewFileBlock, int64(tempSuperblock.S_block_start+(no_bloque)*int32(binary.Size(Structs.Folderblock{}))))
							//Write bloque modificado
							err = Utilities.WriteObject(file, crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{}))))
							if err != nil {
								fmt.Println("Error: ", err)
							}

						}
						fmt.Println("archivo creada")
						creoArchivo = true
						break
					}
				}
				if creoArchivo {
					break
				}
			}
		} else {
			//la posicion es igual a -1
			fmt.Println("Entro en Caso No cabe en bloque")
			//apuntarlo a su nueva direccion

			no_bloque1 := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
			CRRRInode0.I_block[i] = no_bloque1
			tempSuperblock.S_free_blocks_count -= 1
			//CREAR UN nuevo bloque para su apuntador indirecto
			var NewFolderblock Structs.Folderblock
			for i := int32(0); i < 4; i++ {
				NewFolderblock.B_content[i].B_inodo = -1
			}

			//en una posicion libre apuntar al nuevo inodo que se va a crear

			for i, folder := range NewFolderblock.B_content {
				if folder.B_inodo == -1 {
					no_inodo := tempSuperblock.S_inodes_count - tempSuperblock.S_free_inodes_count
					fmt.Println("Numero inodo:", no_inodo)
					tempSuperblock.S_free_inodes_count -= 1
					//Actualiza su bloque apuntador
					NewFolderblock.B_content[i].B_inodo = no_inodo
					copy(NewFolderblock.B_content[i].B_name[:], nameFile)

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

					newInode.I_uid = usuariID
					newInode.I_gid = grupoID
					newInode.I_size = size
					copy(newInode.I_atime[:], formattedDate)
					copy(newInode.I_ctime[:], formattedDate)
					copy(newInode.I_mtime[:], formattedDate)
					copy(newInode.I_type[:], "1") //Carpetas
					copy(newInode.I_perm[:], "664")

					//Rellenar el archivo - de ser mayor a 64  crear mas bloques
					//Divide la cadena en trozos de 64 bytes
					trozos := make([]string, 0)
					inicio := 0

					for i := 0; i < len(data); i++ {
						if i-inicio >= 64 {
							trozos = append(trozos, data[inicio:i])
							inicio = i
						}
					}

					if len(data)-inicio > 0 {
						trozos = append(trozos, data[inicio:])
					}

					for i, trozo := range trozos {
						//Crear el bloque
						var NewFileBlock Structs.Fileblock
						copy(NewFileBlock.B_content[:], trozo)

						//Crear apuntadores
						no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
						newInode.I_block[i] = no_bloque
						//Restar un bloque
						tempSuperblock.S_free_blocks_count -= 1

						//Escribir datos
						//Superbloque
						err := Utilities.WriteObject(file, tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start))
						//bitman inodes
						err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_inode_start+no_inodo))
						//write bitman bloques
						err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_block_start+no_bloque))
						//write inodes
						err = Utilities.WriteObject(file, newInode, int64(tempSuperblock.S_inode_start+no_inodo*int32(binary.Size(Structs.Inode{}))))
						err = Utilities.WriteObject(file, CRRRInode0, int64(tempSuperblock.S_inode_start+indiceO*int32(binary.Size(Structs.Inode{}))))

						//Write bloques
						err = Utilities.WriteObject(file, NewFileBlock, int64(tempSuperblock.S_block_start+(no_bloque)*int32(binary.Size(Structs.Folderblock{}))))
						//Write bloque modificado
						err = Utilities.WriteObject(file, NewFolderblock, int64(tempSuperblock.S_block_start+no_bloque1*int32(binary.Size(Structs.Folderblock{}))))
						if err != nil {
							fmt.Println("Error: ", err)
						}

					}
					fmt.Println("Archivo creada")
					creoArchivo = true
					break
				}
			}
			if creoArchivo {
				break
			}

		}
	}
}
