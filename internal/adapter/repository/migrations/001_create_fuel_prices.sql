CREATE TABLE IF NOT EXISTS fuel_prices (
    id BIGSERIAL PRIMARY KEY,
    price_per_liter DECIMAL(10, 2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_fuel_prices_is_active ON fuel_prices(is_active);
CREATE INDEX IF NOT EXISTS idx_fuel_prices_created_at ON fuel_prices(created_at DESC);