package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"/home/robi/Desktop/DEVPOST/backend/framework"
	"/home/robi/Desktop/DEVPOST/backend/modules"
)

var (
	logMutex sync.Mutex
	logs     []map[string]string
)

func main() {
	fmt.Println("Node Correctness Verification Backend")

	// Initialize the modular framework
	if err := initializeFramework(); err != nil {
		log.Fatalf("Failed to initialize framework: %v", err)
	}

	// Start HTTP server for logs
	http.HandleFunc("/logs", serveLogs)
	fmt.Println("Starting server on :4000")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeFramework() error {
	// Placeholder for initializing the modular framework
	fmt.Println("Initializing modular framework...")

	// Initialize Polkadot module
	polkadotModule := &modules.PolkadotModule{}

	// Create a new challenger with the Polkadot module
	challenger := framework.NewChallenger(polkadotModule)

	// Run challenges
	challenger.RunChallenges()

	return nil
}

func logResult(timestamp, message string) {
	logMutex.Lock()
	defer logMutex.Unlock()
	logs = append(logs, map[string]string{
		"timestamp": timestamp,
		"message":   message,
	})
}

func serveLogs(w http.ResponseWriter, r *http.Request) {
	logMutex.Lock()
	defer logMutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}