-- +migrate Up
ALTER TABLE pool.transaction
ADD COLUMN deposit_count BIGINT;

-- +migrate Down
ALTER TABLE pool.transaction
DROP COLUMN deposit_count;
