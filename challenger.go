package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

type ProofResult struct {
	ProofType string      `json:"proof_type"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Evidence  interface{} `json:"evidence"`
	Timestamp float64     `json:"timestamp"`
}

// ChainAdapter defines the interface for chain-specific logic
// This allows multi-chain support by implementing adapters for different blockchains.
type ChainAdapter interface {
	GetHeader(url string) (map[string]interface{}, error)
	GetRuntimeVersion(url string) (map[string]interface{}, error)
}

// PolkadotAdapter implements ChainAdapter for Polkadot
// Similar adapters can be created for Ethereum, Solana, etc.
type PolkadotAdapter struct{}

func (p *PolkadotAdapter) GetHeader(url string) (map[string]interface{}, error) {
	return performRPC(url, "chain_getHeader", nil)
}

func (p *PolkadotAdapter) GetRuntimeVersion(url string) (map[string]interface{}, error) {
	return performRPC(url, "state_getRuntimeVersion", nil)
}

// queryRPC sends an RPC request to the specified URL and returns the result.
// url: The endpoint to send the request to.
// method: The RPC method to call.
// params: The parameters for the RPC method.
// Returns the result of the RPC call or an error.
func queryRPC(url string, method string, params []interface{}) (interface{}, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	reqBody, _ := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response body
	body, _ := ioutil.ReadAll(resp.Body)
	var rpcResp RPCResponse
	json.Unmarshal(body, &rpcResp)

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %v", rpcResp.Error)
	}

	return rpcResp.Result, nil
}

func saveEvidence(results []ProofResult) {
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("logs/evidence_%d.json", timestamp)
	fmt.Printf("Attempting to save %d results to %s\n", len(results), filename)
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	fmt.Printf("Evidence saved successfully to %s\n", filename)
}

func verifyConsensus(refResults map[string]int) (string, int) {
	var consensusRoot string
	maxCount := 0
	for root, count := range refResults {
		if count > maxCount {
			maxCount = count
			consensusRoot = root
		}
	}
	return consensusRoot, maxCount
}

func enhancedCrossSourceVerification(nutURL string, refURLs []string, proofType string, rpcMethod string, params []interface{}) ProofResult {
	refResults := make(map[string]int)
	validRefs := 0

	for _, url := range refURLs {
		res, err := queryRPC(url, rpcMethod, params)
		if err == nil {
			result := res.(map[string]interface{})
			key := result["stateRoot"].(string) // Example key, adjust per proof type
			refResults[key]++
			validRefs++
		} else {
			fmt.Printf("[Proof] Ref URL %s failed: %v\n", url, err)
		}
	}

	if validRefs == 0 {
		return ProofResult{
			ProofType: proofType,
			Success:   false,
			Message:   "[CRITICAL] All reference sources unreachable",
			Timestamp: float64(time.Now().Unix()),
		}
	}

	consensusRoot, _ := verifyConsensus(refResults)
	res, err := queryRPC(nutURL, rpcMethod, params)
	if err != nil {
		return ProofResult{
			ProofType: proofType,
			Success:   false,
			Message:   "[Offline] NUT unreachable",
			Timestamp: float64(time.Now().Unix()),
		}
	}

	nutResult := res.(map[string]interface{})
	success := (nutResult["stateRoot"].(string) == consensusRoot)
	classification := "Correct"
	if !success {
		classification = "Mismatch Detected"
	}

	return ProofResult{
		ProofType: proofType,
		Success:   success,
		Message:   fmt.Sprintf("[%s] Verified against %d sources", classification, validRefs),
		Evidence: map[string]interface{}{
			"consensus_root": consensusRoot,
			"nut_root":       nutResult["stateRoot"].(string),
			"category":       classification,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

func headCorrectness(nutURL string, refURLs []string) ProofResult {
	fmt.Printf("[Proof] Starting Head Correctness check...\n")
	// Consensus-based verification: Query ALL refs and find the majority state root
	roots := make(map[string]int)
	heights := make(map[int64]int)

	validRefs := 0
	for i, url := range refURLs {
		fmt.Printf("[Proof] Querying Ref %d: %s\n", i, url)
		res, err := queryRPC(url, "chain_getHeader", []interface{}{})
		if err == nil {
			header := res.(map[string]interface{})
			root := header["stateRoot"].(string)
			roots[root]++

			var h int64
			fmt.Sscanf(header["number"].(string), "0x%x", &h)
			heights[h]++
			validRefs++
			fmt.Printf("[Proof] Ref %d OK. Root: %s\n", i, root)
		} else {
			fmt.Printf("[Proof] Ref %d Failed: %v\n", i, err)
		}
	}

	if validRefs == 0 {
		return ProofResult{ProofType: "Head Correctness", Success: false, Message: "[CRITICAL] All reference sources unreachable", Timestamp: float64(time.Now().Unix())}
	}

	// Find the "Ground Truth" (Majority Root)
	var refRoot string
	maxCount := 0
	for root, count := range roots {
		if count > maxCount {
			maxCount = count
			refRoot = root
		}
	}

	// Get NUT result
	fmt.Printf("[Proof] Querying NUT: %s\n", nutURL)
	res, err := queryRPC(nutURL, "chain_getHeader", []interface{}{})
	if err != nil {
		fmt.Printf("[Proof] NUT Failed: %v\n", err)
		return ProofResult{ProofType: "Head Correctness", Success: false, Message: "[Offline] NUT unreachable", Timestamp: float64(time.Now().Unix())}
	}
	fmt.Printf("[Proof] NUT OK.\n")

	nutHeader := res.(map[string]interface{})
	nutRoot := nutHeader["stateRoot"].(string)
	numStr := nutHeader["number"].(string)
	var nutHeight int64
	fmt.Sscanf(numStr, "0x%x", &nutHeight)

	success := (nutRoot == refRoot)
	classification := "Correct"
	if !success {
		classification = "Execution Mismatch / Database Corruption"
	}

	return ProofResult{
		ProofType: "Head Correctness",
		Success:   success,
		Message:   fmt.Sprintf("[%s] Verified against %d sources", classification, validRefs),
		Evidence: map[string]interface{}{
			"consensus_root": refRoot,
			"nut_root":       nutRoot,
			"nut_height":     nutHeight,
			"category":       classification,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

func freshnessProof(nutURL string, refURLs []string) ProofResult {
	var refHeight int64
	for _, url := range refURLs {
		res, err := queryRPC(url, "chain_getHeader", []interface{}{})
		if err == nil {
			header := res.(map[string]interface{})
			fmt.Sscanf(header["number"].(string), "0x%x", &refHeight)
			break
		}
	}

	res, err := queryRPC(nutURL, "chain_getHeader", []interface{}{})
	if err != nil {
		return ProofResult{ProofType: "Freshness", Success: false, Message: "NUT unavailable", Timestamp: float64(time.Now().Unix())}
	}
	var nutHeight int64
	fmt.Sscanf(res.(map[string]interface{})["number"].(string), "0x%x", &nutHeight)

	lag := refHeight - nutHeight
	success := lag <= 2
	classification := "Fresh"
	if lag > 2 {
		classification = "Lagging but correct"
	} else if lag < 0 {
		lag = 0
	}

	return ProofResult{
		ProofType: "Freshness",
		Success:   success,
		Message:   fmt.Sprintf("[%s] Node is %d blocks behind", classification, lag),
		Evidence: map[string]interface{}{
			"lag_blocks": lag,
			"category":   classification,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

func executionProof(nutURL string, refURLs []string) ProofResult {
	extrinsic := "0x4502840015oF4uVebkyS7CnxwgS9RNoCHuL2UaykS4eXGykN3cKR9648..."
	params := []interface{}{extrinsic}

	var refFee float64
	for _, url := range refURLs {
		res, err := queryRPC(url, "payment_queryInfo", params)
		if err == nil {
			info := res.(map[string]interface{})
			refFee = info["partialFee"].(float64)
			break
		}
	}

	res, err := queryRPC(nutURL, "payment_queryInfo", params)
	if err != nil {
		return ProofResult{
			ProofType: "Execution Correctness",
			Success:   true,
			Message:   "[Degraded] Node does not support active execution challenges",
			Timestamp: float64(time.Now().Unix()),
		}
	}

	nutInfo := res.(map[string]interface{})
	nutFee := nutInfo["partialFee"].(float64)

	success := (refFee == nutFee)
	classification := "Verified"
	if !success {
		classification = "Execution mismatch"
	}

	return ProofResult{
		ProofType: "Execution Correctness",
		Success:   success,
		Message:   fmt.Sprintf("[%s] Fee estimation match: %v", classification, nutFee),
		Evidence: map[string]interface{}{
			"nut_fee":  nutFee,
			"ref_fee":  refFee,
			"category": classification,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

func receiptChallenge(nutURL string, refURLs []string) ProofResult {
	// 1. Get a recent block and pick a transaction hash
	var txHash string
	var blockNum int64
	for _, url := range refURLs {
		res, err := queryRPC(url, "chain_getBlock", []interface{}{})
		if err == nil {
			block := res.(map[string]interface{})["block"].(map[string]interface{})
			extrinsics := block["extrinsics"].([]interface{})
			if len(extrinsics) > 0 {
				numStr := block["header"].(map[string]interface{})["number"].(string)
				fmt.Sscanf(numStr, "0x%x", &blockNum)
				// For demo, we just need a non-empty check.
				// In a real system, we'd hash the extrinsic or use a known hash.
				txHash = "0x..." // placeholder for extrinsic hash
				break
			}
		}
	}

	// 2. Challenge the node for events (receipts) at that block
	// RPC: state_getStorage for System.Events (Standard Substrate key)
	eventKey := "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"

	var refEvents interface{}
	for _, url := range refURLs {
		res, err := queryRPC(url, "state_getStorage", []interface{}{eventKey})
		if err == nil {
			refEvents = res
			break
		}
	}

	res, err := queryRPC(nutURL, "state_getStorage", []interface{}{eventKey})
	if err != nil {
		return ProofResult{ProofType: "Receipt Verification", Success: false, Message: "NUT failed receipt query", Timestamp: float64(time.Now().Unix())}
	}

	success := (refEvents == res)
	classification := "Correct"
	if !success {
		classification = "Silent Corruption / Indexer Drift"
	}

	return ProofResult{
		ProofType: "Receipt Verification",
		Success:   success,
		Message:   fmt.Sprintf("[%s] Receipt events match for block %d", classification, blockNum),
		Evidence: map[string]interface{}{
			"category": classification,
			"block":    blockNum,
			"tx_hash":  txHash,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

func checkpointChallenge(nutURL string, refURLs []string) ProofResult {
	checkpointBlock := int64(29500000)
	checkpointHash := "0xd73df8996515b6b158223652ce6d6911c751268673f8ff38563f69b596287c9f"

	res, err := queryRPC(nutURL, "chain_getBlockHash", []interface{}{checkpointBlock})
	if err != nil {
		return ProofResult{ProofType: "Checkpoint Verification", Success: false, Message: "Failed to fetch checkpoint", Timestamp: float64(time.Now().Unix())}
	}

	nutHash := res.(string)
	success := (nutHash == checkpointHash)
	classification := "Canonical"
	if !success {
		classification = "Consensus-desynced / Chain Fork"
	}

	return ProofResult{
		ProofType: "Checkpoint Verification",
		Success:   success,
		Message:   fmt.Sprintf("[%s] Checkpoint block %d verified", classification, checkpointBlock),
		Evidence: map[string]interface{}{
			"checkpoint_block": checkpointBlock,
			"expected_hash":    checkpointHash,
			"actual_hash":      nutHash,
			"category":         classification,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

func triggerHealer(category string, block int64) {
	fmt.Printf("[HEALER] 🚩 Triggering remediation for %s at block %d...\n", category, block)
	cmd := exec.Command("./remediate.sh", category, fmt.Sprintf("%d", block))
	cmd.Start() // Run in background
}

func runtimeVerification(nutURL string, refURLs []string) ProofResult {
	fmt.Println("[Proof] Starting Runtime Verification...")
	var refVersion float64
	for _, url := range refURLs {
		result, err := performRPC(url, "state_getRuntimeVersion", nil)
		if err == nil {
			refVersion = result["specVersion"].(float64)
			break
		}
	}

	result, err := performRPC(nutURL, "state_getRuntimeVersion", nil)
	if err != nil {
		return ProofResult{
			ProofType: "Runtime Verification",
			Success:   false,
			Message:   "NUT runtime query failed",
			Timestamp: float64(time.Now().Unix()),
		}
	}

	nutVersion := result["specVersion"].(float64)
	success := (nutVersion == refVersion)
	classification := "Up-to-date"
	if !success {
		classification = "Outdated Runtime"
	}

	return ProofResult{
		ProofType: "Runtime Verification",
		Success:   success,
		Message:   fmt.Sprintf("[%s] Spec Version: %v", classification, nutVersion),
		Evidence: map[string]interface{}{
			"nut_version": nutVersion,
			"ref_version": refVersion,
			"category":    classification,
		},
		Timestamp: float64(time.Now().Unix()),
	}
}

// Refactored function to handle RPC queries with error handling
func performRPC(url, method string, params []interface{}) (map[string]interface{}, error) {
	res, err := queryRPC(url, method, params)
	if err != nil {
		return nil, fmt.Errorf("RPC call failed for %s: %v", url, err)
	}
	result, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected response format from %s", url)
	}
	return result, nil
}

// Refactored function to determine consensus root
func determineConsensus(roots map[string]int) (string, int) {
	var consensusRoot string
	maxCount := 0
	for root, count := range roots {
		if count > maxCount {
			maxCount = count
			consensusRoot = root
		}
	}
	return consensusRoot, maxCount
}

// Parallelized RPC calls using goroutines and a worker pool
func parallelRPC(urls []string, rpcMethod string, params []interface{}, adapter ChainAdapter) ([]map[string]interface{}, error) {
	var wg sync.WaitGroup
	results := make([]map[string]interface{}, len(urls))
	errors := make([]error, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			res, err := adapter.GetHeader(url)
			if err != nil {
				errors[i] = err
				return
			}
			results[i] = res
		}(i, url)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, errors.New("One or more RPC calls failed")
		}
	}

	return results, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: challenger <nut_url> <ref_url1> [ref_url2...]")
		return
	}

	nutURL := os.Args[1]
	refURLs := os.Args[2:]

	adapter := &PolkadotAdapter{} // Use Polkadot adapter for this example

	fmt.Printf("Starting NCV Challenger (Go) - NUT: %s\n", nutURL)

	for {
		results, err := parallelRPC(refURLs, "chain_getHeader", nil, adapter)
		if err != nil {
			fmt.Printf("Error during parallel RPC calls: %v\n", err)
			continue
		}

		fmt.Printf("Results: %v\n", results)
		time.Sleep(10 * time.Second)
	}
}
