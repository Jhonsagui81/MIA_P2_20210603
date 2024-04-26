package global

import (
	"fmt"
	"strconv"
)

type UserLogin struct {
	User           string
	Pass           string
	UID            string
	GID            string
	disco          string
	particion      string
	partitionStart int
	IdCompleto     string
}

var usuario UserLogin

func LoginValidacion() bool {
	log := false

	if usuario.User == "" {
		//No hay usuari	o logeado
		log = false
	} else {
		//Existe usuario logeado
		log = true
	}
	return log
}

func Deslogearse() {
	usuario.User = ""
	usuario.Pass = ""
	usuario.GID = ""
	usuario.UID = ""
	usuario.disco = ""
	usuario.particion = ""
	usuario.partitionStart = 0
	usuario.IdCompleto = ""
}

func Logear(user string, pass string, UID string, GID string, disco string, particion string, start int, id string) {
	usuario.User = user
	usuario.Pass = pass
	usuario.UID = UID
	usuario.GID = GID
	usuario.disco = disco
	usuario.particion = particion
	usuario.partitionStart = start
	usuario.IdCompleto = id
}

func ShowInfo() {
	fmt.Println(usuario.User)
	fmt.Println(usuario.Pass)
	fmt.Println(usuario.UID)
	fmt.Println(usuario.GID)
	fmt.Println(usuario.disco)
	fmt.Println(usuario.particion)
	fmt.Println(usuario.partitionStart)
}

func ValidaUsuario(user string) bool {
	if usuario.User == "root" {
		return true
	} else {
		return false
	}
}

func InfoDisk() string {
	return usuario.disco
}

func InfoUsuario() int32 {
	idU, err := strconv.Atoi(usuario.UID)
	if err != nil {
		fmt.Println("Error convertir string")
		return -1
	}
	return int32(idU)

}

func DataLogin() (string, string, string) {
	return usuario.User, usuario.Pass, usuario.IdCompleto
}

func InfoGrupo() int32 {
	idG, err := strconv.Atoi(usuario.GID)
	if err != nil {
		fmt.Println("Error convertir string")
		return -1
	}
	return int32(idG)

}

func InfoID() string {
	return usuario.IdCompleto
}

func StartPartition() int {
	return usuario.partitionStart
}

func InfoPartition() int64 {
	index, err := strconv.ParseInt(usuario.particion, 10, 32)
	if err != nil {
		fmt.Println("Error al convertir la cadena:", err)
		return -1
	}

	return index
}

func DeterminarPermisoEscritura(GID string, UID string, p_u string, p_g string, p_o string) bool {
	//Verificar si el usuario actual pertenece al usuario creador o es parte del grupo
	fmt.Println("Entra a determinar permisos =========")
	fmt.Println("GID", GID)
	fmt.Println("UID", UID)
	fmt.Println("permisos: " + p_u)
	fmt.Println("---------------- Usuario info")
	fmt.Println("GidUsuario", usuario.GID)
	fmt.Println("Uidddd", usuario.UID)
	tienePermiso := false
	if usuario.UID == string(UID) {
		//Es el usuario creador
		if p_u == "2" || p_u == "6" || p_u == "3" {
			//Tiene permiso de escritura y lectura
			tienePermiso = true
		}
	} else if usuario.GID == GID {
		//Esta dentro del grupo
		if p_g == "2" || p_g == "6" || p_g == "3" {
			//Tiene permiso de escritura
			tienePermiso = true
		}
	} else {
		tienePermiso = false
		//Se considera como "otros"
		//Solo tiene permisos de "lectura"

	}

	return tienePermiso
}

func DeterminarPermisoLectura(GID string, UID string, p_u string, p_g string, p_o string) bool {
	//Verificar si el usuario actual pertenece al usuario creador o es parte del grupo
	fmt.Println("Entra a determinar permisos =========")
	fmt.Println("GID", GID)
	fmt.Println("UID", UID)
	fmt.Println("permisos: " + p_u)
	fmt.Println("---------------- Usuario info")
	fmt.Println("GidUsuario", usuario.GID)
	fmt.Println("Uidddd", usuario.UID)
	tienePermiso := false
	if usuario.UID == string(UID) {
		//Es el usuario creador
		if p_u == "4" || p_u == "5" || p_u == "6" {
			//Tiene permiso de escritura y lectura
			tienePermiso = true
		}
	} else if usuario.GID == GID {
		//Esta dentro del grupo
		if p_g == "4" || p_g == "5" || p_g == "6" {
			//Tiene permiso de escritura
			tienePermiso = true
		}
	} else {
		if p_o == "4" {
			tienePermiso = true
		}
		//Se considera como "otros"
		//Solo tiene permisos de "lectura"

	}

	return tienePermiso
}
