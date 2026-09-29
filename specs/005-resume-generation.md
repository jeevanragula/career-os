# Specification 005 — Resume Engine

- Status: Draft
- Version: 0.1
- Milestone: M4 Resume Engine

## Purpose

Generate targeted resumes from canonical evidence while preventing factual drift.

## Initial strategies

- Principal / Staff Engineering
- Cloud / Security Engineering
- AI Security / AI Platform

These are presentation strategies over the same canonical evidence.

## Rules

The engine may select verified claims, reorder them by relevance, condense them, rewrite for clarity, and use role-specific terminology when supported by evidence.

It may not invent experience, metrics, technologies, dates, ownership, or expose confidential information.

## Validation

Before application use:

1. Extract factual claims.
2. Map claims to canonical evidence.
3. Reject or flag unsupported claims.
4. Check confidentiality.
5. Record the evidence set and generation metadata.

## Acceptance criteria

- [ ] Resume variants share one canonical evidence source.
- [ ] Material factual claims are traceable.
- [ ] Unsupported claims are blocked or flagged.
- [ ] Confidential claims are excluded from application-safe output.
- [ ] Generated artifacts are versioned.
