'use client'
import styles from '@/app/styles/backend.module.css'
import { useState, useEffect } from 'react'


export default function Backend() {

  const [proposals, setProposals] = useState(null)
  const [events, setEvents] = useState(null)
  const [notifications, setNotifications] = useState(null)
  const [isLoadingProposals, setLoadingProposals] = useState(false)
  const [isLoadingEvents, setLoadingEvents] = useState(false)
  const [isLoadingNotifications, setLoadingNotifications] = useState(false)
  const [isError, setIsError] = useState(false);
  const [decisionsMade, setDecisionsMade] = useState(1)

  useEffect(() => {
    setLoadingProposals(true)
    fetch('/api/c4p/')
      .then((res) => res.json())
      .then((data) => {
        setProposals(data)
        setLoadingProposals(false)
      })
    setLoadingEvents(true)
    fetch('/api/events/')
      .then((res) => res.json())
      .then((data) => {
        setEvents(data)
        setLoadingEvents(false)
      })
    fetch('/api/notifications/')
      .then((res) => res.json())
      .then((data) => {
        setNotifications(data)
        setLoadingNotifications(false)
      })  
  }, [])


  // if (isLoadingProposals) return <p>Loading...</p>
  // if (!proposals) return <p>No Proposals</p>

  // if (isLoadingEvents) return <p>Loading...</p>
  // if (!events) return <p>No Events</p>



  function decide(id, approved) {
    const decision = {
      approved: approved,

    }
    try {
      setLoadingProposals(true);
      fetch('/api/c4p/' + id + "/decide", {
        method: "POST",
        body: JSON.stringify(decision),
        headers: {
          'accept': 'application/json',
        },
      }).then((response) => response.json()).then((data) => {
        setDecisionsMade(decisionsMade + 1)
        setLoadingProposals(false);
      }
      )
    } catch (err) {
      setLoadingProposals(false);
      setIsError(true);
    }

  }


  return (
    <main className={styles.main}>
      <div className="grid content">
        <div className="col full">
          <h1>Backend</h1>
        </div>
      </div>
      <div className="grid content">
        <div className="col third">
          <ul className={styles.tabs}>
            <li className={styles.tabItem}>
              <a >Review Proposals</a>
            </li>
            <li className={styles.tabItem}>
              <a >Notifications</a>
            </li>
            <li className={styles.tabItem}>
              <a >Events</a>
            </li>
          </ul>
          
        </div>

        <div className="col half">
        <h2>Review Proposals (Tab)</h2>
      <div>
        <ul>
          {proposals === null && ((
            <p>No Proposals</p>
          ))}
          {proposals !== null && proposals.map((p) => (
            <li key={p.Id}>{p.Id} - {p.Title} - {p.Description} - {p.Author} - {p.Email}  - {p.Status.Status}  - {p.Approved.toString()}
              <button main onClick={() => decide(p.Id, true)} >Approve</button>
              <button main onClick={() => decide(p.Id, false)}>Reject</button>
            </li>

          ))}
        </ul>
      </div>

      <h2>Notifications (Tab)</h2>
      <div>
        <ul>
          {notifications === null && ((
            <p>No Notifications</p>
          ))}
          {notifications !== null && notifications.map((notif) => (
            <li key={notif.Id}>{notif.EmailText}</li>
          ))}
        </ul>
      </div>

      <h2>Events (Tab)</h2>
      <div>
        <ul>
          {events === null && ((
            <p>No Events</p>
          ))}
          {events !== null && events.map((e) => (
            <li key={e}>{e}</li>
          ))}
        </ul>
      </div>
          
        </div>
      </div>
      
     

    </main>
  )
}
