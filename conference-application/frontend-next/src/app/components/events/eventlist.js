'use client'
import { useState, useEffect } from 'react'
import EventItem from './eventitem'



function EventsList() {
    const [loading, setLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [events, setEvents] = useState([]) // state hook
    const [check, setCheck] = useState(0)

    const fetchData = () => {

        fetch('/api/events/')
            .then((res) => res.json())
            .then((data) => {
                console.log("Fetching Events ...")
                setEvents(data)
                setLoading(false)
            }).catch((error) => {
                console.log(error)
            })
    }

    
    useEffect(() => {

        setLoading(true)
        fetchData();

    }, [])

    useEffect(() => {
        const id = setInterval(() => {
            setLoading(true)
            fetchData();

        }, 3000);
        return () => clearInterval(id);
    }, [check])



    return (
        <div>


            <div>
                {
                    events && events.map((item, index) => (
                        <EventItem
                            key={item.Id}
                            id={item.Id}
                            type={item.Type}
                            payload={item.Payload}

                        />

                    ))
                }
                {
                    events && events.length === 0 && (
                        <span>There are no events.</span>
                    )
                }
            </div>
        </div>
    );

}
export default EventsList;