import Image from 'next/image'
import styles from './styles/home.module.css'
import { Suspense } from 'react'

export async function getAgendaItems() {
  
  const res = await fetch(process.env.REMOTE_URL+"/agenda/");
  if (!res.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error('Failed to fetch data')
  }
  return res.json();
}

export default async function Home() {
  const items = await getAgendaItems();

  return (
    <main className={styles.main}>
      
        <div>
          <Suspense fallback={<div>Loading...</div>}>
            <ul>
              {items.map((item) => (
                <li key={item.Id}>{item.Id} - {item.Title} - {item.Author} - {item.Day} - {item.Time} </li>
              ))}
            </ul>
          </Suspense>
        </div>

       
    </main>
  )
}
