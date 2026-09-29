# Specification 001 — Career Profile / Career Brain

- Status: Draft
- Version: 0.1
- Milestone: M1 Career Brain

## Purpose

Career Brain is the canonical professional knowledge base used by all downstream workflows.

Generated resumes, messages, analyses, and application answers are derived artifacts. They never become the source of truth.

## Core entities

- Person
- Employment
- Project / Product
- Achievement
- Skill
- Patent / Publication
- Education

## Evidence model

Each material claim should support:

| Field | Purpose |
|---|---|
| claim_id | Stable identifier |
| claim | Human-readable statement |
| source_type | Resume, profile, user confirmation, public source, etc. |
| source_reference | Pointer to source |
| verification_status | Verification state |
| confidence | Evidence confidence |
| confidentiality | Publication safety |
| last_verified_at | Verification timestamp |

## Verification states

- verified
- user_asserted
- needs_verification
- rejected

## Confidentiality

- public
- application_safe
- internal
- confidential

## Claim rules

Allowed: rephrase verified claims, condense them, reorder them by relevance, and convert verified work into concise resume language.

Requires review: combining facts, deriving metrics, inferring seniority/scope/impact, or qualitative impact interpretation.

Not allowed: inventing metrics, technologies, responsibilities, ownership, or exposing confidential information.

## Acceptance criteria

- [ ] Canonical person profile exists.
- [ ] Employment history is independent of resume wording.
- [ ] Projects can contain multiple achievements and evidence references.
- [ ] Claims have verification and confidentiality state.
- [ ] Generated artifacts trace claims to evidence.
- [ ] AI-generated claims never silently become verified.
- [ ] Resume strategies reuse canonical evidence without duplicating facts.
