CREATE TABLE outboxes (
    aggregate_id UUID NOT NULL,
    event_name VARCHAR(63) NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,

    PRIMARY KEY (aggregate_id, event_name)
);

-- Index for the relay worker
CREATE INDEX idx_outboxes_unprocessed 
ON outboxes (created_at) 
WHERE processed IS FALSE;
