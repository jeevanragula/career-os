# ADR-005 — Job Source Provenance and Deterministic Deduplication

## Decision

CareerOS will retain source observations separately from canonical jobs.

Deduplication is deterministic and evidence-preserving:
1. canonicalized source URL
2. source + stable external job ID
3. normalized company/title/location/remote fallback

AI-based similarity may propose possible duplicates later, but it is not the primary identity mechanism.

## Consequences

- Multiple providers can represent the same job.
- Historical source observations remain auditable.
- Changed job descriptions become versions.
- Downstream AI analysis operates on a stable canonical job.
