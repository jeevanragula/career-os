# ADR-004 — Human Approval Boundary

- Status: Accepted
- Date: 2026-09-30

## Decision

CareerOS separates preparation from external action.

AI and automation may research, analyze, draft, validate, and prepare an external workflow.

The candidate must explicitly approve before:
- submitting an application
- sending outreach
- publishing a career claim
- accepting a material canonical-profile change
- exposing information above application_safe classification

## Rationale

The system operates on career identity and can create irreversible external side effects. Human approval preserves control and provides a clear accountability boundary.

## Consequence

The workflow engine must represent approval as a durable state/event rather than as an implicit UI interaction.
