package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	api "github.com/salaboy/platforms-on-k8s/conference-application/agenda-service/api"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/pkg/openfeature"
)

var (
	KEY                 = "AGENDAITEMS"
	VERSION             = getEnv("VERSION", "2.0.0")
	SOURCE              = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/v2.0.0/conference-application/agenda-service")
	POD_NAME            = getEnv("POD_NAME", "N/A")
	POD_NAMESPACE       = getEnv("POD_NAMESPACE", "N/A")
	POD_NODENAME        = getEnv("POD_NODENAME", "N/A")
	POD_IP              = getEnv("POD_IP", "N/A")
	POD_SERVICE_ACCOUNT = getEnv("POD_SERVICE_ACCOUNT", "N/A")

	APP_PORT        = getEnv("APP_PORT", "8080")
	STATESTORE_NAME = getEnv("STATESTORE_NAME", "agenda-service-statestore")
	PUBSUB_NAME     = getEnv("PUBSUB_NAME", "conference-pubsub")
	PUBSUB_TOPIC    = getEnv("PUBSUB_TOPIC", "events-topic")
	TENANT_ID       = getEnv("TENANT_ID", "tenant-a")

	FLAGD_HOST = getEnv("FLAGD_HOST", "flagd.default.svc.cluster.local")
)

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

// Event struct to encode events data
type Event struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

// Proposal is a struct to represent a Proposal.
type Proposal struct {
	Id string `json:"id"`
}

