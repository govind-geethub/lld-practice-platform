package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var defaultEvaluator Evaluator

func main() {
	initDB()

	// Initialize the decoupled evaluator pipeline (Answers Change Test B)
	defaultEvaluator = NewDefaultCompositeEvaluator()

	app := fiber.New()

	// Enable CORS for all origins so API can be queried from file:// or dev servers
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Serve static assets both from /static and root /
	app.Static("/static", "./public")
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	// GET /api/problems - List available LLD problems
	app.Get("/api/problems", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, title, difficulty, description, skeleton FROM problems")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		problems := []Problem{}
		for rows.Next() {
			var p Problem
			if err := rows.Scan(&p.ID, &p.Title, &p.Difficulty, &p.Description, &p.Skeleton); err != nil {
				continue
			}
			problems = append(problems, p)
		}
		return c.JSON(problems)
	})

	// POST /api/attempts - Submit design and execute structured evaluation lifecycle
	app.Post("/api/attempts", func(c *fiber.Ctx) error {
		var sub Submission
		if err := c.BodyParser(&sub); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if sub.Format == "" {
			sub.Format = FormatCodeAndNotes
		}

		// Retrieve problem context
		var prob Problem
		err := db.QueryRow("SELECT id, title, difficulty, description, skeleton FROM problems WHERE id = ?", sub.ProblemID).
			Scan(&prob.ID, &prob.Title, &prob.Difficulty, &prob.Description, &prob.Skeleton)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Problem not found"})
		}

		// Step 1: Pre-save attempt in EVALUATING state to prevent data loss on evaluator failure
		// (Section 10 Helping Guide requirement)
		createdAt := time.Now()
		res, err := db.Exec(
			"INSERT INTO attempts (problem_id, code_content, design_notes, status, score, feedback, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			sub.ProblemID, sub.CodeContent, sub.DesignNotes, string(StatusEvaluating), 0, "{}", createdAt,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to record submission: " + err.Error()})
		}
		attemptID, _ := res.LastInsertId()

		// Step 2: Execute evaluator pipeline with Context timeout protection
		evalCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, evalErr := defaultEvaluator.Evaluate(evalCtx, sub, prob)
		if evalErr != nil {
			// Record failure state
			db.Exec("UPDATE attempts SET status = ?, feedback = ? WHERE id = ?",
				string(StatusFailed), `{"error":"Evaluation pipeline failed"}`, attemptID)
			return c.Status(500).JSON(fiber.Map{
				"id":     attemptID,
				"error":  "Evaluation failed: " + evalErr.Error(),
				"status": StatusFailed,
			})
		}

		// Step 3: Update attempt with completed evaluation result
		result.Status = StatusCompleted
		feedbackBytes, _ := json.Marshal(result)

		_, err = db.Exec(
			"UPDATE attempts SET status = ?, score = ?, feedback = ? WHERE id = ?",
			string(StatusCompleted), result.OverallScore, string(feedbackBytes), attemptID,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update attempt results: " + err.Error()})
		}

		return c.JSON(fiber.Map{
			"id":     attemptID,
			"status": StatusCompleted,
			"result": result,
		})
	})

	// GET /api/attempts/:problem_id - Fetch attempt history with parsed rubric results
	app.Get("/api/attempts/:problem_id", func(c *fiber.Ctx) error {
		problemID := c.Params("problem_id")
		rows, err := db.Query(
			"SELECT id, problem_id, code_content, design_notes, status, score, feedback, created_at FROM attempts WHERE problem_id = ? ORDER BY created_at DESC",
			problemID,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		attempts := []Attempt{}
		for rows.Next() {
			var a Attempt
			var feedbackStr string
			var statusStr string
			err := rows.Scan(&a.ID, &a.ProblemID, &a.CodeContent, &a.DesignNotes, &statusStr, &a.Score, &feedbackStr, &a.CreatedAt)
			if err != nil {
				continue
			}
			a.Status = AttemptStatus(statusStr)
			a.Feedback = feedbackStr

			// Unmarshal feedback into result struct if valid
			var res EvaluationResult
			if json.Unmarshal([]byte(feedbackStr), &res) == nil {
				a.Result = &res
			}

			attempts = append(attempts, a)
		}
		return c.JSON(attempts)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
