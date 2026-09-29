# CareerOS

CareerOS is a spec-driven, AI-assisted career operating system for senior engineering job search.

**Core principle: evidence first, automation second.**

## Vision

Maintain a verified career knowledge base and use it to discover opportunities, research companies, analyze job fit, select and tailor resumes, prepare applications, support personalized outreach, track applications, and prepare for interviews.

## SDD workflow

~~~text
Problem -> Specification -> Acceptance Criteria -> Architecture
       -> Implementation -> Tests -> Review -> Merge
~~~

AI may accelerate research, analysis, drafting, and workflow execution, but external actions remain human-approved.

## Roadmap

- M0 Foundation
- M1 Career Brain
- M2 Job Discovery
- M3 Job Analysis
- M4 Resume Engine
- M5 Startup Radar
- M6 Application Assistant
- M7 Outreach
- M8 Interview OS

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
~~~

## Safety and privacy

CareerOS must never fabricate career facts, metrics, employers, technologies, responsibilities, patents, achievements, or qualifications.

Do not store credentials, session cookies, API tokens, passwords, employer/customer confidential information, proprietary source code, or unsafe non-public architecture details.

Generated career artifacts must trace back to canonical evidence.
