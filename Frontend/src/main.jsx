import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import Button from "./components/button.jsx"
import './index.css'
import EntradaTexto from './components/input.jsx'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
    <EntradaTexto></EntradaTexto>
    <Button name="jonas"></Button>
  </React.StrictMode>,
)
