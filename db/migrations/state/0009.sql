-- +migrate Up
ALTER TABLE IF EXISTS state.fork_id DROP CONSTRAINT IF EXISTS fork_id_block_num_fkey;

-- +migrate Down
ALTER TABLE IF EXISTS state.fork_id ADD CONSTRAINT fork_id_block_num_fkey
FOREIGN KEY(block_num)
REFERENCES state.block (block_num) ON DELETE CASCADE;
