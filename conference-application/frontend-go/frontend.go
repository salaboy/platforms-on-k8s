package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func newAgendaServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("http://localhost:8081", w, r)
}

func proxyRequest(serviceUrl string, w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Starting Frontend Go in Port: %s", appPort)

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/agenda/", newAgendaServiceHandler)

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
