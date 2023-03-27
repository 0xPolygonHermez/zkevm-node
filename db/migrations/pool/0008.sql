-- +migrate Up
ALTER TABLE pool.transaction ADD COLUMN failed_reason varchar;

-- +migrate Down
ALTER TABLE pool.transaction DROP COLUMN failed_reason;