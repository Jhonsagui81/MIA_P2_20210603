package reportes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func InitInodos(disco string, id string) string {
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
		fmt.Println("-1")
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
		fmt.Println("-11")
		return ""
	}

	var Inode0 Structs.Inode
	// Read object from bin file
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		fmt.Println("-111")
		return ""
	}

	text += RecorridoInodo(Inode0, file, tempSuperblock)

	return text
}

func RecorridoInodo(Inode Structs.Inode, file *os.File, tempSuperbloque Structs.Superblock) string {
	text := ""
	//PARA INODOS
	indexInodo := fmt.Sprintf("%d", indexI)
	//
	//Para bloques
	// indexBloque := fmt.Sprintf("%d", indexB)

	text += "\tInodo" + indexInodo + " [\n" //Para apuntadores del grafo
	text += "\t\tlabel=<\n"
	text += "\t\t\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	text += "\t\t\t\t<tr><td colspan=\"2\">Inodo " + indexInodo + "</td></tr>\n"

	//Demas atributos de inodos
	//I_uid
	I_uid := fmt.Sprintf("%d", Inode.I_uid)
	text += "\t\t\t\t<tr><td>I_uid</td><td>" + I_uid + "</td></tr>\n"
	//I-gid
	I_gid := fmt.Sprintf("%d", Inode.I_gid)
	text += "\t\t\t\t<tr><td>I_gid</td><td>" + I_gid + "</td></tr>\n"
	//I_size
	I_size := fmt.Sprintf("%d", Inode.I_size)
	text += "\t\t\t\t<tr><td>I_size</td><td>" + I_size + "</td></tr>\n"
	//Iatime
	I_atime := strings.TrimRight(string(Inode.I_atime[:]), "\x00")
	text += "\t\t\t\t<tr><td>I_atime</td><td>" + I_atime + "</td></tr>\n"
	//Ictime
	I_ctime := strings.TrimRight(string(Inode.I_ctime[:]), "\x00")
	text += "\t\t\t\t<tr><td>I_ctime</td><td>" + I_ctime + "</td></tr>\n"
	//Imtime
	I_mtime := strings.TrimRight(string(Inode.I_mtime[:]), "\x00")
	text += "\t\t\t\t<tr><td>I_mtime</td><td>" + I_mtime + "</td></tr>\n"
	//Iblock
	//este for es para construir el grafo del inodo
	for i, block := range Inode.I_block {
		indicePuntero := fmt.Sprintf("%d", i+1)  //Para tabla
		apuntaGrafo := fmt.Sprintf("%d", i)      //Para hacer conexciones con el grafo
		indiceApunta := fmt.Sprintf("%d", block) //indice al que apunta la posicion actual
		text += "\t\t\t\t<tr><td>apt" + indicePuntero + "</td><td port='P" + apuntaGrafo + "'>" + indiceApunta + "</td></tr>\n"
	}
	//itupe
	I_type := strings.TrimRight(string(Inode.I_type[:]), "\x00")
	text += "\t\t\t\t<tr><td>I_type</td><td>" + I_type + "</td></tr>\n"
	//iperm
	I_perm := strings.TrimRight(string(Inode.I_perm[:]), "\x00")
	text += "\t\t\t\t<tr><td>I_perm</td><td>" + I_perm + "</td></tr>\n"
	//Finalizar tabla
	text += "\t\t\t</table>\n"
	text += "\t\t>];\n"

	tipo_doc := strings.TrimRight(string(Inode.I_type[:]), "\x00")
	if tipo_doc == "0" {
		//Iniciar busque de sus apuntadores
		for i, block := range Inode.I_block {
			if block != -1 {
				// indexBloque = fmt.Sprintf("%d", block)
				indexB = block
				if i < 12 { //Pendiente verificar el indice
					//procedimiento para recorrer bloque de carpeta
					var crrFolderBlock Structs.Folderblock
					if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperbloque.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
						return ""
					}
					// //Construir grafo del bloque
					// text += "\tBloque" + indexBloque + " [\n" //Para apuntadores del grafo
					// text += "\t\tlabel=<\n"
					// text += "\t\t\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
					// text += "\t\t\t\t<tr><td colspan=\"2\">Bloque " + indexBloque + "</td></tr>\n"
					// for i, folder := range crrFolderBlock.B_content {
					// 	apuntaGrafoBloque := fmt.Sprintf("%d", i)                         //Para hacer coneciones grafo
					// 	nombre_doc := strings.TrimRight(string(folder.B_name[:]), "\x00") //Nombre de archivo o carpeta
					// 	ApuntaInodo := fmt.Sprintf("%d", folder.B_inodo)                  //inodo al que apunta
					// 	text += "\t\t\t\t<tr><td>" + nombre_doc + "</td><td port='P" + apuntaGrafoBloque + "'>" + ApuntaInodo + "</td></tr>\n"
					// }
					// text += "\t\t\t</table>\n"
					// text += "\t\t>];\n"

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
							text += RecorridoInodo(NextInode, file, tempSuperbloque)

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
					// indexBloque := fmt.Sprintf("%d", indexB)
					//Procedimiento para iterar bloques de texto
					var crrFileBlock Structs.Fileblock
					if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperbloque.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
						return ""
					}
					// //Construir Grafo del bloque
					// text += "\tBloque" + indexBloque + " [\n" //Para apuntadores del grafo
					// text += "\t\tlabel=<\n"
					// text += "\t\t\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
					// text += "\t\t\t\t<tr><td colspan=\"1\">Bloque " + indexBloque + "</td></tr>\n"
					// content := strings.TrimRight(string(crrFileBlock.B_content[:]), "\x00")
					// text += "\t\t\t\t<tr><td>" + content + "</td></tr>\n"
					// // indexB += 1
					// text += "\t\t\t</table>\n"
					// text += "\t\t>];\n"
				}
			}
		}
		return text
	}
	return text
}
