# Job Ingestion Architecture

Job discovery is intentionally split into four layers:

1. Source adapters
2. Normalization
3. Identity and deduplication
4. Career Brain analysis

Source adapters convert an allowed source response into a source observation. Normalization canonicalizes fields while preserving original values. Identity and deduplication determine whether an observation belongs to an existing canonical job. JD parsing, fit analysis and resume selection belong to M3.

CareerOS should prefer official company career pages, authorized job APIs, user-provided feeds, and permitted aggregators/connectors.

It must not bypass authentication, robots/access controls, CAPTCHAs, or contractual restrictions.
