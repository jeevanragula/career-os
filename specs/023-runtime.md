# Spec 023 — Runtime and API

## Goal
Provide a small service shell around CareerOS workflows.

## Requirements
1. Health and readiness endpoints exist independently of external providers.
2. Runtime configuration is environment/secret based.
3. Database credentials and provider secrets never live in repository config.
4. Background workers are disabled by default until their dependencies are configured.
5. API handlers remain thin; workflow logic lives in internal packages.

## Acceptance criteria
- The service starts with only the standard Go runtime.
- Health endpoint does not require database access.
- Production secrets are externalized.
