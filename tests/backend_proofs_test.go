package backend

import (
	"testing"

	"ncv/backend/modules"
)

// These tests now validate that the PolkadotModule implementations return nil error
// for the placeholder logic, ensuring the interface is wired.

func TestHeadCorrectness(t *testing.T) {
	m := &modules.PolkadotModule{}
	if err := m.HeadCorrectness(); err != nil {
		t.Fatalf("HeadCorrectness returned error: %v", err)
	}
}

func TestFreshnessProof(t *testing.T) {
	m := &modules.PolkadotModule{}
	if err := m.FreshnessProof(); err != nil {
		t.Fatalf("FreshnessProof returned error: %v", err)
	}
}

func TestExecutionProof(t *testing.T) {
	m := &modules.PolkadotModule{}
	if err := m.ExecutionCorrectness(); err != nil {
		t.Fatalf("ExecutionCorrectness returned error: %v", err)
	}
}