CREATE TABLE IF NOT EXISTS outreach_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID REFERENCES applications(id),
    recipient_name TEXT,
    recipient_channel TEXT NOT NULL,
    context TEXT NOT NULL,
    draft TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','approved','sent','cancelled')),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS application_timeline_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    event_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes TEXT,
    source TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_outreach_application ON outreach_messages(application_id);
CREATE INDEX IF NOT EXISTS idx_application_timeline_application ON application_timeline_events(application_id, event_at DESC);
