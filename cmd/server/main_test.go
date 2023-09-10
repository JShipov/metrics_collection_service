package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	storage := NewMemStorage()
	handler := updateHandler(storage)

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	resp, err := http.Post(server.URL+"/update/gauge/testMetric/1.23", "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}
	if storage.gauges["testMetric"] != 1.23 {
		t.Fatalf("Expected gauge value 1.23; got %v", storage.gauges["testMetric"])
	}

	resp, err = http.Post(server.URL+"/update/counter/testCounter/10", "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}
	if storage.counters["testCounter"] != 10 {
		t.Fatalf("Expected counter value 10; got %v", storage.counters["testCounter"])
	}
}
