
'use client'
import Button from '../forms/button/button'

function ProposalItem({title, author, id, status, approved, email, description, actionHandler}) {

    const handleAction = (id, status, action) => {

      actionHandler(status, id,action);
    }

    return (
      
      <div>
        <div className="ProposalItem__header">
          <h3>{title}</h3>
          <h5>{author} {email}</h5>
          <div className="ProposalItem__status">
            {status}
          </div>
        </div>
        <div className="ProposalItem__description">
          <p className="p p-s">{description}</p>
        </div>
        {false && (
        <div className="ProposalItem__id">
          {id}
        </div>
        )}
        
        {status && status==="PENDING" && (
          <div className="ProposalItem__actions">
            <div >
              <Button clickHandler={() => handleAction(id, status,"APPROVE")}>Approve</Button>
            </div>
            <div>
              <Button clickHandler={() => handleAction(id, status,"REJECT")}>Reject</Button>
            </div>
            <div>
              <Button clickHandler={() => handleAction(id, status,"ARCHIVE")}>Archive</Button>
            </div>
          </div>
        )}
        
        {status && status==="DECIDED" && (
          <div className="ProposalItem__status-info">

              {approved === true  && (
                <div className="ProposalItem__badge --approved">Approved</div>
              )}
              {approved === false  && (
                <div className="ProposalItem__badge --rejected">Rejected</div>
              )}
              <div>
              <Button clickHandler={() => handleAction(id, status,"ARCHIVE")}>Archive</Button>
            </div>
          </div>
        )}
        {status && status==="ARCHIVED" && (
          <div className="ProposalItem__status-info">
                <div className="ProposalItem__badge --approved">Archived</div>
          </div>
        )}
      </div>
      
    );

}
export default ProposalItem;