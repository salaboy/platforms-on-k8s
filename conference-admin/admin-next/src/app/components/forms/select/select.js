
'use client'
import styles from './select.module.css'




export default function Select({label, id, name, value, children}) {
    
    return (
        <div className={styles.select}>
            <label>{label}</label>
            <div className={styles.selectContainer}>
            <select name={name} id={id}>
                {children}
            </select>
            </div>
            
        </div>
            
     
    
    );
}

