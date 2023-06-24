
'use client'
import Button from '../forms/button/button'
import styles from '@/app/styles/proposals.module.css'

function ProposalItem({title, author, id, status, approved, email, description, actionHandler}) {

    const handleAction = (id, status, action) => {

      actionHandler(status, id,action);
    }

    return (
      
      <div className={`${styles.ProposalItem}  ${status==='PENDING' ? styles.pending : ''}   ${status==='DECIDED' ? styles.decided : ''}  ${status==='REJECT' ? styles.rejected : ''}   ${status==='APPROVE' ? styles.approved : ''}  ${status==='ARCHIVED' ? styles.archived : ''}` }>
        <div className="ProposalItem__header">
          <h4>{title}</h4>
          <div>{author} {email}</div>
          <div className={styles.status}>
            {status}
          </div>
        </div>
        <div className={styles.description}>
          <p className="p --s">{description}</p>
        </div>
        {false && (
        <div className="ProposalItem__id">
          {id}
        </div>
        )}
        
        {status && status==="PENDING" && (
          <div className={styles.actions}>
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
          <div className={styles.actions}>

              {approved === true  && (
                <div className={`${styles.statusTag}  ${styles.approved}`} >Approved</div>
              )}
              {approved === false  && (
                <div className={`${styles.statusTag}  ${styles.rejected}`} >Rejected</div>
              )}
              <div>
              <Button clickHandler={() => handleAction(id, status,"ARCHIVE")}>Archive</Button>
            </div>
          </div>
        )}
        {status && status==="ARCHIVED" && (
          <div className={styles.actions}>
                <div className="ProposalItem__badge --approved">Archived</div>
          </div>
        )}
      </div>
      
    );

}
export default ProposalItem;