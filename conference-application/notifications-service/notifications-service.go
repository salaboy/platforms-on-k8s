package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	kafka "github.com/segmentio/kafka-go"
)

type ServiceInfo struct {
	Name         string
	Version      string
	Source       string
	PodId        string
	PodNamespace string
	PodNodeName  string
}

type Notification struct {
	Id           string
	ProposalId   string
	AgendaItemId string
	Title        string
	EmailTo      string
	Accepted     bool
	EmailSubject string
	EmailBody    string
}

var VERSION = getEnv("VERSION", "1.0.0")
var SOURCE = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/agenda-service")
var POD_ID = getEnv("POD_ID", "N/A")
var POD_NAMESPACE = getEnv("POD_NAMESPACE", "N/A")
var POD_NODENAME = getEnv("POD_NODENAME", "N/A")
var KAFKA_URL = getEnv("KAFKA_URL", "localhost:9094")
var KAFKA_TOPIC = getEnv("KAFKA_TOPIC", "events-topic")

var notifications = []Notification{}

func getAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, notifications)

}

func sendNotificationHandler(kafkaWriter *kafka.Writer) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		notificationJson, err := json.Marshal(notification)
		if err != nil {
			log.Printf("An error occured while marshalling the notification to json: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}

		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("notification-sent-%s", notification.Id)),
			Value: notificationJson,
		}
		err = kafkaWriter.WriteMessages(r.Context(), msg)

		if err != nil {
			log.Printf("An error occured while writing the message to Kafka: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		log.Printf("Notification Sent Event emitted to Kafka: %s", notification)

		// @TODO avoid doing two marshals to json
		respondWithJSON(w, http.StatusOK, notification)
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

	log.Printf("Starting Notifications Service in Port: %s", appPort)

	kafkaAlive := isKafkaAlive(KAFKA_URL, KAFKA_TOPIC)
	if !kafkaAlive {
		log.Printf("Cannot connect to Kafka, restarting until it is healthy.")
		return
	}

	log.Printf("Connecting to Kafka Instance: %s, topic: %s.", KAFKA_URL, KAFKA_TOPIC)
	//https://github.com/segmentio/kafka-go/blob/main/examples/producer-api/main.go
	kafkaWriter := getKafkaWriter(KAFKA_URL, KAFKA_TOPIC)

	log.Printf("Connected to Kafka.")
	defer kafkaWriter.Close()

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/", sendNotificationHandler(kafkaWriter)).Methods("POST")
	r.HandleFunc("/", getAllNotificationsHandler).Methods("GET")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
		var info ServiceInfo = ServiceInfo{
			Name:         "NOTIFICATIONS",
			Version:      VERSION,
			Source:       SOURCE,
			PodId:        POD_ID,
			PodNamespace: POD_NODENAME,
		}
		json.NewEncoder(w).Encode(info)
	})

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+appPort, r)
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

func isKafkaAlive(kafkaURL string, topic string) bool {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	brokers, err := conn.Brokers()

	if err != nil {
		panic(err.Error())
	}

	for _, b := range brokers {
		log.Printf("Available Broker: %s", b)
	}
	if len(brokers) > 0 {
		return true
	} else {
		return false
	}

}
