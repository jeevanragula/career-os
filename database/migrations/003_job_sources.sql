-- CareerOS Job Sources and Discovery Runs
-- Migration 003

CREATE TABLE IF NOT EXISTS job_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    provider TEXT NOT NULL,
    base_url TEXT,
    enabled BOOLEAN NOT NULL DEFAULT true,
    authorization_mode TEXT NOT NULL DEFAULT 'public'
        CHECK (authorization_mode IN ('public','user_authorized','connector')),
    rate_limit_per_min INTEGER NOT NULL DEFAULT 30 CHECK (rate_limit_per_min > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS job_search_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    keywords TEXT[] NOT NULL DEFAULT '{}',
    companies TEXT[] NOT NULL DEFAULT '{}',
    locations TEXT[] NOT NULL DEFAULT '{}',
    remote_modes TEXT[] NOT NULL DEFAULT '{}',
    page_size INTEGER NOT NULL DEFAULT 50 CHECK (page_size > 0),
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS job_discovery_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES job_sources(id),
    query_id UUID NOT NULL REFERENCES job_search_queries(id),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'running' CHECK (status IN ('running','succeeded','partial','failed')),
    observed_count INTEGER NOT NULL DEFAULT 0 CHECK (observed_count >= 0),
    inserted_count INTEGER NOT NULL DEFAULT 0 CHECK (inserted_count >= 0),
    updated_count INTEGER NOT NULL DEFAULT 0 CHECK (updated_count >= 0),
    failed_count INTEGER NOT NULL DEFAULT 0 CHECK (failed_count >= 0),
    error_message TEXT,
    cursor TEXT
);

CREATE INDEX IF NOT EXISTS idx_job_sources_enabled ON job_sources(enabled);
CREATE INDEX IF NOT EXISTS idx_job_search_queries_enabled ON job_search_queries(enabled);
CREATE INDEX IF NOT EXISTS idx_job_discovery_runs_source_started ON job_discovery_runs(source_id, started_at DESC);
