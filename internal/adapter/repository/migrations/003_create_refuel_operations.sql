CREATE TABLE IF NOT EXISTS refuel_operations(
    id BIGSERIAL PRIMARY KEY,

    amount_paid DECIMAL(10, 2) NOT NULL,
    calculated_liters DECIMAL(10, 2) NOT NULL,

    price_per_liter DECIMAL(10, 2) NOT NULL,
    fuel_price_id BIGINT NOT NULL REFERENCES fuel_prices(id),

    counter_before BIGINT NOT NULL,
    counter_after BIGINT NOT NULL,
    counter_state_id BIGINT NOT NULL REFERENCES counter_states(id),

    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP WITH TIME ZONE NULL,
    cancelled_reason TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_refuel_status_date ON refuel_operations(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_refuel_created_at ON refuel_operations(created_at DESC);
-- CREATE INDEX IF NOT EXISTS idx_refuel_fuel_price_id ON refuel_operations(fuel_price_id);
-- CREATE INDEX IF NOT EXISTS idx_refuel_counter_state_id ON refuel_operations(counter_state_id);