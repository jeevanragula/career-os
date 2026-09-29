# Spec 012 — Job Ingestion

## Goal

Create a normalized, evidence-preserving job ingestion layer that can accept jobs from multiple sources without coupling CareerOS to a single provider.

## Requirements

1. Store the original source URL and source identifier.
2. Normalize title, company, location, employment type, remote mode, description and timestamps.
3. Preserve raw source metadata for later inspection.
4. Deduplicate jobs across sources.
5. Keep source provenance when duplicate records are merged.
6. Never infer missing facts as if they came from the source.
7. Support jobs with incomplete fields.
8. Detect materially changed job postings without creating unnecessary duplicates.
9. Allow a job to be marked closed/expired without deleting its historical record.
10. Do not scrape in violation of source terms, authentication boundaries, or access controls.

## Deduplication

Primary identity:
- canonicalized source URL when stable.

Secondary identity:
- normalized company
- normalized title
- normalized location/remote mode
- stable external job identifier when available.

When sources disagree, preserve both source observations and do not silently overwrite the original evidence.

## Acceptance criteria

- The same source URL produces one canonical job.
- Equivalent URLs differing only by tracking parameters deduplicate.
- Stable external IDs deduplicate within a source.
- Different companies with identical titles remain separate.
- Job description changes create a new observation/version, not a new job.
- Every canonical job retains at least one provenance record.
