'use client'
import styles from '@/app/styles/debug.module.css'

function ServiceInfo({ key, name, version, source, podName, nodeName, namespace, podIp, serviceAccount, healthy }) {




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
          <h7>Source: {source}</h7><br/>

          <h7>PodName: {podName}</h7><br/>

          <h7>NodeName: {nodeName}</h7><br/>

          <h7>PodNamespace: {namespace}</h7><br/>

          <h7>PodIP:{podIp}</h7><br/>

          <h7>PodServiceAccount: {serviceAccount}</h7><br/>
        </div>
        <div className={`${styles.statusTag}  ${healthy != null ? styles.healthy : ' '}`}>
          {healthy != null && <>Healthy</>}
        </div>
      </div>


    </div>

  );

}
export default ServiceInfo;