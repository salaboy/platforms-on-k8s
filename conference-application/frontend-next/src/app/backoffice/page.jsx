'use client'

import styles from '@/app/styles/backend.module.css'
import { useState, useEffect } from 'react'

import ProposalList from '../components/proposals/proposallist'
import NotificationList from '../components/notifications/notificationlist'
import EventsList from '../components/events/eventlist'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';


export default function Backoffice() {




  return (
    <main className={styles.main}>

      <Tabs>
        <TabList>
          <Tab>Review Proposals</Tab>
          <Tab>Notifications</Tab>
          <Tab>Events</Tab>
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
      </Tabs>
     
      
            
         
            
          

      



    </main>
  )
}
