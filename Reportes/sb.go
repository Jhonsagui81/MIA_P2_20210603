package reportes

import (
	"PROYECTO1_MIA/Structs"
	"PROYECTO1_MIA/Utilities"
	"fmt"
	"strings"
)

func RecorridoSB(disco string, id string) string {
	text := "" //contenido a retornar
	//existe el disco
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

	//Inicio extraccion de info
	type_system := fmt.Sprintf("%d", tempSuperblock.S_filesystem_type)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>FileSystemType</td>\n"
	text += "\t\t\t<td>" + type_system + "</td>\n"
	text += "\t\t</tr>\n"

	//inodos count
	inodes_count := fmt.Sprintf("%d", tempSuperblock.S_inodes_count)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>InodesCount</td>\n"
	text += "\t\t\t<td>" + inodes_count + "</td>\n"
	text += "\t\t</tr>\n"
	//blocks count
	blocks_count := fmt.Sprintf("%d", tempSuperblock.S_blocks_count)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>BlocksCount</td>\n"
	text += "\t\t\t<td>" + blocks_count + "</td>\n"
	text += "\t\t</tr>\n"
	//free blocks count
	free_Blocks_count := fmt.Sprintf("%d", tempSuperblock.S_free_blocks_count)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Free_BlocksCount</td>\n"
	text += "\t\t\t<td>" + free_Blocks_count + "</td>\n"
	text += "\t\t</tr>\n"
	//free inodes count
	free_Inodes_count := fmt.Sprintf("%d", tempSuperblock.S_free_inodes_count)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Free_InodesCount</td>\n"
	text += "\t\t\t<td>" + free_Inodes_count + "</td>\n"
	text += "\t\t</tr>\n"
	//fecha montaje
	mtime := strings.TrimRight(string(tempSuperblock.S_mtime[:]), "\x00")
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>FechaMontaje</td>\n"
	text += "\t\t\t<td>" + mtime + "</td>\n"
	text += "\t\t</tr>\n"
	//Fecha desmontaje
	umtime := strings.TrimRight(string(tempSuperblock.S_umtime[:]), "\x00")
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Fecha Desmontaje</td>\n"
	text += "\t\t\t<td>" + umtime + "</td>\n"
	text += "\t\t</tr>\n"
	//cantidad montajes
	count_m := fmt.Sprintf("%d", tempSuperblock.S_mnt_count)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Cantidad Montaje</td>\n"
	text += "\t\t\t<td>" + count_m + "</td>\n"
	text += "\t\t</tr>\n"
	//id system
	id_system := fmt.Sprintf("%d", tempSuperblock.S_magic)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Id system file</td>\n"
	text += "\t\t\t<td>" + id_system + "</td>\n"
	text += "\t\t</tr>\n"
	//inode size
	inode_size := fmt.Sprintf("%d", tempSuperblock.S_inode_size)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Inode size</td>\n"
	text += "\t\t\t<td>" + inode_size + "</td>\n"
	text += "\t\t</tr>\n"
	//block size
	block_size := fmt.Sprintf("%d", tempSuperblock.S_block_size)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Block size</td>\n"
	text += "\t\t\t<td>" + block_size + "</td>\n"
	text += "\t\t</tr>\n"
	//prime inodo libre
	fist_inode_free := fmt.Sprintf("%d", tempSuperblock.S_fist_ino)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Primer inodo Libre</td>\n"
	text += "\t\t\t<td>" + fist_inode_free + "</td>\n"
	text += "\t\t</tr>\n"
	//primer bloque libre
	first_block_free := fmt.Sprintf("%d", tempSuperblock.S_first_blo)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Primer Bloque libre</td>\n"
	text += "\t\t\t<td>" + first_block_free + "</td>\n"
	text += "\t\t</tr>\n"
	//Inicio Bitman inodo
	Start_Bit_inodo := fmt.Sprintf("%d", tempSuperblock.S_bm_inode_start)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Inicio Bitman Inodo</td>\n"
	text += "\t\t\t<td>" + Start_Bit_inodo + "</td>\n"
	text += "\t\t</tr>\n"
	//Inicio bitman block
	Start_bit_block := fmt.Sprintf("%d", tempSuperblock.S_bm_block_start)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Inicio Bitman bloques</td>\n"
	text += "\t\t\t<td>" + Start_bit_block + "</td>\n"
	text += "\t\t</tr>\n"
	//inicia Tabla inodo
	Start_tab_inodo := fmt.Sprintf("%d", tempSuperblock.S_inode_start)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Inicia Tabla inodos</td>\n"
	text += "\t\t\t<td>" + Start_tab_inodo + "</td>\n"
	text += "\t\t</tr>\n"
	//inicia tabla blocks
	Start_tab_block := fmt.Sprintf("%d", tempSuperblock.S_block_start)
	text += "\t\t<tr>\n"
	text += "\t\t\t<td>Inicia tabla bloques</td>\n"
	text += "\t\t\t<td>" + Start_tab_block + "</td>\n"
	text += "\t\t</tr>\n"
	//

	return text
}
