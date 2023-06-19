
'use client'
import styles from './textfield.module.css'




export default function Textfield({label, id, name, value}) {
    
    return (
        <div className={styles.textfield}>
            <label>{label}</label>
            <input type="text" id={id} name={name} value={value}/>
        </div>
            
     
    
    );
}

