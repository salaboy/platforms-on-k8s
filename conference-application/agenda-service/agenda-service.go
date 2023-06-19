package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	kafka "github.com/segmentio/kafka-go"
	"golang.org/x/exp/slices"
)

type Proposal struct {
	Id string
}

type AgendaItem struct {
	Id       string
	Proposal Proposal
	Title    string
	Author   string
	Day      string
	Time     string
}

type ServiceInfo struct {
	Name         string
	Version      string
	Source       string
	PodId        string
	PodNamespace string
	PodNodeName  string
}

var rdb *redis.Client
var KEY = "AGENDAITEMS"
var VERSION = getEnv("VERSION", "1.0.0")
var SOURCE = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/agenda-service")
var POD_ID = getEnv("POD_ID", "N/A")
var POD_NAMESPACE = getEnv("POD_NAMESPACE", "N/A")
var POD_NODENAME = getEnv("POD_NODENAME", "N/A")
var REDIS_HOST = getEnv("REDIS_HOST", "localhost")
var REDIS_PORT = getEnv("REDIS_PORT", "6379")
var REDIS_PASSOWRD = getEnv("REDIS_PASSWORD", "")
var KAFKA_URL = getEnv("KAFKA_URL", "localhost:9094")
var KAFKA_TOPIC = getEnv("KAFKA_TOPIC", "events-topic")

func getAgendaByDayHandler(w http.ResponseWriter, r *http.Request) {

}

func getHighlightsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	agendaItemsHashs, err := rdb.HGetAll(ctx, KEY).Result()
	if err != nil {
		panic(err)
	}

	higlights := 4
	min := 0
	max := len(agendaItemsHashs)

	var chosenOnes []int
	counter := 0
	for {
		if len(chosenOnes) == higlights {
			break
		}
		random := rand.Intn(max-min) + min
		if !slices.Contains(chosenOnes, random) {
			chosenOnes = append(chosenOnes, random)
		}

	}
	log.Printf("Chosen ones: %d", chosenOnes)

	counter = 0
	var agendaItems []AgendaItem
	for _, ai := range agendaItemsHashs {
		if slices.Contains(chosenOnes, counter) {
			var agendaItem AgendaItem
			err = json.Unmarshal([]byte(ai), &agendaItem)
			if err != nil {
				log.Printf("There was an error decoding the AgendaItem into the struct: %v", err)
			}
			agendaItems = append(agendaItems, agendaItem)
		}
		counter++
	}

	respondWithJSON(w, http.StatusOK, agendaItems)

}

func getAllAgendaItemsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	agendaItemsHashs, err := rdb.HGetAll(ctx, KEY).Result()
	if err != nil {
		panic(err)
	}
	var agendaItems []AgendaItem

	for _, ai := range agendaItemsHashs {
		var agendaItem AgendaItem
		err = json.Unmarshal([]byte(ai), &agendaItem)
		if err != nil {
			log.Printf("There was an error decoding the AgendaItem into the struct: %v", err)
		}
		agendaItems = append(agendaItems, agendaItem)
	}
	log.Printf("Agenda Items retrieved from Database: %d", len(agendaItems))
	respondWithJSON(w, http.StatusOK, agendaItems)

}

func getAgendaItemByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	agendaItemId := mux.Vars(r)["id"]
	log.Printf("Fetching Agenda Item By Id: %s", agendaItemId)
	agendaItemById, err := rdb.HGet(ctx, KEY, agendaItemId).Result()
	if err != nil {
		panic(err)
	}
	var agendaItem AgendaItem
	err = json.Unmarshal([]byte(agendaItemById), &agendaItem)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
	}
	log.Printf("Agenda Item retrieved from Database: %s", agendaItem)
	respondWithJSON(w, http.StatusOK, agendaItem)
}

func newAgendaItemHandler(kafkaWriter *kafka.Writer) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		var agendaItem AgendaItem
		err := json.NewDecoder(r.Body).Decode(&agendaItem)
		if err != nil {
			log.Printf("There was an error decoding the request body into the struct: %v", err)
		}

		// @TODO: write fail scenario (check for fail string in title return 500)

		agendaItem.Id = uuid.New().String()

		err = rdb.HSetNX(ctx, KEY, agendaItem.Id, agendaItem).Err()
		if err != nil {
			panic(err)
		}

		log.Printf("Agenda Item Stored in Database: %s", agendaItem)

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
		err = kafkaWriter.WriteMessages(r.Context(), msg)

		if err != nil {
			log.Printf("An error occured while writing the message to Kafka: %v", err)
			respondWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		log.Printf("New Agenda Item Event emitted to Kafka: %s", agendaItem)

		// @TODO avoid doing two marshals to json
		respondWithJSON(w, http.StatusOK, agendaItem)
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

	log.Printf("Starting Agenda Service in Port: %s", appPort)

	rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: REDIS_PASSOWRD, // no password set
		DB:       0,              // use default DB
	})

	log.Printf("Connected to Redis.")

	log.Printf("Connecting to Kafka Instance: %s, topic: %s.", KAFKA_URL, KAFKA_TOPIC)
	//https://github.com/segmentio/kafka-go/blob/main/examples/producer-api/main.go
	kafkaWriter := getKafkaWriter(KAFKA_URL, KAFKA_TOPIC)

	log.Printf("Connected to Kafka.")
	defer kafkaWriter.Close()

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/", newAgendaItemHandler(kafkaWriter)).Methods("POST")
	r.HandleFunc("/", getAllAgendaItemsHandler).Methods("GET")
	r.HandleFunc("/highlights", getHighlightsHandler).Methods("GET")
	r.HandleFunc("/{id}", getAgendaItemByIdHandler).Methods("GET")
	r.HandleFunc("/day/{day}", getAgendaByDayHandler).Methods("GET")
	// r.HandleFunc("/{id}", deleteAgendaItemHandler).Methods("DELETE")
	// r.HandleFunc("/", deleteAllHandler).Methods("DELETE")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
		var info ServiceInfo = ServiceInfo{
			Name:         "AGENDA",
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

func (p Proposal) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p Proposal) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	return nil
}

func (ai AgendaItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(ai)
}

func (ai AgendaItem) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &ai); err != nil {
		return err
	}

	return nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
