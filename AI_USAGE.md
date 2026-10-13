# AI_USAGE.md

This document outlines key AI-assisted architectural decisions made during the development of the LLD Practice Platform MVP:

1. **Modular Monolithic Architecture & Tech Stack Selection**
   * *AI Suggestion:* Recommended using Go with the Fiber web framework, structured with a clean separation of concerns across dedicated modules (`main.go`, `evaluator.go`, `database.go`, `models.go`).
   * *Decision:* Accepted. This ensures high maintainability and code readability without the unnecessary complexity or overhead of distributed microservices.

2. **Strict One-to-One Hybrid Domain Evaluation Engine**
   * *AI Suggestion:* Proposed a dual-layer validation mechanism that independently scans both code implementations and architectural notes for problem-specific domain keywords.
   * *Decision:* Accepted. This completely eliminates cross-problem contamination (e.g., mixing a URL Shortener implementation with Splitwise notes) by hard-enforcing a strict one-to-one mapping with zero domain overlap and instant zero-score rejection on mismatches.

3. **Decoupled Multi-File Frontend Design**
   * *AI Suggestion:* Suggested splitting the frontend into standalone static files (`index.html`, `style.css`, `app.js`) served via Fiber's static file handler, integrated with Tailwind CSS (CDN) and Mermaid.js for dynamic UML class diagram generation.
   * *Decision:* Accepted to provide a polished, production-grade user experience while maintaining rapid development velocity.
