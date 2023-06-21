'use client'
import Image from 'next/image'
import styles from './styles/home.module.css'



export default function Home() {
  
  return (
    <main className={styles.main}>
        <section className={`${styles.hero}  ${styles.section} ` }>
          <div className="grid content">
              <div className='col twoThirds'>
                  <h1>Cloud-Native Conf 2023</h1>
              </div>
              <div className='col third '>
                  <p>The flagship conference gathers adopters and technologists from leading Open Source and Cloud-Native communities.</p>
                  <h3>30th February.</h3>
              </div>
          </div>
        </section>
        <section className={`${styles.venue}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col half'>
               
              </div>
              <div className='col half'>
                <h4>The Venue</h4>
                <h2>The Cloud</h2>
               
              </div>
            
          </div>
        </section>
       

        <section className={`${styles.proposals}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col twoThirds'>
                The Call for Proposals is Open. Submit your proposal HERE.
              </div>
             
            
          </div>
        </section>
       
    </main>
  )
}
