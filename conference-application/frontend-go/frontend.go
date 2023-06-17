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

var VERSION = getEnv("VERSION", "1.0.0")
var SOURCE = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-application/frontend-go")
var POD_ID = getEnv("POD_ID", "N/A")
var POD_NAMESPACE = getEnv("POD_NAMESPACE", "N/A")
var POD_NODENAME = getEnv("POD_NODENAME", "N/A")
var AGENDA_SERVICE_URL = getEnv("AGENDA_SERVICE_URL", "http://agenda-service.default.svc.cluster.local")
var C4P_SERVICE_URL = getEnv("C4P_SERVICE_URL", "http://c4p-service.default.svc.cluster.local")
var NOTIFICATION_SERVICE_URL = getEnv("NOTIFICATION_SERVICE_URL", "http://notifications-service.default.svc.cluster.local")

type ServiceInfo struct {
	Name         string
	Version      string
	Source       string
	PodId        string
	PodNamespace string
	PodNodeName  string
}

func agendaServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("api/agenda", AGENDA_SERVICE_URL, w, r)
}

func c4PServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("api/c4p", C4P_SERVICE_URL, w, r)
}

func notificationServiceHandler(w http.ResponseWriter, r *http.Request) {
	proxyRequest("api/notifications", NOTIFICATION_SERVICE_URL, w, r)
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

	log.Printf("Starting Frontend Go in Port: %s", appPort)

	r := mux.NewRouter()

	r.PathPrefix("/api/agenda/").HandlerFunc(agendaServiceHandler)
	r.PathPrefix("/api/c4p/").HandlerFunc(c4PServiceHandler)
	// r.HandleFunc("/agenda/", agendaServiceHandler)
	// r.HandleFunc("/c4p/", c4PServiceHandler)
	// r.HandleFunc("/notifications/", notificationServiceHandler)

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
		var info ServiceInfo = ServiceInfo{
			Name:         "FRONTEND",
			Version:      VERSION,
			Source:       SOURCE,
			PodId:        POD_ID,
			PodNamespace: POD_NODENAME,
		}
		json.NewEncoder(w).Encode(info)
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))

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
