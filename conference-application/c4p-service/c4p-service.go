package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type Proposal struct {
	Id          string
	Title       string
	Description string
	Author      string
	Email       string
	Approved    bool
	Status      ProposalStatus
}

type ProposalStatus struct {
	Status string
}

type ProposalDecision struct {
	Approved bool
}

type ServiceInfo struct {
	Name         string
	Version      string
	Source       string
	PodId        string
	PodNamespace string
	PodNodeName  string
}

var VERSION = getEnv("VERSION", "1.0.0")
var SOURCE = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/c4p-service")
var POD_ID = getEnv("POD_ID", "N/A")
var POD_NAMESPACE = getEnv("POD_NAMESPACE", "N/A")
var POD_NODENAME = getEnv("POD_NODENAME", "N/A")
var POSTGRESQL_HOST = getEnv("POSTGRES_HOST", "localhost")
var POSTGRESQL_PORT = getEnv("POSTGRES_PORT", "5432")
var POSTGRESQL_USERNAME = getEnv("POSTGRES_USERNAME", "postgres")
var POSTGRESQL_PASSOWRD = getEnv("POSTGRES_PASSWORD", "")

var db *sql.DB

func getAllProposalsHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`SELECT * FROM Proposals`)

	if err != nil {
		log.Printf("There was an error executing the query %v", err)
	}

	defer rows.Close()
	var proposals []Proposal
	for rows.Next() {

		var proposal Proposal
		err = rows.Scan(&proposal.Id, &proposal.Title, &proposal.Description, &proposal.Email, &proposal.Author, &proposal.Approved, &proposal.Status.Status)
		if err != nil {
			log.Printf("There was an error scanning the sql rows: %v", err)
		}

		proposals = append(proposals, proposal)

	}

	log.Printf("Proposals retrieved from Database: %d", len(proposals))
	respondWithJSON(w, http.StatusOK, proposals)

}

func decideProposaldHandler(w http.ResponseWriter, r *http.Request) {

	proposalId := mux.Vars(r)["id"]
	var decision ProposalDecision
	err := json.NewDecoder(r.Body).Decode(&decision)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}

	log.Printf("Updating Proposal By Id: %s", proposalId)

	updateStmt := `UPDATE Proposals set Status=$1, Approved=$2 where Id=$3`
	_, err = db.Exec(updateStmt, "DECIDED", decision.Approved, proposalId)
	if err != nil {
		log.Printf("There was an error executing the update query: %v", err)
	}

	rows, err := db.Query(`SELECT * FROM Proposals where id=$1`, proposalId)

	if err != nil {
		log.Printf("There was an error executing the query: %v", err)
	}

	defer rows.Close()

	//@TODO: validate that only one result comes from the query
	var proposal Proposal
	for rows.Next() {
		err = rows.Scan(&proposal.Id, &proposal.Title, &proposal.Description, &proposal.Email, &proposal.Author, &proposal.Approved, &proposal.Status.Status)
		if err != nil {
			log.Printf("There was an error scanning the sql rows: %v", err)
		}
	}

	if decision.Approved {
		log.Printf("Proposal Id: %s was approved!", proposalId)

		log.Printf("Publish Proposal Id: %s to the Conference Agenda", proposalId)
		log.Printf("Email Proposal's Id: %s author about decision", proposalId)

	} else {
		log.Printf("Proposal Id: %s was rejected!", proposalId)
		log.Printf("Email Proposal's Id: %s author about decision", proposalId)
	}

	respondWithJSON(w, http.StatusOK, proposal)
}

func newProposalHandler(w http.ResponseWriter, r *http.Request) {
	var proposal Proposal
	err := json.NewDecoder(r.Body).Decode(&proposal)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}

	proposal.Id = uuid.New().String()

	insertStmt := `insert into Proposals("id", "title", "description", "email", "author", "approved", "status") values($1, $2, $3, $4, $5, $6, $7)`

	_, err = db.Exec(insertStmt, proposal.Id, proposal.Title, proposal.Description, proposal.Email, proposal.Author, false, "PENDING")

	if err != nil {
		log.Printf("An error occured while executing query: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}

	log.Printf("Proposal Stored in Database: %s", proposal)

	respondWithJSON(w, http.StatusOK, proposal)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Starting C4P Service in Port: %s", appPort)

	connStr := "postgresql://" + POSTGRESQL_USERNAME + ":" + POSTGRESQL_PASSOWRD + "@" + POSTGRESQL_HOST + ":" + POSTGRESQL_PORT + "/postgres?sslmode=disable"
	log.Printf("Connecting to Database: %s.", connStr)
	// Connect to database
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to PostgreSQL.")

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/", newProposalHandler).Methods("POST")
	r.HandleFunc("/", getAllProposalsHandler).Methods("GET")
	r.HandleFunc("/{id}/decide", decideProposaldHandler).Methods("POST")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
		var info ServiceInfo = ServiceInfo{
			Name:         "C4P",
			Version:      VERSION,
			Source:       SOURCE,
			PodId:        POD_ID,
			PodNamespace: POD_NODENAME,
		}
		json.NewEncoder(w).Encode(info)
	})

	// Start the server; this is a blocking call
	err = http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
