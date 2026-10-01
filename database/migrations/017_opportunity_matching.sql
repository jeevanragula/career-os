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
ALTER TABLE opportunity_harvest_runs ADD COLUMN IF NOT EXISTS run_key TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_harvest_runs_run_key
    ON opportunity_harvest_runs(run_key) WHERE run_key IS NOT NULL;
