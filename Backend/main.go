package main

import (
	analizer "PROYECTO1_MIA/Analizer"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	//Enrutador
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", )

	analizer.Analize()
}
