-- +migrate Up
ALTER TABLE state.fork_id
DROP COLUMN IF EXISTS block_num;

-- +migrate Down
ALTER TABLE state.fork_id
ADD IF NOT EXISTS block_num BIGINT;
