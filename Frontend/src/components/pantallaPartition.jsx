import React, { useState } from "react";

const Partition = ({ discoData, toggleDatalle }) => {
  // Desestructurar los datos del disco
  const { nombre } = discoData;

  return (
    <div className="w-48 h-32 bg-white rounded-md shadow-md mt-9 m-2 hover:bg-gray-200 hover:shadow-lg">
      <img src='https://c0.klipartz.com/pngpicture/322/736/gratis-png-particion-de-disco-gparted-editor-de-particiones-gnu-parted-tutuapp-brisa.png' alt={nombre} className="w-full h-24 rounded-t-md" />
      <div className="p-2">
        <p className="text-gray-700 font-medium">Paricion: {discoData}</p>
      </div>
    </div>
  );
};

export default Partition;