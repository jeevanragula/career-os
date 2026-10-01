-- Backfill workspace fields on the pre-existing skills table.
-- Migration 014 was already applied in some databases before these columns were added,
-- so keep this as a separate idempotent migration.
ALTER TABLE skills ADD COLUMN IF NOT EXISTS level int NOT NULL DEFAULT 3;
ALTER TABLE skills ADD COLUMN IF NOT EXISTS target_level int;
ALTER TABLE skills ADD COLUMN IF NOT EXISTS evidence_summary text NOT NULL DEFAULT '';
ALTER TABLE skills ADD COLUMN IF NOT EXISTS category text;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'skills_level_range'
    ) THEN
        ALTER TABLE skills
            ADD CONSTRAINT skills_level_range CHECK (level BETWEEN 1 AND 5);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'skills_target_level_range'
    ) THEN
        ALTER TABLE skills
            ADD CONSTRAINT skills_target_level_range
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
