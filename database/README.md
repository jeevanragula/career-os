# CareerOS Database

The initial durable Career Brain model uses PostgreSQL.

## Migration

The first migration is:

database/migrations/001_career_brain.sql

## Rules

- Database state is durable domain state.
- Career claims require explicit verification state.
- Evidence and claims are separate.
- Foreign keys protect canonical relationships.
- Application code must not bypass verification/confidentiality invariants.

The migration intentionally avoids provider-specific features beyond PostgreSQL and pgcrypto.
