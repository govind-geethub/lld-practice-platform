# LLD Practice Platform (MVP)

An engineering practice platform designed to help learners iteratively practice Low-Level Design (LLD), submit class hierarchies and architectural notes, receive explainable rubric-driven feedback, and track design evolution across attempts.

Built in accordance with the **LLD Practice Platform 2-Day Engineering Assignment** and **Candidate Helping Guide**.

---

## 🚀 Key Architectural Highlights

1. **Iterative Practice Loop:**
   $$\text{Choose Problem} \longrightarrow \text{Design \& Code} \longrightarrow \text{Submit} \longrightarrow \text{Explainable Rubric Feedback} \longrightarrow \text{Attempt History} \longrightarrow \text{Iterate}$$
2. **Multi-Dimensional Rubric Engine:**
   Replaces arbitrary scores with a 4-dimensional rubric adhering to the Helping Guide's exact shape:
   $$\text{Criterion} \longrightarrow \text{Score (1-5)} \longrightarrow \text{Evidence} \longrightarrow \text{Concern} \longrightarrow \text{Suggestion} \longrightarrow \text{Confidence}$$
   - **Single Responsibility & Cohesion (SRP)**
   - **Abstraction & Loose Coupling (DIP & ISP)**
   - **Design Patterns & Extensibility (OCP)**
   - **Concurrency Safety & State Management**
3. **Passes the Two Core "Change Tests":**
   - **Change Test A (Format Variation):** Domain model abstracts `SubmissionFormat`, allowing class diagrams or specs to be submitted without altering the core database or attempt lifecycle.
   - **Change Test B (Evaluator Variation):** Evaluators implement the `Evaluator` interface via a `CompositeEvaluator` pipeline, allowing AI models, AST compilers, or human review to be plugged in seamlessly.
4. **Resilient Submission Lifecycle:**
   Submissions are persisted immediately with status `EVALUATING` before computation starts, protected by Go context timeouts, preventing data loss on evaluator failure.
5. **Dynamic UML Class Diagram Generation:**
   Extracts classes and interfaces directly from the candidate's code and renders live interactive Mermaid.js diagrams.

---

## 🛠️ Tech Stack & Decisions

- **Backend:** Go (1.22+) with Fiber web framework (clean modular monolith).
- **Database:** SQLite (file-based persistence with attempt history and status tracking).
- **Frontend:** Vanilla HTML5, Tailwind CSS (CDN), dynamic Mermaid.js UML renderer.
- **Testing:** Go standard library `testing` package covering unit rules, composite pipelines, and edge cases.

---

## 📁 Repository Structure

```
.
├── docs/
│   ├── RESEARCH_NOTE.md   # 1-2 page analysis of learner problem, existing tools, and gaps
│   └── DESIGN_NOTE.md     # Domain model, Change Tests A/B, trade-offs, and architecture
├── public/
│   ├── index.html         # Workspace, rubric scorecards, and attempt history UI
│   ├── app.js             # Client logic, dynamic UML trigger, attempt restoration
│   └── style.css          # Custom styling
├── database.go            # SQLite schema initialization and problem seeding
├── models.go              # Core domain entities, Evaluator interface, Rubric structs
├── evaluator.go           # Deterministic & Composite evaluators + dynamic UML parser
├── evaluator_test.go      # Unit test suite for rubric rules, pipeline, and edge cases
├── main.go                # HTTP handlers and submission state lifecycle
├── AI_USAGE.md            # 5 AI-assisted architectural decisions (what was rejected/accepted & why)
└── README.md              # Project documentation
```

---

## 🚦 Quick Start Guide

### Prerequisites
- [Go (1.22+)](https://golang.org/) installed and available on your system.

### Run Application
```bash
# 1. Download dependencies
go mod tidy

# 2. Start the server
go run .
```

The application will start on **http://localhost:8080**.

### Run Test Suite
```bash
go test -v ./...
```

---

## 📝 Evaluation Deliverables Checklist

- [x] **Research Note:** [`docs/RESEARCH_NOTE.md`](file:///c:/Users/bente/gov/lld-practice-platform/docs/RESEARCH_NOTE.md)
- [x] **Design Note:** [`docs/DESIGN_NOTE.md`](file:///c:/Users/bente/gov/lld-practice-platform/docs/DESIGN_NOTE.md)
- [x] **Working Prototype:** Demonstrable end-to-end loop from problem selection to rubric feedback and attempt history.
- [x] **Tests:** Comprehensive tests in [`evaluator_test.go`](file:///c:/Users/bente/gov/lld-practice-platform/evaluator_test.go).
- [x] **README & AI_USAGE.md:** Full setup guide and engineering discernment documented in [`AI_USAGE.md`](file:///c:/Users/bente/gov/lld-practice-platform/AI_USAGE.md).