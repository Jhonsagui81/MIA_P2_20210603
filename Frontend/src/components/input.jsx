import React, {useState} from "react";
import axios from 'axios';

const EntradaTexto = () => {
  const [Nombre, setNombre] = useState(''); // Estado para el valor del nombre

  const handleSubmit = async () => {

        try {
            const url = "http://localhost:3000/comand"
            // const response = await axios.get(url);
            console.log("Enviando al backend:", {Nombre})
            const response = await axios.post(url, { Nombre }); // Enviar el nombre al backend
            console.log(response)
        } catch (error) {
          console.error('Error al enviar el nombre:', error);
          alert('Error al enviar el nombre.');
        }
    };

    return (
        <div className="flex flex-col items-center">
            <label htmlFor="Nombre" className="text-lg font-bold mb-2">Nombre:</label>
            <input
                type="text"
                id="Nombre"
                name="=Nombre"
                value={Nombre}
                onChange={(event) => setNombre(event.target.value)}
                className="border rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
                type="submit"
                className="mt-4 bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600"
                onClick={handleSubmit}
            >
                Enviar
            </button>
        </div>
    );
};

export default EntradaTexto;