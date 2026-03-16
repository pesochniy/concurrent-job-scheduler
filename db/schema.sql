CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,

    state SMALLINT NOT NULL,
    retry_count INT NOT NULL,
    max_retries INT NOT NULL,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_jobs_state ON jobs(state);