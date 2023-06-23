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
              <p>Lorem ipsum, dolor sit amet consectetur adipisicing elit. Autem nam sit minus quibusdam nisi voluptatem earum eum ipsam sunt consequuntur odit neque libero, modi ut officiis, dignissimos rerum at facere!</p>
            </div>
         
        </div>
      </div>
      
    );

}
export default AgendaItem;