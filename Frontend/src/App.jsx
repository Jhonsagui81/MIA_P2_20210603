import React from 'react';
import EntradaTexto from './components/input';

const Opcion = ({ nombre, onClick }) => {
  return (
    <button className="py-2 px-4 rounded" onClick={onClick}>
      {nombre}
    </button>
  );
};

const PaginaPrincipal = () => {
  const [pantallaActiva, setPantallaActiva] = React.useState('Pantalla1');

  const handlePantallaClick = (pantalla) => {
    setPantallaActiva(pantalla);
  };

  return (
    <div className="w-full flex h-full">
      <div className="w-full max-w-[300px] bg-blue-500 flex flex-col items-stretch">
        <p>asdfhkjashdfkjhas</p>
        <p>Otro div</p>
      </div>
      <div className="w-full bg-gray-200 h-full">
        {/*Se muestra la pantalla seleeccionada utilizando ifs */}
        <h2>contenido</h2>
      </div>
    </div>
  );
};

const App = () => {
  return (
    <div>
      <PaginaPrincipal />
    </div>
  );
};

export default App;

        // {pantallaActiva === 'Pantalla1' && <Pantalla1 />}
        // {pantallaActiva === 'Pantalla2' && <Pantalla2 />}
        // {pantallaActiva === 'Pantalla3' && <Pantalla3 />}
        // <Opcion nombre="Pantalla1" onClick={() => handlePantallaClick('Pantalla1')} />
        // <Opcion nombre="Pantalla2" onClick={() => handlePantallaClick('Pantalla2')} />
        // <Opcion nombre="Pantalla3" onClick={() => handlePantallaClick('Pantalla3')} />