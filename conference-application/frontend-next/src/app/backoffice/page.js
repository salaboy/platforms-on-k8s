'use client'
import styles from '@/app/styles/backend.module.css'
import { useState, useEffect } from 'react'
import ProposalList from '../components/proposals/proposallist'
import NotificationList from '../components/notifications/notificationlist'
import EventsList from '../components/events/eventlist'


export default function Backoffice() {

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
        
          <ProposalList/>
          
        
      </div>

      <h2>Notifications (Tab)</h2>
      <div>
        <NotificationList/>
      </div>

      <h2>Events (Tab)</h2>
      <div>
        <EventsList></EventsList>
      </div>
          
        </div>
      </div>
      
     

    </main>
  )
}
