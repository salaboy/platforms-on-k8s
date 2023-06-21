'use client'
import { useState, useEffect } from 'react'
import ProposalItem from './proposalitem'
import Button from '../forms/button/button'

function ProposalList() {
    const [loading, setLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [decisionsMade, setDecisionsMade] = useState(1)
    const [statusFilter, setStatusFilter] = useState(false)
    const [proposalItems, setProposalItems] = useState([]) // state hook

    const handleApproval = (id, approved) => {
      setLoading(true);
      setIsError(false);
      const data = {
        approved: approved,
      }
      console.log("Decision Made ...")
      fetch('/api/c4p/' + id + "/decide", {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          'accept': 'application/json',
        },
      }).then((response) => response.json()).then((data) => {
        setDecisionsMade(decisionsMade+1)
        setLoading(false);
      }).catch(err => {
          setLoading(false);
          setIsError(true);
        });

    }

    function ApprovalButtons(item){
      if(item.status == "PENDING"){
      return <span><a onClick={() => handleApproval(item.id, true)} disabled={loading}  id="approve" >{loading ? 'Loading...' : 'Approve'}</a> /
      <a onClick={() => handleApproval(item.id,false)} disabled={loading}  id="reject" >{loading ? 'Loading...' : 'Reject'}</a></span>
      }else{
        return <span>No Actions</span>
      }
    }

    function ItemAction(status, id, action){

      if(status == "PENDING"){
        if(action == "APPROVE"){
          handleApproval(id, true)
        }else {
          handleApproval(id, false)
        }
      }
    }

    function PendingFilter(){
      setStatusFilter("PENDING")
    }

    function AllFilter(){
      setStatusFilter("")
    }

    function DecidedFilter(){
      setStatusFilter("DECIDED")
    }

    useEffect(() => {
      setLoading(true)
      var filter = "?status="+statusFilter
      if(statusFilter == ""){
        filter =""
      }
      fetch('/api/c4p/'+filter)
        .then((res) => res.json())
        .then((data) => {
          console.log("Fetching Proposals ...")
          setProposalItems(data)
          setLoading(false)
        }).catch((error) => {
            console.log(error)
          })

    }, [decisionsMade, statusFilter])

    return (
      <div>

            <div className="ProposalList__filters">
              <Button inverted state={statusFilter == "" ? "inactive" : "active"} clickHandler={() => AllFilter()}>All</Button> 
              <Button inverted state={statusFilter == "PENDING" ? "active" : "inactive" } clickHandler={() => PendingFilter()}>Pending</Button>
              <Button inverted state={statusFilter == "DECIDED" ? "active" : "inactive" } clickHandler={() => DecidedFilter()}>Decided</Button>
            </div>
        

        <div>
        {
        proposalItems && proposalItems.map((item,index)=>(
              <ProposalItem
                key={item.Id}
                id={item.Id}
                title={item.Title}
                author={item.Author}
                description={item.Description}
                email={item.Email}
                approved={item.Approved}
                status={item.Status.Status}
                actionHandler={ItemAction}
              />

          ))
        }
        {
          proposalItems && proposalItems.length === 0 && statusFilter === "PENDING" && (
            <span>There are no pending proposals.</span>
          )
        }
        {
          proposalItems && proposalItems.length === 0 && statusFilter === "DECIDED" && (
            <span>There are no decided proposals.</span>
          )
        }
        {
          proposalItems && proposalItems.length === 0 && statusFilter === false && (
            <span>There are no proposals.</span>
          )
        }
        </div>
      </div>
    );

}
export default ProposalList;