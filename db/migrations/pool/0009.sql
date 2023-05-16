-- +migrate Up
CREATE INDEX idx_pool_gas_price_timestamp ON pool.gas_price (timestamp);

-- +migrate Down
DROP INDEX idx_pool_gas_price_timestamp; 
