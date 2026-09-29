# Specification 003 — Job Analysis

- Status: Draft
- Version: 0.1
- Milestone: M3 Job Analysis

## Purpose

Analyze a normalized job against Career Brain without inventing candidate facts.

## Outputs

- role summary
- required and preferred skills
- domain signals
- seniority signals
- architecture/leadership expectations
- candidate evidence matches
- evidence gaps
- application risks
- suggested resume strategy
- questions for human review

## Evidence rules

The analyzer may map requirements to verified evidence, identify missing evidence, and identify transferable experience as explicit interpretation.

It may not create qualifications, turn keyword matches into claimed experience, invent metrics, or silently mark a gap as satisfied.

## Fit representation

Represent fit as explainable requirement-level observations rather than an opaque overall score.

Each observation should identify the job requirement, candidate evidence, verification state, rationale, and confidence.

## Acceptance criteria

- [ ] Candidate-match statements trace to Career Brain evidence.
- [ ] Gaps are explicit.
- [ ] Interpretation is distinct from fact.
- [ ] Analysis is reproducible from recorded inputs and model metadata.
- [ ] Resume strategy is explainable.
