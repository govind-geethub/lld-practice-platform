### Deliverable 6: AI_USAGE.md

```markdown
# AI_USAGE.md

This document outlines key AI-assisted decisions made during the development of the LLD Practice Platform MVP:

1. **Monolithic Architecture & Tech Stack Selection**
   * *AI Suggestion:* Suggested Go with Fiber and SQLite for an instant zero-dependency monolith.
   * *Decision:* Accepted. Avoided unnecessary microservices or Kubernetes overhead as per assignment scope boundaries[cite: 2].
2. **Hybrid Evaluation Engine Design**
   * *AI Suggestion:* Suggested combining deterministic AST/regex pattern matching with structured rubrics for feedback.
   * *Decision:* Accepted. This ensures fast, reliable structural feedback while allowing qualitative assessment notes.
3. **Embedded Tailwind CSS Frontend**
   * *AI Suggestion:* Recommended using Tailwind via CDN for a modern dashboard without heavy frontend build tools (Webpack/Vite).
   * *Decision:* Accepted to maintain rapid iteration speed within the 2-day timeline.