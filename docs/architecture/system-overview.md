# CareerOS System Overview

- Status: Proposed
- Version: 0.1
- Milestone: M0 Foundation

## Logical architecture

~~~text
CareerOS UI
   |
Workflow Layer
   |---- Career Brain
   |---- Opportunity Store
   |---- Application Store
   |
AI / Agent Layer
   |---- Web / Search
   |---- LLM Providers
   |---- Browser Agent
                 |
          Human Approval
                 |
        External Websites
~~~

## Initial technology direction

Favor PostgreSQL for durable structured state, object storage for documents and artifacts, an orchestration layer for workflows, provider adapters for LLM/search/browser/notification integrations, browser automation only where needed, and deterministic validation around AI outputs.

These are directions, not final implementation commitments.

## Data domains

- Career: facts, claims, evidence, skills, projects, publications
- Opportunity: companies, jobs, observations, analysis
- Application: applications, artifacts, approvals, submissions
- Communication: outreach targets, drafts, sent messages, responses
- Interview: preparation, questions, feedback, learning
- AI execution: prompts, model metadata, tool calls, evidence, outputs, validation, audit

## Trust boundaries

1. Canonical evidence
2. External evidence
3. AI interpretation
4. Generated artifact
5. External action

An artifact cannot cross a trust boundary without required validation and/or approval.

## Approval gates

Human approval is required before submitting an application, sending outreach, publishing a career claim externally, accepting material changes to canonical facts, or exposing information above application-safe classification.

## Reliability

Workflow steps should be idempotent where practical. External side effects need explicit state transitions so retries cannot silently duplicate applications, messages, or submissions.

## Security

Never store secrets in Career Brain. Use external secret management, least-privilege integrations, encryption in transit and at rest, separation of internal/application-safe evidence, and an audit trail for external actions.

## Future implementation candidates

Potential components include Go services, Python workers, n8n, PostgreSQL/Supabase, object storage, Playwright, and pluggable LLM providers. No candidate is mandatory until selected by specification or ADR.
