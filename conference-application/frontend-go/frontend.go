package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	kafka "github.com/segmentio/kafka-go"
)

var VERSION = getEnv("VERSION", "1.0.0")
var SOURCE = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/frontend-go")
var POD_NAME = getEnv("POD_NAME", "N/A")
var POD_NAMESPACE = getEnv("POD_NAMESPACE", "N/A")
var POD_NODENAME = getEnv("POD_NODENAME", "N/A")
var POD_IP = getEnv("POD_IP", "N/A")
var POD_SERVICE_ACCOUNT = getEnv("POD_SERVICE_ACCOUNT", "N/A")
var AGENDA_SERVICE_URL = getEnv("AGENDA_SERVICE_URL", "http://agenda-service.default.svc.cluster.local")
var C4P_SERVICE_URL = getEnv("C4P_SERVICE_URL", "http://c4p-service.default.svc.cluster.local")
var NOTIFICATIONS_SERVICE_URL = getEnv("NOTIFICATIONS_SERVICE_URL", "http://notifications-service.default.svc.cluster.local")

var KAFKA_URL = getEnv("KAFKA_URL", "localhost:9094")
var KAFKA_TOPIC = getEnv("KAFKA_TOPIC", "events-topic")
var KAFKA_GROUP_ID = getEnv("KAFKA_GROUP_ID", "app")

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

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

var events = []Event{}

type Event struct {
	Id      int64
	Type    string
	Payload string
}

type Features struct {
	DEBUG_ENABLED     string
	GENERATE_PROPOSAL string
}

var FEATURE_DEBUG_ENABLED = getEnv("FEATURE_DEBUG_ENABLED", "false")

// values: PUBLIC (no filters), GENERATE (Read Only Form - Generate Proposal), GENERATE_ONLY (No Submit until Generated Proposal is created)
var FEATURE_GENERATE_PROPOSAL = getEnv("FEATURE_DEBUG_ENABLED", "GENERATE")

func featureHandler(w http.ResponseWriter, r *http.Request) {
	var features = Features{
		DEBUG_ENABLED:     FEATURE_DEBUG_ENABLED,
		GENERATE_PROPOSAL: FEATURE_GENERATE_PROPOSAL,
	}
	respondWithJSON(w, http.StatusOK, features)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, events)
}

func agendaServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("api/agenda", AGENDA_SERVICE_URL, w, r)
}

func c4PServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("api/c4p", C4P_SERVICE_URL, w, r)
}

func notificationServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("api/notifications", NOTIFICATIONS_SERVICE_URL, w, r)
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    5,    // 5B
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
		MaxWait:     3 * time.Second,
	})
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
func proxyRequest(serviceName string, serviceUrl string, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(requestDump))

	url := fmt.Sprintf("%s%s", serviceUrl, r.RequestURI)
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}

	log.Printf("Proxying request before replace to %s", url)
	// remove the service path
	url = strings.Replace(url, serviceName+"/", "", -1)

	log.Printf("Proxying request to %s", url)

	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for h, val := range resp.Header {
		w.Header()[h] = val
	}

	w.WriteHeader(resp.StatusCode)

	log.Printf("Proxied request response code %s - %d", resp.Status, resp.StatusCode)

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r := mux.NewRouter()

	r.PathPrefix("/api/agenda/").HandlerFunc(agendaServiceHandler)
	r.PathPrefix("/api/c4p/").HandlerFunc(c4PServiceHandler)
	r.PathPrefix("/api/notifications/").HandlerFunc(notificationServiceHandler)
	r.PathPrefix("/api/events/").HandlerFunc(eventsHandler)
	r.PathPrefix("/api/features/").HandlerFunc(featureHandler)
	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
		var info ServiceInfo = ServiceInfo{
			Name:              "FRONTEND",
			Version:           VERSION,
			Source:            SOURCE,
			PodName:           POD_NAME,
			PodNodeName:       POD_NODENAME,
			PodNamespace:      POD_NAMESPACE,
			PodIp:             POD_IP,
			PodServiceAccount: POD_SERVICE_ACCOUNT,
		}
		w.Header().Set(ContentType, ApplicationJson)
		json.NewEncoder(w).Encode(info)
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))

	log.Printf("Connecting to Kafka Instance: %s, topic: %s., group: %s", KAFKA_URL, KAFKA_TOPIC, KAFKA_GROUP_ID)
	reader := getKafkaReader(KAFKA_URL, KAFKA_TOPIC, KAFKA_GROUP_ID)

	kafkaAlive := isKafkaAlive(KAFKA_URL, KAFKA_TOPIC)
	if !kafkaAlive {
		log.Printf("Cannot connect to Kafka, restarting until it is healthy.")
		return
	}

	go consumeFromKafka(reader)

	defer reader.Close()

	log.Printf("Starting Frontend Go in Port: %s", appPort)

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

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

func consumeFromKafka(reader *kafka.Reader) {
	fmt.Println("Consuming Events ...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		var event = Event{
			Id:      m.Offset,
			Type:    string(m.Key),
			Payload: string(m.Value),
		}
		events = append(events, event)
	}
}
