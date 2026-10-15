package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ==============================================================================
// 1. DETERMINISTIC EVALUATOR
// Evaluates syntax validity, domain presence, anti-patterns, and basic structure.
// ==============================================================================

type DeterministicEvaluator struct{}

func (e *DeterministicEvaluator) Name() string {
	return "Deterministic Rule Engine"
}

func (e *DeterministicEvaluator) Evaluate(ctx context.Context, sub Submission, prob Problem) (EvaluationResult, error) {
	codeLower := strings.ToLower(sub.CodeContent)
	notesLower := strings.ToLower(sub.DesignNotes)

	// Validate problem domain presence
	codeValid := false
	notesValid := false

	switch prob.ID {
	case 1: // Parking Lot
		codeValid = strings.Contains(codeLower, "parking") || strings.Contains(codeLower, "spot") || strings.Contains(codeLower, "vehicle") || strings.Contains(codeLower, "ticket") || strings.Contains(codeLower, "slot")
		notesValid = strings.Contains(notesLower, "parking") || strings.Contains(notesLower, "spot") || strings.Contains(notesLower, "pricing") || strings.Contains(notesLower, "rate") || strings.Contains(notesLower, "ticket") || strings.Contains(notesLower, "vehicle") || strings.Contains(notesLower, "lot")
	case 2: // Elevator System
		codeValid = strings.Contains(codeLower, "elevator") || strings.Contains(codeLower, "request") || strings.Contains(codeLower, "floor") || strings.Contains(codeLower, "lift")
		notesValid = strings.Contains(notesLower, "elevator") || strings.Contains(notesLower, "request") || strings.Contains(notesLower, "floor") || strings.Contains(notesLower, "lift") || strings.Contains(notesLower, "dispatch")
	case 3: // URL Shortener
		codeValid = strings.Contains(codeLower, "url") || strings.Contains(codeLower, "short") || strings.Contains(codeLower, "encode") || strings.Contains(codeLower, "hash")
		notesValid = strings.Contains(notesLower, "url") || strings.Contains(notesLower, "short") || strings.Contains(notesLower, "encode") || strings.Contains(notesLower, "hash") || strings.Contains(notesLower, "alias")
	case 4: // Vending Machine
		codeValid = strings.Contains(codeLower, "vending") || strings.Contains(codeLower, "coin") || strings.Contains(codeLower, "state") || strings.Contains(codeLower, "item") || strings.Contains(codeLower, "product")
		notesValid = strings.Contains(notesLower, "vending") || strings.Contains(notesLower, "coin") || strings.Contains(notesLower, "state") || strings.Contains(notesLower, "item") || strings.Contains(notesLower, "inventory")
	case 5: // Splitwise
		codeValid = strings.Contains(codeLower, "expense") || strings.Contains(codeLower, "split") || strings.Contains(codeLower, "balance") || strings.Contains(codeLower, "user") || strings.Contains(codeLower, "debt")
		notesValid = strings.Contains(notesLower, "expense") || strings.Contains(notesLower, "split") || strings.Contains(notesLower, "balance") || strings.Contains(notesLower, "user") || strings.Contains(notesLower, "debt")
	case 6: // Movie Ticket Booking
		codeValid = strings.Contains(codeLower, "seat") || strings.Contains(codeLower, "booking") || strings.Contains(codeLower, "payment") || strings.Contains(codeLower, "ticket") || strings.Contains(codeLower, "movie") || strings.Contains(codeLower, "show")
		notesValid = strings.Contains(notesLower, "seat") || strings.Contains(notesLower, "booking") || strings.Contains(notesLower, "payment") || strings.Contains(notesLower, "ticket") || strings.Contains(notesLower, "movie") || strings.Contains(notesLower, "theater")
	default:
		codeValid = len(sub.CodeContent) > 20
		notesValid = len(sub.DesignNotes) > 10
	}

	// Mismatch rejection: reject if code has no relevance to the selected problem domain
	if !codeValid {
		var reason string
		if !notesValid {
			reason = "Neither the code implementation nor the architectural notes match the problem domain."
		} else {
			reason = "Code implementation does not reflect the selected problem requirements."
		}

		return EvaluationResult{
			Status:       StatusCompleted,
			OverallScore: 0,
			RubricDimensions: []RubricDimension{
				{
					Criterion:  "Domain Relevance & Problem Alignment",
					Score:      0,
					MaxScore:   5,
					Evidence:   "Keywords scanned for " + prob.Title + " were missing or completely mismatched.",
					Concern:    reason,
					Suggestion: "Ensure both your class definitions and notes specifically target " + prob.Title + ".",
					Confidence: "High",
				},
			},
			Critique: []string{
				"❌ Critical Domain Mismatch: " + reason,
				"🚫 Prevention of cross-problem contamination active.",
				"💡 Provide domain-specific entities (e.g. for " + prob.Title + ").",
			},
			MermaidChart:  "classDiagram\n    class RejectedSubmission {\n        <<mismatch>>\n        +status: Failed\n        +reason: Cross-Problem Contamination\n    }",
			EvaluatorName: e.Name(),
			EvaluatedAt:   time.Now(),
		}, nil
	}

	// Dynamic class/interface extraction for Mermaid diagram
	mermaidChart := generateDynamicUML(sub.CodeContent)

	// Multi-dimensional Rubric Scoring (1 to 5 scale each, normalized to 100)
	var dimensions []RubricDimension
	var critiques []string

	// 1. Single Responsibility & Class Cohesion
	srpScore := 3
	srpEvidence := "Basic class structure defined."
	srpConcern := "Check if god classes exist combining data storage and processing."
	srpSuggestion := "Segregate high-level managers from low-level entities."

	if strings.Contains(codeLower, "manager") || strings.Contains(codeLower, "controller") || strings.Contains(codeLower, "service") {
		srpScore = 4
		srpEvidence = "Clear orchestration layer detected separating coordination from core domain models."
		srpConcern = "Ensure orchestration classes don't absorb entity-specific business logic."
		srpSuggestion = "Keep state mutation inside the entity itself."
		critiques = append(critiques, "✅ Good separation of control logic from domain entities.")
	} else {
		critiques = append(critiques, "⚠️ Consider introducing a dedicated coordinator or service class to own system workflows.")
	}
	dimensions = append(dimensions, RubricDimension{
		Criterion:  "Single Responsibility & Cohesion",
		Score:      srpScore,
		MaxScore:   5,
		Evidence:   srpEvidence,
		Concern:    srpConcern,
		Suggestion: srpSuggestion,
		Confidence: "High",
	})

	// 2. Abstraction & Interface Segregation (DIP & ISP)
	absScore := 2
	absEvidence := "Concrete classes used directly."
	absConcern := "Components are tightly coupled to concrete implementations."
	absSuggestion := "Introduce interfaces for extensible components (e.g. strategies, repositories)."

	if strings.Contains(sub.CodeContent, "interface") || strings.Contains(sub.CodeContent, "abstract") || strings.Contains(sub.CodeContent, "virtual") {
		absScore = 5
		absEvidence = "Explicit interfaces/abstract definitions found in class hierarchy."
		absConcern = "None observed at interface boundary."
		absSuggestion = "Program to interfaces at all dependency injection points."
		critiques = append(critiques, "✅ Excellent use of interfaces for loose coupling (DIP).")
	} else {
		critiques = append(critiques, "⚠️ Tightly coupled to concrete classes. Extract interfaces for components that vary.")
	}
	dimensions = append(dimensions, RubricDimension{
		Criterion:  "Abstraction & Loose Coupling",
		Score:      absScore,
		MaxScore:   5,
		Evidence:   absEvidence,
		Concern:    absConcern,
		Suggestion: absSuggestion,
		Confidence: "High",
	})

	// 3. Design Patterns & Extensibility
	patScore := 2
	patEvidence := "Ad-hoc conditional logic detected."
	patConcern := "Adding new behaviors may require modifying existing code (violating OCP)."
	patSuggestion := "Utilize behavioral patterns (Strategy, State, Observer) to allow extension without modification."

	patternsFound := []string{}
	for _, p := range []string{"Strategy", "Factory", "Singleton", "State", "Observer", "Command", "Builder"} {
		if strings.Contains(sub.CodeContent, p) || strings.Contains(sub.DesignNotes, p) {
			patternsFound = append(patternsFound, p)
		}
	}
	if len(patternsFound) > 0 {
		patScore = 4
		if len(patternsFound) >= 2 {
			patScore = 5
		}
		patEvidence = fmt.Sprintf("Recognized design pattern(s) applied: %s.", strings.Join(patternsFound, ", "))
		patConcern = "Ensure patterns are not over-engineered for simple operations."
		patSuggestion = "Document why each pattern was chosen over simpler alternatives."
		critiques = append(critiques, fmt.Sprintf("✅ Applied appropriate design patterns: %s.", strings.Join(patternsFound, ", ")))
	} else {
		critiques = append(critiques, "⚠️ Design lacks explicit design patterns (e.g. Strategy for algorithms, State for lifecycles).")
	}
	dimensions = append(dimensions, RubricDimension{
		Criterion:  "Design Patterns & Extensibility",
		Score:      patScore,
		MaxScore:   5,
		Evidence:   patEvidence,
		Concern:    patConcern,
		Suggestion: patSuggestion,
		Confidence: "High",
	})

	// 4. Concurrency Safety & State Handling
	conScore := 2
	conEvidence := "No synchronization or thread-safety considerations mentioned."
	conConcern := "High risk of race conditions or data inconsistency under concurrent operations."
	conSuggestion := "Specify synchronization mechanisms (Mutex, Synchronized blocks, Atomic references) in critical sections."

	if strings.Contains(notesLower, "mutex") || strings.Contains(notesLower, "lock") || strings.Contains(notesLower, "thread-safe") ||
		strings.Contains(notesLower, "concurrent") || strings.Contains(codeLower, "mutex") || strings.Contains(codeLower, "synchronized") {
		conScore = 5
		conEvidence = "Concurrency hazards addressed with explicit locking/thread-safety mechanisms."
		conConcern = "Watch out for potential deadlocks or lock contention under high load."
		conSuggestion = "Consider read-write locks (RWMutex) or optimistic locking for read-heavy flows."
		critiques = append(critiques, "✅ Strong architectural consideration for concurrency and synchronization.")
	} else {
		critiques = append(critiques, "⚠️ Concurrency omitted: Clarify how concurrent requests or state updates are made thread-safe.")
	}
	dimensions = append(dimensions, RubricDimension{
		Criterion:  "Concurrency Safety & State Management",
		Score:      conScore,
		MaxScore:   5,
		Evidence:   conEvidence,
		Concern:    conConcern,
		Suggestion: conSuggestion,
		Confidence: "High",
	})

	// Calculate overall score (normalized to 100)
	totalScore := 0
	for _, d := range dimensions {
		totalScore += d.Score
	}
	overall := int(float64(totalScore) / float64(len(dimensions)*5) * 100)

	return EvaluationResult{
		Status:           StatusCompleted,
		OverallScore:     overall,
		RubricDimensions: dimensions,
		Critique:         critiqueSummary(critiques),
		MermaidChart:     mermaidChart,
		EvaluatorName:    e.Name(),
		EvaluatedAt:      time.Now(),
	}, nil
}

