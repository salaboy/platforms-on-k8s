'use client'
import styles from '@/app/styles/debug.module.css'

function ServiceInfo({ key, name, version, source, podName, nodeName, namespace, podIp, serviceAccount, healthy }) {




  return (

    <div className={styles.ServiceInfo}>

      <div>
        <div className={styles.header}>
          <h4>
            <a href={source} target='_blank'>
            {name}
            </a>
          </h4>
          <h5>
            {version}
          </h5>
        </div>
        <div className={styles.description}>
          <div className={styles.descriptionItem}>
            <span>
            Pod Name: 
            </span>
            {podName}
          </div>

          <div className={styles.descriptionItem}>
            <span>
            Node Name: 
            </span>
            {nodeName}
          </div>

          <div className={styles.descriptionItem}>
            <span>
            Pod Namespace: 
            </span>
            {namespace}
          </div>

          <div className={styles.descriptionItem}>
            <span>
            Pod IP: 
            </span>
            {podIp}
          </div>

          <div className={styles.descriptionItem}>
            <span>
            Pod Service Account: 
            </span>
            {serviceAccount}
          </div>

        </div>
        <div className={`${styles.statusTag}  ${(healthy != null && healthy) ? styles.healthy : styles.unhealthy }`}>
          {healthy != null && healthy && <>Healthy</>}
          {healthy != null && !healthy && <>Unhealthy</>}
        </div>
      </div>


    </div>

  );

}
export default ServiceInfo;