import React, {useState, useEffect} from "react";
import { BrowserRouter as Router, Route } from 'react-router-dom';
import axios from 'axios'

const Pantalla3 = () => {
    const [showDetalle, setShowDetalle] = useState('pdf');
    const [contenidoPdf, setContenidoPdf] = useState('')
    const [RespuestaBackend, setRespuestaBackend] = useState([])

    useEffect(() => {
        // Carga los datos del backend al montarse la pantalla
        axios.get('http://localhost:3000/reportes') // Ajusta la URL del endpoint
        // .then((response) => console.log(response.data))  
        .then((response) => setRespuestaBackend(response.data))
          
          .catch((error) => console.error(error));
      }, []);

    function handleEstado (contenido) {
        setContenidoPdf(contenido)
        console.log(contenidoPdf)
        const url = "https://quickchart.io/graphviz";

        fetch(url, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            graph: contenido,
          }),
        })
        .then((response) => response.blob())
        .then((blob) => {
          const urlImagen = URL.createObjectURL(blob);
          // Hacer algo con la URL de la imagen
          
          window.open(urlImagen);
        });
    }

    const PdfViewer = ({ pdfBase64 }) => {
        return <PDFViewer pdfURL={`data:application/pdf;base64,${pdfBase64}`} />;
    };

    function obtenerContenido(caso){
        switch (caso){
            case "pdf":
                return (
                    <div className="flex flex-wrap justify-between">
                       {RespuestaBackend.map((parti) => (
                           <div key={parti.Id}>
                             <PDF 
                               nombre={parti.Nombre}
                               imagen={parti.Imagen}
                             />
                               <button type="submit" onClick={() => handleEstado(parti.Content)} className="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600">Abrir PDF</button>
                           </div>
                       ))}
                    </div>);
            case "showContent":
                return (
                    <div className="min-h-screen flex items-center justify-center">
                        <Router>
                          <Route path="/">
                            <ShowPdf dotContent={contenidoPdf} />
                          </Route>
                        </Router>
                    </div>
                );
        }
    }
    
    return (
        <div>
            {obtenerContenido(showDetalle)}
        </div>

    );
};


const PDF = ({ nombre, imagen }) => {
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

  const ShowPdf = ({dotContent}) => {
    const dotUrl = `https://quickchart.io/graphviz?graph=${encodeURIComponent(dotContent)}`;

    return (
        <div>
            <img src={dotUrl} alt="Graph" />
        </div>
    );
  }

export default Pantalla3