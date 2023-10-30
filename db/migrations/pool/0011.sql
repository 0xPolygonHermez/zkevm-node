-- +migrate Up
CREATE TABLE pool.whitelisted (
	addr VARCHAR PRIMARY KEY
);

CREATE INDEX IF NOT EXISTS idx_transaction_from_nonce ON pool.transaction (from_address, nonce);
CREATE INDEX IF NOT EXISTS idx_transaction_wip_status ON pool.transaction (status, is_wip);
CREATE INDEX IF NOT EXISTS idx_transaction_status_gas_price ON pool.transaction (status, gas_price);
CREATE INDEX IF NOT EXISTS idx_transaction_status_received_at ON pool.transaction (status, idx_transaction_status_received_at);
CREATE INDEX IF NOT EXISTS idx_transaction_from_status ON pool.transaction (status, from_address);
CREATE INDEX IF NOT EXISTS idx_transaction_hash ON pool.transaction (hash);

-- +migrate Down
DROP TABLE pool.whitelisted;

DROP INDEX IF EXISTS pool.idx_transaction_from_nonce
DROP INDEX IF EXISTS pool.idx_transaction_wip_status
DROP INDEX IF EXISTS pool.idx_transaction_status_gas_price
DROP INDEX IF EXISTS pool.idx_transaction_status_received_at
DROP INDEX IF EXISTS pool.idx_transaction_from_status
DROP INDEX IF EXISTS pool.idx_transaction_hash
