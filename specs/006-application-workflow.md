# Specification 006 — Application Workflow

- Status: Draft
- Version: 0.1
- Milestone: M6

## Purpose

Prepare applications efficiently while keeping the candidate in control of external submission.

## States

~~~text
discovered -> analyzed -> shortlisted -> preparing -> awaiting_approval
-> submitted -> interview -> offer
                     -> rejected
                     -> withdrawn
~~~

## Approval gate

No workflow may transition to submitted without explicit human approval.

Browser automation may prepare fields, but final submission remains gated.

## Application package

- selected resume
- cover note when appropriate
- application answers
- job analysis
- company research
- outreach draft
- evidence references
- approval record

## Idempotency

Retries must not create duplicate applications or submissions.

## Acceptance criteria

- [ ] State transitions are explicit.
- [ ] Submission requires approval.
- [ ] Submitted artifacts are versioned.
- [ ] Application events are auditable.
- [ ] Retries are safe.
