package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	kafka "github.com/segmentio/kafka-go"
)

var (
	KEY                 = "AGENDAITEMS"
	VERSION             = getEnv("VERSION", "1.0.0")
	SOURCE              = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/agenda-service")
	POD_NAME            = getEnv("POD_NAME", "N/A")
	POD_NAMESPACE       = getEnv("POD_NAMESPACE", "N/A")
	POD_NODENAME        = getEnv("POD_NODENAME", "N/A")
	POD_IP              = getEnv("POD_IP", "N/A")
	POD_SERVICE_ACCOUNT = getEnv("POD_SERVICE_ACCOUNT", "N/A")
	REDIS_HOST          = getEnv("REDIS_HOST", "localhost")
	REDIS_PORT          = getEnv("REDIS_PORT", "6379")
	REDIS_PASSWORD      = getEnv("REDIS_PASSWORD", "")
	KAFKA_URL           = getEnv("KAFKA_URL", "localhost:9094")
	KAFKA_TOPIC         = getEnv("KAFKA_TOPIC", "events-topic")
)

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Proposal struct {
	Id string `json:"id"`
}

type AgendaItem struct {
	Id          string   `json:"id"`
	Proposal    Proposal `json:"proposal"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Archived    bool     `json:"archived"`
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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/", Handler(NewServer()))

	// Add OpenAPI Handler
	OpenAPI(r)

	appPort := getEnv("APP_PORT", "8080")
	err := http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// getEnv returns the value of an environment variable, or a fallback value if
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// respondWithJSON is a helper function to write a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)
	w.Write(response)
}

type server struct {
	KafkaWriter *kafka.Writer
	RedisClient *redis.Client
}

// Get all Agenda Items
func (s server) GetAgendaItems(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	agendaItemsHashs, err := s.RedisClient.HGetAll(ctx, KEY).Result()
	if err != nil {
		panic(err)
	}
	agendaItems := []AgendaItem{}
	for _, ai := range agendaItemsHashs {
		var agendaItem AgendaItem
		err = json.Unmarshal([]byte(ai), &agendaItem)
		if err != nil {
			log.Printf("There was an error decoding the AgendaItem into the struct: %v", err)
		}
		if !agendaItem.Archived {
			agendaItems = append(agendaItems, agendaItem)
		}
	}
	log.Printf("Agenda Items retrieved from Database: %d", len(agendaItems))
	respondWithJSON(w, http.StatusOK, agendaItems)
}

// Create an Agenda Item
func (s server) CreateAgendaItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var agendaItem AgendaItem
	err := json.NewDecoder(r.Body).Decode(&agendaItem)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}

	// @TODO: write fail scenario (check for fail string in title return 500)
	agendaItem.Id = uuid.New().String()

	err = s.RedisClient.HSetNX(ctx, KEY, agendaItem.Id, agendaItem).Err()
	if err != nil {
		panic(err)
	}

	log.Printf("Agenda Item Stored in Database: %v", agendaItem)

	agendaItemJson, err := json.Marshal(agendaItem)
	if err != nil {
		log.Printf("An error occured while marshalling the agenda item to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("new-agenda-item-%s", agendaItem.Id)),
		Value: agendaItemJson,
	}
	err = s.KafkaWriter.WriteMessages(r.Context(), msg)

	if err != nil {
		log.Printf("An error occured while writing the message to Kafka: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("New Agenda Item Event emitted to Kafka: %v", agendaItem)

	// @TODO avoid doing two marshals to json
	respondWithJSON(w, http.StatusOK, agendaItem)
}

// Get an Agenda Item by ID
func (s server) GetAgendaItemById(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	agendaItemId := chi.URLParam(r, "id")
	log.Printf("Fetching Agenda Item By Id: %s", agendaItemId)
	agendaItemById, err := s.RedisClient.HGet(ctx, KEY, agendaItemId).Result()
	if err != nil {
		panic(err)
	}
	var agendaItem AgendaItem
	err = json.Unmarshal([]byte(agendaItemById), &agendaItem)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}
	log.Printf("Agenda Item retrieved from Database: %v", agendaItem)
	respondWithJSON(w, http.StatusOK, agendaItem)
}

// Get Service Info
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
	}
	json.NewEncoder(w).Encode(info)
}

// Archive an Agenda Item by ID
func (s server) ArchiveAgendaItemById(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	agendaItemId := chi.URLParam(r, "id")
	log.Printf("Fetching Agenda Item By Id: %s", agendaItemId)
	agendaItemById, err := s.RedisClient.HGet(ctx, KEY, agendaItemId).Result()
	if err != nil {
		panic(err)
	}
	var agendaItem AgendaItem
	err = json.Unmarshal([]byte(agendaItemById), &agendaItem)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}
	log.Printf("Agenda Item retrieved from Database: %s", agendaItem)
	agendaItem.Archived = true

	err = s.RedisClient.HSet(ctx, KEY, agendaItem.Id, agendaItem).Err()
	if err != nil {
		panic(err)
	}

	log.Printf("Agenda Item Archived in Database: %s", agendaItem)

	agendaItemJson, err := json.Marshal(agendaItem)
	if err != nil {
		log.Printf("An error occured while marshalling the agenda item to json: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("agenda-item-archived-%s", agendaItem.Id)),
		Value: agendaItemJson,
	}
	err = s.KafkaWriter.WriteMessages(r.Context(), msg)

	if err != nil {
		log.Printf("An error occured while writing the message to Kafka: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("Agenda Item Archived Event emitted to Kafka: %s", agendaItem)

	respondWithJSON(w, http.StatusOK, agendaItem)
}

// NewKafkaWriter creates a new Kafka writer.
func NewKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	kafkaAlive := isKafkaAlive(KAFKA_URL, KAFKA_TOPIC)
	if !kafkaAlive {
		log.Printf("Cannot connect to Kafka, restarting until it is healthy.")
		panic("Cannot connect to Kafka")
	}
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// NewRedisClient creates a new Redis client.
func NewRedisClient(redisHost, redisPort, redisPass string) *redis.Client {
	defaultDB := 0
	return redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPass,
		DB:       defaultDB,
	})
}

// NewServer creates a new server.
func NewServer() ServerInterface {
	return server{
		KafkaWriter: NewKafkaWriter(KAFKA_URL, KAFKA_TOPIC),
		RedisClient: NewRedisClient(REDIS_HOST, REDIS_PORT, REDIS_PASSWORD),
	}
}

// OpenAPIHandler returns a handler that serves the OpenAPI spec as JSON.
func OpenAPI(r *chi.Mux) {
	dir := http.Dir("docs")
	fs := http.FileServer(dir)
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", fs))
}

// isKafkaAlive checks if Kafka is alive.
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
