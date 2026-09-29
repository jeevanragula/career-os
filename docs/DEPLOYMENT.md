# CareerOS Deployment Runbook

## 1. Local first run

Requirements: Docker Engine, Docker Compose, Git.

From the repository root:

    cp .env.example .env
    docker compose up --build

Open http://localhost:8080.

Health: http://localhost:8080/healthz

PostgreSQL migrations run automatically when the database volume is initialized. If a development database must be recreated, use `docker compose down -v` followed by `docker compose up --build`. Never use `down -v` for production.

## 2. Configure AI

CareerOS accepts an OpenAI-compatible `/chat/completions` endpoint.

Set these in `.env`:

    AI_BASE_URL=
    AI_API_KEY=
    AI_MODEL=

The same interface can point to a hosted provider or a local model server. The AI layer currently performs job analysis only.

## 3. Discover jobs

The dashboard accepts a provider, provider/company identifier, public API URL, and optional keywords.

For Lever, use the public postings API base and the company site name. For Ashby, use the company's public job-board API endpoint.

CareerOS stores the source observation, canonical job, and job version.

## 4. Production architecture

    Internet -> HTTPS reverse proxy -> CareerOS container -> Managed PostgreSQL

For the first production deployment, use a Docker-capable host with managed PostgreSQL. Kubernetes is optional and should be introduced only when you need its operational features.

## 5. Production requirements

1. Replace the development PostgreSQL password.
2. Put DATABASE_URL and AI credentials in secret management.
3. Use HTTPS.
4. Add authentication before exposing the dashboard publicly.
5. Keep PostgreSQL private.
6. Back up PostgreSQL.
7. Set resource limits and log retention.
8. Never expose provider credentials to browser JavaScript.
9. Keep application submission/outreach disabled until an authorized integration is connected to the exact-payload approval workflow.

## 6. Important MVP boundary

The repository now runs the discovery and AI-analysis path. It is not yet a production application-submission bot. That boundary is intentional: external actions require authorized integrations and human approval.
