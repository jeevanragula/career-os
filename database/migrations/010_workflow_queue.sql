CREATE TABLE IF NOT EXISTS workflow_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','running','succeeded','failed','cancelled')),
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_at TIMESTAMPTZ,
    locked_by TEXT,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_workflow_jobs_ready ON workflow_jobs(status, available_at);
CREATE UNIQUE INDEX IF NOT EXISTS uq_workflow_jobs_active ON workflow_jobs(job_type, entity_type, entity_id) WHERE status IN ('queued','running');

CREATE TABLE IF NOT EXISTS external_action_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action_type TEXT NOT NULL CHECK (action_type IN ('submit_application','send_outreach')),
    entity_id UUID NOT NULL,
    payload_hash TEXT NOT NULL,
    approved_by TEXT NOT NULL,
    approved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(action_type, entity_id, payload_hash)
);
