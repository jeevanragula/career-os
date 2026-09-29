# Career Brain Invariant Tests

These are the initial executable-test requirements for the database implementation.

## Required cases

1. A claim defaults to needs_verification.
2. A claim can reference evidence.
3. A verified claim can retain an audit event.
4. Confidence is constrained to 0..1.
5. Verification state rejects unknown values.
6. Confidentiality rejects unknown values.
7. Deleting evidence referenced by a claim is prevented.
8. Employment deletion is prevented when projects depend on it.
9. Duplicate claim_key is rejected.
10. Duplicate skill name is rejected.

## Security cases

11. Confidential claims are never eligible for application-safe artifact generation.
12. generated_analysis evidence cannot by itself cause a claim to become verified.
13. External action approval is represented outside the canonical claim state and is auditable.
