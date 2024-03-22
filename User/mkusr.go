package user

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Mkusr(userNew string, passNew string, grupoValidar string) {
	//Validar que haya iniciado sesion
	fmt.Println("===== Start Mkusr =====")
	if global.LoginValidacion() {

		PartitionStart := -1

		//Que sea el usuario root
		if global.ValidaUsuario("root") {
			//Esta como root
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
			indexInode := initSearch(file, tempSuperblock, 0, "users.txt")
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
				data := ""
				for i := 0; i < len(crrInode.I_block); i++ {
					if crrInode.I_block[i] != -1 {
						var Fileblock Structs.Fileblock
						if err := Utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[i]*int32(binary.Size(Structs.Fileblock{})))); err != nil {
							return
						}
						data += strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
					}
				}

				fmt.Println("\nque trae data:", data+"\n")
				// Dividir la cadena en líneas
				data = strings.TrimSuffix(data, "\n")
				lines := strings.Split(data, "\n")
				// tamanio := len(lines)

				ExisteUser := false
				// Iterar a través de las líneas
				// count := 0
				//Lista aux para incorporar usuarios
				UsuariosPorGrupo := make(map[string][]Structs.UserTXT)

				for _, linea := range lines {
					fields := strings.Split(linea, ",")
					// Obtener los datos del usuario
					if len(fields) == 5 {

						Id := fields[0]
						Tipo := fields[1]
						Grupo := fields[2]
						User := fields[3]
						Pass := fields[4]

						// Verificar si el grupo existe
						_, ok := UsuariosPorGrupo[Grupo]
						if !ok {
							// Crear una nueva entrada para el grupo
							UsuariosPorGrupo[Grupo] = make([]Structs.UserTXT, 0)
						}

						// Verificar si el usuario ya existe en el grupo
						usuarioExiste := false
						for _, usuarioExistente := range UsuariosPorGrupo[Grupo] {
							if usuarioExistente.User == User {
								usuarioExiste = true
								break
							}
						}

						// Si el usuario no existe en el grupo, agregarlo
						if !usuarioExiste {
							usuario := Structs.UserTXT{Id, Tipo, Grupo, User, Pass}
							UsuariosPorGrupo[Grupo] = append(UsuariosPorGrupo[Grupo], usuario)
						}

					} else {
						// count += 1
						// if count == tamanio-1 {
						// 	break
						// }
						grupo := fields[2]

						_, ok := UsuariosPorGrupo[grupo]
						if !ok {
							// Crear una nueva entrada para el grupo
							UsuariosPorGrupo[grupo] = make([]Structs.UserTXT, 0)
						}
					}

				}

				//verificar si en esos grupo existe el nuevo
				_, ok1 := UsuariosPorGrupo[grupoValidar]
				if ok1 { //existe grupo a ingresar
					// Verificar si el usuario ya existe en el grupo

					usuarioExiste := false
					for _, usuarioExistente := range UsuariosPorGrupo[grupoValidar] {
						if usuarioExistente.User == userNew {
							usuarioExiste = true
							break
						}
					}

					// Si el usuario no existe en el grupo, agregarlo
					if !usuarioExiste {
						usuario := Structs.UserTXT{"0", "U", grupoValidar, userNew, passNew}
						UsuariosPorGrupo[grupoValidar] = append(UsuariosPorGrupo[grupoValidar], usuario)
					} else {
						ExisteUser = true
						fmt.Println("Error: Usuario ya existe, -mkusr")
						fmt.Println("===== END mkgrp =====")
						return
					}

				} else {

					fmt.Println("Error: grupo NO existe")
					fmt.Println("===== END mkgrp =====")
					return

				}
				// if !ExisteGrupo {
				// 	fmt.Println("Error: No existe grupo al que se quiere insertar usuario")
				// 	fmt.Println("===== END Mkusr ======")
				// 	return
				// }

				if !ExisteUser { //No existe el usuario
					newData := ""
					countGrupo := 0
					countUser := 0
					for grupo, usuarios := range UsuariosPorGrupo {
						countGrupo += 1
						correlativo := strconv.Itoa(countGrupo)
						newData += correlativo + ",G," + grupo + "\n"
						//cada usuario
						for _, usuario := range usuarios {
							countUser += 1
							correlativo1 := strconv.Itoa(countUser)
							newData += correlativo1 + ",U," + grupo + "," + usuario.User + "," + usuario.Pass + "\n"
						}
					}
					fmt.Println("##### new data antes de insr\n", newData)
					newData = strings.TrimSuffix(newData, "\n")
					fmt.Println("##### new data apposr\n", newData)
					//Insertarlo en la estructura
					if len(newData) > 64 {

						//se necesitara mas de 1 Fileblock
						numFileblock := len(newData) / 64
						if numFileblock > 11 {
							//Se debera implementar los apuntadores indirectos (pendiente)
						} else {
							//solo con apuntadores directos
							fmt.Println(">64")
							fmt.Println("deberia iterar:", numFileblock)
							diferencia := 64 - (len(newData) % 64)
							cadenaRellena := newData + strings.Repeat(" ", diferencia)
							for i := 0; i <= numFileblock; i++ {

								//resta bloque libres

								//se apunta inodo
								no_bloque := tempSuperblock.S_blocks_count - tempSuperblock.S_free_blocks_count
								fmt.Println("no bloque siguiente:", no_bloque)
								crrInode.I_block[i] = no_bloque
								//Crear bloque

								var newFileblock Structs.Fileblock

								if len(cadenaRellena)%64 == 0 {
									fmt.Println("------- resultado\n", cadenaRellena[:64])
									copy(newFileblock.B_content[:], cadenaRellena[:64])
									cadenaRellena = strings.TrimPrefix(cadenaRellena, cadenaRellena[:64])
								} else {
									fmt.Println("No es modulo 64")
								}

								//Escribir la info en el archivo
								tempSuperblock.S_free_blocks_count -= 1
								// write superblock
								err := Utilities.WriteObject(file, tempSuperblock, int64(PartitionStart))

								// write bitmap blocks
								err = Utilities.WriteObject(file, byte(1), int64(tempSuperblock.S_bm_block_start+(no_bloque-1)))
								//Write  inode
								crrInode.I_block[i] = (no_bloque - 1)
								err = Utilities.WriteObject(file, crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))) //Inode 1

								// write blocks
								err = Utilities.WriteObject(file, newFileblock, int64(tempSuperblock.S_block_start+(int32(no_bloque-1)*int32(binary.Size(Structs.Fileblock{})))))

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
						fmt.Print("se escribe:", newData)
						// fmt.Print("El bloque es:,", no_bloque-1)
						copy(Fileblock.B_content[:], newData)

						// write blocks
						err = Utilities.WriteObject(file, Fileblock, int64(tempSuperblock.S_block_start+(int32(binary.Size(Structs.Fileblock{})))))

						fmt.Println("Se agrego usuario")
						if err != nil {
							fmt.Println("Error: Escritura de bloques -mkgrp", err)
						}
					}
				}

				//Buscar el archivo user.txt
				//Verificar que el numero grupo no exista
				//Verificar el correlativo
				//insertar de nuevo el archivo user.txt -> comenzar a validar los 64 bytes de cada bloque

			} else {
				fmt.Println("Error: No se encontro archivo user.txt -mkusr")
			}
		} else {
			fmt.Println("Error: debe iniciar sesion como root -mkgrp")
		}
	} else {
		fmt.Println("Erro: Debe iniciar sesion para usar este comando -mkgrp")
	}
	fmt.Println("===== END Mkusr =====")
}
