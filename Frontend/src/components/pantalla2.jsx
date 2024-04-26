import React, {useState, useEffect} from "react";
import axios from 'axios'
import Partition from "./pantallaPartition";

const Pantalla2 = ({ip}) => {
    //recorrer el arreglo de discos
    const [discos, setDiscos] = useState([]);
    //Disco seleccionado 
    const [selectedDisk, setSelectDisk] = useState(null);
    //Estado a mostrar en pantalla 2
    const [showDetalle, setShowDetalle] = useState("disco");
    //Parition seleccionadad para login
    const [selectedPartition, setPartition] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [archivos, setArchivos] = useState([]);
    //Navegacion en el sistema
    const [ruta, setRuta] = useState('/');

    //Para pedir discos 
    useEffect(() => {
        // Carga los datos del backend al montarse la pantalla
        axios.get(`http://${ip}:3000/discos`) // Ajusta la URL del endpoint
        // .then((response) => console.log(response.data))  
        .then((response) => setDiscos(response.data))
          
          .catch((error) => console.error(error));
      }, []);

    //Para Pedir particiones  
    const handleCardClick = (nombreDisco) => {
        console.log(`**Se hizo clic en la tarjeta con nombre: ${nombreDisco}**`);
        axios.post(`http://${ip}:3000/partitions`, {Nombre: nombreDisco})
            .then((response) => {
                console.log('Petición POST exitosa:', response.data);
                setShowDetalle("particion")
                setSelectDisk(response.data)
            })
            .catch((error) => {
                console.error('Error en la peticion POST:', error)
            })
    }

    //Para hacer el post del login
    function handleLoginClick(){
        
        

        //Construir el comando
        const comand = "login -user="+username+" -pass="+password+" -id="+selectedPartition

        //Hacer peticion post
        axios.post(`http://${ip}:3000/login`, {Nombre: comand})
        .then((response) =>{
            setShowDetalle("sistema")
            setArchivos(response.data)
            console.log('Petición POST exitosa:', response.data);
            //Respondera con los archivos y carpetas de la carpeta root-> Guardarla en un arreglo como discos

        })
        .catch((error) => {
            console.error('Error en la Peticion POST', error)
        })
    }

    function handleDesconectar(){
        axios.post(`http://${ip}:3000/comand`, {Nombre: "logout"})
        .then((response) => {
            setShowDetalle("disco")

        })
        .catch((error) => {
            console.error('Error en la peticion POST', error)
        })
    }

    function agregarRuta(subdirectorio) {
        setRuta((ruta) => ruta + subdirectorio+"/");
      }
    
    function quitarRuta() {
      ruta.slice(0, -1);
      setRuta((prevruta) => prevruta.slice(0, ruta.lastIndexOf('/')));
    }

    //Para abrir carpetas o archivos
    function AbrirArchivo(nombre){
        
        agregarRuta(nombre)
        console.log(ruta+nombre)
        //preparar la ruta peticion (ya seria ruta -> hacer un console.log)
        //hacer una peticion mandandole la ruta
        axios.post(`http://${ip}:3000/sistema`, {Nombre: ruta+nombre})
        .then((response) => {
            setArchivos(response.data)
            response.data.forEach((obj) => {
                if (obj.Content == "") {
                  console.log('No hay contenido disponible en este objeto.');
                  setShowDetalle("sistema")

                } else {
                  console.log('Contenido:', obj.Content);
                  setShowDetalle("Contenido")

                }
            });
        })
        .catch((error) => {
            console.error('Error en la peticion POST', error)
        })
        //verificar si la respuesta.data.contenido es vacio -> si es vacio crear el caso mostrarTexto 
        //si no es vacion actualizar setArchivos y mostrar caso sistema
    }

    function handleRegresar(){
        const nuevaRuta1 = ruta.slice(0, -1);

// Asigna el valor de la variable temporal a "ruta"
        setRuta(nuevaRuta1);

        // Realiza la segunda operación de slice y almacena el resultado en otra variable temporal
        const nuevaRuta2 = nuevaRuta1.slice(0, nuevaRuta1.lastIndexOf('/'));

        // Asigna el valor de la segunda variable temporal a "ruta"
        

        if (nuevaRuta2.length === 0) {
            console.log(ruta)
            setRuta("/")
            axios.post(`http://${ip}:3000/sistema`, {Nombre: "/"})
            .then ((response) => {
                setArchivos(response.data)
            
                setShowDetalle("sistema")

            })
            .catch((error) => {
                console.error('Error en la peticion POST', error)
            })
        } else {
            console.log(ruta)
            axios.post(`http://${ip}:3000/sistema`, {Nombre: nuevaRuta2})
            .then ((response) => {
                setArchivos(response.data)
            
                setShowDetalle("sistema")

            })
            .catch((error) => {
                console.error('Error en la peticion POST', error)
            })
        }
        
        const rutaabsoluta = nuevaRuta2+"/"
        setRuta(rutaabsoluta);
    }

    //Establece que pestana mostrar 
    const handleEstado = (estado, namePartition) => {
        setShowDetalle(estado)
        setPartition(namePartition)
    }


    function obtenerContenido (caso) {
        switch (caso) {
            case "disco":

                return (<div className="flex flex-wrap justify-between">
                    {discos.map((disco) => (
                        <div key={disco.Id}>
                          <Tarjeta 
                            nombre={disco.Nombre} 
                            imagen="./disco2.jpeg" 
                            />
                            <button type="submit" onClick={() => handleCardClick(disco.Nombre)} className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600">Detalles</button>
                    
                        </div>
                    ))}
                </div>);
                
            case "particion":
                return (<div className="flex flex-wrap justify-between">
                     <button
                            onClick={() => handleEstado("disco", "")}
                            className="absolute top-2 right-2 text-gray-600 hover:text-gray-800"
                          >
                            Regresar
                          </button>
                    {selectedDisk.map((parti) => (
                        <div key={parti.Id}>
                          <Partition 
                            discoData={parti.Nombre}
                          />
                            <button type="submit" onClick={() => handleEstado("login", parti.Nombre)} className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600">Logear</button>
                        </div>
                    ))}
                </div>);
            
            case "login":
                return (
                    <div className="bg-gray-100 min-h-screen flex items-center justify-center">
                      <div className="bg-white p-8 rounded-lg shadow-md w-96 relative">
                        <button
                          onClick={() => handleEstado("particion", "")}
                          className="absolute top-2 right-2 text-gray-600 hover:text-gray-800"
                        >
                          Regresar
                        </button>
                        <h2 className="text-2xl font-semibold mb-4">¡Bienvenido!</h2>
                        <p className="text-gray-600 mb-6">Estás iniciando sesión en la particion {selectedPartition}</p>
                        <div className="mb-4">
                          <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                            Usuario
                          </label>
                          <input
                            type="text"
                            id="username"
                            name="username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                            placeholder="Ingresa tu usuario"
                          />
                        </div>
                        <div className="mb-4">
                          <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                            Contraseña
                          </label>
                          <input
                            type="password"
                            id="password"
                            name="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                            placeholder="Ingresa tu contraseña"
                          />
                        </div>
                        <button
                          onClick={handleLoginClick}
                          className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600"
                        >
                          Iniciar sesión
                        </button>
                      </div>
                    </div>
                );
            case "sistema":
                return (
                <div className="flex flex-wrap justify-between">
                    <div className="w-full">
                        <input
                          type="text"
                          value={ruta}
                          readOnly
                          className="border rounded-md p-2 px-2 w-1/2"
                        />
                        <button
                          onClick={() => handleRegresar()}
                          className="bg-gray-300 text-gray-700 font-medium px-4 py-2 rounded-md hover:bg-gray-400"
                        >
                          Regresar
                        </button>
                    </div>
                  
                    
                    {archivos.map((archivo) => (
                        <div key={archivo.Id}>
                          <Archivos 
                            nombre={archivo.Nombre} 
                            imagen={archivo.Imagen} 
                            />
                            <button type="submit" onClick={() => AbrirArchivo(archivo.Nombre)} className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600">Abrir</button>
                    
                        </div>
                    ))}
                    <div className="fixed bottom-0 right-0 p-4 bg-gray-200">
                        <button
                          onClick={() => handleDesconectar()}
                          className="bg-red-500 text-white font-medium px-4 py-2 rounded-md hover:bg-red-600"
                        >
                          Desconectar
                        </button>
                    </div>

                </div>);
            case "Contenido":
                return(
                    <div>
                         <div className="w-full">
                            <input
                              type="text"
                              value={ruta}
                              readOnly
                              className="border rounded-md p-2 px-2 w-1/2"
                            />
                            <button
                              onClick={() => handleRegresar()}
                              className="bg-gray-300 text-gray-700 font-medium px-4 py-2 rounded-md hover:bg-gray-400"
                            >
                              Regresar
                            </button>
                        </div>
                        {archivos.map((archivo)=>(
                            <div class="bg-white p-4 rounded-md shadow-md h-1/2 flex flex-col">
                                <textarea
                                  class="flex-grow border rounded-md p-2"
                                  placeholder={archivo.Content}
                                ></textarea>
                            </div>
                        ))}
                         <div className="fixed bottom-0 right-0 p-4 bg-gray-200">
                            <button
                              onClick={() => handleDesconectar()}
                              className="bg-red-500 text-white font-medium px-4 py-2 rounded-md hover:bg-red-600"
                            >
                              Desconectar
                            </button>
                        </div>
                    </div>
                );
            default:
                return <p>No aplica caso</p>
        }
    }

    return (
        <div>
            {obtenerContenido(showDetalle)}
        </div>
        
    );
}

const Tarjeta = ({ nombre, imagen }) => {
    return (
      <div className="bg-white rounded-md shadow-md w-48 h-32 m-2 hover:bg-gray-200 hover:shadow-lg">
        <img
          src="https://cdn-icons-png.flaticon.com/512/689/689331.png"
          alt={imagen}
          className="w-full h-24 rounded-t-md"
        />
        <div className="p-2">
          <p className="text-gray-700 font-medium">{nombre}</p>
        </div>
      </div>
    );
  };
  
  const Archivos = ({ nombre, imagen }) => {
    return (
      <div className="bg-white rounded-md shadow-md w-48 h-32 m-2 hover:bg-gray-200 hover:shadow-lg">
        <img
          src={imagen}
          alt={imagen}
          className="w-full h-24 rounded-t-md"
        />
        <div className="p-2">
          <p className="text-gray-700 font-medium">{nombre}</p>
        </div>
      </div>
    );
  };

const Contenido = ({contenido}) => {
    return(
        <div class="bg-white p-4 rounded-md shadow-md h-1/2 flex flex-col">
          <textarea
            class="flex-grow border rounded-md p-2"
            placeholder={contenido}
          ></textarea>
        </div>
    )
}

export default Pantalla2