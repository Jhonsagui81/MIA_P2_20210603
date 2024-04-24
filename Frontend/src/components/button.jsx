

const Button = (props) =>{
    const {name} = props;
    return (
        <>
            <button className="bg-primary px-5 py-2 text-white" >{name}</button>
        </>
    )
}

const Opcion = ({ nombre, onClick }) => {
    return (
      <button className="py-2 px-4 rounded" onClick={onClick}>
        {nombre}
      </button>
    );
  };

export default Button