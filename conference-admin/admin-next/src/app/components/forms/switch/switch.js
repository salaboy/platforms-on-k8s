
'use client'
import styles from './switch.module.css'




export default function Switch({label, id, name, value}) {
    
    return (
        <div className={styles.switch} id={id}>
            <label>{label}
           
                <input type="checkbox" value={value}/>
                <span class={styles.slider}></span>
            </label>
        </div>
            
     
    
    );
}

