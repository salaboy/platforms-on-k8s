package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	api "github.com/salaboy/platforms-on-k8s/conference-application/notifications-service/api"
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

type Notification struct {
	Id           string `json:"id"`
	ProposalId   string `json:"proposalId"`
	AgendaItemId string `json:"agendaItemId"`
	Title        string `json:"title"`
	EmailTo      string `json:"emailTo"`
	Accepted     bool   `json:"accepted"`
	EmailSubject string `json:"emailSubject"`
	EmailBody    string `json:"emailBody"`
}

type EventsEnabled struct {
	AgendaService        bool `json:"agenda-service"`
	NotificationsService bool `json:"notifications-service"`
	C4PService           bool `json:"c4p-service"`
}

func (s Notification) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

var (
	VERSION             = getEnv("VERSION", "2.0.0")
	SOURCE              = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/v2.0.0/conference-application/notifications-service")
	POD_NAME            = getEnv("POD_NAME", "N/A")
	POD_NAMESPACE       = getEnv("POD_NAMESPACE", "N/A")
	POD_NODENAME        = getEnv("POD_NODENAME", "N/A")
	POD_IP              = getEnv("POD_IP", "N/A")
	POD_SERVICE_ACCOUNT = getEnv("POD_SERVICE_ACCOUNT", "N/A")
	APP_PORT            = getEnv("APP_PORT", "8080")
	PUBSUB_NAME         = getEnv("PUBSUB_NAME", "conference-pubsub")
	PUBSUB_TOPIC        = getEnv("PUBSUB_TOPIC", "events-topic")
	TENANT_ID           = getEnv("TENANT_ID", "tenant-a")
	notifications       = []Notification{}
	FLAGD_HOST          = getEnv("FLAGD_HOST", "flagd.default.svc.cluster.local")
)

// respondWithJSON is a helper function to write a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	chiServer := NewChiServer()

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+APP_PORT, chiServer)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// getEnv returns the value of an environment variable, or a fallback value if not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// NewChiServer creates a new Chi server.
func NewChiServer() *chi.Mux {
	log.Printf("Starting Notifications Service in Port: %s", APP_PORT)

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

// OpenApi returns a handler that serves the OpenAPI spec as JSON.
func OpenAPI(r *chi.Mux) {
	fs := http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH") + "/docs/"))
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", fs))
}

// server implements the api.ServerInterface
type server struct {
	APIClient     dapr.Client
	FeatureClient *openfeature.Client
}

// NewServer creates a new api.ServerInterface instance.
func NewServer() api.ServerInterface {
	openfeature.SetProvider(flagd.NewProvider(
		flagd.WithHost(FLAGD_HOST),
		flagd.WithPort(8013),
	))

	openfeatureClient := openfeature.NewClient("notifications-service")

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	return &server{
		APIClient:     client,
		FeatureClient: openfeatureClient,
	}
}

// GetAllNotifications returns all notifications.
func (s *server) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, notifications)
}

func (s *server) areEventsEnabled() bool {
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

	return eventsEnabledStruct.NotificationsService
}

// CreateNotification creates a new notification.
func (s *server) CreateNotification(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}

	notification.Id = uuid.New().String()

	// Here you connect to email service and send the following Payload
	var status string
	var mood string
	if notification.Accepted {
		status = "accepted"
		mood = "happy"
	} else {
		status = "rejected"
		mood = "sad"
	}
	bodyAccepted := fmt.Sprintf("We will send further instructions closer to your presentation date. \n\t You session has been published into the conference website: %s", notification.AgendaItemId)
	bodyRejected := "We hope you submit again next year. \n\t Here is a %20 discount for this year ticket: \"PLATFORMS-ON-K8S\"."

	bodyThanks := "Thanks and see you soon. \n\n\t - The Conference Organizers -"

	emailTo := notification.EmailTo
	subject := fmt.Sprintf("Your proposal  %s was %s", notification.Title, status)
	notification.EmailSubject = subject
	body := fmt.Sprintf("We are %s to inform that your proposal: `%s` with title: `%s` was %s", mood, notification.ProposalId, notification.Title, status)

	if notification.Accepted {
		body = fmt.Sprintf("%s \n\t %s \n", body, bodyAccepted)
	} else {
		body = fmt.Sprintf("%s \n\t %s \n", body, bodyRejected)
	}
	body = fmt.Sprintf("%s \n\t %s  \n", body, bodyThanks)
	notification.EmailBody = body

	log.Printf("\n To: %s \n Subject: %s \n Body: %s \n", emailTo, subject, body)

	notifications = append(notifications, notification)

	log.Printf("Notification Sent: %v", notification)

	eventsEnabled := s.areEventsEnabled()
	log.Printf("Events Enabled? %s", eventsEnabled)
	if eventsEnabled {
		notificationJson, err := json.Marshal(notification)
		if err != nil {
			log.Printf("An error occured while marshalling the notification to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		event := Event{
			Id:      uuid.New().String(),
			Type:    "notification-sent",
			Payload: string(notificationJson),
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

		log.Printf("Notification Sent - Event published: %v", notification)
	}

	// @TODO avoid doing two marshals to json
	respondWithJSON(w, http.StatusOK, notification)
}

// GetServiceInfo returns service information.
func (s *server) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	var info ServiceInfo = ServiceInfo{
		Name:              "NOTIFICATIONS",
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
