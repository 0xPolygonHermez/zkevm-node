-- +migrate Up
CREATE TABLE IF NOT EXISTS pool.whitelisted (
	addr VARCHAR PRIMARY KEY
);

CREATE INDEX IF NOT EXISTS idx_transaction_from_nonce ON pool.transaction (from_address, nonce);
CREATE INDEX IF NOT EXISTS idx_transaction_status ON pool.transaction (status);
CREATE INDEX IF NOT EXISTS idx_transaction_hash ON pool.transaction (hash);

-- +migrate Down
DROP TABLE pool.whitelisted;

DROP INDEX IF EXISTS pool.idx_transaction_from_nonce;
DROP INDEX IF EXISTS pool.idx_transaction_status;
DROP INDEX IF EXISTS pool.idx_transaction_hash;
