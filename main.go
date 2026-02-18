package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	Name string `json:"server_name"`
	IP string `json:"ip_address"`
	Status bool `json:"is_online"`
}

type Switch struct {
	Hostname string
	PortCount int
	IsManaged bool
}

type PromptGenerator interface {
	GeneratePrompt() string
}

type Document struct {
	Text string `json:"text"`
}

var (	
		Servers = []Server{{Name: "dc1", IP: "192.168.1.2", Status: true}, {Name: "fs1", IP: "192.168.1.3", Status: true}, {Name: "web1", IP: "192.168.1.3", Status: false}}
 		ServersMutex sync.Mutex
		Switches = []Switch{}
		SwitchesMutex sync.Mutex
		Documents = []Document{}
		DocumentsMutex sync.Mutex
	)

func (s Server) GeneratePrompt() string {
	var operationStatus string
	if s.Status {
		operationStatus = "Operating Normally"
		} else { 
			operationStatus = "Maintenance Needed"
		}
		result := fmt.Sprintf("Server %s: %s\n", s.Name, operationStatus)
		return result
	}

func (s Switch) GeneratePrompt() string {
	result := fmt.Sprintf("Switch %s has %d ports.\nManaged: %t", s.Hostname, s.PortCount, s.IsManaged)
	return result
}
	
func healthCheck (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"status": "up", "system": "healthy"}
	json.NewEncoder(w).Encode(data)
}

func handleServers (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		ServersMutex.Lock()
		defer ServersMutex.Unlock()
		json.NewEncoder(w).Encode(Servers)
		return
	}

	if r.Method == http.MethodPost {
		var newServer Server
		err := json.NewDecoder(r.Body).Decode(&newServer)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		
		ServersMutex.Lock()
		Servers = append(Servers, newServer)
		ServersMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newServer)
		fmt.Println("New server added successfully")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func handleDocuments (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(Documents)
		return
	}

	if r.Method == http.MethodPost {
		var newDocument Document
		err := json.NewDecoder(r.Body).Decode(&newDocument)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		Documents = append(Documents, newDocument)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newDocument)
		fmt.Println("New document added successfully")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}


func main() {

	fmt.Println("Starting Netscribe server...")
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/servers", handleServers)
	http.HandleFunc("/documents", handleDocuments)
	http.ListenAndServe(":8080", nil)
}