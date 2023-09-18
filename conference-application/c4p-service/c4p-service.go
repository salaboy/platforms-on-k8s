package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/salaboy/platforms-on-k8s/conference-application/c4p-service/api"

	kafka "github.com/segmentio/kafka-go"
)

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Proposal struct {
	Id          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Author      string         `json:"author"`
	Email       string         `json:"email"`
	Approved    bool           `json:"approved"`
	Status      ProposalStatus `json:"status"`
}

// Event struct to encode events data
type Event struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

func (p *Proposal) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

type ProposalStatus struct {
	Status string `json:"status"`
}

type ProposalDecision struct {
	Approved bool `json:"approved"`
}

type ProposalRef struct {
	Id string `json:"id"`
}

type AgendaItem struct {
	Id          string      `json:"id"`
	Proposal    ProposalRef `json:"proposal"`
	Title       string      `json:"title"`
	Author      string      `json:"author"`
	Description string      `json:"description"`
}

type DecisionResponse struct {
	ProposalId string     `json:"proposalId"`
	AgendaItem AgendaItem `json:"agendaItem"`
	Proposal   Proposal   `json:"proposal"`
	Decision   bool       `json:"decision"`
}

type Notification struct {
	ProposalId   string `json:"proposalId"`
	AgendaItemId string `json:"agendaItemId"`
	Title        string `json:"title"`
	EmailTo      string `json:"emailTo"`
	Accepted     bool   `json:"accepted"`
}

type ServiceInfo struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	Source            string `json:"source"`
	PodName           string `json:"podName"`
	PodNamespace      string `json:"podNamespace"`
	PodNodeName       string `json:"podNodeName"`
	PodIp             string `json:"podIp"`
	PodServiceAccount string `json:"podServiceAccount"`
}

type Config struct {
	AgendaServiceUrl        string
	NotificationsServiceUrl string
}

var (
	Version            = getEnv("VERSION", "1.0.0")
	Source             = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/c4p-service")
	PodName            = getEnv("POD_NAME", "N/A")
	PodNamespace       = getEnv("POD_NAMESPACE", "N/A")
	PodNodeName        = getEnv("POD_NODENAME", "N/A")
	PodIp              = getEnv("POD_IP", "N/A")
	PodServiceAccount  = getEnv("POD_SERVICE_ACCOUNT", "N/A")
	PostgresqlHost     = getEnv("POSTGRES_HOST", "localhost")
	PostgresqlPort     = getEnv("POSTGRES_PORT", "5432")
	PostgresqlUsername = getEnv("POSTGRES_USERNAME", "postgres")
	PostgresqlPassword = getEnv("POSTGRES_PASSWORD", "postgres")

	KafkaUrl   = getEnv("KAFKA_URL", "localhost:9094")
	KafkaTopic = getEnv("KAFKA_TOPIC", "events-topic")
	AppPort    = getEnv("APP_PORT", "8080")

	AgendaServiceUrl        = getEnv("AGENDA_SERVICE_URL", "http://agenda-service.default.svc.cluster.local")
	NotificationsServiceUrl = getEnv("NOTIFICATIONS_SERVICE_URL", "http://notifications-service.default.svc.cluster.local")
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func isKafkaAlive(kafkaURL string, topic string) bool {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)
	if err != nil {
		panic(any(err.Error()))
	}
	defer conn.Close()

	brokers, err := conn.Brokers()

	if err != nil {
		panic(any(err.Error()))
	}

	for _, b := range brokers {
		log.Printf("Available Broker: %v", b)
	}
	if len(brokers) > 0 {
		return true
	} else {
		return false
	}

}

// server
type server struct {
	KafkaWriter             *kafka.Writer
	DB                      *sql.DB
	AgendaServiceURL        string
	NotificationsServiceURL string
}

