-- +migrate Up
ALTER TABLE pool.blocked ADD COLUMN block_reason VARCHAR;

-- +migrate Down
ALTER TABLE pool.blocked DROP COLUMN block_reason;
