package framework

import (
	"fmt"
)

// Challenger handles node challenges and verification.
// It uses a ChainModule to perform chain-specific logic.
type Challenger struct {
	ChainModule ChainModule
}

// ChainModule defines the interface for chain-specific logic.
// Implementations of this interface must provide methods for:
// - HeadCorrectness: Verifying the state root at a specific block height.
// - ExecutionCorrectness: Simulating transactions and verifying results.
// - FreshnessProof: Monitoring block arrival times.
type ChainModule interface {
	HeadCorrectness() error
	ExecutionCorrectness() error
	FreshnessProof() error
}

// NewChallenger creates a new Challenger instance.
// module: The ChainModule implementation to use.
// Returns a pointer to the created Challenger instance.
func NewChallenger(module ChainModule) *Challenger {
	return &Challenger{
		ChainModule: module,
	}
}

// RunChallenges runs all the defined challenges in the ChainModule.
// Logs the results of each challenge to the console.
func (c *Challenger) RunChallenges() {
	fmt.Println("Running challenges...")

	if err := c.ChainModule.HeadCorrectness(); err != nil {
		fmt.Printf("Head Correctness failed: %v\n", err)
	}

	if err := c.ChainModule.ExecutionCorrectness(); err != nil {
		fmt.Printf("Execution Correctness failed: %v\n", err)
	}

	if err := c.ChainModule.FreshnessProof(); err != nil {
		fmt.Printf("Freshness Proof failed: %v\n", err)
	}
}