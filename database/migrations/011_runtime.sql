CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Ensure canonical URL is unique when present. Historical observations remain separate.
CREATE UNIQUE INDEX IF NOT EXISTS uq_jobs_canonical_url
    ON jobs(canonical_url) WHERE canonical_url IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_jobs_last_seen ON jobs(last_seen_at DESC);
