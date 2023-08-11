'use client'
import { useState, useEffect } from 'react'
import ProposalItem from './proposalitem'
import Button from '../forms/button/button'
import styles from '@/app/styles/proposals.module.css'

function ProposalList() {
  const [loading, setLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [statusFilter, setStatusFilter] = useState("PENDING")
  const [proposalItems, setProposalItems] = useState([]) // state hook

  const handleApproval = (id, approved) => {
    setLoading(true);
    setIsError(false);
    const data = {
      approved: approved,
    }
    console.log("Decision Made ...")
    fetch('/api/c4p/proposals/' + id + "/decide/", {
      method: "POST",
      body: JSON.stringify(data),
      headers: {
        'accept': 'application/json',
      },
    }).then((response) => response.json()).then((data) => {
      var filter = "?status=" + statusFilter
      if (statusFilter == "") {
        filter = ""
      }
      fetchData(filter)
      setLoading(false);
    }).catch(err => {
      setLoading(false);
      setIsError(true);
    });

  }

  const handleArchive = (id) => {
    setLoading(true);
    setIsError(false);
    console.log("Archiving Proposal ..." + id)
    fetch('/api/c4p/proposals/' + id , {
      method: "DELETE",
      headers: {
        'accept': 'application/json',
      },
    }).then((response) => response.json()).then(() => {
      var filter = "?status=" + statusFilter
      if (statusFilter == "") {
        filter = ""
      }
      fetchData(filter)
      setLoading(false);
    }).catch(err => {
      console.log(err);
      setLoading(false);
      setIsError(true);
    });

  }


  function ItemAction(status, id, action) {
    console.log("status: " + status + " - id: " + id + " - action: " + action)
    if (status == "PENDING") {
      if (action == "APPROVE") {
        handleApproval(id, true)
      } else {
        handleApproval(id, false)
      } 
    }
    if (action == "ARCHIVE"){
      handleArchive(id)
    }
  }

  function PendingFilter() {
    setStatusFilter("PENDING")
  }

  function AllFilter() {
    setStatusFilter("")
  }

  function DecidedFilter() {
    setStatusFilter("DECIDED")
  }
  function ArchivedFilter() {
    setStatusFilter("ARCHIVED")
  }

  const fetchData = (filter) => {
    console.log("Fetching Proposals ... (" + filter + ").")
    fetch('/api/c4p/proposals/' + filter)
      .then((res) => res.json())
      .then((data) => {
        setProposalItems(data)
        setLoading(false)
      }).catch((error) => {
        console.log(error)
      })
  }
  useEffect(() => {
    setLoading(true)
    var filter = "?status=" + statusFilter
    if (statusFilter == "") {
      filter = ""
    }
    fetchData(filter)

  }, [ statusFilter])

  return (
    <div className={styles.ProposalList}>

      <div className={styles.ProposalList_Filters}>
        <div className={styles.container}>
          <div className={styles.filterLabel}>Filter By: </div>
          
          <div className={`${statusFilter == "" ? styles.inactive : styles.active}  ${ styles.filter }` }  onClick={() => AllFilter()}>All</div>
          <div className={`${statusFilter == "PENDING" ? styles.inactive : styles.active}   ${ styles.filter }` }  onClick={() => PendingFilter()}>Pending</div>
          <div className={`${statusFilter == "DECIDED" ? styles.inactive : styles.active}   ${ styles.filter }` }  onClick={() => DecidedFilter()}>Decided</div>
          <div className={`${statusFilter == "ARCHIVED" ? styles.inactive : styles.active}   ${ styles.filter }` }  onClick={() => ArchivedFilter()}>Archived</div>
        </div>
      </div>

      <div className={styles.ProposalList_Items}>
        {
          proposalItems && proposalItems.map((item, index) => (
            <ProposalItem
              key={item.id}
              id={item.id}
              title={item.title}
              author={item.author}
              description={item.description}
              email={item.email}
              approved={item.approved}
              status={item.status.status}
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