package UtilitiesInodes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func InitSearch(path string, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("======Start INITSEARCH======")
	fmt.Println("path:", path)
	// path = "/ruta/nueva"

	// split the path by /
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:] //corta el / inicial

	// fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 Structs.Inode
	// Read object from bin file
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	fmt.Println("======End INITSEARCH======")

	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// Para elimincar espacion entre palabras (no estoy seguro)
func pop(s *[]string) string {
	lastIndex := 0 // Cambiar el índice a 0
	last := (*s)[lastIndex]
	*s = (*s)[lastIndex+1:] // Cambiar el orden de las operaciones
	return last
}

func SarchInodeByPath(StepsPath []string, Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("======Start SARCHINODEBYPATH======")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFolderBlock Structs.Folderblock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {
					// fmt.Println("Folder found======")
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if strings.Contains(string(folder.B_name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_inodo
						} else {
							fmt.Println("NextInode======")
							var NextInode Structs.Inode
							// Read object from bin file
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return -1
							}
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}

			} else {
				//CASO INDIRECTO
			}
		}
		index++
	}

	fmt.Println("======End SARCHINODEBYPATH======")
	return -1
}

func GetInodeFileData(Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) string {
	fmt.Println("======Start GETINODEFILEDATA======")
	index := int32(0)
	// define content as a string
	content := ""

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFileBlock Structs.Fileblock
				// Read object from bin file
				if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
					return ""
				}

				content += strings.TrimRight(string(crrFileBlock.B_content[:]), "\x00")

			} else {
				//CASO INDIRECTO
			}
		}
		index++
	}

	fmt.Println("======End GETINODEFILEDATA======")
	return content
}
