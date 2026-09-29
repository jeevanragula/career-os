# Spec 022 — Operations and Learning Loop

## Goal
Make CareerOS continuously useful while keeping automation bounded and auditable.

## Requirements
1. Queue work by durable entity/version IDs.
2. Make workers idempotent and retryable.
3. Capture model, prompt, source, and timing metadata for AI jobs.
4. Never retry an external action unless an approval is still valid for the exact payload hash.
5. Track application outcomes and use them as learning signals.
6. Learning produces proposed strategy changes, not silent changes to canonical claims.
7. Provide health metrics for discovery freshness, analysis backlog, application state, and failed jobs.

## Acceptance criteria
- Duplicate queue requests coalesce for active work.
- External actions require exact approval binding.
- Outcome analysis cannot mutate Career Brain verification without a separate evidence workflow.
