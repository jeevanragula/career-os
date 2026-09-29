# Spec 014 — AI Job Analysis

## Goal
Turn an ingested job into structured requirements and an evidence-backed fit analysis.

## Requirements
1. Preserve the exact job version analyzed.
2. Extract responsibilities, required skills, preferred skills, seniority, location, work mode, compensation when explicitly present, and application constraints.
3. Every extracted item records source text or a source span.
4. Match requirements only against Career Brain claims with allowed confidentiality.
5. Distinguish matched evidence, partial evidence, missing evidence, and needs verification.
6. AI-generated analysis is not career evidence.
7. Produce an explanation and confidence per requirement.
8. Never invent requirements absent from the job posting.

## Acceptance criteria
- Analysis is reproducible against a specific job version.
- Every match references claim IDs.
- Confidential claims cannot be used for application-safe output.
- Missing evidence is represented explicitly rather than guessed.
