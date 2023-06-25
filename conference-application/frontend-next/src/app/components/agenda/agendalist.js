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
        "title": "Cached Entry",
        "author": "Cached Author",
        "time": "1pm",
        "day": "Monday",
        "description": "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Id officia doloribus, molestiae, mollitia quia maiores velit consequuntur dolorem labore beatae, porro aliquam quis! Quasi commodi aperiam, assumenda rem molestiae porro."
    }]
    
    const fetchData = () => {
        console.log("Querying /agenda/")
        fetch('/api/agenda/')
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
        fetch('/api/agenda/' + id , {
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
                        name={item.Title}
                        time={item.Time}
                        key={index}
                        id={item.Id}
                        description={item.Description}
                        author={item.Author}
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