import React, {useState} from "react";
import axios from 'axios'

const Pantalla1 = ({ip}) => {
    const [Nombre, setNombre] = useState('');
    const [RespuestaBackend, setRespuestaBackend] = useState('')

    const handleSubmit = async () => {
        try{
            const url = `http://${ip}:3000/input`
            console.log("Enviando al backend:", {Nombre})
            const response = await axios.post(url, {Nombre});
            console.log(response)
            setRespuestaBackend(Nombre)
        } catch {
            console.error('Error al enviar el nombre:', error);
            alert('Error al enviar el nombre.');
        }
    }

      //http://localhost:3000

    return (
        <div className="p-4">
          <h2>Contenido de la Pantalla 1</h2>
          <div className="w-full bg-gray-200 border border-gray-300 rounded-md p-2 mb-4" style={{ minHeight: '85vh' , whiteSpace: 'pre-wrap' }}>
            <p className="text-gray-700">{RespuestaBackend}</p>
          </div>
          <div className="flex w-full">
          <textarea
            className="flex-grow resize-none border border-gray-300 rounded-md p-2"
            placeholder="Escribe aquí"
            id="Nombre"
            name="Nombre"
            value={Nombre}
            onChange={(event) => setNombre(event.target.value)}
          />
            <button 
                type="submit"
                className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600"
                onClick={handleSubmit}
                >Enviar</button>
          </div>
        </div>
      );
};

export default Pantalla1