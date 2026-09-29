# Agent Orchestration Architecture

Event-driven flow:

Discovery -> Job Version -> Analysis -> Fit Evidence -> Resume Draft -> Application Package -> Human Approval -> Submission -> Timeline -> Interview -> Outcome -> Learning

The Career Brain is the authoritative evidence store. Agents read canonical claims and write proposed artifacts, analysis, and workflow events. Only explicit verification workflows can change claim verification state.

## Execution model
- Stateless agent workers where practical.
- PostgreSQL is the workflow/source-of-truth store.
- Object storage can hold generated resume artifacts.
- A queue separates discovery, analysis, generation, and external-action workers.
- Every job is idempotent using deterministic entity/version identifiers.

## Human approval boundary
The system may automate preparation. Submission and outbound communication remain blocked until an approval event exists for that exact action payload.
