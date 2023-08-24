-- #2324 speed up delete batches
-- +migrate Up
CREATE INDEX IF NOT EXISTS receipt_block_num_idx ON state.receipt USING btree (block_num);

-- +migrate Down
DROP INDEX IF EXISTS state.receipt_block_num_idx;
