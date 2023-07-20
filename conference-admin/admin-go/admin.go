package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/api/types/v1alpha1"
	clientV1alpha1 "github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/clientset/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	VERSION             = getEnv("VERSION", "1.0.0")
	SOURCE              = getEnv("SOURCE", "https://github.com/salaboy/platforms-on-k8s/tree/main/conference-admin/admin-go")
	POD_NAME            = getEnv("POD_NAME", "N/A")
	POD_NAMESPACE       = getEnv("POD_NAMESPACE", "N/A")
	POD_NODENAME        = getEnv("POD_NODENAME", "N/A")
	POD_IP              = getEnv("POD_IP", "N/A")
	POD_SERVICE_ACCOUNT = getEnv("POD_SERVICE_ACCOUNT", "N/A")
	kubeconfig          string
)

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type EnvironmentSimple struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	InstallInfra bool   `json:"installInfra"`
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
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

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
		panic(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/", getAllEnvironmentsHandler(clientSet)).Methods("GET")
	r.HandleFunc("/api/", createEnvironmentHandler(clientSet)).Methods("POST")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/service/info", func(w http.ResponseWriter, r *http.Request) {
		var info ServiceInfo = ServiceInfo{
			Name:              "ADMIN",
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

	// Start the server; this is a blocking call
	err = http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

func getAllEnvironmentsHandler(clientSet *clientV1alpha1.ConferenceAdminV1Alpha1Client) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		environments, err := clientSet.Environments("default").List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		fmt.Printf("environments found: %+v\n", environments)

		log.Printf("Environments retrieved from Kube API: %d", len(environments.Items))
		respondWithJSON(w, http.StatusOK, environments.Items)
	})

}

func createEnvironmentHandler(clientSet *clientV1alpha1.ConferenceAdminV1Alpha1Client) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				CompositionSelector: v1alpha1.CompositionSelector{
					MatchLabels: map[string]string{
						"type": env.Type,
					},
				},
				Parameters: v1alpha1.Parameters{
					InstallInfra: env.InstallInfra,
				},
			},
		}
		result, err := clientSet.Environments("default").Create(fullEnv)
		if err != nil {
			panic(err)
		}

		respondWithJSON(w, http.StatusOK, result)
	})

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
