package analizer

import (
	filesystem "PROYECTO1_MIA/FileSystem"
	User "PROYECTO1_MIA/User"
	user "PROYECTO1_MIA/User"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Global Variables
var ContadorArchivos int
var Letra_Disco = 'A'
var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func Analize() {
	for true {
		var input string
		fmt.Println("Ingrese comando: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getCommandAndParams(input)

		if command != "" {
			fmt.Println("Command: ", command, "Params: ", params)
			AnalyzeCommnad(strings.ToLower(command), params)
		}
	}
}

func AnalyzeCommnad(command string, params string) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "unmount") {
		fn_unmount(params)
	} else if strings.Contains(command, "mount") {
		fn_mount(params)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params)
	} else if strings.Contains(command, "login") {
		fn_login(params)
	} else if strings.Contains(command, "logout") {
		fn_logout()
	} else if strings.Contains(command, "mkgrp") {
		fn_mkgrp(params)
	} else if strings.Contains(command, "rmgrp") {
		fn_rmgrp(params)
	} else if strings.Contains(command, "mkusr") {
		fn_mkusr(params)
	} else if strings.Contains(command, "rmusr") {
		fn_rmusr(params)
	} else if strings.Contains(command, "rep") {
		fn_rep(params)
	} else if strings.Contains(command, "execute") {
		fn_execute(params)
	} else if strings.Contains(command, "#") {
		fmt.Println("Es un comentario")
	} else if strings.Contains(command, "pause") {
		fn_pause()
	} else if strings.Contains(command, "cat") {
		fn_cat(params)
	} else if strings.Contains(command, "mkdir") {
		fn_mkdir(params)
	} else {
		fmt.Println("Error: Command not found")
	}

}

func AnalizerType(ruta_archivo string) {
	// Expresión regular para "execute path direccion"

	//Extrae la ruta del archivo a ejecutar
	// Abrir el archivo
	fmt.Println("ruta :", ruta_archivo)
	file, err := os.Open(ruta_archivo)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Crea un escáner para leer el archivo línea por línea
	scanner := bufio.NewScanner(file)
	// Itera sobre cada línea del archivo
	for scanner.Scan() {
		line := scanner.Text() // Obtiene la línea actual
		//Expresion para eliminar comentarios
		regex := regexp.MustCompile(`#.*S`)
		lineaSinComentario := regex.ReplaceAllString(line, "")

		command, params := getCommandAndParams(lineaSinComentario)
		if command != "" {
			fmt.Println("Command: ", command, "Params: ", params)
			AnalyzeCommnad(strings.ToLower(command), params)
		}
	}
}

func fn_mkdir(params string) {
	path := ""
	r := false
	// var re = regexp.MustCompile(`-(\w+)|\s-`)
	slide := strings.Split(params, "-")
	for _, param := range slide {
		if param == "r " {
			r = true
		}
	}
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "path":
			path = flagValue
		case "r":
			r = true
		default:
			fmt.Println("Error: Flag not found")
		}

	}
	filesystem.Mkdir(path, r)
}

func fn_cat(params string) {
	path := ""
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "file":
			path = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}

	}
	filesystem.Cat(path)

}

func fn_pause() {
	fmt.Println("====== Start Pause ======")
	fmt.Println("Presione ENTER para continuar")
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada del usuario:", err)
		return
	}
	fmt.Println("===== END Pause =====")
}

func fn_execute(params string) {
	path := ""
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "path":
			path = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}

	}
	AnalizerType(path)
}

func fn_mkdisk(params string) {
	// Define flags
	size := 0
	unit := "m"
	fit := "ff"
	flags := true
	// Parse the flags

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			size = sizeValue
		case "unit":
			flagValue = strings.ToLower(flagValue)
			unit = flagValue
		case "fit":
			flagValue = strings.ToLower(flagValue)
			fit = flagValue
		default:
			flags = false
		}
	}

	// Call the function
	if flags {
		mkdisk(size, unit, fit)
	} else {
		fmt.Println("Error: Parametros no validos -MKDISK")
	}

}

func fn_rmdisk(params string) {
	letter := ""
	matches := re.FindAllStringSubmatch(params, -1)
	flags := true

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "driveletter":
			flagValue = strings.ToUpper(flagValue)
			letter = flagValue
		default:
			flags = false
		}

	}
	if flags {
		rmdisk(letter)
	} else {
		fmt.Println("Error: Parametros no validos. -RMDISK")
	}

}

