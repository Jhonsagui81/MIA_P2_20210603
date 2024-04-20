package user

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"fmt"
	"strings"
)

func Rmgrp(name string) {
	//
	fmt.Println("===== Start Rmgrp ===== ")
	if global.LoginValidacion() {
		PartitionStart := -1
		if global.ValidaUsuario("root") {
			//Estraigo el id del disco
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
							PartitionStart = int(TempMBR.Partitions[i].Start)
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

			// initSearch /users.txt -> regresa no Inodo
			// initSearch -> 1
			//Mando a bucar el archivo user.txt
			indexInode := UtilitiesInodes.InitSearch("/users.txt", file, tempSuperblock)
			//file, superbloque, posicion de inodo inicial, archivo/carpeta a buscar

			//recupero el indo en indexInode -> concateno todos sus bloques porque este inode es tipo archivo
			//(aunque hay que validar para los bloque en i = 13,14,15 que son indirectos)
			if indexInode != -1 {
				//Recupero el Inodo que indica indexInode
				var crrInode Structs.Inode
				if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil {
					return
				}

				// getInodeFileData -> Iterate the I_Block n concat the data
				//Falta crear funcion para iterar todos los bloques del inodo archivo y concatenar
				//Validar apuntadores indirectos

				//Repuero el bloque del archivo, para obtener la data
				data := UtilitiesInodes.GetInodeFileData(crrInode, file, tempSuperblock)

				fmt.Println("que trae data:", data+"\n")
				// Dividir la cadena en líneas
				data = strings.TrimSuffix(data, "\n")
				lines := strings.Split(data, "\n")

				// tamanio := len(lines)
				ExisteGrupo := false
				// count := 0
				// tamanio := len(lines)

				for i, line := range lines {
					fields := strings.Split(line, ",")
					// Obtener los datos del usuario
					id := fields[0]
					Grupo := fields[2]
					if len(fields) == 3 {
						if Grupo == name {
							//Se encontro grupo a eliminar
							if id != "0" {
								lines[i] = "0,G," + Grupo
								ExisteGrupo = true
							} else {
								fmt.Println("Error: Este grupo ya fue eliminado")
								fmt.Println("===== END Rmgrp =====")
								return
							}
						}
					}
					//Si existe grupo es true todos los usuarios de ese grupo hacerlos id = 0
					if ExisteGrupo {
						if len(fields) == 5 {
							user := fields[3]
							pass := fields[4]
							if Grupo == name {
								lines[i] = "0,U," + Grupo + "," + user + "," + pass
							}
						}
					}

				}

				if ExisteGrupo {
					//Si existe el grupo
					//Buscar la fila donde esta el grupo a eliminar

					newData := ""
					for _, line := range lines {
						newData += line + "\n"
					}

					newData = strings.TrimSuffix(newData, "\n")
					fmt.Print("se escribe:", newData)

					//crear slide de cadena
					trozos := make([]string, 0)
					inicio := 0

					for i := 0; i < len(newData); i++ {
						if i-inicio >= 64 {
							trozos = append(trozos, newData[inicio:i])
							inicio = i
						}
					}

					if len(newData)-inicio > 0 {
						trozos = append(trozos, newData[inicio:])
					}

					//Insertarlo en la estructura
					if len(newData) > 64 {

						//se necesitara mas de 1 Fileblock
						numFileblock := len(newData) / 64
						if numFileblock > 12 {
							//Se debera implementar los apuntadores indirectos (pendiente)
						} else {
							//solo con apuntadores directos
							fmt.Println(">64")
							fmt.Println("deberia iterar:", numFileblock)
							// diferencia := 64 - (len(newData) % 64)
							no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
							if no_bloque < int32(numFileblock+1) {
								tempSuperblock.S_free_blocks_count -= 1
							}
							// cadenaRellena := newData + strings.Repeat(" ", diferencia)
							for i, trozo := range trozos {

								//resta bloque libres

								//se apunta inodo
								fmt.Println("no bloque siguiente:", i+1)

								//Crear bloque

								var newFileblock Structs.Fileblock

								// if len(cadenaRellena)%64 == 0 {
								// 	copy(newFileblock.B_content[:], cadenaRellena[:64])
								// 	cadenaRellena = strings.TrimPrefix(cadenaRellena, cadenaRellena[:64])
								// } else {
								// 	fmt.Println("No es modulo 64")
								// }
								copy(newFileblock.B_content[:], trozo)

								//Escribir la info en el archivo
								// write superblock
								err := Utilities.WriteObject(file, tempSuperblock, int64(PartitionStart))

								// write bitmap blocks
								err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_block_start+(no_bloque-1)))

								//Write  inode

								crrInode.I_block[i] = (int32(i) + 1)

								err = Utilities.WriteObject(file, crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))) //Inode 1

								// write blocks
								err = Utilities.WriteObject(file, newFileblock, int64(tempSuperblock.S_block_start+(int32(i+1)*int32(binary.Size(Structs.Fileblock{})))))

								if err != nil {
									fmt.Println("Error: Escritura de bloques -mkgrp", err)
								}
								fmt.Println("Escritura correcta")
							}
						}
					} else {
						fmt.Println("<64")
						var Fileblock Structs.Fileblock
						// no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
						//Se inserta en el mismo fileblock
						// fmt.Print("El bloque es:,", no_bloque-1)
						copy(Fileblock.B_content[:], newData)

						// write blocks
						err = Utilities.WriteObject(file, Fileblock, int64(tempSuperblock.S_block_start+(int32(binary.Size(Structs.Fileblock{})))))

						fmt.Println("Escritura correcta")
						if err != nil {
							fmt.Println("Error: Escritura de bloques -mkgrp", err)
						}

					}

					fmt.Println("Grupo Eliminado")
				} else {
					fmt.Println("Error: El grupo que pretende eliminar no existe")
				}
			} else {
				fmt.Println("Error: No se encontro archivo user.txt -rmgrp")
			}
		} else {
			fmt.Println("Error: debe iniciar sesion como root -rmgrp")
		}
	} else {
		fmt.Println("Erro: Debe iniciar sesion para usar este comando -rmgrp")
	}
	fmt.Println("===== END rmgrp ======")
}
