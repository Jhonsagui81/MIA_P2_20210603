import React, {useState, useEffect} from "react";
import axios from 'axios'
import Partition from "./pantallaPartition";

const Pantalla2 = () => {
    //recorrer el arreglo de discos
    const [discos, setDiscos] = useState([]);
    //Disco seleccionado 
    const [selectedDisk, setSelectDisk] = useState(null);
    //Estado a mostrar en pantalla 2
    const [showDetalle, setShowDetalle] = useState("disco");
    //Parition seleccionadad para login
    const [selectedPartition, setPartition] = useState('')

    //Para pedir discos 
    useEffect(() => {
        // Carga los datos del backend al montarse la pantalla
        axios.get('http://localhost:3000/discos') // Ajusta la URL del endpoint
        // .then((response) => console.log(response.data))  
        .then((response) => setDiscos(response.data))
          
          .catch((error) => console.error(error));
      }, []);

    //Para Pedir particiones  
    const handleCardClick = (nombreDisco) => {
        console.log(`**Se hizo clic en la tarjeta con nombre: ${nombreDisco}**`);
        axios.post('http://localhost:3000/partitions', {Nombre: nombreDisco})
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
    function handleLoginSubmit(event){
        event.prevenDefault();
        const username = event.target.username.value;
        const password = event.target.password.value;
        const partition = selectedPartition;

        //Construir el comando
        const comand = "login -user="+username+" -pass="+password+" -id="+partition

        //Hacer peticion post
        axios.post('http://localhost:3000/login', {Nombre: comand})
        .then((response) =>{
            setShowDetalle("sistema")
            //Respondera con los archivos y carpetas de la carpeta root-> Guardarla en un arreglo como discos

        })
        .catch((error) => {
            console.error('Error en la Peticion POST', error)
        })
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
                    <form onSubmit={}>
                        <div className="bg-gray-100 min-h-screen flex items-center justify-center">
                            <div className="bg-white p-8 rounded-lg shadow-md w-96">
                              <button
                                onClick={() => handleEstado("particion", "")}
                                className="absolute top-2 right-2 text-gray-600 hover:text-gray-800"
                              >
                                Regresar
                              </button>
                              <h2 className="text-2xl font-semibold mb-4">¡Bienvenido!</h2>
                              <p className="text-gray-600 mb-6">Estás iniciando sesión en la particion {selectedPartition}</p>
                              <form>
                                <div className="mb-4">
                                  <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                                    Usuario
                                  </label>
                                  <input
                                    type="text"
                                    id="username"
                                    name="username"
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
                                    className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                                    placeholder="Ingresa tu contraseña"
                                  />
                                </div>
                                <button
                                  type="submit"
                                  className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600"
                                >
                                  Iniciar sesión
                                </button>
                              </form>
                            </div>
                        </div>
                    </form>
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
  

export default Pantalla2