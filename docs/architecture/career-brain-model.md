# Career Brain Domain Model

## Relationships

~~~text
Person
  |
  +-- Employment
  |     |
  |     +-- Project
  |            |
  |            +-- Achievement
  |
  +-- Skill
  +-- Patent / Publication
  +-- Education

Claim
  |
  +-- references -> Person / Employment / Project / Achievement / Skill / Patent / Education
  +-- supported_by -> Evidence

Generated Artifact
  |
  +-- derived_from -> Claim
~~~

## Domain invariants

1. A claim without evidence may be needs_verification but cannot be verified.
2. generated_analysis is never sufficient evidence for verification.
3. Confidential claims cannot appear in application-safe artifacts.
4. Generated artifacts retain source claim references.
5. Removing evidence does not silently delete a claim; it makes the claim needs_verification.
6. Changing a verified claim creates an auditable change event.
7. Resume variants share canonical claims.

## Suggested relational tables

- people
- employments
- projects
- achievements
- skills
- claims
- evidence
- claim_evidence
- project_skills
- achievement_claims
- artifact_claims
- career_change_events

Large source documents belong in object storage; the database stores metadata and references.
