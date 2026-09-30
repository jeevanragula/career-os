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

CREATE TABLE IF NOT EXISTS job_source_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    source_url TEXT NOT NULL,
    evidence_type TEXT NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    UNIQUE(job_id, source_url, evidence_type)
);

CREATE INDEX IF NOT EXISTS idx_harvest_runs_started ON opportunity_harvest_runs(started_at DESC);
CREATE INDEX IF NOT EXISTS idx_job_source_evidence_job ON job_source_evidence(job_id);
