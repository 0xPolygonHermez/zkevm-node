-- +migrate Up
ALTER TABLE state.proof
ADD COLUMN prover_id VARCHAR;
UPDATE state.proof
SET prover_id = prover;
UPDATE state.proof
SET prover = NULL;
ALTER TABLE state.proof
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE state.proof
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE state.proof
ADD COLUMN generating_since TIMESTAMP WITH TIME ZONE;
UPDATE state.proof
SET generating_since = NOW() WHERE generating IS TRUE;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS generating;

CREATE INDEX IF NOT EXISTS transaction_l2_block_num_idx ON state.transaction (l2_block_num);
CREATE INDEX IF NOT EXISTS l2block_batch_num_idx ON state.l2block (batch_num);
CREATE INDEX IF NOT EXISTS l2block_received_at_idx ON state.l2block (received_at);
CREATE INDEX IF NOT EXISTS batch_timestamp_idx ON state.batch ("timestamp");
CREATE INDEX IF NOT EXISTS log_tx_hash_idx ON state.log (tx_hash);
CREATE INDEX IF NOT EXISTS log_address_idx ON state.log (address);

ALTER TABLE state.virtual_batch
ADD COLUMN sequencer_addr VARCHAR DEFAULT '0x0000000000000000000000000000000000000000';

CREATE TABLE IF NOT EXISTS state.trusted_reorg
(
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    batch_num BIGINT,
    reason    VARCHAR NOT NULL
);

-- +migrate Down
UPDATE state.proof
SET prover = prover_id;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS prover_id;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS created_at;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS updated_at;
ALTER TABLE state.proof
ADD COLUMN  generating BOOLEAN DEFAULT FALSE;
UPDATE state.proof
SET generating = TRUE WHERE generating_since IS NOT NULL;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS generating_since;

DROP INDEX IF EXISTS state.transaction_l2_block_num_idx;
DROP INDEX IF EXISTS state.l2block_batch_num_idx;
DROP INDEX IF EXISTS state.l2block_received_at_idx;
DROP INDEX IF EXISTS state.batch_timestamp_idx;
DROP INDEX IF EXISTS state.log_tx_hash_idx;
DROP INDEX IF EXISTS state.log_address_idx;

ALTER TABLE state.virtual_batch
DROP COLUMN IF EXISTS sequencer_addr;

DROP TABLE IF EXISTS state.trusted_reorg;
