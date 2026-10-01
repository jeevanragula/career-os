-- Opportunity matching and harvest-run extensions.
-- Keep the prerequisite harvest table idempotent for legacy databases.
CREATE TABLE IF NOT EXISTS opportunity_harvest_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    candidates_seen INTEGER NOT NULL DEFAULT 0,
    pages_fetched INTEGER NOT NULL DEFAULT 0,
    jobs_observed INTEGER NOT NULL DEFAULT 0,
    jobs_ingested INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'running'
        CHECK (status IN ('running','completed','failed')),
    error TEXT
);

CREATE INDEX IF NOT EXISTS idx_harvest_runs_started
    ON opportunity_harvest_runs(started_at DESC);

CREATE TABLE IF NOT EXISTS opportunity_matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    score NUMERIC(6,3) NOT NULL DEFAULT 0,
    matched_skills JSONB NOT NULL DEFAULT '[]'::jsonb,
    gaps JSONB NOT NULL DEFAULT '[]'::jsonb,
    rationale TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'new'
        CHECK (status IN ('new','reviewed','saved','dismissed')),
    computed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(job_id)
);

CREATE INDEX IF NOT EXISTS idx_opportunity_matches_score
    ON opportunity_matches(status, score DESC, computed_at DESC);

ALTER TABLE opportunity_harvest_runs
    ADD COLUMN IF NOT EXISTS run_key TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_harvest_runs_run_key
    ON opportunity_harvest_runs(run_key)
    WHERE run_key IS NOT NULL;
