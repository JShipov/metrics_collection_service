package main

import (
	"testing"
)

func TestGatherMetrics(t *testing.T) {
	metrics := gatherMetrics()

	if _, exists := metrics["Alloc"]; !exists {
		t.Fatalf("Expected metric 'Alloc' to be present")
	}

	if val, exists := metrics["RandomValue"]; !exists || val.(float64) < 0 || val.(float64) > 1 {
		t.Fatalf("Expected metric 'RandomValue' to be between 0 and 1")
	}
}
