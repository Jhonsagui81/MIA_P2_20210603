package Utilities

import (
	"PROYECTO1_MIA/Structs"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func BusquedaArchivo(letter string) bool {
	archivos, err := ioutil.ReadDir("./MIA/P1")
	if err != nil {
		fmt.Println("Error: No se puedo acceder a los discos ", err)
	}

	name_file := letter + ".dsk"

	encontrado := false
	for _, archivo := range archivos {
		if archivo.Name() == name_file {
			//Encontro el disco buscado
			encontrado = true
			break
		}
	}

	return encontrado
}

func RecorridoParticionesDisco(file *os.File, nombre_comparar string, TempMBR Structs.MRB) bool {
	Se_repite := false
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			nombre_particion := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
			if nombre_particion == nombre_comparar {
				//Los nombres de las particiones se repiten
				Se_repite = true
				return Se_repite
			}
		}
	}
	return Se_repite
}

func ConteoParticiones(TempMBR Structs.MRB) bool {
	Disco_lleno := true
	for i := 0; i < 4; i++ {
		// fmt.Println("size particion:", TempMBR.Partitions[i].Size)
		if TempMBR.Partitions[i].Size == 0 {
			Disco_lleno = false
		}
	}
	// fmt.Println("retornar:", Disco_lleno)
	return Disco_lleno
}

func ParticionExtendida(TempMBR Structs.MRB) bool {
	existe := false
	for i := 0; i < 4; i++ {
		if string(TempMBR.Partitions[i].Type[:]) == "e" {
			//Ya existe una particion extendida
			existe = true
		}
	}
	return existe
}

func InsertarParticion(file *os.File, TempMBR *Structs.MRB, size int, unit string, name string, type_ string, fit string) bool {
	var count = 0
	var gap = int32(0)
	aceptacion := false
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			count++
			//Donde inicio + su tamano = el inicio de la siguiente particion
			gap = TempMBR.Partitions[i].Start + TempMBR.Partitions[i].Size
		}
	}

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size == 0 {
			if gap+int32(size) < TempMBR.MbrSize {
				TempMBR.Partitions[i].Size = int32(size)
				if count == 0 {
					TempMBR.Partitions[i].Start = int32(binary.Size(TempMBR))
				} else {
					TempMBR.Partitions[i].Start = gap
				}
				copy(TempMBR.Partitions[i].Name[:], name)
				// copy(TempMBR.Partitions[i].Fit[:], fit)
				copy(TempMBR.Partitions[i].Status[:], "0")
				copy(TempMBR.Partitions[i].Type[:], type_)
				copy(TempMBR.Partitions[i].Fit[:], fit)
				TempMBR.Partitions[i].Correlative = int32(count + 1)
				aceptacion = true
				break
			}
		}
	}
	// Overwrite the MBR
	if err := WriteObject(file, TempMBR, 0); err != nil {
		return false
	}
	return aceptacion
}

