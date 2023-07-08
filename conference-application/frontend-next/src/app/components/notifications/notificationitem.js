
'use client'
import styles from "@/app/styles/notifications.module.css"
import { useState } from 'react'

function NotificationItem({id, title, emailTo,  emailSubject, emailBody, approved}) {
  const [open, setOpen] = useState(false) // state hook
  const handleOpen = () => {
    if(open){
      setOpen(false);
    }else {
      setOpen(true);
    }
  }
  
    return (
      
      <div onClick={() => handleOpen()} className={`${styles.NotificationItem}  ${open ? styles.open : ' '} ` }>
        <div className={styles.openTag}>
          {!open && (
            <>Click for details</>
          )}
          {open && (
            <>Close</>
          )}
        </div>
        <div className={styles.header}>
          <h5> <span>Proposal:</span> {title}</h5>
          
          <div className={styles.headerStatus}>
            {approved && (
              <div className={`${styles.headerStatusTag}  ${ styles.approved } ` }>Approved</div>
            )}
            {!approved && (
              <div className={`${styles.headerStatusTag}  ${ styles.rejected } ` }>Rejected</div>
            )}
          </div>
        </div>
        <div className={styles.description}>
          <div className={styles.descriptionTo}>
            <span>To:</span> {emailTo}
          </div>

          <div className={styles.descriptionSubject}>
            <span> Subject:</span> {emailSubject}
          </div>
          <div className={styles.descriptionBody}>
            <p>
             {emailBody}
            </p>
          </div>
      
        </div>
        
      </div>
      
    );

}
export default NotificationItem;