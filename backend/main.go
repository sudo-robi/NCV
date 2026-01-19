package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"ncv/backend/framework"
	"ncv/backend/modules"
)

// LogEntry matches the frontend schema
// proof_type, success, message, evidence (optional), timestamp (seconds)
type LogEntry struct {
	ProofType string      `json:"proof_type"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Evidence  interface{} `json:"evidence,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

var (
	logMutex sync.Mutex
	logs     []LogEntry
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
	fmt.Println("Initializing modular framework...")

	polkadotModule := &modules.PolkadotModule{}
	challenger := framework.NewChallenger(polkadotModule)

	// Also run via the framework to exercise that path (prints to stdout)
	challenger.RunChallenges()

	// Run challenges and append placeholder logs that align with frontend schema
	start := time.Now()
	if err := polkadotModule.HeadCorrectness(); err != nil {
		appendLog("Head Correctness", false, fmt.Sprintf("Head Correctness failed: %v", err), nil)
	} else {
		appendLog("Head Correctness", true, "Head Correctness check passed.", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	}

	start = time.Now()
	if err := polkadotModule.ExecutionCorrectness(); err != nil {
		appendLog("Execution Correctness", false, fmt.Sprintf("Execution Correctness failed: %v", err), nil)
	} else {
		appendLog("Execution Correctness", true, "Execution Correctness check passed.", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	}

	start = time.Now()
	if err := polkadotModule.FreshnessProof(); err != nil {
		appendLog("Freshness", false, fmt.Sprintf("Freshness Proof failed: %v", err), nil)
	} else {
		appendLog("Freshness", true, "Freshness Proof check passed.", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	}

	return nil
}

func appendLog(proofType string, success bool, message string, evidence interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	logs = append(logs, LogEntry{
		ProofType: proofType,
		Success:   success,
		Message:   message,
		Evidence:  evidence,
		Timestamp: time.Now().Unix(),
	})
}

func serveLogs(w http.ResponseWriter, r *http.Request) {
	// Basic CORS to allow CRA dev server
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	logMutex.Lock()
	defer logMutex.Unlock()
	_ = json.NewEncoder(w).Encode(logs)
}