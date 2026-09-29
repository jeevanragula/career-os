# Spec 016 — Application Assistant

## Goal
Prepare complete application packages while keeping submission under explicit human control.

## Requirements
1. Package includes selected resume, job analysis, application answers, and evidence references.
2. AI may draft answers from verified/application-safe claims.
3. Missing required information is surfaced as a blocking question.
4. Submission is a separate state transition requiring explicit user approval.
5. Credentials and session cookies are never stored in the Career Brain.
6. Every submission attempt is auditable.

## Acceptance criteria
- No external submit occurs from a draft or prepared state.
- Approval records actor and timestamp.
- Rejected or withdrawn applications remain historical.
