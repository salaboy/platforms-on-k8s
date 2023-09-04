
'use client'
import styles from './switch.module.css'




export default function Switch({label, id, name, value}) {
    
    return (
        <div className={styles.switch} >
            <label>{label}
           
                <input type="checkbox" value={value} id={id}/>
                <span class={styles.slider}></span>
            </label>
        </div>
            
     
    
    );
}

