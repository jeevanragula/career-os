-- Repair discovery tables for databases that recorded legacy migrations
-- before the migration runner was introduced.
CREATE TABLE IF NOT EXISTS discovery_runs (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 started_at timestamptz NOT NULL DEFAULT now(),
 finished_at timestamptz,
 status text NOT NULL DEFAULT 'running',
 query jsonb NOT NULL DEFAULT '{}'::jsonb,
 discovered_count int NOT NULL DEFAULT 0,
 error text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS discovered_companies (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 name text NOT NULL,
 canonical_domain text NOT NULL,
 homepage_url text NOT NULL DEFAULT '',
 career_url text NOT NULL DEFAULT '',
 status text NOT NULL DEFAULT 'discovered' CHECK(status IN ('discovered','verified','rejected')),
 first_seen_at timestamptz NOT NULL DEFAULT now(),
 last_seen_at timestamptz NOT NULL DEFAULT now(),
 evidence jsonb NOT NULL DEFAULT '[]'::jsonb,
 UNIQUE(canonical_domain)
);

CREATE TABLE IF NOT EXISTS discovery_evidence (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 company_id uuid NOT NULL REFERENCES discovered_companies(id) ON DELETE CASCADE,
 signal text NOT NULL,
 category text NOT NULL,
 source_url text NOT NULL,
 evidence text NOT NULL DEFAULT '',
 weight int NOT NULL DEFAULT 1,
 discovered_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_discovered_companies_status
 ON discovered_companies(status,last_seen_at DESC);

CREATE INDEX IF NOT EXISTS idx_discovery_evidence_company
 ON discovery_evidence(company_id,discovered_at DESC);
