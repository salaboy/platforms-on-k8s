
'use client'
import styles from "@/app/styles/environments.module.css"
import { useState } from 'react'
import Button from '../forms/button/button';

function EnvironmentItem({id, name, type,  installInfra, debug, status, synced, vclusterRef, secretRef, handleDelete}) {
  const [open, setOpen] = useState(false) // state hook
  const handleOpen = () => {
    if(open){
      setOpen(false);
    }else {
      setOpen(true);
    }
  }

  
  
    return (
      
      <div onClick={() => handleOpen()} className={`${styles.EnvironmentItem}  ${open ? styles.open : ' '} ` }>
        <div className={styles.openTag}>
          {!open && (
            <>Click for details</>
          )}
          {open && (
            <>Close</>
          )}
        </div>
        <div className={styles.header}>
          <h5> <span>Environment:</span> {name}</h5>
          
          <div className={styles.headerStatusTags}>
            
              <div className={`${styles.headerStatusTag}  ${ styles.approved } ` }>Synced: {synced}</div>
            
              <div className={`${styles.headerStatusTag}  ${ styles.approved } ` }>Ready: {status} </div>
          
          </div> 
        </div>
        <div className={styles.description}>
          <div className={styles.descriptionTo}>
            <span>Type:</span> {type}
          </div>

          <div className={styles.descriptionSubject}>
            <span> Install Infrastructure:</span> {installInfra.toString()}
          </div>
          <div className={styles.descriptionSubject}>
            <span> Frontend Debug:</span> {debug.toString()}
          </div>

          <div className={styles.descriptionSubject}>
            {status != "True" && (<p>Waiting for the Environment to be Ready.</p>) }
            {status == "True" && vclusterRef != null && (<p>Connect to this environment running <b>`vcluster connect {vclusterRef} --server https://localhost:8443 -- zsh`</b> </p>)}
            {status == "True" && secretRef != null && (<p>Use the secret called <b>`{secretRef}`</b> to connect</p>)}
            
          </div>
          <div className={styles.descriptionAction}>
              <Button clickHandler={() => handleDelete(name)}>Delete</Button>
          </div>
      
        </div>
        
      </div>
      
    );

}
export default EnvironmentItem;