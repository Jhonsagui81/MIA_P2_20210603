package Utilities

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Funtion to create bin file
func CreateFile(name string) error {
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Create file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name) //si no existe el archivo, lo crea
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funtion to open bin file in read/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	// fmt.Println("Error lecura??: ", err)
	return nil
}

func NumeroRandom() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(5-1) + 1
}

func CrearGrafo(text string, pathDot string, pathPDF string) {
	err := CreateFile(pathDot)
	if err != nil {
		fmt.Println("Error: Creacion archivo para grafo MBR ", err)
		return
	}

	//apertura de archivo
	archivo, err := os.OpenFile(pathDot, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("apertura:", err)
		return
	}
	defer archivo.Close()

	//Escritura de archivo
	// fmt.Println(text)
	_, err = io.WriteString(archivo, text)
	if err != nil {
		fmt.Println("Escritura:", err)
		return
	}

	//Conversion
	// fmt.Println("dot:", pathDot)
	// fmt.Println("pdf:", pathPDF)
	cmd := exec.Command("dot", "-Tpdf", pathDot, "-o", pathPDF)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Transformacion:", err)
		fmt.Println(string(stderr.Bytes()))
		return
	}
}
