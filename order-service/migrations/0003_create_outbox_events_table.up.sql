CREATE TABLE outbox_events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(64) NOT NULL,
    order_id UUID NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    processed_at TIMESTAMP
);

CREATE INDEX idx_outbox_status ON outbox_events(status);
CREATE INDEX idx_outbox_order_id ON outbox_events(order_id);
