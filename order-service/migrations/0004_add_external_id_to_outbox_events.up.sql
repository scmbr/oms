ALTER TABLE outbox_events
ADD COLUMN external_id UUID UNIQUE;

CREATE UNIQUE INDEX idx_outbox_events_external_id ON outbox_events(external_id);