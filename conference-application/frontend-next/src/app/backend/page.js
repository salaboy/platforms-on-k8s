import styles from '@/app/styles/backend.module.css'
import { Suspense } from 'react'

export async function getProposals() {
  
  const res = await fetch(process.env.REMOTE_URL+"/c4p/");
  if (!res.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error('Failed to fetch data')
  }
  return res.json();
}

export default async function Backend() {
  
  const proposals = await getProposals();

  return (
    <main className={styles.main}>
      <h1>Backend</h1>
      <h2>Review Proposals (Tab)</h2>
      <div>
          <Suspense fallback={<div>Loading...</div>}>
            <ul>
              {proposals.map((p) => (
                <li key={p.Id}>{p.Id} - {p.Title} - {p.Description} - {p.Author} - {p.Email}  - {p.Approved} </li>
              ))}
            </ul>
          </Suspense>
        </div>

      <h2>Notifications (Tab)</h2> 
      (TBD)

      <h2>Events (Tab)</h2> 
      (TBD)
    </main>
  )
}
