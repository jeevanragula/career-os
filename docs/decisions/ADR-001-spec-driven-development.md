# ADR-001 — Specification-Driven Development

- Status: Accepted
- Date: 2026-09-30

## Context

CareerOS combines AI agents, external integrations, automation, sensitive career data, and externally visible actions. Without explicit specifications, behavior can drift as prompts, models, integrations, and implementations change.

## Decision

CareerOS uses Specification-Driven Development (SDD):

~~~text
Problem
  -> Specification
  -> Acceptance Criteria
  -> Architecture / Design
  -> Implementation
  -> Tests
  -> Review
  -> Merge
~~~

## Rules

1. Major features do not start with implementation alone.
2. Specifications contain explicit acceptance criteria.
3. Implementations are traceable to specifications.
4. Tests verify acceptance criteria where applicable.
5. Material behavior changes require a specification update.
6. Prompts that materially affect behavior are versioned.
7. External actions require explicit human approval.
8. Generated career content is validated against canonical evidence.
9. Material trust/data/side-effect architecture changes should have an ADR.
10. Experimental spikes cannot become production behavior until incorporated into an accepted specification.

## Consequences

Benefits include reduced behavior drift, clearer review, better auditability, provider replacement, and safer external automation.

Trade-offs include additional documentation and maintenance of specifications as behavior evolves.
