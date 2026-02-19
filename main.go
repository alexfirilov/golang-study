package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"context"
	"os"
	"os/signal"
	"syscall"
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

var AIUpdates = make(chan string)

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

func ProcessPrompt(prompt string) {
	fmt.Printf("AI agent received prompt: [%s]\n", prompt)
	time.Sleep(5 * time.Second)
	AIUpdates <- "COMMAND_PARSE_SUCCESS: " + prompt
}

func InfrastructureWorker() {
	for {
		processedPrompt := <-AIUpdates 
        fmt.Printf("Worker received AI instruction: [%s]\n", processedPrompt)
	}
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
		DocumentsMutex.Lock()
		defer DocumentsMutex.Unlock()
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

		DocumentsMutex.Lock()
		Documents = append(Documents, newDocument)
		DocumentsMutex.Unlock()
		go ProcessPrompt(newDocument.Text)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newDocument)
		fmt.Println("New document added successfully")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}


func main() {

	fmt.Println("Starting Netscribe server...")
	srv := &http.Server {
		Addr: ":8080",
		Handler: nil,
	}
	
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/servers", handleServers)
	http.HandleFunc("/documents", handleDocuments)
	go InfrastructureWorker()

	go func() { srv.ListenAndServe() }()
	
	fmt.Println("Server started successfully!")
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutdown signal received. Gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}