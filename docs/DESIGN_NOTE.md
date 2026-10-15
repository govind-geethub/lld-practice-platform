# Design Note: Low-Level Design (LLD) Practice Platform

**Author:** Govind Kumar  
**Date:** October 2026  
**Assignment:** 2-Day Engineering Prototype  

---

## 1. System Overview & Practice Loop

The platform is designed around the core learner iteration cycle:
$$\text{Choose Problem} \longrightarrow \text{Design Implementation \& Notes} \longrightarrow \text{Submit} \longrightarrow \text{Receive Rubric Scorecard} \longrightarrow \text{Inspect Dynamic UML} \longrightarrow \text{Iterate Attempt}$$

Rather than building an LMS or heavyweight assessment suite, this prototype prioritizes **explainable, evidence-backed architectural feedback** and **rapid re-attempt learning**.

---

## 2. Core Domain Model & Responsibilities

The domain is organized into decoupled entities reflecting distinct responsibilities:

```
┌─────────────────┐       1..*       ┌───────────────────┐
│     Problem     │ ◄─────────────── │      Attempt      │
│  - ID           │                  │  - ID             │
│  - Title        │                  │  - Status         │
│  - Difficulty   │                  │  - Submission     │
│  - Skeleton     │                  │  - Result         │
└─────────────────┘                  └─────────┬─────────┘
                                               │
                                               │ Evaluates
                                               ▼
                                     ┌───────────────────┐
                                     │  Evaluator (IF)   │
                                     └─────────▲─────────┘
                                               │
                      ┌────────────────────────┴────────────────────────┐
                      │                                                 │
            ┌───────────────────┐                             ┌───────────────────┐
            │ DeterministicEval │                             │  CompositeEval    │
            │  (Rules & UML)    │                             │ (Chained Pipeline)│
            └───────────────────┘                             └───────────────────┘
```

### Key Entities:
1. **`Problem`**: Represents an LLD scenario (e.g. *Parking Lot*, *Movie Booking System*) with constraints, requirements, and starting skeleton code.
2. **`Submission`**: Encapsulates the learner's solution payload, supporting format abstraction (`FormatCodeAndNotes`, `FormatClassDiagram`, `FormatSpecification`).
3. **`Attempt`**: Models an attempt within the learning loop, owning an explicit lifecycle:
   $$\text{SUBMITTED} \longrightarrow \text{EVALUATING} \longrightarrow \text{COMPLETED} \quad (\text{or } \text{FAILED})$$
4. **`Evaluator` (Interface)**: The architectural contract decoupling evaluation logic from HTTP handlers and data access.
5. **`RubricDimension`**: Structured critique item adhering to the specification:
   $$\text{Criterion} \longrightarrow \text{Score} \longrightarrow \text{Evidence} \longrightarrow \text{Concern} \longrightarrow \text{Suggestion} \longrightarrow \text{Confidence}$$

---

## 3. Addressing the Assignment's Two Critical "Change Tests"

The design was tested against the two explicit change scenarios from the Candidate Helping Guide:

### Change Test A: Changing Submission Format
> *"Today the learner submits text/code. Later the platform supports a class diagram. How much of your domain model changes?"*

* **Design Solution:** The domain defines `SubmissionFormat` as an abstraction. The `Submission` entity holds:
  ```go
  type Submission struct {
      ProblemID   int              `json:"problem_id"`
      Format      SubmissionFormat `json:"format"`
      CodeContent string           `json:"code_content"`
      DesignNotes string           `json:"design_notes"`
      DiagramData string           `json:"diagram_data,omitempty"`
  }
  ```
* **Impact of Change:** Adding graphical diagram submissions requires **zero changes** to `Problem`, `Attempt`, or the database schema; only a new evaluator strategy or parser is registered to inspect `DiagramData`.

---

### Change Test B: Adding New Evaluation Strategies
> *"Today feedback comes from one evaluator. Later you add a rule-based evaluator, human review, or an LLM. Can you add it without rewriting the practice flow?"*

* **Design Solution:** Open/Closed Principle applied via the `Evaluator` interface and `CompositeEvaluator` pattern:
  ```go
  type Evaluator interface {
      Name() string
      Evaluate(ctx context.Context, sub Submission, prob Problem) (EvaluationResult, error)
  }
  ```
* **Impact of Change:** To introduce an OpenAI or Gemini evaluator:
  1. Implement `type AIEvaluator struct{}` satisfying `Evaluator`.
  2. Register it in `NewDefaultCompositeEvaluator()`.
  3. The HTTP handler in `main.go` and database code remain completely untouched.

---

## 4. Addressing Core Design Questions (Section 4)

### Q1: What does a learner actually need to provide for an attempt to be meaningful?
- **Classes & Interfaces (Code / Skeleton):** Concrete evidence of abstraction, encapsulation, and type hierarchy.
- **Architectural Notes:** Explains intent, non-obvious trade-offs (e.g. concurrency locks vs. optimistic locking), and pattern choices that cannot be inferred from syntax alone.

### Q2: What makes feedback useful when multiple valid solutions exist?
- **Multi-Dimensional Rubric:** Evaluating designs across standardized engineering dimensions (SRP, Loose Coupling, Concurrency, Extensibility) rather than matching against a single "golden" solution.
- **Evidence-Based Feedback:** Citing specific classes or methods the candidate created, explaining the structural hazard, and suggesting the exact pattern or refactoring to apply.

### Q3: What is deterministic vs. what benefits from an LLM?
- **Deterministic:** Domain keyword validation, syntax checking, anti-pattern detection (god classes, circular dependencies), and AST-based class extraction for dynamic UML generation.
- **LLM-Assisted:** Qualitative evaluation of trade-offs, naming clarity, domain nuance, and edge-case scenario reasoning.

### Q4: Handling slow or failing evaluations (Scale & Reliability)
- **Pre-Evaluation Persistence:** The submission is immediately written to the database with `status = 'EVALUATING'` before the evaluator is invoked. If an evaluator crashes, times out, or runs out of tokens, candidate code is never lost.
- **Context Timeouts:** The evaluation pipeline is protected by Go's `context.WithTimeout(ctx, 10*time.Second)` to guarantee the server never hangs.
- **State Transition:** Status updates to `COMPLETED` on success or `FAILED` with retry capability on error.

---

## 5. Key Engineering Trade-offs

| Decision | Chosen Approach | Alternative Considered | Rationale |
| :--- | :--- | :--- | :--- |
| **Architecture** | Go Monolith (Fiber) | Microservices (API + Worker) | For a 2-day MVP, a modular monolith minimizes operational complexity and deployment friction while offering microsecond-level performance. |
| **Storage** | SQLite with WAL mode | PostgreSQL / MongoDB | SQLite provides zero-dependency local persistence, atomic transactions, and effortless portability. |
| **UML Generation** | Dynamic client-side Mermaid.js | Server-side Graphviz/PlantUML | Zero external binary dependencies; parses class hierarchies dynamically on the frontend directly from code. |
