-- +migrate Up
ALTER TABLE pool.txs ADD COLUMN failed_counter DECIMAL(78, 0) DEFAULT 0;
CREATE INDEX idx_failed_counter ON pool.txs (failed_counter);

-- +migrate Down
DROP INDEX pool.idx_failed_counter;
ALTER TABLE pool.txs DROP COLUMN failed_counter;
