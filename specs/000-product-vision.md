# Specification 000 — Product Vision

- Status: Accepted
- Version: 0.1
- Milestone: M0 Foundation

## 1. Problem

Senior engineering job search is fragmented across job boards, company career sites, startup databases, professional networks, resumes, spreadsheets, outreach channels, and interview preparation materials.

The result is repeated manual research, inconsistent career narratives, weak traceability between claims and evidence, and time lost on low-fit opportunities.

## 2. Product goal

CareerOS provides one system of record for verified career evidence and an AI-assisted workflow for the complete senior-engineering job-search lifecycle.

The system should improve search quality and execution speed without turning the candidate into a passive approver of opaque automation.

## 3. Primary user

The initial user is a senior/principal engineering candidate with deep experience in cloud-native distributed systems, Kubernetes and public cloud, cybersecurity and cloud security, data security and data platforms, AI platforms and AI security, and architecture/technical leadership.

The model must not assume any particular employment history, metric, title, or technology unless represented in verified career evidence.

## 4. Product principles

### P1 — Evidence before generation
Every externally useful career claim should trace to canonical evidence.

### P2 — Human approval for external actions
Submitting applications, sending outreach, publishing career claims, or making material changes to canonical facts requires explicit human approval.

### P3 — Facts and interpretation are separate
Documented facts, AI-derived analysis, predictions, recommendations, and user decisions must be represented separately.

### P4 — Reproducibility
Important AI outputs should record inputs, specification/prompt version, model/provider, and relevant evidence references.

### P5 — Provider independence
Web search, LLM, browser automation, storage, and orchestration providers should be replaceable behind interfaces.

### P6 — Minimize proprietary information
CareerOS should retain enough evidence to support accurate job search while excluding employer/customer confidential material.

### P7 — Deterministic controls around AI
AI-generated artifacts must pass deterministic validation before becoming eligible for external use.

## 5. Capabilities

- C1 Career Brain
- C2 Opportunity Discovery
- C3 Job and Company Intelligence
- C4 Fit Analysis
- C5 Resume Engine
- C6 Application Assistant
- C7 Outreach
- C8 Interview OS
- C9 Learning Loop

## 6. Non-goals

CareerOS will not:
- perform fully autonomous mass applications
- fabricate credentials, achievements, experience, or metrics
- submit an application without human approval
- send outreach without human approval
- scrape sources in violation of applicable terms or access controls
- store employer/customer confidential information as career evidence
- replace the candidate's judgment about career decisions

## 7. Functional requirements

- FR-1 Canonical profile
- FR-2 Evidence traceability
- FR-3 Job normalization
- FR-4 Deduplication
- FR-5 Analysis provenance
- FR-6 Human approval
- FR-7 Fact validation
- FR-8 Privacy classification
- FR-9 Versioned specifications
- FR-10 Replaceable integrations

## 8. Non-functional requirements

Security, auditability, idempotent ingestion, reliable bounded retries, structured observability, AI cost controls, explicit retention policies, and testability of deterministic workflow components.

## 9. MVP boundary

1. Career Brain
2. Job ingestion
3. Job normalization and deduplication
4. Job analysis
5. Resume strategy selection
6. Human-approved application preparation
7. Application tracking

## 10. M0 acceptance criteria

- [x] Product scope documented
- [x] Principles documented
- [x] MVP boundaries documented
- [x] Capabilities identified
- [x] Privacy/confidentiality explicit
- [x] Human approval explicit
- [x] SDD workflow documented
- [x] Initial architecture direction documented
- [x] Career Brain data model initiated
