-- +migrate Up
ALTER TABLE pool.blocked ADD COLUMN block_reason VARCHAR;
CREATE INDEX idx_pool_gas_price_timestamp ON pool.gas_price (timestamp);

-- +migrate Down
ALTER TABLE pool.blocked DROP COLUMN block_reason;
DROP INDEX IF EXISTS idx_pool_gas_price_timestamp;
