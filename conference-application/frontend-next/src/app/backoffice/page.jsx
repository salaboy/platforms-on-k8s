'use client'

import styles from '@/app/styles/backoffice.module.css'
import { useState, useEffect } from 'react'

import ProposalList from '../components/proposals/proposallist'
import NotificationList from '../components/notifications/notificationlist'
import EventsList from '../components/events/eventlist'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import AgendaList from '../components/agenda/agendalist'
import Debug from '../components/debug/debug'



export default function Backoffice() {
  const [isLoading, setLoading] = useState(false)
  const [features, setFeatures] = useState('') // state hook

  const fetchFeatures = () => {
    setLoading(true);
    console.log("Querying /api/features/")
    fetch('/api/features/')
      .then((res) => res.json())
      .then((data) => {
        console.log("Features Data: " + data)
        setFeatures(data)
        setLoading(false)
      }).catch((error) => {
        setFeatures({})
        console.log(error)
      })
  };

  useEffect(() => {                           // side effect hook
    setLoading(true)
    fetchFeatures()
   
  }, [])


  return (
    <main className={styles.main}>
      <div className={`${styles.hero} `}>
        <div className={`grid content noMargin`}>
          <div className="col full">
            <h3>Backoffice</h3>
          </div>
        </div>
      </div>
      <div className={`${styles.BackofficeContent} `}>
        <div className={`grid content noMargin`}>
          <div className="col full">
            <div className={`${styles.tabs} `}>
              <Tabs>
                <TabList>
                  <Tab>Review Proposals</Tab>
                  <Tab>Agenda Items</Tab>
                  <Tab>Notifications</Tab>
                  <Tab>Events</Tab>
                  {features.DebugEnabled == "true" && (<Tab>Debug</Tab>)}
                </TabList>

                <TabPanel>
                  <ProposalList></ProposalList>
                </TabPanel>
                <TabPanel>
                  <AgendaList admin="true" />
                </TabPanel>
                <TabPanel>
                  <NotificationList />
                </TabPanel>
                <TabPanel>
                  <EventsList />
                </TabPanel>
                {features.DebugEnabled == "true" && (<TabPanel>
                  <Debug />
                </TabPanel>)}
              </Tabs>
            </div>
          </div>
        </div>
      </div>



    </main>
  )
}
