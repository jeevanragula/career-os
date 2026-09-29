CREATE TABLE IF NOT EXISTS job_analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_version_id UUID NOT NULL REFERENCES job_versions(id),
    summary TEXT NOT NULL,
    overall_confidence NUMERIC(5,4) CHECK (overall_confidence >= 0 AND overall_confidence <= 1),
    model_provider TEXT,
    model_name TEXT,
    prompt_version TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS job_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id UUID NOT NULL REFERENCES job_analyses(id) ON DELETE CASCADE,
    category TEXT NOT NULL,
    requirement_text TEXT NOT NULL,
    priority TEXT NOT NULL CHECK (priority IN ('required','preferred','context')),
    source_quote TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS requirement_matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requirement_id UUID NOT NULL REFERENCES job_requirements(id) ON DELETE CASCADE,
    claim_id UUID NOT NULL REFERENCES claims(id),
    status TEXT NOT NULL CHECK (status IN ('matched','partial','unmatched','needs_verification')),
    confidence NUMERIC(5,4) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    rationale TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_job_analyses_version ON job_analyses(job_version_id);
CREATE INDEX IF NOT EXISTS idx_job_requirements_analysis ON job_requirements(analysis_id);
CREATE INDEX IF NOT EXISTS idx_requirement_matches_claim ON requirement_matches(claim_id);
