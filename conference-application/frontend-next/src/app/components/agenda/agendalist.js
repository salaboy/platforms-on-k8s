'use client'
import styles from '@/app/styles/agenda.module.css'
import { useState, useEffect } from 'react'
import AgendaItem from './agendaitem'

function AgendaList(props) {

    
    const [isError, setIsError] = useState(false);
    const [agendaItems, setAgendaItems] = useState('') // state hook
    const {day, highlights, admin} = props;
    const [isLoading, setLoading] = useState(false)
    
    const mockAgendaItems = [{
        "id": "ABC-123",
        "title": "Cached Entry",
        "author": "Cached Author",
        "description": "Cached Content"
    }]
    
    const fetchData = () => {
        console.log("Querying /agenda/agenda-items/")
        fetch('/api/agenda/agenda-items/')
        .then((res) => res.json())
        .then((data) => {
            setAgendaItems(data)
            setLoading(false)
        }).catch((error) => {    
            setAgendaItems(mockAgendaItems)
            console.log(error)
        })
    };

    const handleArchive = (id) => {
        setLoading(true);
        setIsError(false);
        console.log("Archiving Agenda Item ..." + id)
        fetch('/api/agenda/agenda-items/' + id , {
          method: "DELETE",
          headers: {
            'accept': 'application/json',
          },
        }).then((response) => response.json()).then(() => {
          fetchData()
          setLoading(false);
        }).catch(err => {
          console.log(err);
          setLoading(false);
          setIsError(true);
        });
    
      }



    useEffect(() => {                           // side effect hook
        setLoading(true)
        fetchData()
  
    }, [setAgendaItems])

    return (
        <div>
            <div className={`${styles.agendaList}  ${admin ? styles.backoffice : ' '} ` }>
                {agendaItems && agendaItems.length > 0 && agendaItems.map((item, index) => (

                    <AgendaItem
                        name={item.title}
                        key={index}
                        id={item.id}
                        description={item.description}
                        author={item.author}
                        admin={admin}
                        handleArchive={handleArchive}
                    />


                ))}
                {agendaItems && agendaItems.length == 0 && (
                    <p>
                            There are no confirmed talks just yet.
                    </p>
                )}
            </div>

        </div>
    );

}

export default AgendaList;