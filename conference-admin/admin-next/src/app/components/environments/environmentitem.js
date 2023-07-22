
'use client'
import styles from "@/app/styles/environments.module.css"
import { useState } from 'react'

function EnvironmentItem({id, name, type,  installInfra, debug}) {
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
          
          {/* <div className={styles.headerStatus}>
            {approved && (
              <div className={`${styles.headerStatusTag}  ${ styles.approved } ` }>Approved</div>
            )}
            {!approved && (
              <div className={`${styles.headerStatusTag}  ${ styles.rejected } ` }>Rejected</div>
            )}
          </div> */}
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
      
        </div>
        
      </div>
      
    );

}
export default EnvironmentItem;