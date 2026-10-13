package main

import (
	"time"
)

type Problem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Difficulty  string `json:"difficulty"`
	Description string `json:"description"`
	Skeleton    string `json:"skeleton"`
}

type Attempt struct {
	ID          int       `json:"id"`
	ProblemID   int       `json:"problem_id"`
	CodeContent string    `json:"code_content"`
	DesignNotes string    `json:"design_notes"`
	Score       int       `json:"score"`
	Feedback    string    `json:"feedback"`
	CreatedAt   time.Time `json:"created_at"`
}

type EvaluationResult struct {
	OverallScore int      `json:"overall_score"`
	SolidScore   int      `json:"solid_score"`
	PatternScore int      `json:"pattern_score"`
	Concurrency  int      `json:"concurrency_score"`
	Critique     []string `json:"critique"`
	MermaidChart string   `json:"mermaid_chart"`
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
