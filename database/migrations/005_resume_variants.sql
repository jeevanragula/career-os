CREATE TABLE IF NOT EXISTS resume_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES jobs(id),
    strategy TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    model_provider TEXT,
    model_name TEXT,
    prompt_version TEXT,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','reviewed','approved','archived')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS resume_variant_claims (
    variant_id UUID NOT NULL REFERENCES resume_variants(id) ON DELETE CASCADE,
    claim_id UUID NOT NULL REFERENCES claims(id),
    PRIMARY KEY (variant_id, claim_id)
);
