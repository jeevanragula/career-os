# CareerOS Deployment Runbook

## Local

```bash
cp .env.example .env
# set CAREEROS_AUTH_USER and CAREEROS_AUTH_PASSWORD
docker compose up --build
```

Open **http://localhost:8080**. Basic Auth protects the dashboard when credentials are configured.

If the local database volume predates migration 014:

```bash
make reset
```

This deletes only the local Docker database volume. Never run it against production data.

## Production

CareerOS is packaged as one container containing the Go API and dashboard. Put it behind HTTPS and connect it to managed PostgreSQL.

Required environment:

```env
DATABASE_URL=postgres://...
CAREEROS_AUTH_USER=...
CAREEROS_AUTH_PASSWORD=...
AI_BASE_URL=https://provider.example/v1
AI_API_KEY=...
AI_MODEL=...
```

For the first deployment, a Docker-capable VM/container platform is sufficient. Keep the database on a private network and store secrets in the platform's secret manager.

### Startup

The container serves:

- `/` — Career Command Center
- `/api/dashboard` — workspace snapshot
- `/api/agents` — agent registry
- `/api/agents/run` — run one agent
- `/api/discover` — public job discovery
- `/api/analyze` — evidence-backed job analysis
- `/api/resume` — tailored resume generation
- `/api/applications` — application preparation
- `/healthz` — liveness
- `/readyz` — database readiness

## Security

- Basic Auth is optional but should be enabled for any public deployment until a stronger identity layer is added.
- Never expose `DATABASE_URL` or `AI_API_KEY` to the browser.
- External application submission and outreach remain human-approved.
- Keep public portfolio data separate from private CareerOS data.

## Backups

Back up PostgreSQL and test restore procedures before treating CareerOS as the system of record.
