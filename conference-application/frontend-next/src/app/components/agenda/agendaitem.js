'use client'
import styles from '@/app/styles/agenda.module.css'
import { useState, useContext } from 'react'
import Button from '../forms/button/button'

function AgendaItem({key, id, name, day, time, author, description, admin, handleArchive}) {
    const [open, setOpen] = useState(false) // state hook

    const handleAction = (id) => {
      handleArchive(id);
    }

    const handleOpen = () => {
      if(open){
        setOpen(false);
      }else {
        setOpen(true);
      }


    }

    return (
      
      <div onClick={() => handleOpen()} className={`${styles.agendaItem}  ${open ? styles.open : ' '} ` }>
        <div className={styles.openTag}>
          {!open && (
            <>Click for details</>
          )}
          {open && (
            <>Close</>
          )}
        </div>
        <div className="AgendaItem__date">
          <div className="AgendaItem__day">
            {day}
          </div>
          <div className="AgendaItem__time">
            {time}
          </div>
        </div>
        <div className="AgendaItem__data">
          <h4>{name}</h4>
          <p className="p p-s"> {author}</p>
         
            <div className={styles.description} >
              <p>{description}</p>
            </div>
         
        </div>
          {admin && (
            <Button clickHandler={() => handleAction(id)}>Archive</Button>
          )}
      </div>
      
    );

}
export default AgendaItem;