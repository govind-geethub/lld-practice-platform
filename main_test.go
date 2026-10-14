package main

import (
	"context"
	"testing"
)

func TestEvaluateDesignLegacyCompatibility(t *testing.T) {
	eval := &DeterministicEvaluator{}
	sub := Submission{
		ProblemID:   1,
		CodeContent: "class ParkingLot { interface SlotSelector {} }",
		DesignNotes: "Used Strategy pattern for pricing and handled concurrent thread-safe mutex locks.",
	}
	prob := Problem{
		ID:    1,
		Title: "Parking Lot",
	}

	result, err := eval.Evaluate(context.Background(), sub, prob)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

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