// AgendaItem is a struct to represent an Agenda Item.
type AgendaItem struct {
	Id          string   `json:"id,omitempty"`
	Proposal    Proposal `json:"proposal"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Archived    bool     `json:"archived"`
}

type EventsEnabled struct {
	AgendaService        bool `json:"agenda-service"`
	NotificationsService bool `json:"notifications-service"`
	C4PService           bool `json:"c4p-service"`
}

// MarshalBinary is a custom marshaler for AgendaItem.
func (s AgendaItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// ServiceInfo is a struct to represent a Service Info describing the service and the pod it is running on.
type ServiceInfo struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	Source            string `json:"source"`
	PodName           string `json:"podName"`
	PodNamespace      string `json:"podNamespace"`
	PodNodeName       string `json:"podNodeName"`
	PodIp             string `json:"podIp"`
	PodServiceAccount string `json:"podServiceAccount"`
	EventsEnabled     bool   `json:"eventsEnabled"`
}

// main is the entrypoint of the application.
func main() {
	chiServer := NewChiServer()
	err := http.ListenAndServe(":"+APP_PORT, chiServer)
	log.Printf("Starting Agenda Service in Port: %s", APP_PORT)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// getEnv returns the value of an environment variable, or a fallback value if it is not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// respondWithJSON is a helper function to write a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)
	w.Write(response)
}

// errorHandler is a helper function to write an error response.
func errorHandler(w http.ResponseWriter, statusCode int, message string) {
	mapResponse := map[string]interface{}{
		"message": message,
	}
	respondWithJSON(w, statusCode, mapResponse)
}

// server is the API server struct that implements api.ServerInterface.
type server struct {
	APIClient     dapr.Client
	FeatureClient *openfeature.Client
}

// GetAgendaItems is a handler to get all Agenda Items.
func (s server) GetAgendaItems(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	agendaItemsStateItem, err := s.APIClient.GetState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), nil)
	if err != nil {
		log.Printf("An error occured while getting agenda items from store: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}
	agendaItems := []AgendaItem{}
	err = json.Unmarshal(agendaItemsStateItem.Value, &agendaItems)
	if err != nil {
		log.Printf("There was an error decoding the AgendaItems into the struct array: %v", err)
	}

	// Let's remove the Archived Agenda Items
	for k, v := range agendaItems {
		if v.Archived {
			agendaItems = RemoveIndex(agendaItems, k)
		}
	}

	log.Printf("Agenda Items retrieved from Database: %d", len(agendaItems))
	respondWithJSON(w, http.StatusOK, agendaItems)
}

func RemoveIndex(s []AgendaItem, index int) []AgendaItem {
	ret := make([]AgendaItem, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

// CreateAgendaItem is a handler to create an Agenda Item.
func (s server) CreateAgendaItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var agendaItem AgendaItem
	err := json.NewDecoder(r.Body).Decode(&agendaItem)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		errorHandler(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := s.APIClient.GetState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), nil)
	if err != nil {
		log.Printf("An error occured while getting agenda items from store: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}
	agendaItems := []AgendaItem{}
	if result != nil {
		json.Unmarshal(result.Value, &agendaItems)
	}

	agendaItem.Id = uuid.New().String()
	log.Printf("Creating Agenda Item: %v", agendaItem)
	agendaItems = append(agendaItems, agendaItem)

	jsonData, err := json.Marshal(agendaItems)
	if err != nil {
		log.Printf("An error occured while marshalling agenda items to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.APIClient.SaveState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), jsonData, nil); err != nil {
		panic(err)
	}

	log.Printf("Agenda Item Stored in Database: %v", agendaItem)
	eventsEnabled := s.areEventsEnabled()
	log.Printf("Events Enabled? %s", eventsEnabled)
	if eventsEnabled {
		agendaItemJson, err := json.Marshal(agendaItem)
		if err != nil {
			log.Printf("An error occured while marshalling the agenda item to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		event := Event{
			Id:      uuid.New().String(),
			Type:    "new-agenda-item",
			Payload: string(agendaItemJson),
		}
		eventJson, err := json.Marshal(event)
		if err != nil {
			log.Printf("An error occured while marshalling the event to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		//@TODO: add tenant to PUBSUB_TOPIC
		if err := s.APIClient.PublishEvent(ctx, PUBSUB_NAME, PUBSUB_TOPIC, eventJson); err != nil {
			log.Printf("An error occured while publishing the event: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		log.Printf("New Agenda Item Event emitted: %v", agendaItem)
	}

	respondWithJSON(w, http.StatusCreated, agendaItem)
}

// GetAgendaItemById is a handler to get an Agenda Item by ID.
func (s server) GetAgendaItemById(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	agendaItemId := chi.URLParam(r, "id")
	log.Printf("Fetching Agenda Item By Id: %s", agendaItemId)

	result, err := s.APIClient.GetState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), nil)
	if err != nil {
		log.Printf("An error occured while getting agenda items from store: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}
	agendaItems := []AgendaItem{}
	if result != nil {
		json.Unmarshal(result.Value, &agendaItems)
	}

	agendaItemById := getAgendaItemById(agendaItems, agendaItemId)
	if agendaItemById.Id == "" {
		errorHandler(w, http.StatusNotFound, "Agenda Item not found")
	}

	log.Printf("Agenda Item retrieved from Database: %v", agendaItemById)
	respondWithJSON(w, http.StatusOK, agendaItemById)
}

func getAgendaItemById(agendaItems []AgendaItem, agendaItemId string) AgendaItem {
	var agendaItemById AgendaItem
	for _, v := range agendaItems {
		if v.Id == agendaItemId {
			agendaItemById = v
			break
		}
	}
	return agendaItemById
}

// GetServiceInfo is a handler to get service info.
func (s server) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	var info ServiceInfo = ServiceInfo{
		Name:              "AGENDA",
		Version:           VERSION,
		Source:            SOURCE,
		PodName:           POD_NAME,
		PodNodeName:       POD_NODENAME,
		PodNamespace:      POD_NAMESPACE,
		PodIp:             POD_IP,
		PodServiceAccount: POD_SERVICE_ACCOUNT,
		EventsEnabled:     s.areEventsEnabled(),
	}
	w.Header().Set(ContentType, ApplicationJson)
	json.NewEncoder(w).Encode(info)
}

// ArchiveAgendaItemById is a handler to archive an Agenda Item by ID.
func (s server) ArchiveAgendaItemById(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	agendaItemId := chi.URLParam(r, "id")
	log.Printf("Fetching Agenda Item By Id: %s", agendaItemId)

	result, err := s.APIClient.GetState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), nil)
	if err != nil {
		log.Printf("An error occured while getting agenda items from store: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
	}
	agendaItems := []AgendaItem{}
	if result != nil {
		json.Unmarshal(result.Value, &agendaItems)
	}

	// Archiving Agenda Item
	var archivedAgendaItem AgendaItem
	for k, v := range agendaItems {
		if v.Id == agendaItemId {
			v.Archived = true
			agendaItems = RemoveIndex(agendaItems, k)
			archivedAgendaItem = v
			break
		}
	}

	agendaItems = append(agendaItems, archivedAgendaItem)

	agendaItemsJson, err := json.Marshal(agendaItems)
	if err != nil {
		log.Printf("An error occured while marshalling the agenda items to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.APIClient.SaveState(ctx, STATESTORE_NAME, fmt.Sprintf("%s-%s", TENANT_ID, KEY), agendaItemsJson, nil); err != nil {
		panic(err)
	}

	log.Printf("Agenda Item Archived in Database: %v", archivedAgendaItem)
	eventsEnabled := s.areEventsEnabled()
	log.Printf("Events Enabled? %s", eventsEnabled)
	if eventsEnabled {
		agendaItemJson, err := json.Marshal(archivedAgendaItem)
		if err != nil {
			log.Printf("An error occured while marshalling the agenda item to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		event := Event{
			Id:      uuid.New().String(),
			Type:    "agenda-item-archived",
			Payload: string(agendaItemJson),
		}
		eventJson, err := json.Marshal(event)
		if err != nil {
			log.Printf("An error occured while marshalling the event to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		//@TODO: add tenant to PUBSUB_TOPIC
		if err := s.APIClient.PublishEvent(ctx, PUBSUB_NAME, PUBSUB_TOPIC, eventJson); err != nil {
			log.Printf("An error occured while publishing the event: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		log.Printf("Agenda Item Archived Event published: %v", archivedAgendaItem)
	}

	respondWithJSON(w, http.StatusNoContent, archivedAgendaItem)
}

func (s server) areEventsEnabled() bool {
	ctx := context.Background()
	eventsEnabled, err := s.FeatureClient.ObjectValue(ctx, "eventsEnabled", EventsEnabled{}, openfeature.EvaluationContext{})
	if err != nil {
		log.Println("failed to find Feature Flag `eventsEnabled`.")
		return false
	}

	jsonData, err := json.Marshal(eventsEnabled)
	if err != nil {
		log.Fatal(err)
	}

	var eventsEnabledStruct EventsEnabled
	err = json.Unmarshal(jsonData, &eventsEnabledStruct)
	if err != nil {
		log.Fatal(err)
	}

	return eventsEnabledStruct.AgendaService
}

// NewServer creates a new api.ServerInterface.
func NewServer() api.ServerInterface {
	openfeature.SetProvider(flagd.NewProvider(
		flagd.WithHost(FLAGD_HOST),
		flagd.WithPort(8013),
	))

	openfeatureClient := openfeature.NewClient("agenda-service")

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	return &server{
		APIClient:     client,
		FeatureClient: openfeatureClient,
	}
}

// NewChiServer creates a new *chi.Mux server.
func NewChiServer() *chi.Mux {
	log.Printf("Starting Agenda Service in Port: %s", APP_PORT)
	// create new router
	r := chi.NewRouter()

	// add middlewares
	r.Use(middleware.Logger)

	// create new server
	server := NewServer()

	// add openapi spec
	OpenAPI(r)

	// add routes
	r.Mount("/", api.Handler(server))

	// add health check
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	return r
}

// OpenAPIHandler returns a handler that serves the OpenAPI documentation.
func OpenAPI(r *chi.Mux) {
	fs := http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH") + "/docs/"))
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", fs))
}