// GetProposals gets all proposals.
func (s server) GetProposals(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	var query = "SELECT id, title, description, email, author, approved, status FROM Proposals p"
	if status != "" {
		query = fmt.Sprintf("%s where p.status=$1", query)
	}
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = s.DB.Query(query, status)
	} else {
		rows, err = s.DB.Query(query)
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

// CreateProposal creates a new proposal.
func (s server) CreateProposal(w http.ResponseWriter, r *http.Request) {
	var proposal Proposal
	err := json.NewDecoder(r.Body).Decode(&proposal)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	proposal.Status = ProposalStatus{Status: "PENDING"}
	proposal.Id = uuid.New().String()

	insertStmt := `insert into Proposals("id", "title", "description", "email", "author", "approved", "status") values($1, $2, $3, $4, $5, $6, $7)`

	_, err = s.DB.Exec(insertStmt, proposal.Id, proposal.Title, proposal.Description, proposal.Email, proposal.Author, false, "PENDING")

	if err != nil {
		log.Printf("An error occured while executing query: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Proposal Stored in Database: %v", proposal)

	proposalJson, err := json.Marshal(proposal)
	if err != nil {
		log.Printf("An error occured while marshalling the proposal to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	event := Event{
		Id:      uuid.New().String(),
		Type:    "new-proposal",
		Payload: string(proposalJson),
	}

	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("An error occured while marshalling the event for the proposal to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("new-proposal-%s", proposal.Id)),
		Value: eventJson,
	}
	err = s.KafkaWriter.WriteMessages(r.Context(), msg)

	if err != nil {
		log.Printf("An error occured while writing the message to Kafka: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("New Proposal Event emitted to Kafka: %v", proposal)
	respondWithJSON(w, http.StatusOK, proposal)
}

// DeleteProposal deletes a proposal.
func (s server) DeleteProposal(w http.ResponseWriter, r *http.Request, proposalId string) {

	var decision ProposalDecision
	err := json.NewDecoder(r.Body).Decode(&decision)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}

	log.Printf("Archiving Proposal By Id: %s", proposalId)

	updateStmt := `UPDATE Proposals set Status=$1 where Id=$2`
	_, err = s.DB.Exec(updateStmt, "ARCHIVED", proposalId)
	if err != nil {
		log.Printf("There was an error executing the update query: %v", err)
	}
	rows, err := s.DB.Query(`SELECT * FROM Proposals where id=$1`, proposalId)

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

	proposalJson, err := json.Marshal(proposal)
	if err != nil {
		log.Printf("An error occured while marshalling the proposal to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	event := Event{
		Id:      uuid.New().String(),
		Type:    "proposal-archived",
		Payload: string(proposalJson),
	}

	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("An error occured while marshalling the event for the proposal to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("proposal-archived-%s", proposal.Id)),
		Value: eventJson,
	}
	err = s.KafkaWriter.WriteMessages(r.Context(), msg)

	if err != nil {
		log.Printf("An error occured while writing the message to Kafka: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("Proposal Archived Event emitted to Kafka: %v", proposal)

	respondWithJSON(w, http.StatusOK, proposal)
}

// DecideProposal updates the status of a proposal.oa
func (s server) DecideProposal(w http.ResponseWriter, r *http.Request, proposalId string) {
	var decision ProposalDecision
	err := json.NewDecoder(r.Body).Decode(&decision)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}

	log.Printf("Updating Proposal By Id: %s", proposalId)

	updateStmt := `UPDATE Proposals set Status=$1, Approved=$2 where Id=$3`
	_, err = s.DB.Exec(updateStmt, "DECIDED", decision.Approved, proposalId)
	if err != nil {
		log.Printf("There was an error executing the update query: %v", err)
	}

	rows, err := s.DB.Query(`SELECT id, title, description, email, author, approved, status FROM Proposals where id=$1`, proposalId)

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
			Description: proposal.Description,
			Author:      proposal.Author,
		}
		agendaItemJson, err := json.Marshal(agendaItem)
		if err != nil {
			log.Printf("There was an error marshalling the Agenda Item to JSON: %v", err)
		}
		r, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/", s.AgendaServiceURL, "agenda-items"), bytes.NewBuffer(agendaItemJson))
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

		event := Event{
			Id:      uuid.New().String(),
			Type:    "proposal-approved",
			Payload: string(proposalJson),
		}

		eventJson, err := json.Marshal(event)
		if err != nil {
			log.Printf("An error occured while marshalling the event for the proposal to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("proposal-approved-%s", proposal.Id)),
			Value: eventJson,
		}
		err = s.KafkaWriter.WriteMessages(r.Context(), msg)

		if err != nil {
			log.Printf("An error occured while writing the message to Kafka: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		log.Printf("Proposal Approved Event emitted to Kafka: %v", proposal)

	} else {

		log.Printf("Proposal Id: %s was rejected!", proposalId)

	}
	log.Printf("Sending Notification to Proposal's author: %s author about decision", proposal.Email)

	notification := Notification{
		ProposalId:   decisionResponse.ProposalId,
		AgendaItemId: decisionResponse.AgendaItem.Id,
		Title:        decisionResponse.Proposal.Title,
		EmailTo:      decisionResponse.Proposal.Email,
		Accepted:     decisionResponse.Proposal.Approved,
	}

	notificationJson, err := json.Marshal(notification)
	if err != nil {
		log.Printf("An error occured while marshalling the proposal to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	r, err = http.NewRequest("POST", fmt.Sprintf("%s/%s/", s.NotificationsServiceURL, "notifications"), bytes.NewBuffer(notificationJson))
	if err != nil {
		log.Printf("There was an error creating the request to the Notifications Service: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		log.Printf("There was an error submitting the request to the Notifications Service: %v", err)
	} else {
		defer res.Body.Close()
	}

	respondWithJSON(w, http.StatusOK, proposal)
}

func (s server) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	var info = ServiceInfo{
		Name:              "C4P",
		Version:           Version,
		Source:            Source,
		PodName:           PodName,
		PodNodeName:       PodNodeName,
		PodNamespace:      PodNamespace,
		PodIp:             PodIp,
		PodServiceAccount: PodServiceAccount,
	}
	w.Header().Set(ContentType, ApplicationJson)
	json.NewEncoder(w).Encode(info)
}

func main() {
	r := NewChiServer(&Config{
		AgendaServiceUrl:        AgendaServiceUrl,
		NotificationsServiceUrl: NotificationsServiceUrl,
	})

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+AppPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// NewChiServer creates a new *chi.Mux server.
func NewChiServer(conf *Config) *chi.Mux {
	// create new chi router
	r := chi.NewRouter()

	// add logger middleware
	r.Use(middleware.Logger)

	log.Printf("Starting C4P Service in Port: %s", AppPort)

	fmt.Println("Environment Variables:")
	fmt.Println("Agenda Service URL: ", conf.AgendaServiceUrl)
	fmt.Println("Notifications Service URL: ", conf.NotificationsServiceUrl)

	// connect to database
	db := NewDB()

	// check if database is alive
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to PostgreSQL.")

	// create a new *kafka.Writer using the given broker addresses and configuration
	kafkaWriter := NewKafkaWriter(KafkaUrl, KafkaTopic)

	// Create a new server
	server := NewServer(kafkaWriter, db, conf.AgendaServiceUrl, conf.NotificationsServiceUrl)
	OpenAPI(r)

	// mount the API on the server
	r.Mount("/", api.Handler(server))

	// add health check
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// add openapi spec

	return r
}

// OpenAPI OpenAPIHandler returns a handler that serves the OpenAPI documentation.
func OpenAPI(r *chi.Mux) {
	fs := http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH") + "/docs/"))
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", fs))
}

// NewServer creates a new api.ServerInterface.
func NewServer(kafkaWriter *kafka.Writer, db *sql.DB, agendaServiceURL string, notificationsServiceURL string) api.ServerInterface {
	return &server{
		KafkaWriter:             kafkaWriter,
		DB:                      db,
		AgendaServiceURL:        agendaServiceURL,
		NotificationsServiceURL: notificationsServiceURL,
	}
}

// NewKafkaWriter creates a new *kafka.Write.
func NewKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	log.Printf("Connecting to Kafka Instance: %s, topic: %s.", KafkaUrl, KafkaTopic)
	kafkaAlive := isKafkaAlive(KafkaUrl, KafkaTopic)
	if !kafkaAlive {
		log.Printf("Cannot connect to Kafka, restarting until it is healthy.")
		panic(any("Cannot connect to Kafka"))
	}

	log.Printf("Connected to Kafka.")

	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewDB() *sql.DB {
	connStr := "postgresql://" + PostgresqlUsername + ":" + PostgresqlPassword + "@" + PostgresqlHost + ":" + PostgresqlPort + "/postgres?sslmode=disable"
	log.Printf("Connecting to Database: %s.", connStr)
	// Connect to database

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
