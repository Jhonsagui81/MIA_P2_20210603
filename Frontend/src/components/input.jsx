import React, {useState} from "react";
import axios from 'axios';

const EntradaTexto = () => {
  const [nombre, setNombre] = useState(''); // Estado para el valor del nombre

  const handleSubmit = async (event) => {
        event.preventDefault(); // Evitar la recarga de la página

        try {
          const respuesta = await axios.post('/api/enviarNombre', { nombre }); // Enviar el nombre al backend
          console.log('Nombre enviado:', respuesta.data); // Mostrar la respuesta del backend
          alert('Nombre enviado correctamente!');
        } catch (error) {
          console.error('Error al enviar el nombre:', error);
          alert('Error al enviar el nombre.');
        }
    };

    return (
        <div className="flex flex-col items-center">
            <label htmlFor="nombre" className="text-lg font-bold mb-2">Nombre:</label>
            <input
                type="text"
                id="nombre"
                name="nombre"
                value={nombre}
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