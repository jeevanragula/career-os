CREATE TABLE IF NOT EXISTS resume_variant_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID NOT NULL REFERENCES resume_variants(id) ON DELETE CASCADE,
    format TEXT NOT NULL CHECK (format IN ('markdown','text')),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_resume_variant_content_variant ON resume_variant_content(variant_id);
