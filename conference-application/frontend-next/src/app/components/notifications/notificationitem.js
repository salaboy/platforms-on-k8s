
'use client'


function NotificationItem({id, title, emailTo,  emailSubject, emailBody, approved}) {

    return (
      
      <div>
        <div className="ProposalItem__header">
          <h3>Proposal: {title}</h3>
          
          <div className="ProposalItem__status">
            Approved? {approved}
          </div>
        </div>
        <div className="ProposalItem__description">
          To: <p className="p p-s">{emailTo}</p>
          Subject: <p className="p p-s">{emailSubject}</p>
          Body: <p className="p p-s">{emailBody}</p>
        </div>
        
      </div>
      
    );

}
export default NotificationItem;