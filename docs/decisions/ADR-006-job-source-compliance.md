# ADR-006 — Compliant Job Source Access

## Decision
CareerOS will use documented/public job interfaces and explicitly authorized connectors. It will not bypass authentication, CAPTCHAs, robots restrictions, rate limits, or terms of service.

## Rationale
Evidence must be durable and automation operationally safe. Provider-neutral adapters keep source failures isolated and replaceable.

## Consequences
- Some sources require a user-connected integration.
- Coverage may be lower than unrestricted scraping.
- Every observation retains provenance and retrieval metadata.
- Provider-specific adapters can be added without changing the canonical job model.
