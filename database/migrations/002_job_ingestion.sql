-- CareerOS Job Ingestion schema
-- Migration 002

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company TEXT NOT NULL,
    title TEXT NOT NULL,
    canonical_url TEXT,
    location TEXT,
    remote_mode TEXT,
    employment_type TEXT,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'open'
        CHECK (status IN ('open','closed','expired','unknown')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS job_observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    source TEXT NOT NULL,
    source_job_id TEXT,
    source_url TEXT NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    raw_title TEXT,
    raw_company TEXT,
    raw_location TEXT,
    raw_description TEXT,
    raw_metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS job_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    content_hash TEXT NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    title TEXT NOT NULL,
    company TEXT NOT NULL,
    location TEXT,
    remote_mode TEXT,
    employment_type TEXT,
    description TEXT NOT NULL,
    UNIQUE (job_id, content_hash)
);

CREATE INDEX IF NOT EXISTS idx_jobs_company_title ON jobs(company, title);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_job_observations_job ON job_observations(job_id);
CREATE INDEX IF NOT EXISTS idx_job_observations_source_id ON job_observations(source, source_job_id);
CREATE INDEX IF NOT EXISTS idx_job_versions_job ON job_versions(job_id);
