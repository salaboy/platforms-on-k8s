package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/api"
	"log"
	"net/http"
	"os"

	"github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/api/types/v1alpha1"
	clientV1alpha1 "github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/clientset/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Version           = getEnv("VERSION", "1.0.0")
	Source            = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-admin/admin-go")
	PodName           = getEnv("POD_NAME", "N/A")
	PodNamespace      = getEnv("POD_NAMESPACE", "N/A")
	PodNodeName       = getEnv("POD_NODENAME", "N/A")
	PodIp             = getEnv("POD_IP", "N/A")
	PodServiceAccount = getEnv("POD_SERVICE_ACCOUNT", "N/A")
	AppPort           = getEnv("APP_PORT", "8080")
	KoDataPath        = getEnv("KO_DATA_PATH", "kodata")
	kubeconfig        string
)

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Frontend struct {
	Debug bool `json:"debug"`
}

type Parameters struct {
	Type         string   `json:"type"`
	InstallInfra bool     `json:"installInfra"`
	Frontend     Frontend `json:"frontend"`
}

type EnvironmentSimple struct {
	Name       string     `json:"name"`
	Parameters Parameters `json:"parameters"`
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

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	r := NewChiServer()

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+AppPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// NewChiServer returns a new chi router with the API routes configured.
func NewChiServer() *chi.Mux {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(any(err))
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config)

	if err != nil {
		panic(any(err))
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	fs := http.FileServer(http.Dir(KoDataPath))

	server := NewServer(clientSet)

	OpenAPI(r)

	r.Mount("/api/", api.Handler(server))
	r.Handle("/*", http.StripPrefix("/", fs))

	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	log.Printf("Starting Conference Admin API in port: %s", AppPort)

	return r
}

// server represents the Conference Admin API server.
type server struct {
	ClientSet *clientV1alpha1.ConferenceAdminV1Alpha1Client
}

// ListEnvironments returns a list of Environments.
func (s *server) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	environments, err := s.ClientSet.Environments("default").List(metav1.ListOptions{})
	if err != nil {
		panic(any(err))
	}

	fmt.Printf("environments found: %+v\n", environments)

	log.Printf("Environments retrieved from Kube API: %d", len(environments.Items))
	respondWithJSON(w, http.StatusOK, environments.Items)
}

// CreateEnvironment creates a new Environment.
func (s *server) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var env EnvironmentSimple
	err := json.NewDecoder(r.Body).Decode(&env)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	fullEnv := &v1alpha1.Environment{
		ObjectMeta: metav1.ObjectMeta{
			Name: env.Name,
		},
		Spec: v1alpha1.EnvironmentSpec{
			WriteConnectionSecretToRef: v1alpha1.WriteConnectionSecretToRef{
				Name: env.Name,
			},
			CompositionSelector: v1alpha1.CompositionSelector{
				MatchLabels: map[string]string{
					"type": env.Parameters.Type,
				},
			},
			Parameters: v1alpha1.Parameters{
				InstallInfra: env.Parameters.InstallInfra,
				Frontend: v1alpha1.Frontend{
					Debug: env.Parameters.Frontend.Debug,
				},
			},
		},
	}
	result, err := s.ClientSet.Environments("default").Create(fullEnv)
	if err != nil {
		panic(any(err))
	}

	respondWithJSON(w, http.StatusOK, result)
}

// DeleteEnvironment deletes an environment.
func (s *server) DeleteEnvironment(w http.ResponseWriter, r *http.Request, id string) {
	err := s.ClientSet.Environments("default").Delete(id, metav1.DeleteOptions{})
	if err != nil {
		panic(any(err))
	}
}

// GetServiceInfo returns information about the service.
func (s *server) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	var info = ServiceInfo{
		Name:              "ADMIN",
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

// NewServer returns an instance of the api.ServerInterface implementation.
func NewServer(client *clientV1alpha1.ConferenceAdminV1Alpha1Client) api.ServerInterface {
	return &server{
		ClientSet: client,
	}
}

// OpenAPI OpenAPIHandler returns a handler that serves the OpenAPI documentation.
func OpenAPI(r *chi.Mux) {
	fs := http.FileServer(http.Dir(KoDataPath + "/docs/"))
	r.Handle("/openapi/*", http.StripPrefix("/openapi/", fs))
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
