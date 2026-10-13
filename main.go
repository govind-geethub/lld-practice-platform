package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	initDB()
	app := fiber.New()

	app.Static("/static", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	app.Get("/api/problems", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, title, difficulty, description, skeleton FROM problems")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var problems []Problem
		for rows.Next() {
			var p Problem
			rows.Scan(&p.ID, &p.Title, &p.Difficulty, &p.Description, &p.Skeleton)
			problems = append(problems, p)
		}
		return c.JSON(problems)
	})

	app.Post("/api/attempts", func(c *fiber.Ctx) error {
		type Submission struct {
			ProblemID   int    `json:"problem_id"`
			CodeContent string `json:"code_content"`
			DesignNotes string `json:"design_notes"`
		}

		var sub Submission
		if err := c.BodyParser(&sub); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		result := evaluateDesign(sub.ProblemID, sub.CodeContent, sub.DesignNotes)
		feedbackBytes, _ := json.Marshal(result)
		createdAt := time.Now()

		res, err := db.Exec("INSERT INTO attempts (problem_id, code_content, design_notes, score, feedback, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			sub.ProblemID, sub.CodeContent, sub.DesignNotes, result.OverallScore, string(feedbackBytes), createdAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		id, _ := res.LastInsertId()
		return c.JSON(fiber.Map{
			"id":     id,
			"result": result,
		})
	})

	app.Get("/api/attempts/:problem_id", func(c *fiber.Ctx) error {
		problemID := c.Params("problem_id")
		rows, err := db.Query("SELECT id, problem_id, code_content, design_notes, score, feedback, created_at FROM attempts WHERE problem_id = ? ORDER BY created_at DESC", problemID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var attempts []Attempt
		for rows.Next() {
			var a Attempt
			rows.Scan(&a.ID, &a.ProblemID, &a.CodeContent, &a.DesignNotes, &a.Score, &a.Feedback, &a.CreatedAt)
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
