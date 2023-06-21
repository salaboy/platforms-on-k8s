'use client'
import { useState, useEffect } from 'react'
import NotificationItem from './notificationitem'


function NotificationList() {
    const [loading, setLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [notifications, setNotifications] = useState([]) // state hook
    const [check, setCheck] = useState(0)

    useEffect(() => {
        const id = setInterval(() => {
        setLoading(true)
    
        fetch('/api/notifications/')
            .then((res) => res.json())
            .then((data) => {
            console.log("Fetching Notifications ...")
            setNotifications(data)
            setLoading(false)
            }).catch((error) => {
                console.log(error)
            })
        }, 3000);
            return () => clearInterval(id);
        }, [check])



    return (
      <div>

            
        

        <div>
        {
        notifications && notifications.map((item,index)=>(
              <NotificationItem
                key={item.Id}
                id={item.Id}
                title={item.Title}
                author={item.Author}
                emailText={item.EmailText}
                email={item.Email}
                approved={item.Approved}
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