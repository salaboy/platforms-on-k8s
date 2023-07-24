'use client'
import { useState, useEffect } from 'react'
import EnvironmentItem from './environmentitem'
import styles from "@/app/styles/environments.module.css"



function EnvironmentList() {
    const [loading, setLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [environments, setEnvironments] = useState([]) // state hook
    const [check, setCheck] = useState(0)

    const handleDelete = (name) => {
      setLoading(true);
      setIsError(false);
      console.log("Deleting Environment ..." + name)
      fetch('/api/environments/' + name , {
        method: "DELETE",
        headers: {
          'accept': 'application/json',
        },
      }).then((response) => response.json()).then(() => {
        fetchData(filter)
        setLoading(false);
      }).catch(err => {
        console.log(err);
        setLoading(false);
        setIsError(true);
      });
  
    }

    const fetchData = () => {
        fetch('/api/environments/')
          .then((res) => res.json())
          .then((data) => {
          console.log("Fetching Environments ...")
          setEnvironments(data)
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

            
        

        <div className={styles.EnvironmentList}>
        {
        environments && environments.map((item,index)=>(
              <EnvironmentItem
                key={item.metadata.uid}
                id={item.metadata.uid}
                name={item.metadata.name}
                type={item.spec.compositionSelector.matchLabels.type}
                installInfra={item.spec.parameters.installInfra}
                debug={item.spec.parameters.frontend.debug}
                status={item.status.conditions[1].status}
                synced={item.status.conditions[0].status}
                vclusterRef={item.spec.resourceRef.name}
                secretRef={item.spec.writeConnectionSecretToRef.name}
                handleDelete={handleDelete}
              />

          ))
        }
        {
          environments && environments.length === 0 && (
            <span>There are no environments.</span>
          )
        }
        </div>
      </div>
    );

}
export default EnvironmentList;