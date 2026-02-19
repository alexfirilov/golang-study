package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type (
	Server struct {
		Name   string `json:"server_name"`
		IP     string `json:"ip_address"`
		Status bool   `json:"is_online"`
	}

	Document struct {
		Text string `json:"text"`
	}

	API struct {
		Servers        []Server `json:"servers"`
		ServersMutex   sync.Mutex
		Documents      []Document `json:"documents"`
		DocumentsMutex sync.Mutex
		AIUpdates      chan string
		Wg             sync.WaitGroup
	}
)

func (api *API) InfrastructureWorker() {
	for {
		processedPrompt := <-api.AIUpdates
		fmt.Printf("Worker received AI instruction: [%s]\n", processedPrompt)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"status": "up", "system": "healthy"}
	json.NewEncoder(w).Encode(data)
}

func (api *API) HandleServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		api.ServersMutex.Lock()
		defer api.ServersMutex.Unlock()
		json.NewEncoder(w).Encode(api.Servers)
		return
	}

	if r.Method == http.MethodPost {
		var newServer Server
		err := json.NewDecoder(r.Body).Decode(&newServer)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		api.ServersMutex.Lock()
		api.Servers = append(api.Servers, newServer)
		api.ServersMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newServer)
		fmt.Println("New server added successfully")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (api *API) HandleDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		api.DocumentsMutex.Lock()
		defer api.DocumentsMutex.Unlock()
		json.NewEncoder(w).Encode(api.Documents)
		return
	}

	if r.Method == http.MethodPost {
		var newDocument Document
		err := json.NewDecoder(r.Body).Decode(&newDocument)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		api.DocumentsMutex.Lock()
		api.Documents = append(api.Documents, newDocument)
		api.DocumentsMutex.Unlock()
		api.Wg.Add(1)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newDocument)
		fmt.Println("New document added successfully")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
