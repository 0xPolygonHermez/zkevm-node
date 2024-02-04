-- +migrate Up
ALTER TABLE pool.transaction
    ADD COLUMN l2_hash VARCHAR UNIQUE,
    ADD COLUMN used_sha256_hashes INTEGER DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_transaction_l2_hash ON pool.transaction (l2_hash);

-- +migrate Down
DROP INDEX IF EXISTS pool.idx_transaction_l2_hash;
ALTER TABLE pool.transaction
    DROP COLUMN l2_hash,
    DROP COLUMN used_sha256_hashes;

