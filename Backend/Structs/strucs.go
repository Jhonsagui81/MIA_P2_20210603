package Structs

type MRB struct {
	MbrSize      int32
	CreationDate [10]byte
	Signature    int32
	Fit          [2]byte
	Partitions   [4]Partition
}

type Partition struct {
	Status      [1]byte //1 -> montada, 0 -> Desmontada
	Type        [1]byte //P o E
	Fit         [2]byte // B o F o W
	Start       int32   //Inicia la particion
	Size        int32   //Tamano de la particion
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}

type EBR struct {
	Part_mount [1]byte
	Part_fit   [2]byte
	Part_start int32 //Donde incia
	Part_s     int32 //Tamano de particion
	Part_next  int32 //Byte donde incia el siguiente EBR
	Part_name  [16]byte
}

//  =============================================================

type Superblock struct {
	S_filesystem_type   int32    //Identifica sistema de archivos
	S_inodes_count      int32    //numero total inodos
	S_blocks_count      int32    //numero total bloques
	S_free_blocks_count int32    //bloques libres
	S_free_inodes_count int32    //inodos libres
	S_mtime             [17]byte //fecha fue montado
	S_umtime            [17]byte //fecha fue desmontado
	S_mnt_count         int32    //cuantas veces monto
	S_magic             int32    //identifica sitem file 0xEF56
	S_inode_size        int32    //size inodo
	S_block_size        int32    //size bloque
	S_fist_ino          int32    //first inodo free
	S_first_blo         int32    //first block free
	S_bm_inode_start    int32    //inicio bitmap inodo
	S_bm_block_start    int32    //inicio bitmap bloques
	S_inode_start       int32    //inicio tabla inodos
	S_block_start       int32    //inicio tabla bloques
}

//  =============================================================

type Inode struct {
	I_uid   int32 //UsuarioID usuario propietario de file o folder
	I_gid   int32 //Grupo ID
	I_size  int32
	I_atime [17]byte
	I_ctime [17]byte
	I_mtime [17]byte
	I_block [15]int32 //array de apuntador que tiene inodo
	I_type  [1]byte   //1=archivo \n 0=carpeta
	I_perm  [3]byte   //Perimiso de archivo/carpeta UGO
}

//  =============================================================

type Fileblock struct { //Bloques de archivos
	B_content [64]byte
}

//  =============================================================

type Content struct {
	B_name  [12]byte //Nombre carpeta archivo
	B_inodo int32    //Apuntador
}

type Folderblock struct { //Bloque de carpetas
	B_content [4]Content
}

//  =============================================================

type Pointerblock struct { //Bloque de apuntadores (indirectos)
	B_pointers [16]int32
}

//  =============================================================

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [17]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}

type UserTXT struct {
	Id    string
	Tipo  string
	Grupo string
	User  string
	Pass  string
}

type Command struct {
	Nombre string `json: "Nombre"`
	Id     int    `json:"ID"`
}

type Discos []Command

type RespuestaFron struct {
	Respuesta string `json: "Respuesta"`
}
