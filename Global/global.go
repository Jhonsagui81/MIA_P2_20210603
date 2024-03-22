package global

type UserLogin struct {
	User      string
	Pass      string
	disco     string
	particion string
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
	usuario.disco = ""
	usuario.particion = ""
}

func Logear(user string, pass string, disco string, particion string) {
	usuario.User = user
	usuario.Pass = pass
	usuario.disco = disco
	usuario.particion = particion
}

func ValidaUsuario(user string) bool {
	if usuario.User == "root" {
		return true
	} else {
		return false
	}
}

func InfoDisk() string {
	return usuario.particion
}
