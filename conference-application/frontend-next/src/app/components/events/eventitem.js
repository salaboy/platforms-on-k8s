
'use client'


function EventItem({id, type, payload}) {

    return (
      
      <div>
        <div className="ProposalItem__header">
          <h3>{id}</h3>
          <h5>{type}</h5>
          {/* Maybe render using: https://www.npmjs.com/package/react-json-pretty */}
          <div className="ProposalItem__status">
            {payload}
          </div>
        </div>
        
        
      </div>
      
    );

}
export default EventItem;