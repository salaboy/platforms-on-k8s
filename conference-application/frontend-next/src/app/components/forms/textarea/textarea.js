
'use client'
import styles from './textarea.module.css'




export default function Textarea({label, id, name, value}) {
    
    return (
        <div className={styles.textarea}>
            <label>{label}</label>
            <textarea name={name}  id={id}  cols="30" rows="4" value={value}></textarea>
            
        </div>
            
     
    
    );
}

