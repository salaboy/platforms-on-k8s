'use client'


function ServiceInfo({key, name, version, source, podName, nodeName, namespace, healthy}) {
   

   

    return (
      
      <div>
        
        <div className="AgendaItem__date">
          <div className="AgendaItem__day">
            {name}
          </div>
          <div className="AgendaItem__time">
            {version}
          </div>
          <div className="AgendaItem__time">
            {source}
          </div>
          <div className="AgendaItem__time">
            {podName}
          </div>
          <div className="AgendaItem__time">
            {nodeName}
          </div>
          <div className="AgendaItem__time">
            {namespace}
          </div>
          <div className="AgendaItem__time">
            {healthy != null && (healthy.toString())}
          </div>
        </div>
       
         
      </div>
      
    );

}
export default ServiceInfo;