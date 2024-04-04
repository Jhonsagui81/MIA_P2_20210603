package reportes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// Cuentan los inodos y bloques globales
var indexI int32 //---- deberian ser globales
var indexB int32

func TreeSystem(disco string, id string) string {
	text := ""
	indexI = 0
	indexB = 0
	filepath := "./MIA/P1/" + disco + ".dsk"
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		return ""
	}
	defer file.Close()
	//recoger  informacion del disco
	var TempMBR Structs.MRB
	// Read object from bin file
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return ""
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
					return ""
				}
				break
			}
		}
	}

	if PartitionStart != -1 {
		fmt.Println("Se encontro la particion")
	} else {
		fmt.Println("Partition not found")
		return ""
	}

	//Recupera el super bloque
	var tempSuperblock Structs.Superblock
	// Read object from bin file
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[PartitionStart].Start)); err != nil {
		return ""
	}

	var Inode0 Structs.Inode
	// Read object from bin file
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return ""
	}

	//Crea las tablas
	text += RecorridoSystem(Inode0, file, tempSuperblock)
	//Crea las conexciones
	indexI = 0
	indexB = 0
	text += ConexionSystem(Inode0, file, tempSuperblock)
	return text
}

func ConexionSystem(Inode Structs.Inode, file *os.File, tempSuperbloque Structs.Superblock) string {
	text := ""
	//Reiniciar apuntadores

	//para concatenar
	indexInodo := fmt.Sprintf("%d", indexI)
	//
	//Para bloques
	indexBloque := fmt.Sprintf("%d", indexB)

	tipo_doc := strings.TrimRight(string(Inode.I_type[:]), "\x00")
	if tipo_doc == "0" {
		//Iniciar busque de sus apuntadores
		for i, block := range Inode.I_block {
			if block != -1 {
				text += "\tInodo" + indexInodo + ":"
				indexBloque = fmt.Sprintf("%d", block)
				indexB = block
				if i < 12 { //Pendiente verificar el indice
					//procedimiento para recorrer bloque de carpeta
					apuntaGrafo := fmt.Sprintf("%d", i) //Para hacer conexciones con el grafo
					// indiceApunta := fmt.Sprintf("%d", block) //indice al que apunta la posicion actual
					text += "P" + apuntaGrafo + " -> "
					var crrFolderBlock Structs.Folderblock
					if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperbloque.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
						return ""
					}
					//Construir grafo del bloque
					text += "Bloque" + indexBloque + ";\n" //Para apuntadores del grafo
					//para llamar al siguiente inodo
					for i, folder1 := range crrFolderBlock.B_content {
						nombre_docu := strings.TrimRight(string(folder1.B_name[:]), "\x00")
						if nombre_docu != "." && nombre_docu != ".." && nombre_docu != "" {
							//Descartamos primeros dos inodos
							var NextInode Structs.Inode
							// Read object from bin file
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperbloque.S_inode_start+folder1.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return ""
							}
							indiceBoque := fmt.Sprintf("%d", i)
							text += "\tBloque" + indexBloque + ":P" + indiceBoque + " -> "
							indexI = folder1.B_inodo
							indexInodo := fmt.Sprintf("%d", folder1.B_inodo)
							text += "Inodo" + indexInodo + ";\n"
							text += ConexionSystem(NextInode, file, tempSuperbloque)

						}
					}

				}
			}
		}
	} else if tipo_doc == "1" {
		//Iniciar busque de sus apuntadores
		for i, block := range Inode.I_block {
			if block != -1 {
				text += "\tInodo" + indexInodo + ":"
				indexB = block //Para apuntadores del grafo
				if i < 12 {    //Pendiente verificar el indice

					apuntaGrafo := fmt.Sprintf("%d", i) //Para hacer conexciones con el grafo
					// indiceApunta := fmt.Sprintf("%d", block) //indice al que apunta la posicion actual
					text += "P" + apuntaGrafo + " -> "

					indexBloque := fmt.Sprintf("%d", indexB)
					//Procedimiento para iterar bloques de texto
					var crrFileBlock Structs.Fileblock
					if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperbloque.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
						return ""
					}
					//Construir Grafo del bloque
					text += "Bloque" + indexBloque + ";\n" //Para apuntadores del graf

				}
			}
		}

	}

	return text
}

