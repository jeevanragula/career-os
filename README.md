# CareerOS

CareerOS is a private, evidence-first **personal career operating system** for senior engineering careers.

It is intentionally separate from the public portfolio repository. CareerOS is the private system of record and action workspace; the portfolio is a downstream public presentation surface.

## What works now

- Career Brain with claims, evidence, verification and confidentiality boundaries
- Canonical public job ingestion with provenance, normalization, versioning and deduplication
- Lever and Ashby discovery adapters
- AI job analysis grounded only in application-safe career claims
- Evidence-grounded tailored Markdown resume generation
- Application preparation records with explicit human approval boundaries
- Persistent Career OS workspace for goals, skills, learning, networking, startup ideas, tasks and recommendations
- 14 specialized Career OS agents with a shared orchestration contract
- Command-center dashboard showing opportunities, actions, skills, learning, network and startup tracks
- Agent recommendations can be accepted/dismissed and tasks can be completed from the dashboard
- Optional HTTP Basic Authentication for a private deployment
- Docker Compose local runtime with PostgreSQL
- Go test/validation workflow in GitHub Actions

## Agent system

The dashboard exposes these agents:

1. Career Intelligence
2. Opportunity
3. Skill Gap
4. Learning
5. Project
6. Resume
7. Interview
8. Networking
9. Personal Brand
10. Startup
11. Market Intelligence
12. Achievement
13. Portfolio
14. Career Strategy

Agents analyze the canonical Career Brain and workspace. They may propose recommendations and tasks, but they cannot silently change verified career facts, publish content, submit applications, or send outreach.

## Core loop

```text
Capture experience / opportunity / learning
                ↓
           Career Brain
                ↓
       ┌────────┼─────────┐
       ↓        ↓         ↓
  Opportunity Learning  Strategy
       ↓        ↓         ↓
    Analysis  Projects  Networking
       └────────┼─────────┘
                ↓
          Action Queue
                ↓
       Human-approved actions
```

## Run locally

Requirements: Docker Engine, Docker Compose, Git.

```bash
git clone https://github.com/jeevanragula/career-os.git
cd career-os
cp .env.example .env
# edit .env and set a private password
docker compose up --build
```

Open **http://localhost:8080**. The browser will request the configured Basic Auth credentials.

Health remains available at **http://localhost:8080/healthz**.

If you already have an old local PostgreSQL volume, run `make reset` once so migration 014 creates the new workspace tables. Do not use `make reset` against production data.

## Configure AI

CareerOS accepts an OpenAI-compatible `/chat/completions` endpoint:

```env
AI_BASE_URL=https://your-provider.example/v1
AI_API_KEY=...
AI_MODEL=...
```

The same interface can point to a compatible hosted provider or a local model gateway. Without AI configuration, the dashboard still loads and job discovery still works; agent buttons explain that AI configuration is required.

## Job workflow

1. Set your career targets in the dashboard.
2. Run **Discover opportunities**; CareerOS chooses public signal sources itself.
3. CareerOS resolves public career pages, harvests postings, preserves source evidence, versions changes, and computes a transparent skill-match score.
4. Review the Opportunity Inbox and run **Analyze** on roles worth deeper evaluation.
5. Generate an evidence-grounded tailored resume.
6. Create an application preparation record.
7. Complete the actual external submission yourself after reviewing the package.

CareerOS deliberately does not auto-submit applications or send outreach.

## Deploy

The simplest production shape is:

```text
HTTPS / auth
    ↓
CareerOS container
    ↓
Managed PostgreSQL
```

Use any Docker-capable host (for example a VM/container platform) and supply:

- `DATABASE_URL`
- `CAREEROS_AUTH_USER`
- `CAREEROS_AUTH_PASSWORD`
- `AI_BASE_URL`
- `AI_API_KEY`
- `AI_MODEL`

Keep PostgreSQL private and HTTPS enabled. Back up the database. Never put provider keys in browser code.

For a first deployment, keep the dashboard and API together in the same container. Kubernetes is optional later; it is not required for the MVP.

## Repository structure

```text
specs/              Product and capability specifications
docs/architecture/  System architecture
career-profile/     Canonical career evidence
resumes/            Resume strategies and artifacts
agents/             Agent contracts
database/            PostgreSQL migrations
internal/            Go domain, AI, jobs, career workspace and API packages
cmd/                 Executable entry points
web/                 CareerOS command-center UI
```

## Safety and privacy

CareerOS must never fabricate career facts, metrics, employers, technologies, responsibilities, patents, achievements or qualifications.

Do not store credentials, session cookies, API tokens, passwords, employer/customer confidential information, proprietary source code or unsafe non-public architecture details.

Generated career artifacts must trace back to canonical evidence.

External actions remain human-approved.

## Autonomous opportunity discovery

You do not maintain a company or job-source list.

CareerOS discovers opportunity signals from multiple public ecosystems, including major job portals, startup ecosystems, Gartner Peer Insights/vendor ecosystems, CNCF/cloud-native communities, public conference sponsorships, cloud events, security conferences/competitions, and engineering communities.

The discovery flow is:

1. Start with the user's target-role and technology profile.
2. Search across multiple independent signal classes.
3. Deduplicate discovered companies/domains.
4. Preserve the evidence that caused each company to be discovered.
5. Resolve and verify the company's own career page.
6. Extract public roles and feed verified roles into the existing job-analysis pipeline.
7. Analyze opportunities against the user's verified career evidence.

A discovery signal is **not** an employer rating. It is only a reason to investigate a company or role.

The current web-search implementation uses Tavily Search API. Configure only:

TAVILY_API_KEY=...

You do not need to configure individual job portals or companies.


## Scheduled autonomous discovery

GitHub Actions runs the opportunity harvester daily. Add these **repository Actions secrets**:

- `DATABASE_URL` — PostgreSQL connection string for the CareerOS database.
- `TAVILY_API_KEY` — Tavily Search API key.

Optional repository Actions variables:

- `CAREEROS_DISCOVERY_ROLES`
- `CAREEROS_DISCOVERY_DOMAINS`
- `CAREEROS_DISCOVERY_LOCATIONS`

The workflow is intentionally fail-closed when either required secret is missing. GitHub stores Actions secrets encrypted and exposes them to workflows only when referenced by the workflow.

For an existing PostgreSQL/Supabase database, apply all migrations through `017_opportunity_matching.sql` before enabling the scheduled harvester. The repository also contains `cmd/career-migrate` for idempotent migration execution on databases that have not previously been migrated.

The harvester records harvest-run status, job source evidence, job versions, opportunity matches, and freshness. It does **not** submit applications, send outreach, or publish content.
