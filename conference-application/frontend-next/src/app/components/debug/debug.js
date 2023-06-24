'use client'
import ServiceInfo from './serviceinfo'

import { useState, useEffect } from 'react'


function Debug() {
    const [isLoading, setLoading] = useState(false)
    const [frontendServiceInfo, setFrontendServiceInfo] = useState('') // state hook
    const [c4pServiceInfo, setC4pServiceInfo] = useState('') // state hook
    const [agendaServiceInfo, setAgendaServiceInfo] = useState('') // state hook
    const [notificationsServiceInfo, setNotificationsServiceInfo] = useState('') // state hook

    const mockServiceInfo = {
        "Name": "N/A",
        "PodId": "N/A",
        "PodNamespace": "N/A",
        "PodNodeName": "N/A",
        "Source": "N/A",
        "Version": "N/A",
        "Healthy": false
    }

    const fetchFrontendServiceInfo = () => {
        setLoading(true);
        console.log("Querying service/info")
        fetch('service/info')
            .then((res) => res.json())
            .then((data) => {
                data.Healthy = true;
                setFrontendServiceInfo(data)
                setLoading(false)
            }).catch((error) => {
                setFrontendServiceInfo(mockServiceInfo)
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
                setAgendaServiceInfo(mockServiceInfo)
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
                setC4pServiceInfo(mockServiceInfo)
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
                setNotificationsServiceInfo(mockServiceInfo)
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
            <ServiceInfo
                name={frontendServiceInfo.Name}
                version={frontendServiceInfo.Version}
                key="frontend"
                source={frontendServiceInfo.Source}
                podName={frontendServiceInfo.PodName}
                nodeName={frontendServiceInfo.NodeName}
                namespace={frontendServiceInfo.Namespace}
                healthy={frontendServiceInfo.Healthy}
            />

            <ServiceInfo
                name={c4pServiceInfo.Name}
                version={c4pServiceInfo.Version}
                key="c4p"
                source={c4pServiceInfo.Source}
                podName={c4pServiceInfo.PodName}
                nodeName={c4pServiceInfo.NodeName}
                namespace={c4pServiceInfo.Namespace}
                healthy={c4pServiceInfo.Healthy}
            />
            <ServiceInfo
                name={agendaServiceInfo.Name}
                version={agendaServiceInfo.Version}
                key="agenda"
                source={agendaServiceInfo.Source}
                podName={agendaServiceInfo.PodName}
                nodeName={agendaServiceInfo.NodeName}
                namespace={agendaServiceInfo.Namespace}
                healthy={agendaServiceInfo.Healthy}
            />
            <ServiceInfo
                name={notificationsServiceInfo.Name}
                version={notificationsServiceInfo.Version}
                key="notifications"
                source={notificationsServiceInfo.Source}
                podName={notificationsServiceInfo.PodName}
                nodeName={notificationsServiceInfo.NodeName}
                namespace={notificationsServiceInfo.Namespace}
                healthy={notificationsServiceInfo.Healthy}
            />
        </div>
    );

}

export default Debug;