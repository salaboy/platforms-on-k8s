'use client'
import styles from '@/app/styles/agenda.module.css'
import { useState, useContext } from 'react'

function AgendaItem({key, name, day, time, author, description}) {
    const [open, setOpen] = useState(false) // state hook

    const handleOpen = () => {
      if(open){
        setOpen(false);
      }else {
        setOpen(true);
      }


    }

    return (
      
      <div onClick={() => handleOpen()} className={styles.agendaItem}>
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
          {open && (
            <div className="AgendaItem__description">
              <p>{description}</p>
            </div>
          )}
        </div>
      </div>
      
    );

}
export default AgendaItem;