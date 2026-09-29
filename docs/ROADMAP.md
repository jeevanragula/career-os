# CareerOS Implementation Roadmap

## Completed foundation

- M0: repository, product vision, architecture, SDD, privacy and approval principles.
- M1: Career Brain schema, evidence model, claims, validation, public evidence reconciliation.
- M2: canonical job model, source observations, deterministic normalization/deduplication, source configuration, discovery runs, provider contract, compliant provider notes.
- M3: job-analysis contract, requirement extraction model, evidence-backed requirement matching.
- M4: resume variant contract and claim traceability.
- M5: startup radar contract and provenance model.
- M6: application workflow schema with explicit approval boundary.
- M7: outreach draft/send lifecycle and application timeline.
- M8: interview plans, questions, practice sessions, and learning-loop operations.

## Next implementation wave

1. Add concrete Lever and Ashby public-posting adapters using fixture tests.
2. Add persistence/repository services around PostgreSQL migrations.
3. Add queue worker and scheduler with bounded retries and idempotency.
4. Add LLM provider abstraction and structured JSON outputs for job analysis.
5. Add resume rendering and artifact storage.
6. Add application-form adapters only where the provider permits and only behind exact human approval.
7. Add company research and startup-source adapters.
8. Add dashboard/API and authentication.
9. Add observability, metrics, tracing, and operational alerts.

The system is intentionally designed so provider integrations and model choices can evolve without changing the canonical Career Brain.
