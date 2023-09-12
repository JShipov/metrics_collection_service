package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/mux"
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
		vars := mux.Vars(r)
		metricType := vars["metricType"]
		metricName := vars["metricName"]
		metricValue := vars["metricValue"]

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

		fmt.Fprint(w, "Metric updated")
	}
}

func valueHandler(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		metricType := vars["metricType"]
		metricName := vars["metricName"]

		switch metricType {
		case "gauge":
			if value, ok := storage.gauges[metricName]; ok {
				fmt.Fprint(w, value)
				return
			}
		case "counter":
			if value, ok := storage.counters[metricName]; ok {
				fmt.Fprint(w, value)
				return
			}
		}

		http.Error(w, "Metric not found", http.StatusNotFound)
	}
}

func listMetrics(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Metrics</h1>")
		fmt.Fprintf(w, "<h2>Gauges</h2>")
		for name, value := range storage.gauges {
			fmt.Fprintf(w, "<div>%s: %f</div>", name, value)
		}
		fmt.Fprintf(w, "<h2>Counters</h2>")
		for name, value := range storage.counters {
			fmt.Fprintf(w, "<div>%s: %d</div>", name, value)
		}
	}
}

var addr string

func init() {
	flag.StringVar(&addr, "a", getEnv("ADDRESS", "localhost:8080"), "HTTP server address")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("Unknown arguments:", flag.Args())
		os.Exit(1)
	}

	r := mux.NewRouter()
	storage := NewMemStorage()

	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", updateHandler(storage)).Methods("POST")
	r.HandleFunc("/value/{metricType}/{metricName}", valueHandler(storage)).Methods("GET")
	r.HandleFunc("/", listMetrics(storage)).Methods("GET")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("Caught signal %s: shutting down.\n", sig)
			os.Exit(0)
		}
	}()

	fmt.Printf("Server running on http://%s\n", addr)
	http.ListenAndServe(addr, r)
}
