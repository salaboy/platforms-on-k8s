
'use client'


function EventItem({id, type, payload}) {

    return (
      
      <div>
        <div className="ProposalItem__header">
          <h3>{title}</h3>
          <h5>{email}</h5>
          <div className="ProposalItem__status">
            {approved}
          </div>
        </div>
        <div className="ProposalItem__description">
          <p className="p p-s">{emailText}</p>
        </div>
        
      </div>
      
    );

}
export default EventItem;