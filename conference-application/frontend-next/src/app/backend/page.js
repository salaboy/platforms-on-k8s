'use client'
import styles from '@/app/styles/backend.module.css'
import { useState, useEffect } from 'react'


export default function Backend() {

  const [data, setData] = useState(null)
  const [isLoading, setLoading] = useState(false)
  const [isError, setIsError] = useState(false);
  const [decisionsMade, setDecisionsMade] = useState(1)

  useEffect(() => {
    setLoading(true)
    fetch('/api/c4p')
      .then((res) => res.json())
      .then((data) => {
        setData(data)
        setLoading(false)
      })
  }, [decisionsMade])
 
  if (isLoading) return <p>Loading...</p>
  if (!data) return <p>No Proposals</p>

  

  function decide(id, approved){
    const decision = {
      approved: approved,
     
    }
    try{
      setLoading(true);
      fetch('/api/c4p/'+id+"/decide", {
        method: "POST",
        body: JSON.stringify(decision),
        headers: {
          'accept': 'application/json',
        },
      }).then((response) => response.json()).then((data) => {
          setDecisionsMade(decisionsMade+1)
          setLoading(false);
        }
      )
    }catch(err){
        setLoading(false);
        setIsError(true);
    }
  
  }

  


  return (
    <main className={styles.main}>
      <h1>Backend</h1>
      <h2>Review Proposals (Tab)</h2>
      <div>
            <ul>
              {data.map((p) => (
                <li key={p.Id}>{p.Id} - {p.Title} - {p.Description} - {p.Author} - {p.Email}  - {p.Status.Status}  - {p.Approved} 
                  <button main onClick={() => decide(p.Id, true)} >Approve</button>
                  <button main onClick={() => decide(p.Id, false)}>Reject</button>
                </li>
                
              ))}
            </ul>
        </div>

      <h2>Notifications (Tab)</h2> 
      (TBD)

      <h2>Events (Tab)</h2> 
      (TBD)
    </main>
  )
}
