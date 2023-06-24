
'use client'
import Button from '../forms/button/button'
import styles from '@/app/styles/proposals.module.css'

function ProposalItem({title, author, id, status, approved, email, description, actionHandler}) {

    const handleAction = (id, status, action) => {

      actionHandler(status, id,action);
    }

    return (
      
      <div className={`${styles.ProposalItem}  ${status==='PENDING' ? styles.pending : ''}   ${status==='REJECT' ? styles.rejected : ''}   ${status==='APPROVE' ? styles.approved : ''}  ${status==='ARCHIVE' ? styles.archived : ''}` }>
        <div className="ProposalItem__header">
          <h4>{title}</h4>
          <div>{author} {email}</div>
          <div className="ProposalItem__status">
            {status}
          </div>
        </div>
        <div className="ProposalItem__description">
          <p className="p --s">{description}</p>
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