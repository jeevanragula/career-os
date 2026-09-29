# Spec 015 — Evidence-Grounded Resume Engine

## Goal
Generate targeted resume variants from canonical claims without creating unsupported career facts.

## Requirements
1. A resume variant references a target job and analyzed requirements.
2. Every substantive bullet maps to one or more canonical claims.
3. Only application-safe claims may be rendered.
4. User-approved wording may be retained as an artifact but does not become evidence automatically.
5. Resume generation must preserve chronology and employer truth.
6. Tailoring may reorder or emphasize evidence but may not fabricate metrics, technologies, ownership, or outcomes.
7. Each variant records model and prompt versions.

## Acceptance criteria
- Unsupported bullets are rejected by validation.
- A variant can be traced to claim IDs and job requirements.
- Multiple variants share the same canonical career facts.
