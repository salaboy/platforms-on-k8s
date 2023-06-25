'use client'
import styles from '@/app/styles/debug.module.css'

function ServiceInfo({key, name, version, source, podName, nodeName, namespace, healthy}) {
   

   

    return (
      
      <div className={styles.ServiceInfo}>
        
        <div>
          <div className={styles.header}>
            <h4>
              {name}
            </h4>
            <h5>
              {version}
            </h5>
          </div>
          <div className={styles.description}>
            {source}
         
            {podName}
         
            {nodeName}
         
            {namespace}
          </div>
          <div  className={`${styles.statusTag}  ${healthy != null ? styles.healthy : ' '}`}>
            {healthy != null && <>Healthy</>}
          </div>
        </div>
       
         
      </div>
      
    );

}
export default ServiceInfo;