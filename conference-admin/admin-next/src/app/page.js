'use client'

import styles from '@/app/styles/backoffice.module.css'
import { useState, useEffect } from 'react'

import EnvironmentList from './components/environments/environmentlist'
import NewEnvironment from './components/newenvironment/newenvironment'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';

export default function Backoffice() {
  const [isLoading, setLoading] = useState(false)



  useEffect(() => {                           // side effect hook
    setLoading(true)
    setLoading(false)
   
  }, [])


  return (
    <main className={styles.main}>
      <div className={`${styles.hero} `}>
        <div className={`grid content noMargin`}>
          <div className="col full">
            <h3>Admin</h3>
          </div>
        </div>
      </div>
      <div className={`${styles.BackofficeContent} `}>
        <div className={`grid content noMargin`}>
          <div className="col full">
            <div className={`${styles.tabs} `}>
              <Tabs>
                <TabList>
                  
                  
                  <Tab>Environments</Tab>
                  <Tab>New Environment</Tab>
                </TabList>

                
                <TabPanel>
                  <EnvironmentList />
                </TabPanel>
                <TabPanel>
                  <NewEnvironment />
                </TabPanel>
              </Tabs>
            </div>
          </div>
        </div>
      </div>



    </main>
  )
}
