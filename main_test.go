package main

import (
	"testing"
)

func TestEvaluateDesign(t *testing.T) {
	code := "class ParkingLot { interface SlotSelector {} }"
	notes := "Used Strategy pattern for pricing and handled concurrent thread-safe mutex locks."

	// Pass problemID = 1 as the first argument
	result := evaluateDesign(1, code, notes)

	if result.OverallScore < 70 {
		t.Errorf("Expected passing score, got %d", result.OverallScore)
	}
	if len(result.Critique) == 0 {
		t.Errorf("Expected critique output, got empty list")
	}
	if result.MermaidChart == "" {
		t.Errorf("Expected generated mermaid chart, got empty string")
	}
}
