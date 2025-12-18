CREATE TABLE outbox_events(
    id SERIAL PRIMARY KEY,
    external_id UUID NOT NULL UNIQUE,
    event_type VARCHAR(64) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(20) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(16) NOT NULL CHECK (status IN('PENDING','SENT','FAILED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ

);
CREATE INDEX idx_outbox_status ON outbox_events(status);
CREATE INDEX idx_outbox_external_id ON outbox_events(external_id);
