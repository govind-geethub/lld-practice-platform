package main

import (
	"strings"
)

func evaluateDesign(problemID int, code, notes string) EvaluationResult {
	codeLower := strings.ToLower(code)
	notesLower := strings.ToLower(notes)

	// Strict Independent Validation for Code and Notes
	codeValid := false
	notesValid := false

	switch problemID {
	case 1: // Parking Lot
		codeValid = strings.Contains(codeLower, "parking") || strings.Contains(codeLower, "spot") || strings.Contains(codeLower, "vehicle")
		notesValid = strings.Contains(notesLower, "parking") || strings.Contains(notesLower, "spot") || strings.Contains(notesLower, "pricing")
	case 2: // Elevator System
		codeValid = strings.Contains(codeLower, "elevator") || strings.Contains(codeLower, "request") || strings.Contains(codeLower, "floor")
		notesValid = strings.Contains(notesLower, "elevator") || strings.Contains(notesLower, "request") || strings.Contains(notesLower, "floor")
	case 3: // URL Shortener
		codeValid = strings.Contains(codeLower, "url") || strings.Contains(codeLower, "short") || strings.Contains(codeLower, "encode")
		notesValid = strings.Contains(notesLower, "url") || strings.Contains(notesLower, "short") || strings.Contains(notesLower, "encode")
	case 4: // Vending Machine
		codeValid = strings.Contains(codeLower, "vending") || strings.Contains(codeLower, "coin") || strings.Contains(codeLower, "state")
		notesValid = strings.Contains(notesLower, "vending") || strings.Contains(notesLower, "coin") || strings.Contains(notesLower, "state")
	case 5: // Splitwise
		codeValid = strings.Contains(codeLower, "expense") || strings.Contains(codeLower, "split") || strings.Contains(codeLower, "balance")
		notesValid = strings.Contains(notesLower, "expense") || strings.Contains(notesLower, "split") || strings.Contains(notesLower, "balance")
	case 6: // Movie Ticket Booking
		codeValid = strings.Contains(codeLower, "seat") || strings.Contains(codeLower, "booking") || strings.Contains(codeLower, "payment")
		notesValid = strings.Contains(notesLower, "seat") || strings.Contains(notesLower, "booking") || strings.Contains(notesLower, "payment")
	}

	// Hard rejection if either the code or the notes mismatch the selected problem domain
	if !codeValid || !notesValid {
		var reason string
		if !codeValid && !notesValid {
			reason = "Both implementation code and architectural notes belong to a different domain."
		} else if !codeValid {
			reason = "Implementation code does not match the selected problem domain."
		} else {
			reason = "Architectural notes belong to a different problem domain."
		}

		return EvaluationResult{
			OverallScore: 0,
			SolidScore:   0,
			PatternScore: 0,
			Concurrency:  0,
			Critique: []string{
				"❌ Critical Domain Mismatch: " + reason,
				"🚫 Strict One-to-One Enforcement: Mixing code from one problem with notes from another is prohibited.",
				"💡 Action Required: Ensure both your code and notes explicitly target the selected problem.",
			},
			MermaidChart: "classDiagram\n    class RejectedSubmission {\n        <<mismatch>>\n        +status: Failed\n        +reason: Cross-Problem Contamination\n    }",
		}
	}

	// Normal Evaluation Pipeline for valid matching submissions
	solid := 70
	patterns := 65
	concurrency := 60
	critique := []string{}

	if strings.Contains(code, "interface") || strings.Contains(code, "abstract") || strings.Contains(code, "virtual") {
		solid += 25
		critique = append(critique, "✅ Excellent use of interfaces/abstractions for decoupling.")
	} else {
		critique = append(critique, "⚠️ Consider introducing interfaces to adhere to the Dependency Inversion Principle.")
	}

	if strings.Contains(code, "Strategy") || strings.Contains(code, "Factory") || strings.Contains(code, "Singleton") || strings.Contains(code, "State") {
		patterns += 30
		critique = append(critique, "✅ Detected recognized design patterns.")
	} else {
		critique = append(critique, "⚠️ Explicitly highlight design patterns.")
	}

	if strings.Contains(notesLower, "mutex") || strings.Contains(notesLower, "concurrent") || strings.Contains(notesLower, "thread-safe") || strings.Contains(notesLower, "lock") || strings.Contains(codeLower, "mutex") || strings.Contains(codeLower, "lock") {
		concurrency += 35
		critique = append(critique, "✅ Strong architectural consideration for thread-safety and concurrency.")
	} else {
		critique = append(critique, "⚠️ Elaborate on concurrency management in your design notes.")
	}

	if solid > 100 {
		solid = 100
	}
	if patterns > 100 {
		patterns = 100
	}
	if concurrency > 100 {
		concurrency = 100
	}

	overall := (solid + patterns + concurrency) / 3

	mermaid := "classDiagram\n    class CoreManager {\n        +getInstance()\n        +processRequest()\n    }\n    class StrategyInterface {\n        <<interface>>\n        +execute()\n    }\n    CoreManager --> StrategyInterface : uses"

	return EvaluationResult{
		OverallScore: overall,
		SolidScore:   solid,
		PatternScore: patterns,
		Concurrency:  concurrency,
		Critique:     critique,
		MermaidChart: mermaid,
	}
}
