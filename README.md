# CareerOS

CareerOS is a spec-driven, AI-assisted career operating system for senior engineering job search.

**Core principle: evidence first, automation second.**

## Current implementation

- Career Brain with claims, evidence, verification, confidentiality, and validation
- Canonical job ingestion with provenance, normalization, versioning, and deduplication
- Compliant public Lever and Ashby job-posting adapters
- Discovery source/query/run persistence and workflow queue foundation
- AI job-analysis and evidence-matching contracts
- Evidence-grounded resume variant model
- Startup radar, application, outreach, timeline, and interview foundations
- Exact-payload human approval boundary for external actions
- Runnable Go HTTP service with health/readiness endpoints

## SDD workflow

~~~text
Problem -> Specification -> Acceptance Criteria -> Architecture
       -> Implementation -> Tests -> Review -> Merge
~~~

AI may accelerate research, analysis, drafting, and workflow execution, but external actions remain human-approved.

## Roadmap

- M0 Foundation — complete
- M1 Career Brain — complete
- M2 Job Discovery — runnable public-provider discovery + persistence
- M3 Job Analysis — runnable AI analysis path
- M4 Resume Engine — runnable AI Markdown resume generation
- M5 Startup Radar — persistence foundation
- M6 Application Assistant — runnable preparation workflow
- M7 Outreach — workflow foundation; sending intentionally gated
- M8 Interview OS — workflow foundation; practice UI next

See `docs/ROADMAP.md` and `docs/DEPLOYMENT.md` for the implementation and run instructions.\n\n## Runnable MVP\n\nThe current MVP can discover public Lever/Ashby postings, persist jobs in PostgreSQL, run evidence-backed AI job analysis, generate a tailored Markdown resume, and create an application preparation record. External submission and outreach remain human-approved and are not automated.

## Repository structure

~~~text
specs/              Product and capability specifications
docs/architecture/  System architecture
docs/decisions/     Architecture decision records
career-profile/     Canonical career evidence
resumes/            Resume strategies and artifacts
agents/             Agent specifications and implementations
application/        Application workflow assets
infrastructure/     Runtime and deployment assets
internal/            Go domain, provider, workflow, and API packages
cmd/                Executable entry points
~~~

## Safety and privacy

CareerOS must never fabricate career facts, metrics, employers, technologies, responsibilities, patents, achievements, or qualifications.

Do not store credentials, session cookies, API tokens, passwords, employer/customer confidential information, proprietary source code, or unsafe non-public architecture details.

Generated career artifacts must trace back to canonical evidence.
