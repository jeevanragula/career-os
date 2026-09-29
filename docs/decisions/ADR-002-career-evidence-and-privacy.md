# ADR-002 — Career Evidence and Privacy Model

- Status: Accepted
- Date: 2026-09-30

## Context

CareerOS needs rich evidence for accurate job-search artifacts, but professional experience can contain confidential employer and customer information.

## Decision

CareerOS models career claims separately from evidence and classifies every material claim by verification status and confidentiality.

The system distinguishes canonical evidence, application-safe claims, internal-only evidence, and confidential evidence.

AI-generated text is never evidence by itself.

## Verification states

- verified
- user_asserted
- needs_verification
- rejected

The system must not silently promote a claim between states.

## Confidentiality levels

- public
- application_safe
- internal
- confidential

Downstream artifacts may only consume claims allowed by their target publication level.

## Consequence

The metadata is richer than a conventional resume database, but it prevents factual drift and accidental confidential disclosure.
