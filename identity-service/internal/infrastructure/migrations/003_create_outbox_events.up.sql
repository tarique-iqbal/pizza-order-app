CREATE TABLE outbox_events (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    aggregate_id  UUID NOT NULL,
    event_name    VARCHAR(64) NOT NULL,
    payload       JSONB NOT NULL,
    status        VARCHAR(16) NOT NULL DEFAULT 'pending',
    attempts      INTEGER NOT NULL DEFAULT 0,
    locked_until  TIMESTAMPTZ,
    last_error    TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at  TIMESTAMPTZ,

    CONSTRAINT outbox_events_status_check
    CHECK (status IN ('pending', 'processing', 'processed', 'failed'))
);

-- critical index for worker scanning
CREATE INDEX idx_outbox_events_relay
ON outbox_events (status, created_at ASC)
WHERE status IN ('pending', 'processing');
