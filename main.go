package main

import "fmt"
// import "encoding/json"
import "net/http"

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
	w.Write([]byte("Netscribe System: Online"))
}

func main() {
	http.HandleFunc("/health", healthCheck)
	fmt.Println("Starting Netscribe server...")
	http.ListenAndServe(":8080", nil)
}