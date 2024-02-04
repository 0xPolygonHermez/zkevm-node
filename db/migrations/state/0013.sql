-- +migrate Up
ALTER TABLE state.exit_root
    ADD COLUMN IF NOT EXISTS prev_block_hash BYTEA DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS l1_info_root BYTEA DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS l1_info_tree_index BIGINT DEFAULT NULL UNIQUE;
CREATE INDEX IF NOT EXISTS idx_exit_root_l1_info_tree_index ON state.exit_root (l1_info_tree_index);

ALTER TABLE state.transaction
    ADD COLUMN IF NOT EXISTS l2_hash VARCHAR UNIQUE;

CREATE INDEX IF NOT EXISTS idx_transaction_l2_hash ON state.transaction (l2_hash);

ALTER TABLE state.batch
    ADD COLUMN IF NOT EXISTS wip BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE state.virtual_batch
    ADD COLUMN IF NOT EXISTS timestamp_batch_etrog TIMESTAMP WITH TIME ZONE NULL,
    ADD COLUMN IF NOT EXISTS l1_info_root VARCHAR;

-- +migrate Down
ALTER TABLE state.exit_root
    DROP COLUMN IF EXISTS prev_block_hash,
    DROP COLUMN IF EXISTS l1_info_root,
    DROP COLUMN IF EXISTS l1_info_tree_index;
DROP INDEX IF EXISTS state.idx_exit_root_l1_info_tree_index;

ALTER TABLE state.transaction
    DROP COLUMN IF EXISTS l2_hash;
    
DROP INDEX IF EXISTS state.idx_transaction_l2_hash;

ALTER TABLE state.batch
    DROP COLUMN IF EXISTS wip;

ALTER TABLE state.virtual_batch
    DROP COLUMN IF EXISTS timestamp_batch_etrog,
    DROP COLUMN IF EXISTS l1_info_root;