func InsertaLogica(file *os.File, TempMBR *Structs.MRB, size int, name string, fit string) {
	// var auxEBR int            //Posicion donde acaba la extendida
	var continuar bool = true //Salir del while de busqueda de siguiente EBR
	for i := 0; i < 4; i++ {
		if string(TempMBR.Partitions[i].Type[:]) == "e" {
			//tamano de extendida
			fmt.Println("La Particion expandida", i)
			tamano := TempMBR.Partitions[i].Start + TempMBR.Partitions[i].Size
			var WriteNextEBR int = 0
			SaltosEBR := int(TempMBR.Partitions[i].Start)

			for continuar {
				var TempEBR Structs.EBR

				if err := ReadObject(file, &TempEBR, int64(SaltosEBR)); err != nil {
					return
				} //En caso si Reconozca algo en el archivo EBR

				// // Read object from bin filefo
				// fmt.Println("-----Entro al for Datos del TempEBR recuperado-------- ")
				// fmt.Println("inicia:", TempEBR.Part_start)
				// fmt.Println("Tamano:", TempEBR.Part_s)
				// fmt.Println("Next:", TempEBR.Part_next)
				// fmt.Println("SaltosEBR:", SaltosEBR)
				// fmt.Println("startExtende:", TempMBR.Partitions[i].Start)
				// fmt.Println("-----Ent-------- ")
				if TempEBR.Part_next == 0 { //Como no encontro nada entonces se crea EBR
					// fmt.Println("No entra a if??")
					if int(SaltosEBR)+size < int(tamano) {
						//Rellenar el struct y meterlo en la posicion start de la Particion actual
						// var newEBR Structs.EBR
						copy(TempEBR.Part_mount[:], "0")
						copy(TempEBR.Part_fit[:], fit)
						// int32(binary.Size(TempMBR)
						TempEBR.Part_start = (int32(SaltosEBR) + int32(binary.Size(TempEBR)))
						TempEBR.Part_s = int32(size)
						TempEBR.Part_next = int32(-1)
						copy(TempEBR.Part_name[:], name)
						//Escribir la estructura en el archivo
						if err := WriteObject(file, TempEBR, int64(SaltosEBR)); err != nil {
							return
						}
						continuar = false
						fmt.Println("Se introdujo EBR")
					}
				} else if TempEBR.Part_next == -1 {
					fmt.Println("entra a -1")
					//Estamos en el ultimo EBR y se debe escribir su siguiente
					WriteNextEBR = int(TempEBR.Part_start) + int(TempEBR.Part_s)
					// auxEBR = int(final)
					//en la posicion writeNextEBR voy a meter el nuevo EBR (particion logica)
					//Crear nuevo EBR
					if int(WriteNextEBR)+int(binary.Size(TempEBR))+size < int(tamano) { //Valida si la particion cabe dentro de la extendida
						TempEBR.Part_next = int32(WriteNextEBR) //Tengo que sobreescribir esta posicion para actualizar cambios start - sizeEBR
						//Rellenar el struct y meterlo en la posicion start de la Particion actual
						var newEBR Structs.EBR
						copy(newEBR.Part_mount[:], "0")
						copy(newEBR.Part_fit[:], fit)
						// int32(binary.Size(TempMBR)
						newEBR.Part_start = int32(WriteNextEBR) + int32(binary.Size(newEBR))
						newEBR.Part_s = int32(size)
						newEBR.Part_next = int32(-1)
						copy(newEBR.Part_name[:], name) //Aqui si debe validar si se repite nombre de particion logica
						//Escribir la estructura en el archivo
						if err := WriteObject(file, newEBR, int64(WriteNextEBR)); err != nil {
							return
						}
						posicion_anterior := int32(TempEBR.Part_start) - int32(binary.Size(TempEBR))
						if err := WriteObject(file, TempEBR, int64(posicion_anterior)); err != nil {
							return
						}
						continuar = false
						fmt.Println("Se va a crear particion logica anidada")
						fmt.Println("")
					} else {
						fmt.Println("No cabe la particion")
					}

				} else if TempEBR.Part_next > 0 {
					fmt.Println("Entra a saltos")
					SaltosEBR = int(TempEBR.Part_next) //posicion de siguiente EBR
				}
			}
		}
	}
}

func EliminarParticion(name string, disco string) bool {
	//Buscar el disco
	var confirm_delete string //letra S/n para eliminar
	var posicion_par int      //Posicion de la particion a eliminar
	salio_mal := true         //Estado del proceso
	existe_disco := BusquedaArchivo(disco)
	if existe_disco {
		//existe el disco
		filepath := "./MIA/P1/" + disco + ".dsk"
		file, err := OpenFile(filepath)
		if err != nil {
			return salio_mal
		}
		defer file.Close()
		//recoger  informacion del disco
		var TempMBR Structs.MRB
		// Read object from bin file
		if err := ReadObject(file, &TempMBR, 0); err != nil {
			return salio_mal
		}

		//Buscar el nombre de la particion
		Se_repite := false
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				nombre_particion := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
				if nombre_particion == name {
					//Los nombres de las particiones se repiten
					Se_repite = true
					posicion_par = i

				}
			}
		}
		//validar la particion
		if Se_repite {
			//La particion existe
			fmt.Println("->Desea Eliminar la particion: " + name + ", del disco (S/N)? ")
			_, err := fmt.Scan(&confirm_delete)

			//valida lectura de consola
			if err != nil {
				fmt.Println("Error: confirmacion eliminacion particion", err)
				return salio_mal
			}

			//Valida si la respusta es valida
			aux_confir := strings.ToLower(confirm_delete)
			if aux_confir != "s" && aux_confir != "n" {
				fmt.Println("Error: opcion incorrecta fdisk -delete ", err)
				return salio_mal
			}

			if aux_confir == "s" {
				//Se elimina particion
				TempMBR.Partitions[posicion_par].Size = 0
				TempMBR.Partitions[posicion_par].Start = 0
				TempMBR.Partitions[posicion_par].Correlative = 0
				reflect.ValueOf(&TempMBR.Partitions[posicion_par].Status).Elem().Set(reflect.Zero(reflect.TypeOf(TempMBR.Partitions[posicion_par].Status)))
				reflect.ValueOf(&TempMBR.Partitions[posicion_par].Type).Elem().Set(reflect.Zero(reflect.TypeOf(TempMBR.Partitions[posicion_par].Type)))
				reflect.ValueOf(&TempMBR.Partitions[posicion_par].Fit).Elem().Set(reflect.Zero(reflect.TypeOf(TempMBR.Partitions[posicion_par].Fit)))
				reflect.ValueOf(&TempMBR.Partitions[posicion_par].Name).Elem().Set(reflect.Zero(reflect.TypeOf(TempMBR.Partitions[posicion_par].Name)))
				reflect.ValueOf(&TempMBR.Partitions[posicion_par].Id).Elem().Set(reflect.Zero(reflect.TypeOf(TempMBR.Partitions[posicion_par].Id)))

				salio_mal = false
				//Mandar a escribirlo de nuevo al disco
				if err := WriteObject(file, TempMBR, 0); err != nil {
					return salio_mal
				}

			} else {
				//No se eliminar
				fmt.Println("Cancela proceso de eliminacion de particion")
				return salio_mal
			}

		} else {
			//NO existe el nombre de la particion en el disco
			fmt.Println("Error: No existe la Particion a eliminar ")
			return salio_mal
		}

	} else {
		//El disco no existe
		fmt.Println("Erro: No existe el disco donde quiere eliminar la particion")
		return salio_mal
	}
	return salio_mal
}

