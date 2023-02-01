-- +migrate Up
CREATE INDEX IF NOT EXISTS transaction_l2_block_num_idx ON state.transaction (l2_block_num);
CREATE INDEX IF NOT EXISTS l2block_batch_num_idx ON state.l2block (batch_num);
CREATE INDEX IF NOT EXISTS l2block_received_at_idx ON state.l2block (received_at);
CREATE INDEX IF NOT EXISTS batch_timestamp_idx ON state.batch ("timestamp");
CREATE INDEX IF NOT EXISTS log_tx_hash_idx ON state.log (tx_hash);
CREATE INDEX IF NOT EXISTS log_address_idx ON state.log (address);

-- +migrate Down
DROP INDEX IF EXISTS state.transaction_l2_block_num_idx;
DROP INDEX IF EXISTS state.l2block_batch_num_idx;
DROP INDEX IF EXISTS state.l2block_received_at_idx;
DROP INDEX IF EXISTS state.batch_timestamp_idx;
DROP INDEX IF EXISTS state.log_tx_hash_idx;
DROP INDEX IF EXISTS state.log_address_idx;
