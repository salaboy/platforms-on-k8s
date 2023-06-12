package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
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

var rdb *redis.Client
var KEY = "AGENDAITEMS"

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
				log.Fatalln("There was an error decoding the AgendaItem into the struct")
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
			log.Fatalln("There was an error decoding the AgendaItem into the struct")
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
		log.Fatalln("There was an error decoding the request body into the struct")
	}
	log.Printf("Agenda Item retrieved from Database: %s", agendaItem)
	respondWithJSON(w, http.StatusOK, agendaItem)
}

func newAgendaItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var agendaItem AgendaItem
	err := json.NewDecoder(r.Body).Decode(&agendaItem)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	// @TODO: write fail scenario (check for fail string in title return 500)

	agendaItem.Id = uuid.New().String()

	err = rdb.HSetNX(ctx, KEY, agendaItem.Id, agendaItem).Err()
	if err != nil {
		panic(err)
	}

	log.Printf("Agenda Item Stored in Database: %s", agendaItem)

	// @TODO avoid doing two marshals to json
	respondWithJSON(w, http.StatusOK, agendaItem)
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

	log.Printf("Starting Agenda Service in Port: %s", appPort)

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Printf("Connected to Redis.")

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/", newAgendaItemHandler).Methods("POST")
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
