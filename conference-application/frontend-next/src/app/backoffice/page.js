'use client'
import styles from '@/app/styles/backend.module.css'
import { useState, useEffect } from 'react'
import ProposalList from '../components/proposals/proposallist'


export default function Backend() {

  const [proposals, setProposals] = useState(null)
  const [events, setEvents] = useState(null)
  const [notifications, setNotifications] = useState(null)
  const [isLoadingProposals, setLoadingProposals] = useState(false)
  const [isLoadingEvents, setLoadingEvents] = useState(false)
  const [isLoadingNotifications, setLoadingNotifications] = useState(false)
  const [isError, setIsError] = useState(false);
  const [decisionsMade, setDecisionsMade] = useState(1)
  const [check, setCheck] = useState(0)

  useEffect(() => {
    const id = setInterval(() => {
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
        setCheck(check + 1)
    }, 3000);
    return () => clearInterval(id);
  }, [check])

 





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
          <ProposalList/>
          
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
