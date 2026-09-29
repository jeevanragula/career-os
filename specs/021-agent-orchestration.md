# Spec 021 — CareerOS Agent Orchestration

## Goal
Coordinate specialized AI agents around one canonical Career Brain and auditable workflow state.

## Agents
- Discovery Agent: finds and normalizes permitted job observations.
- Company Research Agent: gathers public company facts and source provenance.
- Job Analysis Agent: extracts requirements and maps them to claims.
- Resume Agent: drafts evidence-grounded variants.
- Application Agent: prepares forms and answers but cannot submit without approval.
- Outreach Agent: drafts personalized messages but cannot send without approval.
- Interview Agent: generates preparation plans and practice sessions.
- Learning Agent: analyzes outcomes and proposes changes to search/resume strategy.

## Shared rules
1. Agents do not write directly to canonical verified claims.
2. Every artifact has source references and model/prompt versions.
3. External actions are explicit state transitions requiring approval.
4. Agents cannot bypass provider access controls.
5. Agent output is not evidence by itself.
6. Conflicts are surfaced for human review rather than silently resolved.
