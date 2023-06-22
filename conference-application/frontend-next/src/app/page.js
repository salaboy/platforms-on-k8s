'use client'
import Image from 'next/image'
import styles from './styles/home.module.css'
import Button from './components/forms/button/button'


export default function Home() {
  
  return (
    <main className={styles.main}>
        <section className={`${styles.hero}  ${styles.section} ` }>
          <div className="grid content">
              <div className='col twoThirds'>
                  <h1>Cloud-Native Conf 2023</h1>
                  <h3>30th February.</h3>
              </div>
              <div className='col third '>
                  <p>The flagship conference gathers adopters and technologists from leading Open Source and Cloud-Native communities.</p>
                  
                  <Button link="/agenda/">Explore the Schedule</Button>
              </div>
          </div>
        </section>
        <section className={`${styles.venue}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col half'>
               
              </div>
              <div className='col half'>
                <h5>The Venue</h5>
                <h2>The Cloud</h2>
               
              </div>
            
          </div>
        </section>
       

        <section className={`${styles.proposals}  ${styles.section} ` }>
          <div className="grid content" >
              
              <div className='col third '>
                <h3>The Call for Proposals is Open</h3>
                <p>Lorem ipsum, dolor sit amet consectetur adipisicing elit. Autem, adipisci dolore eum, molestiae aut quidem asperiores culpa quae error harum suscipit repellendus excepturi delectus labore officiis! In saepe reiciendis rerum!</p>
                <Button link="/proposals/">Submit your proposal</Button> 
              </div>
             
            
          </div>
        </section>
       
    </main>
  )
}
