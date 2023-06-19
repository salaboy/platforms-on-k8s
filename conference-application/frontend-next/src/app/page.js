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
        <section className={`${styles.hero}  ${styles.section} ` }>
          <div className="grid content">
              <div className='col half'>
                  <h1>Cloud Con 2023</h1>
              </div>
              <div className='col third '>
                  <p>The flagship conference gathers adopters and technologists from leading open source and cloud native communities.</p>
                  <h2>18 to 21, April.</h2>
              </div>
          </div>
        </section>
        <section className={`${styles.venue}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col half'>
               
              </div>
              <div className='col half'>
                <h4>The Venue</h4>
                <h2>RAI Amsterdam Convention Centre.</h2>
                <p>Europaplein 24, 1078 GZ. Amsterdam, The Netherlands.</p>
              </div>
            
          </div>
        </section>
        <section className={`${styles.agenda}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col twoThirds'>
                <ul>
                  {data.map((item) => (
                    <li key={item.Id}>{item.Id} - {item.Title} - {item.Author} - {item.Day} - {item.Time} </li>
                  ))}
                </ul>
              </div>
              <div className='col third'>
                <h2>Agenda</h2>
                
              </div>
            
          </div>
        </section>

        <section className={`${styles.speakers}  ${styles.section} ` }>
          <div className="grid content" >
              <div className='col third'>
                <h2>Speakers</h2>
                
              </div>
              <div className='col half'>
                Highlighted Speakers list
              </div>
             
            
          </div>
        </section>

        <section className={`${styles.proposals}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col twoThirds'>
                Submit your proposal
              </div>
             
            
          </div>
        </section>

        <section className={`${styles.sponsors}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col full'>
                List of sponsors
              </div>
             
            
          </div>
        </section>

       
    </main>
  )
}
