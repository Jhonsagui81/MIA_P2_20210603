package user

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"encoding/binary"
	"fmt"
	"strings"
)

func Rmusr(userNew string) {
	fmt.Println("===== Start Rmusr ======")
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

				fmt.Println("que trae data:", data+"\n")
				data = strings.TrimSuffix(data, "\n")
				// Dividir la cadena en líneas
				lines := strings.Split(data, "\n")

				ExisteUser := false

				for i, linea := range lines {
					fields := strings.Split(linea, ",")
					// Obtener los datos del usuario
					if len(fields) == 5 {
						id := fields[0]
						grupo := fields[2]
						user := fields[3]
						pass := fields[4]

						if user == userNew {
							//Se encontro el usuario a elimincar
							if id != "0" {
								ExisteUser = true
								lines[i] = "0,U," + grupo + "," + user + "," + pass
							} else {
								fmt.Println("Error: usuario ya eliminado")
								fmt.Println("===== END RMUSR =====")
								return
							}
						}
					}
				}

				if ExisteUser {
					//Si existe el grupo
					//Reconstruir cadena
					newData := ""
					for _, line := range lines {
						newData += line + "\n"
					}

					newData = strings.TrimSuffix(newData, "\n")
					//Insertarlo en la estructura
					fmt.Print("se escribe:", newData)

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

								//Crear bloque

								var newFileblock Structs.Fileblock

								if len(cadenaRellena)%64 == 0 {
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
						// fmt.Print("El bloque es:,", no_bloque-1)
						copy(Fileblock.B_content[:], newData)

						// write blocks
						err = Utilities.WriteObject(file, Fileblock, int64(tempSuperblock.S_block_start+(int32(binary.Size(Structs.Fileblock{})))))

						fmt.Println("Escritura correcta")
						if err != nil {
							fmt.Println("Error: Escritura de bloques -mkgrp", err)
						}

					}

					fmt.Println("Usuario Eliminado")
				} else {
					fmt.Println("Error: El usuario que pretende eliminar no existe")
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

	fmt.Println("===== End Rmusr =====")
}