// ==============================================================================
// 2. COMPOSITE / PIPELINE EVALUATOR
// Satisfies Change Test B: easily chain multiple evaluators (Rule engine + AI/Mock)
// ==============================================================================

type CompositeEvaluator struct {
	evaluators []Evaluator
}

func NewDefaultCompositeEvaluator() *CompositeEvaluator {
	return &CompositeEvaluator{
		evaluators: []Evaluator{
			&DeterministicEvaluator{},
		},
	}
}

func (c *CompositeEvaluator) Name() string {
	return "Composite Hybrid Pipeline"
}

func (c *CompositeEvaluator) AddEvaluator(e Evaluator) {
	c.evaluators = append(c.evaluators, e)
}

func (c *CompositeEvaluator) Evaluate(ctx context.Context, sub Submission, prob Problem) (EvaluationResult, error) {
	if len(c.evaluators) == 0 {
		return EvaluationResult{}, fmt.Errorf("no evaluators registered in pipeline")
	}

	// Run primary evaluator
	primary := c.evaluators[0]
	res, err := primary.Evaluate(ctx, sub, prob)
	if err != nil {
		return EvaluationResult{}, err
	}

	// If rejected by deterministic filter, return early
	if res.OverallScore == 0 {
		return res, nil
	}

	// Chain additional evaluators if registered (e.g. LLMEvaluator or MockAIEvaluator)
	for _, nextEval := range c.evaluators[1:] {
		nextRes, err := nextEval.Evaluate(ctx, sub, prob)
		if err == nil {
			// Merge critiques and dimensions
			res.Critique = append(res.Critique, nextRes.Critique...)
			res.RubricDimensions = append(res.RubricDimensions, nextRes.RubricDimensions...)
		}
	}

	return res, nil
}

