'use client'

import styles from '@/app/styles/backoffice.module.css'
import { useState, useEffect } from 'react'

import ProposalList from '../components/proposals/proposallist'
import NotificationList from '../components/notifications/notificationlist'
import EventsList from '../components/events/eventlist'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import AgendaList from '../components/agenda/agendalist'



export default function Backoffice() {




  return (
    <main className={styles.main}>
      <div className={`${styles.hero} ` }>
        <div className={ `grid content noMargin`}>
          <div className="col full">
            <h3>Backoffice</h3>
          </div>
        </div>
      </div>
      <div className={`${styles.BackofficeContent} ` }>
      <div className={ `grid content noMargin`}>
        <div className="col full">
          <div className={`${styles.tabs} ` }>
          <Tabs>
            <TabList>
              <Tab>Review Proposals</Tab>
              <Tab>Notifications</Tab>
              <Tab>Events</Tab>
              <Tab>Agenda Items</Tab>
            </TabList>

            <TabPanel>
              <ProposalList></ProposalList>
            </TabPanel>
            <TabPanel>
              <NotificationList />
            </TabPanel>
            <TabPanel>
              <EventsList />
            </TabPanel>
            <TabPanel>
              <AgendaList admin="true" />
            </TabPanel>
          </Tabs>
          </div>
        </div>
      </div>
      </div>
     


    </main>
  )
}
