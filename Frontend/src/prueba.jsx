import React, { useState } from 'react';
import Pantalla1 from './components/pantalla1';
import Pantalla2 from './components/pantalla2';
import Pantalla3 from './components/pantalla3';

const Menu = () => {
  const [pantallaActiva, setPantallaActiva] = useState('Pantalla1');

  const handlePantallaClick = (pantalla) => {
    setPantallaActiva(pantalla);
  };

  return (
    <div className="flex  h-screen">
      <div className="w-1/4 bg-blue-500 flex flex-col items-stretch order-first">
        <button onClick={() => handlePantallaClick('Pantalla1')} className="text-white text-lg font-medium py-2 px-4 border-b border-transparent hover:bg-blue-600">Pantalla 1</button>
        <button onClick={() => handlePantallaClick('Pantalla2')} className="text-white text-lg font-medium py-2 px-4 border-b border-transparent hover:bg-blue-600">Pantalla 2</button>
        <button onClick={() => handlePantallaClick('Pantalla3')} className="text-white text-lg font-medium py-2 px-4 hover:bg-blue-600">Pantalla 3</button>
      </div>
      <div className="w-full bg-gray-200 flex-grow">
        {/* Se muestra la pantalla seleccionada utilizando condicionales */}
        {pantallaActiva === 'Pantalla1' && <Pantalla1 />}
        {pantallaActiva === 'Pantalla2' && <Pantalla2 />}
        {pantallaActiva === 'Pantalla3' && <Pantalla3 />}
      </div>
    </div>
  );
};





const Salida = () => {
  return (
    <div>
      <Menu />
    </div>
  );
};



  export default Salida;

