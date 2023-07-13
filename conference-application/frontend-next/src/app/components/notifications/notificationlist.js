'use client'
import { useState, useEffect } from 'react'
import NotificationItem from './notificationitem'
import styles from "@/app/styles/notifications.module.css"



function NotificationList() {
    const [loading, setLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [notifications, setNotifications] = useState([]) // state hook
    const [check, setCheck] = useState(0)

    const fetchData = () => {
        fetch('/api/notifications/notifications/')
                   .then((res) => res.json())
                   .then((data) => {
                   console.log("Fetching Notifications ...")
                   setNotifications(data)
                   setLoading(false)
                  
                   }).catch((error) => {
                       console.log(error)
                   })
               }

    useEffect(() => {
         const id = setInterval(() => {
        setLoading(true)
    
        fetchData()
         }, 3000);
             return () => clearInterval(id);
        }, [check])

    useEffect(() => {
       
        setLoading(true)
        fetchData()
       
       
        }, [])



    return (
      <div >

            
        

        <div className={styles.NotificationList}>
        {
        notifications && notifications.map((item,index)=>(
              <NotificationItem
                key={item.id}
                id={item.id}
                title={item.title}
                emailTo={item.emailTo}
                emailBody={item.emailBody}
                emailSubject={item.emailSubject}
                approved={item.accepted}
              />

          ))
        }
        {
          notifications && notifications.length === 0 && (
            <span>There are no notifications.</span>
          )
        }
        </div>
      </div>
    );

}
export default NotificationList;