func Add_Espacio(disco string, name string, unit string, add int) bool {
	agrego := false      //Estado de esta funcion, si logro su funcion = true
	positivo := false    //Valida si se agrega o quita espacio
	var posicion_par int //Posicion de la particion afectada

	confir_disco := BusquedaArchivo(disco)
	if confir_disco {
		//Encontro Disco
		//existe el disco
		filepath := "./MIA/P1/" + disco + ".dsk"
		file, err := OpenFile(filepath)
		if err != nil {
			fmt.Println("Error: No se pudo leer el disco  fdisk -add")
			return agrego
		}
		defer file.Close()
		//recoger  informacion del disco
		var TempMBR Structs.MRB
		// Read object from bin file
		if err := ReadObject(file, &TempMBR, 0); err != nil {
			fmt.Println("Error: No se pudo leer el MBR  fdisk -add")
			return agrego
		}

		//Buscar el nombre de la particion
		Se_repite := false
		for i := 0; i < 4; i++ {
			if TempMBR.Partitions[i].Size != 0 {
				nombre_particion := strings.TrimRight(string(TempMBR.Partitions[i].Name[:]), "\x00")
				if nombre_particion == name {
					//Los nombres de las particiones se repiten
					Se_repite = true
					posicion_par = i

				}
			}
		}
		//validar la particion
		if Se_repite {
			//Existe particion
			//1.pasar la cantidad de add a k o M
			if unit == "k" {
				add = add * 1024
			} else {
				add = add * 1024 * 1024
			}
			//Validar si es negativo o positivo
			if add > 0 {
				positivo = true
			} // caso contrario false

			nuevo_espacio := TempMBR.Partitions[posicion_par].Start + TempMBR.Partitions[posicion_par].Size + int32(add)
			//Validar si cabe o se puede quitar
			if positivo {
				//se agrega espacio
				if nuevo_espacio < TempMBR.MbrSize {
					//Si cabe en el disco
					TempMBR.Partitions[posicion_par].Size += int32(add)
					//Rescribir el MBR
					if err := WriteObject(file, TempMBR, 0); err != nil {
						return agrego
					}
					agrego = true
				} else {
					//No cabe en el disco
					fmt.Println("Error: Espacio insuficiente fdisk -add")
					return agrego
				}
			} else {
				//Se quia espacio
				if nuevo_espacio > TempMBR.Partitions[posicion_par].Start {
					//Si puedo quitar esa cantidad
					TempMBR.Partitions[posicion_par].Size += int32(add)
					//Rescribir el MBR
					if err := WriteObject(file, TempMBR, 0); err != nil {
						fmt.Println("Error: No se pudo reescribir MBR  fdisk -add")
						return agrego
					}
					agrego = true
				} else {
					fmt.Println("Error: No se puede restar esa cantidad  fdisk -add -")
					return agrego
				}
			}

		} else {
			//NO existe la particion
			fmt.Println("Error: No se encontro la particion  fdisk -add")
			return agrego
		}
	} else {
		//No encontro disco
		fmt.Println("Error: No se encontro el disco  fdisk -add")
		return agrego
	}

	return agrego
}
