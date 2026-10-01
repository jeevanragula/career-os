-- Backfill workspace fields and opportunity matching schema on existing databases.
-- Migrations 014/017 may have been skipped by an older init volume.
-- This migration is intentionally idempotent.

ALTER TABLE skills ADD COLUMN IF NOT EXISTS level int NOT NULL DEFAULT 3;
ALTER TABLE skills ADD COLUMN IF NOT EXISTS target_level int;
ALTER TABLE skills ADD COLUMN IF NOT EXISTS evidence_summary text NOT NULL DEFAULT '';
ALTER TABLE skills ADD COLUMN IF NOT EXISTS category text;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'skills_level_range') THEN
        ALTER TABLE skills ADD CONSTRAINT skills_level_range CHECK (level BETWEEN 1 AND 5);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'skills_target_level_range') THEN
        ALTER TABLE skills ADD CONSTRAINT skills_target_level_range
            CHECK (target_level IS NULL OR target_level BETWEEN 1 AND 5);
    END IF;
END $$;

INSERT INTO skills(name, category, level, target_level, evidence_summary)
VALUES
    ('Distributed Systems','architecture',5,5,'Architected and built cloud-native distributed systems and security platforms.'),
    ('Kubernetes','platform',5,5,'Extensive production Kubernetes/EKS and operator/platform experience.'),
    ('Cloud Security','security',5,5,'Built CNAPP, CSPM, DSPM and cloud posture capabilities.'),
    ('AI Security','security',4,5,'Built AI security/asset-management capabilities across models, agents, MCP and identities.'),
    ('AI Agent Architecture','ai',3,5,'Hands-on agent, MCP, ADK and A2A-compatible chatbot work.'),
    ('Technical Leadership','leadership',4,5,'Architect-level ownership across products and cross-functional technical decisions.'),
    ('Data Platforms','data',5,5,'Spark, Delta Lake, Databricks, Snowflake, ClickHouse and large-scale data pipelines.')
ON CONFLICT (name) DO UPDATE SET
    category = COALESCE(NULLIF(skills.category, ''), EXCLUDED.category),
    evidence_summary = COALESCE(NULLIF(skills.evidence_summary, ''), EXCLUDED.evidence_summary);

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
