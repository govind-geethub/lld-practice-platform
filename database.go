package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./lld_platform.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS problems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		difficulty TEXT,
		description TEXT,
		skeleton TEXT
	);
	CREATE TABLE IF NOT EXISTS attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		problem_id INTEGER,
		code_content TEXT,
		design_notes TEXT,
		status TEXT DEFAULT 'COMPLETED',
		score INTEGER,
		feedback TEXT,
		created_at DATETIME
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	// Ensure status column exists if table was previously created without it
	db.Exec("ALTER TABLE attempts ADD COLUMN status TEXT DEFAULT 'COMPLETED'")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count)
	if count == 0 {
		seedProblems()
	}
}
