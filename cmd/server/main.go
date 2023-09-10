package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func updateHandler(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		if len(parts) != 5 || parts[1] != "update" {
			http.Error(w, "Invalid request path", http.StatusNotFound)
			return
		}

		metricType := parts[2]
		metricName := parts[3]
		metricValue := parts[4]

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.gauges[metricName] = value

		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.counters[metricName] += value

		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Metric updated")

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	storage := NewMemStorage()

	http.HandleFunc("/update/", updateHandler(storage))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("Caught signal %s: shutting down.\n", sig)
			os.Exit(0)
		}
	}()

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
