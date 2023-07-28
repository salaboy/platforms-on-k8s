package main

import (
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

	dapr "github.com/dapr/go-sdk/client"
	"github.com/salaboy/platforms-on-k8s/conference-application/c4p-service/api"
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

type Event struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

type Event struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
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

var (
	Version            = getEnv("VERSION", "2.0.0")
	Source             = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/v2.0.0/conference-application/c4p-service")
	PodName            = getEnv("POD_NAME", "N/A")
	PodNamespace       = getEnv("POD_NAMESPACE", "N/A")
	PodNodeName        = getEnv("POD_NODENAME", "N/A")
	PodIp              = getEnv("POD_IP", "N/A")
	PodServiceAccount  = getEnv("POD_SERVICE_ACCOUNT", "N/A")
	PostgresqlHost     = getEnv("POSTGRES_HOST", "localhost")
	PostgresqlPort     = getEnv("POSTGRES_PORT", "5432")
	PostgresqlUsername = getEnv("POSTGRES_USERNAME", "postgres")
	PostgresqlPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	PubSubName         = getEnv("PUBSUB_NAME", "conference-pubsub")
	PubSubTopic        = getEnv("PUBSUB_TOPIC", "events-topic")
	TenantId           = getEnv("TENANT_ID", "tenant-a")
	AppPort            = getEnv("APP_PORT", "8080")
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

// server
type server struct {
	APIClient               *dapr.Client
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

	ctx := context.Background()
	var proposal Proposal
	err := json.NewDecoder(r.Body).Decode(&proposal)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	proposal.Id = uuid.New().String()
	proposal.Status = ProposalStatus{Status: "PENDING"}

	insertStmt := `insert into Proposals("id", "title", "description", "email", "author", "approved", "status") values($1, $2, $3, $4, $5, $6, $7)`

	_, err = s.DB.Exec(insertStmt, proposal.Id, proposal.Title, proposal.Description, proposal.Email, proposal.Author, false, proposal.Status)

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
	event := Event{
		Id:      uuid.New().String(),
		Type:    "new-proposal",
		Payload: string(proposalJson),
	}
	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("An error occured while marshalling the event to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	//@TODO: add tenant to PUBSUB_TOPIC
	if err := s.APIClient.PublishEvent(ctx, PubSubName, PubSubTopic, eventJson); err != nil {
		log.Printf("An error occured while publishing the event: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("New Proposal Event published: %s", proposal)
	respondWithJSON(w, http.StatusOK, proposal)

}

// DeleteProposal deletes a proposal.
func (s server) DeleteProposal(w http.ResponseWriter, r *http.Request, proposalId string) {
	ctx := context.Background()
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
		log.Printf("An error occured while marshalling the event to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.APIClient.PublishEvent(ctx, PubSubName, PubSubTopic, eventJson); err != nil {
		log.Printf("An error occured while publishing the event: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Proposal Archived Event published: %s", proposal)

	respondWithJSON(w, http.StatusOK, proposal)
}

// DecideProposal updates the status of a proposal.oa
func (s server) DecideProposal(w http.ResponseWriter, r *http.Request, proposalId string) {
	ctx := context.Background()
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

		content := &dapr.DataContent{
			ContentType: "application/json",
			Data:        agendaItemJson,
		}

		res, err := s.APIClient.InvokeMethodWithContent(ctx, "agenda-service", "agenda-items/", "POST", content)
		if err != nil {
			log.Printf("There was an error creating the request to the Agenda Item Service: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
		}

		log.Printf("Response from calling agenda-service %s: ", res)

		var agendaItemResponse AgendaItem
		err = json.Unmarshal(res, &agendaItemResponse)
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

		//@TODO: add tenant to PUBSUB_TOPIC
		event := Event{
			Id:      uuid.New().String(),
			Type:    "proposal-approved",
			Payload: string(proposalJson),
		}
		eventJson, err := json.Marshal(event)
		if err != nil {
			log.Printf("An error occured while marshalling the event to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		if err := s.APIClient.PublishEvent(ctx, PubSubName, PubSubTopic, eventJson); err != nil {
			log.Printf("An error occured while publishing the event: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		log.Printf("Proposal Approved Event published: %s", proposal)

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
	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        notificationJson,
	}

	resp, err := s.APIClient.InvokeMethodWithContent(ctx, "notifications-service", "notifications/", "POST", content)
	if err != nil {
		log.Printf("There was an error creating the request to the Notifications Service: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}
	log.Printf("Response from calling notification service %s", resp)

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
	r := NewChiServer()

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+AppPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// NewChiServer creates a new *chi.Mux server.
func NewChiServer() *chi.Mux {
	// create new chi router
	r := chi.NewRouter()

	// add logger middleware
	r.Use(middleware.Logger)

	log.Printf("Starting C4P Service in Port: %s", AppPort)

	// connect to database
	db := NewDB()

	// check if database is alive
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to PostgreSQL.")

	APIClient, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	// Create a new server
	server := NewServer(APIClient, db)
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

func NewServer(daprClient *dapr.Client, db *sql.DB) api.ServerInterface {
	return &server{
		APIClient: daprClient,
		DB:        db,
	}
}
