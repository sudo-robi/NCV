package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRPCQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"jsonrpc":"2.0","result":{"number":"0x1","stateRoot":"0xabc"},"id":1}`)
	}))
	defer server.Close()

	res, err := queryRPC(server.URL, "chain_getHeader", []interface{}{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	header := res.(map[string]interface{})
	if header["stateRoot"] != "0xabc" {
		t.Errorf("Expected 0xabc, got %v", header["stateRoot"])
	}
}
