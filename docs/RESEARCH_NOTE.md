# Research Note: Low-Level Design (LLD) Practice & Evaluation Platform

**Candidate:** Govind Kumar  
**Date:** October 2026  
**Assignment:** 2-Day LLD Practice Platform Prototype  

---

## 1. The Core Learner Problem

Practicing Data Structures & Algorithms (DSA) has a clear, deterministic loop: a learner writes code against input tests, receives instantaneous green/red execution results, and optimizes time/space complexity.

In contrast, **practicing Low-Level Design (LLD) is fraught with ambiguity**:
1. **Subjective & Multi-Faceted Solutions:** Unlike algorithmic problems with unique deterministic outputs, an LLD problem (such as *Parking Lot*, *Elevator System*, or *Movie Booking*) has dozens of valid designs. A candidate can implement an Elevator using a State Pattern or a Command Pattern; both can be effective depending on trade-offs.
2. **Evaluation Void:** Candidates frequently write classes, interfaces, and design notes, but remain fundamentally unsure whether:
   * Class responsibilities adhere to the Single Responsibility Principle (SRP) or exhibit "God Class" anti-patterns.
   * Coupling between abstractions and concrete implementations is sufficiently loose (Dependency Inversion Principle).
   * Concurrency, thread-safety, and state mutation were handled soundly.
3. **Absence of an Iterative Learning Loop:** Existing resources treat LLD as static reading. Candidates read a pre-baked reference solution, attempt to memorize class diagrams, and fail to develop the iterative architectural muscle needed to refine a design based on targeted critique.

---

## 2. Competitive Landscape & Gaps in Existing Tools

| Platform / Tool | Typical Workflow | Feedback Provided | Critical Gap / Failure Mode |
| :--- | :--- | :--- | :--- |
| **LeetCode / HackerRank** | Run code against automated unit test suites. | Pass/Fail boolean and runtime performance (ms). | Completely unsuitable for LLD. Evaluates only algorithmic outputs, ignoring class structure, interfaces, coupling, or design patterns. |
| **Educative / Grokking LLD** | Static reading of reference class diagrams and code solutions. | Static text and pre-drawn diagrams. | Zero interactive practice. The learner cannot submit their own variations or receive tailored feedback on their trade-offs. |
| **Generic LLMs (ChatGPT / Claude)** | Freeform chat prompt with pasted code or notes. | Unconstrained text critique and arbitrary scores (e.g. "8/10"). | **High inconsistency and hallucinations**. Lacks a standardized rubric; provides ungrounded scores without citing concrete evidence or pointing out domain mismatch. |
| **GitHub Open-Source Repositories** | Browse curated implementations of classic LLD problems. | None (static reference). | Encourages rote memorization rather than active problem-solving or iterative improvement. |

---

## 3. Key Product Insights & Gaps Identified

From our research, three fundamental gaps must be bridged for an LLD practice tool to be effective:
1. **The need for structured, rubric-based evaluation:** Feedback cannot be a single arbitrary number or an unconstrained "Is this good?" LLM response. Useful evaluation must dissect designs across standard engineering dimensions: **Single Responsibility & Cohesion**, **Abstraction & Loose Coupling (DIP/ISP)**, **Design Patterns & Extensibility**, and **Concurrency & State Safety**.
2. **Evidence-Based Explanations:** For feedback to be actionable, every critique must cite **concrete evidence** directly from the candidate's submission, explain the underlying architectural concern, and propose a specific next step for their subsequent attempt.
3. **Progressive Learning Loop (Attempt History):** The platform must retain attempt history so learners can compare Attempt #1 to Attempt #2, observing how extracting an interface or adding synchronization directly improves their design score.

---

## 4. Product Direction & MVP Scope

Our prototype focuses strictly on the **Learner Journey**:
$$\text{Choose Problem} \longrightarrow \text{Think \& Design} \longrightarrow \text{Submit} \longrightarrow \text{Receive Rubric Feedback} \longrightarrow \text{Review Evolution} \longrightarrow \text{Iterate}$$

### Key Architectural Choices:
- **Clean Go Monolith:** High-performance, zero-overhead modular architecture using Go and Fiber, avoiding unnecessary microservice complexity.
- **Hybrid Decoupled Evaluation Pipeline:** 
  - *Deterministic Rule Engine:* Instantly rejects cross-problem contamination, validates required domain models, and parses dynamic Mermaid UML diagrams directly from student code.
  - *Pluggable Evaluator Interface:* Decouples evaluation logic from HTTP handlers, enabling seamless extension to external LLMs or human reviewers.
- **Pre-Evaluation State Persistence:** Submissions are persisted immediately with status `EVALUATING` before computation begins, ensuring zero data loss during timeouts or pipeline failures.
