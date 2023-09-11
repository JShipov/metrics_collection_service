package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandlers(t *testing.T) {
	storage := NewMemStorage()
	r := mux.NewRouter()

	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", updateHandler(storage)).Methods("POST")
	r.HandleFunc("/value/{metricType}/{metricName}", valueHandler(storage)).Methods("GET")
	r.HandleFunc("/", listMetrics(storage)).Methods("GET")

	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := http.Post(fmt.Sprintf("%s/update/gauge/testGauge/42.0", server.URL), "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	resp, err = http.Get(fmt.Sprintf("%s/value/gauge/testGauge", server.URL))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	expected := "42"
	if string(body) != expected {
		t.Fatalf("Expected body to be %v; got %v", expected, string(body))
	}
}
