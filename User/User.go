package user

import (
	global "PROYECTO1_MIA/Global"
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"PROYECTO1_MIA/UtilitiesInodes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Login(user string, pass string, id string) {
	fmt.Println("===== Start Login ======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id", id)

	logeado := global.LoginValidacion()
	if logeado {
		fmt.Println("Debe deslogearse para iniciar sesion")
		fmt.Println("===== END login ======")
		return
	}

	driveletter := string(id[0])
	partition := string(id[1])
	partitionStart := 0

	// Open bin file
	filepath := "./MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: El disco no existe -> login")
		return
	}

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
					partitionStart = int(TempMBR.Partitions[i].Start)
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
		data = strings.TrimSuffix(data, "\n")
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		// Iterar a través de las líneas
		var Uid string
		var grupoLogeado string
		existeUsuario := false
		for _, line := range lines {
			// Imprimir cada línea
			fields := strings.Split(line, ",")

			if len(fields) != 3 {
				UID := fields[0]
				grupo := fields[2]
				username := fields[3]
				password := fields[4]

				if username == user && password == pass {
					existeUsuario = true
					Uid = UID
					grupoLogeado = grupo
					break
				}
			}
		}

		if existeUsuario {
			for _, line := range lines {
				// Imprimir cada línea
				fields := strings.Split(line, ",")

				if len(fields) == 3 {
					GIDActual := fields[0]
					grupoActual := fields[2]

					if grupoActual == grupoLogeado {
						global.Logear(user, pass, Uid, GIDActual, driveletter, partition, partitionStart, id)
						fmt.Println("logeo exitoso...")
					}
				}
			}
		} else {
			fmt.Println("Error: Usuario no registrado")
		}

	} else {
		fmt.Println("Error: No se encontro el archivo, users.txt")
	}

	// Close bin file
	defer file.Close()

	fmt.Println("====== END login ======")
}

func initSearch(file *os.File, tempSuperblock Structs.Superblock, posicion int, dato string) int32 {
	indexInode := int32(-1)

	//Recupero el inodo
	var TempInodeRaiz Structs.Inode
	if err := Utilities.ReadObject(file, &TempInodeRaiz, int64(tempSuperblock.S_inode_start+(int32(posicion)*int32(binary.Size(Structs.Inode{}))))); err != nil {
		return indexInode
	}

	//Itera sus apuntadores I_block
	for i := int32(0); i < 15; i++ {
		if TempInodeRaiz.I_block[i] != -1 {
			tipo_inodo := strings.TrimRight(string(TempInodeRaiz.I_type[:]), "\x00")
			if tipo_inodo == "0" {
				posicion_block := TempInodeRaiz.I_block[i]
				//Recupero el bloque que indique el inodo.I_block en posicion i
				var TempBlock Structs.Folderblock
				if err := Utilities.ReadObject(file, &TempBlock, int64(tempSuperblock.S_block_start+(int32(posicion_block)*int32(binary.Size(Structs.Folderblock{}))))); err != nil {
					return indexInode
				}

				//Itero los cuatro apuntadores del bloque de carpetas
				for i := int32(0); i < 4; i++ {
					B_name := strings.TrimRight(string(TempBlock.B_content[i].B_name[:]), "\x00")
					//comparar carpeta/archivo buscado con el nombre de apuntador
					if B_name == dato {
						//Si coinciden - retornar el inodo
						indexInode = TempBlock.B_content[i].B_inodo
						break
					}
				}
			} //else { Es un inodo para archivo, no puedo buscar carpetas}
		}
	}
	return indexInode
}

func SearchFileOrFolder(file *os.File, tempSuperblock Structs.Superblock, posicion int, dato string) int32 {
	indexInode := int32(-1)

	//Recupero el inodo
	var TempInodeRaiz Structs.Inode
	if err := Utilities.ReadObject(file, &TempInodeRaiz, int64(tempSuperblock.S_inode_start+(int32(posicion)*int32(binary.Size(Structs.Inode{}))))); err != nil {
		return indexInode
	}

	//Itera sus apuntadores I_block
	for i := int32(0); i < 15; i++ {
		if TempInodeRaiz.I_block[i] != -1 {
			tipo_inodo := strings.TrimRight(string(TempInodeRaiz.I_type[:]), "\x00")
			if tipo_inodo == "0" {
				posicion_block := TempInodeRaiz.I_block[i]
				//Recupero el bloque que indique el inodo.I_block en posicion i
				var TempBlock Structs.Folderblock
				if err := Utilities.ReadObject(file, &TempBlock, int64(tempSuperblock.S_block_start+(int32(posicion_block)*int32(binary.Size(Structs.Folderblock{}))))); err != nil {
					return indexInode
				}

				//Itero los cuatro apuntadores del bloque de carpetas
				for i := int32(0); i < 4; i++ {
					B_name := strings.TrimRight(string(TempBlock.B_content[i].B_name[:]), "\x00")
					//comparar carpeta/archivo buscado con el nombre de apuntador
					if B_name == dato {
						//Si coinciden - retornar el inodo
						indexInode = TempBlock.B_content[i].B_inodo
						break
					} else {
						//buscar la siguiente posicion -> llamar a la misma funcion
						//if B_name ==  "." or ".." continue
						//Para no ciclar la recursividad
						initSearch(file, tempSuperblock, int(TempBlock.B_content[i].B_inodo), dato)
					}
				}
			} //else { Es un inodo para archivo, no puedo buscar carpetas}
		}
	}
	return indexInode
}
