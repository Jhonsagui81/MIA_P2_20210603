package user

import (
	global "PROYECTO1_MIA/Global"
	"fmt"
)

func Logout() {
	fmt.Println("===== Start logout ======")
	if global.LoginValidacion() {
		//Existe usuario logeado
		global.Deslogearse()
		fmt.Println("Deslogeo correcto")
	} else {
		//No existe usuario logeado
		fmt.Println("Error: No hay usuario logeado -logout")
	}
	fmt.Println("===== End logout ======")
}
