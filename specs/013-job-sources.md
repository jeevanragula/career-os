# Spec 013 — Job Sources and Discovery Runs

## Goal
Define compliant, auditable job-source configuration and discovery execution without coupling CareerOS to a single job board.

## Requirements
1. A source declares provider type, endpoint, authorization mode, and access constraints.
2. Discovery queries are explicit and versioned.
3. Each discovery run records timestamps, source, query, counts, and failure state.
4. Adapters return source observations and never invent missing fields.
5. Adapters use documented/public or explicitly authorized interfaces.
6. No CAPTCHA bypass, credential sharing, robots/terms bypass, or stealth scraping.
7. Discovery is idempotent.
8. Closed/expired observations remain historical.
9. Source failures are isolated and retryable.
10. Human approval is required before external application/outreach actions.

## Acceptance criteria
- Sources can be enabled/disabled without code changes.
- Runs are auditable from source, query, timestamps, and outcome.
- Adapter failures do not corrupt prior jobs.
- Every observation has provenance.
- Sources can be independently rate-limited or paused.