func fn_fdisk(params string) {
	size := 0
	unit := "k"
	letter := ""
	name := ""
	type_ := "p"
	fit := "wf"
	delete := ""
	add := 0
	flags := true

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			size = sizeValue
		case "unit":
			flagValue = strings.ToLower(flagValue)
			unit = flagValue

		case "driveletter":
			flagValue = strings.ToUpper(flagValue)
			letter = flagValue

		case "name":
			flagValue = strings.ToLower(flagValue)
			name = flagValue

		case "type":
			flagValue = strings.ToLower(flagValue)
			type_ = flagValue
		case "fit":
			flagValue = strings.ToLower(flagValue)
			fit = flagValue
		case "delete":
			flagValue = strings.ToLower(flagValue)
			delete = flagValue
		case "add":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			add = sizeValue
		default:
			flags = false
		}
	}
	if flags {
		fdisk(size, unit, letter, name, type_, fit, delete, add)
	} else {
		fmt.Println("Error: Parametros no validos. - FDISK")
	}
}

func fn_mount(params string) {
	letter := ""
	name := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter":
			flagValue = strings.ToUpper(flagValue)
			letter = flagValue
		case "name":
			flagValue = strings.ToLower(flagValue)
			name = flagValue
		default:
			flags = false
		}
	}
	if flags {
		mount(letter, name)
	} else {
		fmt.Println("Error: Parametros no validos. -MOUNT")
	}
}

func fn_unmount(params string) {
	id := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id":
			flagValue = strings.ToUpper(flagValue)
			id = flagValue
		default:
			flags = false
		}
	}
	if flags {
		Unmount(id)
	} else {
		fmt.Println("Error: Parametros no validos. -UNMOUNT")
	}
}

func fn_mkfs(params string) {
	id := ""
	type_ := "full"
	fs := "2fs"
	flags := true

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id":
			flagValue = strings.ToUpper(flagValue)
			id = flagValue
		case "type":
			flagValue = strings.ToLower(flagValue)
			type_ = flagValue
		case "fs":
			flagValue = strings.ToLower(flagValue)
			fs = flagValue
		default:
			flags = false
		}
	}

	if flags {
		filesystem.Mkfs(id, type_, fs)
	} else {
		fmt.Println("Error: Parametros no validos. -MKFS")
	}
}

func fn_login(params string) {
	user := ""
	pass := ""
	id := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "user":
			user = flagValue
		case "pass":
			pass = flagValue
		case "id":
			flagValue = strings.ToUpper(flagValue)
			id = flagValue
		default:
			flags = false
		}

	}

	if flags {
		User.Login(user, pass, id)
	} else {
		fmt.Println("Error: parametros no validos. -LOGIN")
	}
}

func fn_logout() {
	user.Logout()
}

func fn_mkgrp(params string) {
	name := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "name":
			name = flagValue
		default:
			flags = false
		}

	}

	if flags {
		User.Mkgrp(name)
	} else {
		fmt.Println("Error: Parametros no validos. -MKGRP")
	}
}

func fn_rmgrp(params string) {
	name := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "name":
			name = flagValue
		default:
			flags = false
		}

	}
	if flags {
		user.Rmgrp(name)
	} else {
		fmt.Println("Error: Parametros no validos. -RMGRP")
	}
}

func fn_mkusr(params string) {
	newUser := ""
	newPass := ""
	Group := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "user":
			newUser = flagValue
		case "pass":
			newPass = flagValue
		case "grp":
			Group = flagValue
		default:
			flags = false
		}

	}
	if flags {
		user.Mkusr(newUser, newPass, Group)
	} else {
		fmt.Println("Error: Parametros no validos. -MKUSR")
	}
}

func fn_rmusr(params string) {
	user1 := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"") // Removing the double
		switch flagName {
		case "user":
			user1 = flagValue
		default:
			flags = false
		}

	}
	if flags {
		user.Rmusr(user1)
	} else {
		fmt.Println("Error: Parametros no validos. -RMUSR")
	}
}

func fn_rep(params string) {
	name := ""
	path := ""
	id := ""
	ruta := ""
	flags := true
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name":
			flagValue = strings.ToLower(flagValue)
			name = flagValue
		case "path":
			path = flagValue
		case "id":
			flagValue = strings.ToUpper(flagValue)
			id = flagValue
		case "ruta":
			flagValue = strings.ToLower(flagValue)
			ruta = flagValue
		default:
			flags = false
		}
	}

	// fmt.Println("Ruta:", ruta)
	if flags {
		rep(name, path, id, ruta)
	} else {
		fmt.Println("Error: Parametros no validos. -REP")
	}
}
