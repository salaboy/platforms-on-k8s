'use client'
import ServiceInfo from './serviceinfo'
import styles from '@/app/styles/debug.module.css'

import { useState, useEffect } from 'react'


function Debug() {
    const [isLoading, setLoading] = useState(false)
    const [frontendServiceInfo, setFrontendServiceInfo] = useState('') // state hook
    const [c4pServiceInfo, setC4pServiceInfo] = useState('') // state hook
    const [agendaServiceInfo, setAgendaServiceInfo] = useState('') // state hook
    const [notificationsServiceInfo, setNotificationsServiceInfo] = useState('') // state hook

    const [check, setCheck] = useState(0)

    const mockFrontendServiceInfo = {
        "Name": "FRONTEND",
        "PodId": "N/A",
        "PodNamespace": "N/A",
        "PodNodeName": "N/A",
        "PodName": "N/A",
        "PodServiceAccount": "N/A",
        "Source": "N/A",
        "Version": "N/A",
        "Healthy": false
    }

    const mockAgendaServiceInfo = {
        "Name": "AGENDA",
        "PodId": "N/A",
        "PodNamespace": "N/A",
        "PodNodeName": "N/A",
        "PodName": "N/A",
        "PodServiceAccount": "N/A",
        "Source": "N/A",
        "Version": "N/A",
        "Healthy": false
    }

    const mockC4PServiceInfo = {
        "Name": "C4P",
        "PodId": "N/A",
        "PodNamespace": "N/A",
        "PodNodeName": "N/A",
        "PodName": "N/A",
        "PodServiceAccount": "N/A",
        "Source": "N/A",
        "Version": "N/A",
        "Healthy": false
    }

    const mockNotificationsServiceInfo = {
        "Name": "NOTIFICATIONS",
        "PodId": "N/A",
        "PodNamespace": "N/A",
        "PodNodeName": "N/A",
        "PodName": "N/A",
        "ServiceAccount": "N/A",
        "Source": "N/A",
        "Version": "N/A",
        "Healthy": false
    }

    useEffect(() => {
        const id = setInterval(() => {
            setLoading(true)
            fetchFrontendServiceInfo()
            fetchAgendaServiceInfo()
            fetchC4PServiceInfo()
            fetchNotificationsServiceInfo()

        }, 3000);
        return () => clearInterval(id);
    }, [check])

    const fetchFrontendServiceInfo = () => {
        setLoading(true);
        console.log("Querying service/info")
        fetch('/service/info')
            .then((res) => res.json())
            .then((data) => {
                data.Healthy = true;
                setFrontendServiceInfo(data)
                setLoading(false)
            }).catch((error) => {
                setFrontendServiceInfo(mockFrontendServiceInfo)
                console.log(error)
            })
    };

    const fetchAgendaServiceInfo = () => {
        setLoading(true);
        console.log("Querying /api/agenda/service/info")
        fetch('/api/agenda/service/info')
            .then((res) => res.json())
            .then((data) => {
                data.Healthy = true;
                setAgendaServiceInfo(data)
                setLoading(false)
            }).catch((error) => {
                setAgendaServiceInfo(mockAgendaServiceInfo)
                console.log(error)
            })
    };

    const fetchC4PServiceInfo = () => {
        setLoading(true);
        console.log("Querying /api/c4p/service/info")
        fetch('/api/c4p/service/info')
            .then((res) => res.json())
            .then((data) => {
                data.Healthy = true;
                setC4pServiceInfo(data)
                setLoading(false)
            }).catch((error) => {
                setC4pServiceInfo(mockC4PServiceInfo)
                console.log(error)
            })
    };
    const fetchNotificationsServiceInfo = () => {
        setLoading(true);
        console.log("Querying /api/notifications/service/info")
        fetch('/api/notifications/service/info')
            .then((res) => res.json())
            .then((data) => {
                data.Healthy = true;
                setNotificationsServiceInfo(data)
                setLoading(false)
            }).catch((error) => {
                setNotificationsServiceInfo(mockNotificationsServiceInfo)
                console.log(error)
            })
    };


    useEffect(() => {                           // side effect hook
        setLoading(true)
        fetchFrontendServiceInfo()
        fetchAgendaServiceInfo()
        fetchC4PServiceInfo()
        fetchNotificationsServiceInfo()
    }, [])

    return (
        <div className={styles.DebugList}>
            <ServiceInfo
                name={frontendServiceInfo.Name}
                version={frontendServiceInfo.Version}
                key="frontend"
                source={frontendServiceInfo.Source}
                podIp={frontendServiceInfo.PodIp}
                podName={frontendServiceInfo.PodName}
                nodeName={frontendServiceInfo.PodNodeName}
                namespace={frontendServiceInfo.PodNamespace}
                serviceAccount={frontendServiceInfo.PodServiceAccount}
                healthy={frontendServiceInfo.Healthy}
            />

            <ServiceInfo
                name={c4pServiceInfo.Name}
                version={c4pServiceInfo.Version}
                key="c4p"
                source={c4pServiceInfo.Source}
                podName={c4pServiceInfo.PodName}
                podIp={c4pServiceInfo.PodIp}
                nodeName={c4pServiceInfo.PodNodeName}
                namespace={c4pServiceInfo.PodNamespace}
                serviceAccount={c4pServiceInfo.PodServiceAccount}
                healthy={c4pServiceInfo.Healthy}
            />
            <ServiceInfo
                name={agendaServiceInfo.Name}
                version={agendaServiceInfo.Version}
                key="agenda"
                source={agendaServiceInfo.Source}
                podIp={agendaServiceInfo.PodIp}
                podName={agendaServiceInfo.PodName}
                nodeName={agendaServiceInfo.PodNodeName}
                namespace={agendaServiceInfo.PodNamespace}
                serviceAccount={agendaServiceInfo.PodServiceAccount}
                healthy={agendaServiceInfo.Healthy}
            />
            <ServiceInfo
                name={notificationsServiceInfo.Name}
                version={notificationsServiceInfo.Version}
                key="notifications"
                source={notificationsServiceInfo.Source}
                podIp={notificationsServiceInfo.PodIp}
                podName={notificationsServiceInfo.PodName}
                nodeName={notificationsServiceInfo.PodNodeName}
                namespace={notificationsServiceInfo.PodNamespace}
                serviceAccount={notificationsServiceInfo.PodServiceAccount}
                healthy={notificationsServiceInfo.Healthy}
            />
        </div>
    );

}

export default Debug;