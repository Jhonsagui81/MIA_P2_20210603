package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// Endpoint para recibir el nombre
	router.HandleFunc("/comand", InputCommand)

	fmt.Println("Servidor iniciado en el puerto 3000")
	http.ListenAndServe(":3000", router)
}

func InputCommand(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Si llesdsdgo el comando")
}
