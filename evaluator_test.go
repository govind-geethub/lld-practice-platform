package main

import (
	"context"
	"testing"
	"time"
)

func TestDeterministicEvaluator_PassingDesign(t *testing.T) {
	eval := &DeterministicEvaluator{}
	sub := Submission{
		ProblemID:   1,
		CodeContent: "class ParkingLot { interface ParkingStrategy {} class CarManager {} }",
		DesignNotes: "Used Strategy pattern for pricing calculations and sync.Mutex for concurrent thread-safe slot allocation.",
	}
	prob := Problem{
		ID:    1,
		Title: "Parking Lot",
	}

	res, err := eval.Evaluate(context.Background(), sub, prob)
	if err != nil {
		t.Fatalf("Unexpected error from evaluation: %v", err)
	}

	if res.OverallScore < 70 {
		t.Errorf("Expected score >= 70, got %d", res.OverallScore)
	}

	if len(res.RubricDimensions) != 4 {
		t.Errorf("Expected 4 rubric dimensions, got %d", len(res.RubricDimensions))
	}

	if res.MermaidChart == "" {
		t.Errorf("Expected generated dynamic Mermaid chart, got empty string")
	}
}

func TestDeterministicEvaluator_DomainMismatchRejection(t *testing.T) {
	eval := &DeterministicEvaluator{}
	// Submitting URL shortener code for Parking lot problem
	sub := Submission{
		ProblemID:   1,
		CodeContent: "class URLShortener { String encode() {} }",
		DesignNotes: "Hashing with base62 encoding.",
	}
	prob := Problem{
		ID:    1,
		Title: "Parking Lot",
	}

	res, err := eval.Evaluate(context.Background(), sub, prob)
	if err != nil {
		t.Fatalf("Unexpected error from evaluation: %v", err)
	}

	if res.OverallScore != 0 {
		t.Errorf("Expected 0 score on domain mismatch, got %d", res.OverallScore)
	}

	if len(res.Critique) == 0 {
		t.Errorf("Expected rejection critique message")
	}
}

func TestCompositeEvaluator_PipelineChaining(t *testing.T) {
	composite := NewDefaultCompositeEvaluator()

	sub := Submission{
		ProblemID:   2,
		CodeContent: "class ElevatorController { interface DispatchStrategy {} }",
		DesignNotes: "Handles floor requests with thread-safe Mutex locks.",
	}
	prob := Problem{
		ID:    2,
		Title: "Elevator System",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := composite.Evaluate(ctx, sub, prob)
	if err != nil {
		t.Fatalf("Pipeline failed: %v", err)
	}

	if res.Status != StatusCompleted {
		t.Errorf("Expected status COMPLETED, got %s", res.Status)
	}
}

func TestDynamicUMLParser(t *testing.T) {
	code := `
class ParkingLot {
}
interface PricingStrategy {
}
class CarManager {
}
`
	uml := generateDynamicUML(code)
	if uml == "" {
		t.Fatalf("Expected non-empty UML")
	}
	if !contains(uml, "class ParkingLot") {
		t.Errorf("Expected ParkingLot in UML")
	}
	if !contains(uml, "<<interface>>") {
		t.Errorf("Expected interface annotation in UML")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && stringIndex(s, substr) >= 0))
}

func stringIndex(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