func RecorridoSystem(Inode Structs.Inode, file *os.File, tempSuperbloque Structs.Superblock) string {
	text := ""
	//PARA INODOS
	indexInodo := fmt.Sprintf("%d", indexI)
	//
	//Para bloques
	indexBloque := fmt.Sprintf("%d", indexB)

	text += "\tInodo" + indexInodo + " [\n" //Para apuntadores del grafo
	text += "\t\tlabel=<\n"
	text += "\t\t\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	text += "\t\t\t\t<tr><td colspan=\"2\">Inodo " + indexInodo + "</td></tr>\n"

	//este for es para construir el grafo del inodo
	for i, block := range Inode.I_block {
		indicePuntero := fmt.Sprintf("%d", i+1)  //Para tabla
		apuntaGrafo := fmt.Sprintf("%d", i)      //Para hacer conexciones con el grafo
		indiceApunta := fmt.Sprintf("%d", block) //indice al que apunta la posicion actual
		text += "\t\t\t\t<tr><td>apt" + indicePuntero + "</td><td port='P" + apuntaGrafo + "'>" + indiceApunta + "</td></tr>\n"
	}
	text += "\t\t\t</table>\n"
	text += "\t\t>];\n"

	tipo_doc := strings.TrimRight(string(Inode.I_type[:]), "\x00")
	if tipo_doc == "0" {
		//Iniciar busque de sus apuntadores
		for i, block := range Inode.I_block {
			if block != -1 {
				indexBloque = fmt.Sprintf("%d", block)
				indexB = block
				if i < 12 { //Pendiente verificar el indice
					//procedimiento para recorrer bloque de carpeta
					var crrFolderBlock Structs.Folderblock
					if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperbloque.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
						return ""
					}
					//Construir grafo del bloque
					text += "\tBloque" + indexBloque + " [\n" //Para apuntadores del grafo
					text += "\t\tlabel=<\n"
					text += "\t\t\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
					text += "\t\t\t\t<tr><td colspan=\"2\">Bloque " + indexBloque + "</td></tr>\n"
					for i, folder := range crrFolderBlock.B_content {
						apuntaGrafoBloque := fmt.Sprintf("%d", i)                         //Para hacer coneciones grafo
						nombre_doc := strings.TrimRight(string(folder.B_name[:]), "\x00") //Nombre de archivo o carpeta
						ApuntaInodo := fmt.Sprintf("%d", folder.B_inodo)                  //inodo al que apunta
						text += "\t\t\t\t<tr><td>" + nombre_doc + "</td><td port='P" + apuntaGrafoBloque + "'>" + ApuntaInodo + "</td></tr>\n"
					}
					text += "\t\t\t</table>\n"
					text += "\t\t>];\n"

					//para llamar al siguiente inodo
					// coincide := false
					for _, folder1 := range crrFolderBlock.B_content {
						nombre_docu := strings.TrimRight(string(folder1.B_name[:]), "\x00")
						if nombre_docu != "." && nombre_docu != ".." && nombre_docu != "" {
							//Descartamos primeros dos inodos
							// coincide = true
							var NextInode Structs.Inode
							// Read object from bin file
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperbloque.S_inode_start+folder1.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return ""
							}
							indexI = folder1.B_inodo
							text += RecorridoSystem(NextInode, file, tempSuperbloque)

						}
					}
					// if !coincide {
					// 	return text
					// }
				}

			}
		}
		return text
	} else if tipo_doc == "1" {
		//Iniciar busque de sus apuntadores
		for i, block := range Inode.I_block {
			if block != -1 {
				indexB = block
				if i < 12 { //Pendiente verificar el indice
					indexBloque := fmt.Sprintf("%d", indexB)
					//Procedimiento para iterar bloques de texto
					var crrFileBlock Structs.Fileblock
					if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperbloque.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
						return ""
					}
					//Construir Grafo del bloque
					text += "\tBloque" + indexBloque + " [\n" //Para apuntadores del grafo
					text += "\t\tlabel=<\n"
					text += "\t\t\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
					text += "\t\t\t\t<tr><td colspan=\"1\">Bloque " + indexBloque + "</td></tr>\n"
					content := strings.TrimRight(string(crrFileBlock.B_content[:]), "\x00")
					text += "\t\t\t\t<tr><td>" + content + "</td></tr>\n"
					// indexB += 1
					text += "\t\t\t</table>\n"
					text += "\t\t>];\n"
				}
			}
		}
		return text
	}
	return text
}
