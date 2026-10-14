package main

import (
	"context"
	"time"
)

// AttemptStatus models the explicit lifecycle requested in Section 10 of the Candidate Helping Guide.
type AttemptStatus string

const (
	StatusSubmitted  AttemptStatus = "SUBMITTED"
	StatusEvaluating AttemptStatus = "EVALUATING"
	StatusCompleted  AttemptStatus = "COMPLETED"
	StatusFailed     AttemptStatus = "FAILED"
)

// SubmissionFormat models the submission representation (Addresses Change Test A:
// "Today the learner submits text/code. Later the platform supports a class diagram. How much of your domain model changes?")
type SubmissionFormat string

const (
	FormatCodeAndNotes SubmissionFormat = "CODE_AND_NOTES"
	FormatClassDiagram SubmissionFormat = "CLASS_DIAGRAM"
	FormatSpecification SubmissionFormat = "SPECIFICATION"
)

type Problem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Difficulty  string `json:"difficulty"`
	Description string `json:"description"`
	Skeleton    string `json:"skeleton"`
}

// Submission encapsulates candidate input payload with format abstraction.
type Submission struct {
	ProblemID   int              `json:"problem_id"`
	Format      SubmissionFormat `json:"format,omitempty"` // defaults to FormatCodeAndNotes
	CodeContent string           `json:"code_content"`
	DesignNotes string           `json:"design_notes"`
	DiagramData string           `json:"diagram_data,omitempty"`
}

// RubricDimension defines the exact structured feedback shape recommended in the Helping Guide:
// criterion -> score -> evidence -> concern -> suggestion -> confidence
type RubricDimension struct {
	Criterion  string `json:"criterion"`  // e.g. "Single Responsibility & Cohesion"
	Score      int    `json:"score"`      // 1-5 scale (or percentage)
	MaxScore   int    `json:"max_score"`  // 5
	Evidence   string `json:"evidence"`   // What was found in candidate submission
	Concern    string `json:"concern"`    // What is missing or problematic
	Suggestion string `json:"suggestion"` // Actionable advice to improve
	Confidence string `json:"confidence"` // High, Medium, Low
}

// EvaluationResult aggregates multi-dimensional rubric feedback and dynamic UML visualization.
type EvaluationResult struct {
	Status           AttemptStatus     `json:"status"`
	OverallScore     int               `json:"overall_score"`
	RubricDimensions []RubricDimension `json:"rubric_dimensions"`
	Critique         []string          `json:"critique"`
	MermaidChart     string            `json:"mermaid_chart"`
	EvaluatorName    string            `json:"evaluator_name"`
	EvaluatedAt      time.Time         `json:"evaluated_at"`
}

// Attempt models a persisted practice attempt supporting the iterative learning loop.
type Attempt struct {
	ID          int              `json:"id"`
	ProblemID   int              `json:"problem_id"`
	CodeContent string           `json:"code_content"`
	DesignNotes string           `json:"design_notes"`
	Status      AttemptStatus    `json:"status"`
	Score       int              `json:"score"`
	Feedback    string           `json:"feedback"`
	Result      *EvaluationResult `json:"result,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

// Evaluator is the core domain interface answering Change Test B:
// "Today feedback comes from one evaluator. Later you add a rule-based evaluator, human review, or an LLM.
// Can you add it without rewriting the practice flow?"
type Evaluator interface {
	Name() string
	Evaluate(ctx context.Context, sub Submission, prob Problem) (EvaluationResult, error)
}

func seedProblems() {
	problems := []Problem{
		{
			Title:       "Parking Lot",
			Difficulty:  "Medium",
			Description: "Design a parking lot supporting multiple floors, vehicle types (Car, Bike, Truck), hourly pricing strategies, and parking spot allocation.",
			Skeleton:    "class ParkingLot {\n  // TODO: Add floors, capacity, entry/exit gates\n}\n\ninterface PricingStrategy {\n  calculateCost(Ticket t): float\n}",
		},
		{
			Title:       "Elevator System",
			Difficulty:  "Hard",
			Description: "Design an intelligent elevator dispatch system for a multi-story building handling concurrent request queues, direction optimization, and door states.",
			Skeleton:    "class ElevatorController {\n  // TODO: Manage multiple elevators, requests, and dispatch algorithms\n}",
		},
		{
			Title:       "URL Shortener",
			Difficulty:  "Medium",
			Description: "Design a URL shortener service that encodes long URLs into short unique hashes, handles custom aliases, and tracks click analytics.",
			Skeleton:    "class URLService {\n  // TODO: Implement hashing and URL redirection mapping\n}\n\ninterface ShorteningStrategy {\n  String encode(long id);\n}",
		},
		{
			Title:       "Vending Machine",
			Difficulty:  "Medium",
			Description: "Design a state-driven vending machine managing coin insertion, product inventory, state transitions, and change return.",
			Skeleton:    "class VendingMachine {\n  // TODO: Manage states (Idle, HasMoney, Dispensing)\n}\n\ninterface VendingState {\n  void selectItem();\n}",
		},
		{
			Title:       "Splitwise",
			Difficulty:  "Medium",
			Description: "Design an expense-sharing application to track individual balances, handle equal/exact/percentage splits, and simplify debts.",
			Skeleton:    "class ExpenseManager {\n  // TODO: Add expenses, compute balances, simplify debts\n}\n\ninterface SplitStrategy {\n  List<Split> split(double amount, List<User> users);\n}",
		},
		{
			Title:       "Movie Ticket Booking System",
			Difficulty:  "Hard",
			Description: "Design a concurrent movie ticket booking platform supporting theaters, shows, seat categories, temporary locks, and payment processing.",
			Skeleton:    "class BookingManager {\n  // TODO: Handle concurrent seat selection, temporary locks, and payments\n}\n\ninterface PaymentStrategy {\n  boolean processPayment(double amount);\n}",
		},
	}

	for _, p := range problems {
		db.Exec("INSERT INTO problems (title, difficulty, description, skeleton) VALUES (?, ?, ?, ?)",
			p.Title, p.Difficulty, p.Description, p.Skeleton)
	}
}
