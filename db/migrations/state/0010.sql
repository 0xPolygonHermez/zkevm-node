-- +migrate Up
CREATE INDEX IF NOT EXISTS l2block_block_hash_idx ON state.l2block (block_hash);

-- +migrate Down
DROP INDEX IF EXISTS state.l2block_block_hash_idx;
