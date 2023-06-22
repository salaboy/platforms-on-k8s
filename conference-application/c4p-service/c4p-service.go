package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	kafka "github.com/segmentio/kafka-go"
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

type ProposalRef struct {
	Id string
}

type AgendaItem struct {
	Id       string
	Proposal ProposalRef
	Title    string
	Author   string
	Day      string
	Time     string
}

type DecisionResponse struct {
	ProposalId string
	AgendaItem AgendaItem
	Proposal   Proposal
	Decision   bool
}

type Notification struct {
	ProposalId   string
	AgendaItemId string
	Title        string
	EmailTo      string
	Accepted     bool
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
var POSTGRESQL_PASSOWRD = getEnv("POSTGRES_PASSWORD", "postgres")
var AGENDA_SERVICE_URL = getEnv("AGENDA_SERVICE_URL", "http://agenda-service.default.svc.cluster.local")
var NOTIFICATIONS_SERVICE_URL = getEnv("NOTIFICATIONS_SERVICE_URL", "http://notifications-service.default.svc.cluster.local")

var KAFKA_URL = getEnv("KAFKA_URL", "localhost:9094")
var KAFKA_TOPIC = getEnv("KAFKA_TOPIC", "events-topic")

var db *sql.DB

func getAllProposalsHandler(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")
	var query = "SELECT * FROM Proposals p"
	if status != "" {
		query = fmt.Sprintf("%s where p.status=$1", query)
	}
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = db.Query(query, status)
	} else {
		rows, err = db.Query(query)
	}

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

func decideProposaldHandler(kafkaWriter *kafka.Writer) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
		var decisionResponse DecisionResponse
		decisionResponse.ProposalId = proposalId
		decisionResponse.Decision = decision.Approved
		decisionResponse.Proposal = proposal
		if decision.Approved {
			log.Printf("Proposal Id: %s was approved!", proposalId)
			log.Printf("Publish Proposal Id: %s to the Conference Agenda", proposalId)
			agendaItem := AgendaItem{
				Title: proposal.Title,
				Proposal: ProposalRef{
					Id: proposal.Id,
				},
				Author: proposal.Author,
				Day:    "",
				Time:   "",
			}
			agendaItemJson, err := json.Marshal(agendaItem)
			if err != nil {
				log.Printf("There was an error marshalling the Agenda Item to JSON: %v", err)
			}
			r, err := http.NewRequest("POST", AGENDA_SERVICE_URL, bytes.NewBuffer(agendaItemJson))
			if err != nil {
				log.Printf("There was an error creating the request to the Agenda Item Service: %v", err)
			}
			r.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			res, err := client.Do(r)
			if err != nil {
				log.Printf("There was an error submitting the request to the Agenda Item Service: %v", err)
			}
			defer res.Body.Close()
			var agendaItemResponse AgendaItem
			err = json.NewDecoder(res.Body).Decode(&agendaItemResponse)
			if err != nil {
				log.Printf("There was an error decoding the request body into the struct: %v", err)
			}
			decisionResponse.AgendaItem = agendaItemResponse

			proposalJson, err := json.Marshal(proposal)
			if err != nil {
				log.Printf("An error occured while marshalling the proposal to json: %v", err)
				respondWithJSON(w, http.StatusInternalServerError, err)
				return
			}
			msg := kafka.Message{
				Key:   []byte(fmt.Sprintf("proposal-approved-%s", proposal.Id)),
				Value: proposalJson,
			}
			err = kafkaWriter.WriteMessages(r.Context(), msg)

			if err != nil {
				log.Printf("An error occured while writing the message to Kafka: %v", err)
				respondWithJSON(w, http.StatusInternalServerError, err)
				return
			}
			log.Printf("Proposal Approved Event emitted to Kafka: %s", proposal)

		} else {

			log.Printf("Proposal Id: %s was rejected!", proposalId)

		}
		log.Printf("Sending Notification to Proposal's author: %s author about decision", proposal.Email)

		notification := Notification{
			ProposalId:   decisionResponse.ProposalId,
			AgendaItemId: decisionResponse.AgendaItem.Id,
			Title:        decisionResponse.AgendaItem.Title,
			EmailTo:      decisionResponse.AgendaItem.Author,
			Accepted:     decisionResponse.Proposal.Approved,
		}

		notificationJson, err := json.Marshal(notification)
		if err != nil {
			log.Printf("An error occured while marshalling the proposal to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		r, err = http.NewRequest("POST", NOTIFICATIONS_SERVICE_URL, bytes.NewBuffer(notificationJson))
		if err != nil {
			log.Printf("There was an error creating the request to the Notifications Service: %v", err)
		}
		r.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			log.Printf("There was an error submitting the request to the Agenda Item Service: %v", err)
		}
		defer res.Body.Close()

		respondWithJSON(w, http.StatusOK, proposal)
	})
}

func newProposalHandler(kafkaWriter *kafka.Writer) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var proposal Proposal
		err := json.NewDecoder(r.Body).Decode(&proposal)
		if err != nil {
			log.Printf("There was an error decoding the request body into the struct: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		proposal.Id = uuid.New().String()

		insertStmt := `insert into Proposals("id", "title", "description", "email", "author", "approved", "status") values($1, $2, $3, $4, $5, $6, $7)`

		_, err = db.Exec(insertStmt, proposal.Id, proposal.Title, proposal.Description, proposal.Email, proposal.Author, false, "PENDING")

		if err != nil {
			log.Printf("An error occured while executing query: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		log.Printf("Proposal Stored in Database: %s", proposal)

		proposalJson, err := json.Marshal(proposal)
		if err != nil {
			log.Printf("An error occured while marshalling the proposal to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("new-proposal-%s", proposal.Id)),
			Value: proposalJson,
		}
		err = kafkaWriter.WriteMessages(r.Context(), msg)

		if err != nil {
			log.Printf("An error occured while writing the message to Kafka: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		log.Printf("New Proposal Event emitted to Kafka: %s", proposal)
		respondWithJSON(w, http.StatusOK, proposal)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
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

	log.Printf("Connecting to Kafka Instance: %s, topic: %s.", KAFKA_URL, KAFKA_TOPIC)
	//https://github.com/segmentio/kafka-go/blob/main/examples/producer-api/main.go
	kafkaWriter := getKafkaWriter(KAFKA_URL, KAFKA_TOPIC)

	log.Printf("Connected to Kafka.")
	defer kafkaWriter.Close()

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/", newProposalHandler(kafkaWriter)).Methods("POST")
	r.HandleFunc("/", getAllProposalsHandler).Methods("GET")
	r.HandleFunc("/{id}/decide", decideProposaldHandler(kafkaWriter)).Methods("POST")

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
