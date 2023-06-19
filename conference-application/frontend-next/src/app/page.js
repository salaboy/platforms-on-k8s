'use client'
import Image from 'next/image'
import styles from './styles/home.module.css'
import { useState, useEffect } from 'react'



export default function Home() {
  
  const [data, setData] = useState(null)
  const [isLoading, setLoading] = useState(false)
 
  useEffect(() => {
    setLoading(true)
    fetch('/api/agenda/')
      .then((res) => res.json())
      .then((data) => {
        setData(data)
        setLoading(false)
      })
  }, [])
 
  if (isLoading) return <p>Loading...</p>
  if (!data) return <p>No Proposals</p>

  return (
    <main className={styles.main}>
      
        <div className="grid">
            <div className='col half'>
              <ul>
                {data.map((item) => (
                  <li key={item.Id}>{item.Id} - {item.Title} - {item.Author} - {item.Day} - {item.Time} </li>
                ))}
              </ul>
            </div>
          
        </div>

       
    </main>
  )
}
