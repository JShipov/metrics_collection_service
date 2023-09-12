package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverEndpoint = "http://localhost:8080"
)

var pollCount int64

func main() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics := gatherMetrics()
			go reportMetrics(metrics)
		}
	}
}

func gatherMetrics() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]interface{}{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": float64(m.GCCPUFraction),
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
		"RandomValue":   rand.Float64(),
		"PollCount":     pollCount,
	}
	pollCount++

	return metrics
}

func reportMetrics(metrics map[string]interface{}) {
	for name, value := range metrics {
		var metricType, metricValue string

		switch v := value.(type) {
		case float64:
			metricType = "gauge"
			metricValue = fmt.Sprintf("%f", v)
		case int64:
			metricType = "counter"
			metricValue = fmt.Sprintf("%d", v)
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s", serverEndpoint, metricType, name, metricValue)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Printf("Failed to send metric %s: %v\n", name, err)
			continue
		}
		resp.Body.Close()
	}
}
