-- CareerOS Career Brain schema
-- Migration 001
-- PostgreSQL

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS people (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name TEXT NOT NULL,
    professional_positioning TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS employments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE RESTRICT,
    employer TEXT NOT NULL,
    title TEXT,
    start_date DATE,
    end_date DATE,
    verification_status TEXT NOT NULL DEFAULT 'needs_verification'
        CHECK (verification_status IN ('needs_verification','user_asserted','verified','rejected')),
    confidentiality TEXT NOT NULL DEFAULT 'internal'
        CHECK (confidentiality IN ('public','application_safe','internal','confidential')),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    employment_id UUID REFERENCES employments(id) ON DELETE RESTRICT,
    summary TEXT,
    role TEXT,
    verification_status TEXT NOT NULL DEFAULT 'needs_verification'
        CHECK (verification_status IN ('needs_verification','user_asserted','verified','rejected')),
    confidentiality TEXT NOT NULL DEFAULT 'internal'
        CHECK (confidentiality IN ('public','application_safe','internal','confidential')),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
    claim TEXT NOT NULL,
    verification_status TEXT NOT NULL DEFAULT 'needs_verification'
        CHECK (verification_status IN ('needs_verification','user_asserted','verified','rejected')),
    confidentiality TEXT NOT NULL DEFAULT 'internal'
        CHECK (confidentiality IN ('public','application_safe','internal','confidential')),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    category TEXT,
    verification_status TEXT NOT NULL DEFAULT 'needs_verification'
        CHECK (verification_status IN ('needs_verification','user_asserted','verified','rejected')),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_type TEXT NOT NULL,
    source_reference TEXT NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    confidentiality TEXT NOT NULL DEFAULT 'internal'
        CHECK (confidentiality IN ('public','application_safe','internal','confidential')),
    integrity_status TEXT NOT NULL DEFAULT 'unreviewed'
        CHECK (integrity_status IN ('unreviewed','reviewed','superseded')),
    content_hash TEXT,
    verification_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_key TEXT NOT NULL UNIQUE,
    text TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID,
    verification_status TEXT NOT NULL DEFAULT 'needs_verification'
        CHECK (verification_status IN ('needs_verification','user_asserted','verified','rejected')),
    confidence NUMERIC(5,4) CHECK (confidence IS NULL OR (confidence >= 0 AND confidence <= 1)),
    confidentiality TEXT NOT NULL DEFAULT 'internal'
        CHECK (confidentiality IN ('public','application_safe','internal','confidential')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS claim_evidence (
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    evidence_id UUID NOT NULL REFERENCES evidence(id) ON DELETE RESTRICT,
    relationship TEXT NOT NULL DEFAULT 'supports',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (claim_id, evidence_id)
);

CREATE TABLE IF NOT EXISTS career_change_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    actor_type TEXT NOT NULL,
    reason TEXT,
    before_state JSONB,
    after_state JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_employments_person ON employments(person_id);
CREATE INDEX IF NOT EXISTS idx_projects_employment ON projects(employment_id);
CREATE INDEX IF NOT EXISTS idx_achievements_project ON achievements(project_id);
CREATE INDEX IF NOT EXISTS idx_claims_status ON claims(verification_status);
CREATE INDEX IF NOT EXISTS idx_claims_confidentiality ON claims(confidentiality);
CREATE INDEX IF NOT EXISTS idx_evidence_source ON evidence(source_type);
