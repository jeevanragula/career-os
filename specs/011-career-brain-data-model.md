# Specification 011 — Career Brain Data Model

- Status: Accepted
- Version: 0.1
- Milestone: M1 Career Brain

## Purpose

Define the durable domain model for canonical career facts, evidence, verification, and publication safety.

## Design principles

1. Canonical facts are independent of resume wording.
2. Every material claim has verification and confidentiality metadata.
3. Evidence is immutable by default; corrections create new evidence versions.
4. AI output is an analysis artifact, never canonical evidence.
5. Generated artifacts reference claims rather than copying facts into a second source of truth.

## Core entities

Person, Employment, Project, Achievement, Skill, Patent/Publication, Education, Evidence, and Claim.

### Evidence fields

- evidence_id
- source_type
- source_reference
- captured_at
- confidentiality
- integrity_status

Allowed source types include user_statement, resume, linkedin, portfolio, public_documentation, public_article, patent_record, employment_record, and generated_analysis.

A generated analysis cannot independently verify a canonical claim.

### Claim fields

- claim_id
- text
- entity_type
- entity_ref
- verification_status
- confidence
- confidentiality
- evidence_refs
- created_at
- updated_at

## Verification states

- needs_verification
- user_asserted
- verified
- rejected

A claim cannot become verified solely because an AI model generated or repeated it.

## Confidentiality levels

- public
- application_safe
- internal
- confidential

Lowering confidentiality requires human review.

## Generated artifact linkage

Resume bullets, cover notes, outreach messages, and application answers should retain machine-readable references to their source claims.

## Acceptance criteria

- [x] Canonical entities are defined.
- [x] Claim and evidence are separate.
- [x] Verification states are explicit.
- [x] Confidentiality is explicit.
- [x] AI output cannot independently verify claims.
- [x] Generated artifacts can reference claims.
- [x] Evidence can be audited.
