'use client'
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
      
      <div onClick={() => handleOpen()}>
        <div className="AgendaItem__date">
          <div className="AgendaItem__day">
            {day}
          </div>
          <div className="AgendaItem__time">
            {time}
          </div>
        </div>
        <div className="AgendaItem__data">
          <h3>{name}</h3>
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