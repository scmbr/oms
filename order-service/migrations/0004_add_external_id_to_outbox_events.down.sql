DROP INDEX IF EXISTS idx_outbox_events_external_id;

ALTER TABLE outbox_events
DROP COLUMN IF EXISTS external_id;
