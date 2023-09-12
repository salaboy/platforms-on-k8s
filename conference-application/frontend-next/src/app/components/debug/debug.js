'use client'
import ServiceInfo from './serviceinfo'
import styles from '@/app/styles/debug.module.css'

import { useState, useEffect } from 'react'
import JSONPretty from 'react-json-pretty';

function Debug() {
    const [isLoading, setLoading] = useState(false)
    const [frontendServiceInfo, setFrontendServiceInfo] = useState('') // state hook
    const [c4pServiceInfo, setC4pServiceInfo] = useState('') // state hook
    const [agendaServiceInfo, setAgendaServiceInfo] = useState('') // state hook
    const [notificationsServiceInfo, setNotificationsServiceInfo] = useState('') // state hook

    const [check, setCheck] = useState(0)

    
    const [features, setFeatures] = useState('') // state hook

    const fetchFeatures = () => {
        setLoading(true);
        console.log("Querying /api/features/")
        fetch('/api/features/')
        .then((res) => res.json())
        .then((data) => {
            console.log("Features Data: " + data)
            setFeatures(data)
            setLoading(false)
        }).catch((error) => {
            setFeatures({})
            console.log(error)
        })
    };


    const mockFrontendServiceInfo = {
        "name": "FRONTEND",
        "podId": "N/A",
        "podNamespace": "N/A",
        "podNodeName": "N/A",
        "podName": "N/A",
        "podServiceAccount": "N/A",
        "source": "N/A",
        "version": "N/A",
        "healthy": false,
        "eventsEnabled": false
    }

    const mockAgendaServiceInfo = {
        "name": "AGENDA",
        "podId": "N/A",
        "podNamespace": "N/A",
        "podNodeName": "N/A",
        "podName": "N/A",
        "podServiceAccount": "N/A",
        "source": "N/A",
        "version": "N/A",
        "healthy": false,
        "eventsEnabled": false

    }

    const mockC4PServiceInfo = {
        "name": "C4P",
        "podId": "N/A",
        "podNamespace": "N/A",
        "podNodeName": "N/A",
        "podName": "N/A",
        "podServiceAccount": "N/A",
        "source": "N/A",
        "version": "N/A",
        "healthy": false,
        "eventsEnabled": false
    }

    const mockNotificationsServiceInfo = {
        "name": "NOTIFICATIONS",
        "podId": "N/A",
        "podNamespace": "N/A",
        "podNodeName": "N/A",
        "podName": "N/A",
        "serviceAccount": "N/A",
        "source": "N/A",
        "version": "N/A",
        "healthy": false,
        "eventsEnabled": false
    }

    useEffect(() => {
        const id = setInterval(() => {
            setLoading(true)
            fetchFrontendServiceInfo()
            fetchAgendaServiceInfo()
            fetchC4PServiceInfo()
            fetchNotificationsServiceInfo()
            fetchFeatures()
        }, 3000);
        return () => clearInterval(id);
    }, [check])

    const fetchFrontendServiceInfo = () => {
        setLoading(true);
        console.log("Querying /api/service/info")
        fetch('/api/service/info')
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
        <div>
            <div>
                <h3>Feature Flags</h3>
                <JSONPretty id="json-pretty" data={JSON.stringify(features)}></JSONPretty>
            </div>
            <div>
            <h3>Services</h3>
            
                <div className={styles.DebugList}>
                    
                    <ServiceInfo
                        name={frontendServiceInfo.name}
                        version={frontendServiceInfo.version}
                        key="frontend"
                        source={frontendServiceInfo.source}
                        podIp={frontendServiceInfo.podIp}
                        podName={frontendServiceInfo.podName}
                        nodeName={frontendServiceInfo.podNodeName}
                        namespace={frontendServiceInfo.podNamespace}
                        serviceAccount={frontendServiceInfo.podServiceAccount}
                        healthy={frontendServiceInfo.healthy}
                        eventsEnabled={frontendServiceInfo.eventsEnabled}
                    />

                    <ServiceInfo
                        name={c4pServiceInfo.name}
                        version={c4pServiceInfo.version}
                        key="c4p"
                        source={c4pServiceInfo.source}
                        podName={c4pServiceInfo.podName}
                        podIp={c4pServiceInfo.podIp}
                        nodeName={c4pServiceInfo.podNodeName}
                        namespace={c4pServiceInfo.podNamespace}
                        serviceAccount={c4pServiceInfo.podServiceAccount}
                        healthy={c4pServiceInfo.healthy}
                        eventsEnabled={c4pServiceInfo.eventsEnabled}
                    />
                    <ServiceInfo
                        name={agendaServiceInfo.name}
                        version={agendaServiceInfo.version}
                        key="agenda"
                        source={agendaServiceInfo.source}
                        podIp={agendaServiceInfo.podIp}
                        podName={agendaServiceInfo.podName}
                        nodeName={agendaServiceInfo.podNodeName}
                        namespace={agendaServiceInfo.podNamespace}
                        serviceAccount={agendaServiceInfo.podServiceAccount}
                        healthy={agendaServiceInfo.healthy}
                        eventsEnabled={agendaServiceInfo.eventsEnabled}
                    />
                    <ServiceInfo
                        name={notificationsServiceInfo.name}
                        version={notificationsServiceInfo.version}
                        key="notifications"
                        source={notificationsServiceInfo.source}
                        podIp={notificationsServiceInfo.podIp}
                        podName={notificationsServiceInfo.podName}
                        nodeName={notificationsServiceInfo.podNodeName}
                        namespace={notificationsServiceInfo.podNamespace}
                        serviceAccount={notificationsServiceInfo.podServiceAccount}
                        healthy={notificationsServiceInfo.healthy}
                        eventsEnabled={notificationsServiceInfo.eventsEnabled}
                    />
                </div>
            </div>
        </div>
    );

}

export default Debug;