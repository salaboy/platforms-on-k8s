
'use client'
import styles from '@/app/styles/events.module.css'
import { useState } from 'react'
import JSONPretty from 'react-json-pretty';

function EventItem({id, type, payload}) {
  const [open, setOpen] = useState(false) // state hook
  const handleOpen = () => {
    if(open){
      setOpen(false);
    }else {
      setOpen(true);
    }
  }

    return (
      
      <div onClick={() => handleOpen()} className={`${styles.EventItem}  ${open ? styles.open : ' '} ` }>
        <div className={styles.openTag}>
          {!open && (
            <>Click for details</>
          )}
          {open && (
            <>Close</>
          )}
        </div>
        <div className={styles.header}>
          <h5><span>#{id}</span>  {type}</h5>
        </div>
          {/* Maybe render using: https://www.npmjs.com/package/react-json-pretty */}
          <div className={styles.description}>
            <div className={styles.codeContainer}>
              <JSONPretty id="json-pretty" data={payload}></JSONPretty>
            </div>
          </div>
        
        
        
      </div>
      
    );

}
export default EventItem;