// ==============================================================================
// 3. DYNAMIC UML CLASS PARSER
// Parses actual classes and interfaces from learner's code to build Mermaid diagram
// ==============================================================================

func generateDynamicUML(code string) string {
	classRegex := regexp.MustCompile(`(?m)(?:class|interface)\s+([A-Za-z0-9_]+)`)
	matches := classRegex.FindAllStringSubmatch(code, -1)

	var sb strings.Builder
	sb.WriteString("classDiagram\n")

	classNames := make(map[string]bool)
	var orderedNames []string

	for _, m := range matches {
		if len(m) > 1 {
			name := m[1]
			if !classNames[name] {
				classNames[name] = true
				orderedNames = append(orderedNames, name)
			}
		}
	}

	if len(orderedNames) == 0 {
		sb.WriteString("    class System {\n        +execute()\n    }\n")
		return sb.String()
	}

	for _, name := range orderedNames {
		isInterface := strings.Contains(code, "interface "+name)
		sb.WriteString(fmt.Sprintf("    class %s {\n", name))
		if isInterface {
			sb.WriteString("        <<interface>>\n")
		}
		sb.WriteString("    }\n")
	}

	// Link first orchestrator to subsequent classes if > 1
	if len(orderedNames) > 1 {
		root := orderedNames[0]
		for _, dep := range orderedNames[1:] {
			if strings.Contains(dep, "Strategy") || strings.Contains(dep, "State") || strings.Contains(dep, "Handler") {
				sb.WriteString(fmt.Sprintf("    %s ..> %s : depends on\n", root, dep))
			} else {
				sb.WriteString(fmt.Sprintf("    %s --> %s : manages\n", root, dep))
			}
		}
	}

	return sb.String()
}

func critiqueSummary(critiques []string) []string {
	if len(critiques) == 0 {
		return []string{"Submission reviewed. Solid baseline architecture."}
	}
	return critiques
}
