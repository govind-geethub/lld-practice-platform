# AI_USAGE.md

This document outlines key AI-assisted architectural decisions made during the design and development of the **LLD Practice Platform**, explicitly documenting what the AI suggested, what was accepted or rejected, and the engineering rationale.

---

### Decision 1: Evaluation Architecture — Fixed Multi-Dimensional Rubric vs. Freeform LLM Prompt

* **AI Suggestion:** Initially suggested delegating the entire evaluation to an external LLM with a generic prompt: *"Evaluate this LLD solution on a scale of 1-100 and give feedback on how good the design is."*
* **Our Decision: REJECTED.**
* **Engineering Rationale:**
  - An unconstrained LLM prompt violates the core principles outlined in the Candidate Helping Guide. It produces non-deterministic scores, suffers from hallucinations, and fails to cite concrete code evidence.
  - Furthermore, relying solely on an external API introduces external latency, rate limiting, and failure risks during interview demos.
* **Adopted Alternative:**
  - Implemented a **Hybrid Pipeline Architecture**: A deterministic pre-validator that rejects cross-problem contamination and god-class patterns, coupled with a **formal multi-dimensional rubric** adhering to the Helping Guide's exact shape:
    $$\text{Criterion} \longrightarrow \text{Score (1-5)} \longrightarrow \text{Evidence} \longrightarrow \text{Concern} \longrightarrow \text{Suggestion} \longrightarrow \text{Confidence}$$

---

### Decision 2: Extensibility Architecture — Monolithic Function vs. Pluggable `Evaluator` Interface

* **AI Suggestion:** Suggested writing a single monolithic evaluation function inside `main.go` using regular expressions for quick prototyping.
* **Our Decision: REJECTED.**
* **Engineering Rationale:**
  - Directly violates **Change Test B** from the assignment: *"Today feedback comes from one evaluator. Later you add a rule-based evaluator or human review. Can you add it without rewriting the practice flow?"*
  - A monolithic function tightly couples HTTP request handling to evaluation rules, making the codebase unmaintainable and failing LLD design grading.
* **Adopted Alternative:**
  - Defined the `Evaluator` interface and a `CompositeEvaluator` (Go implementation of the Strategy and Composite patterns).
  - Chaining a new evaluator (e.g. AST parser or OpenAI evaluator) requires zero modifications to HTTP handlers or database persistence.

---

### Decision 3: Resilience & State Lifecycle — Synchronous Execution vs. Pre-Evaluation Persistence

* **AI Suggestion:** Suggested executing evaluation synchronously in the POST handler and returning the JSON result immediately without recording intermediate status.
* **Our Decision: REJECTED.**
* **Engineering Rationale:**
  - If the evaluation pipeline times out or crashes, the learner's submitted code is permanently lost.
  - Section 10 of the Candidate Helping Guide explicitly requires:
    > *"Store the submission before evaluation starts so it is not lost on evaluator failure. Give evaluation a clear state such as Submitted → Evaluating → Completed / Failed."*
* **Adopted Alternative:**
  - Modeled an explicit state machine: `SUBMITTED` $\rightarrow$ `EVALUATING` $\rightarrow$ `COMPLETED` / `FAILED`.
  - Submissions are immediately saved to SQLite with status `EVALUATING` before the evaluation pipeline runs under a 10-second Go `context.WithTimeout`.

---

### Decision 4: Dynamic Architecture Diagramming — Static Chart vs. AST/Regex Extraction

* **AI Suggestion:** Suggested returning a static, pre-rendered Mermaid diagram of a generic "Manager -> Strategy" pattern for every problem.
* **Our Decision: REJECTED.**
* **Engineering Rationale:**
  - A static diagram does not provide personalized feedback; learners cannot see their own classes and relationships visualized.
* **Adopted Alternative:**
  - Implemented `generateDynamicUML()`, a lightweight parser that extracts actual `class`, `interface`, and inheritance tokens from the student's submission and renders **their actual architecture** live in Mermaid.js.

---

### Decision 5: Tech Stack Selection — Lightweight Monolith vs. Distributed Microservices

* **AI Suggestion:** Suggested spinning up a Go API gateway, Python evaluation worker service, Redis message broker, and PostgreSQL database.
* **Our Decision: REJECTED.**
* **Engineering Rationale:**
  - The assignment brief states: *"This is primarily an LLD/domain-design exercise. Do not spend the majority of your time on Kubernetes, microservices, multi-region deployment, sharding... A simple monolith is completely acceptable."*
  - Building distributed microservices would distract from the core domain model and introduce needless failure modes.
* **Adopted Alternative:**
  - Built a clean, modular Go monolith using Fiber and SQLite with embedded static assets. It boots in milliseconds, has zero external server dependencies, and showcases clean object-oriented domain boundaries.
