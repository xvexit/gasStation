CREATE TABLE IF NOT EXISTS counter_states(
    id BIGSERIAL PRIMARY KEY,
    current_value INTEGER NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_counter_states_updated_at ON counter_states(updated_at DESC);