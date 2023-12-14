-- +migrate Up
ALTER TABLE state.exit_root
    ADD COLUMN prev_block_hash BYTEA DEFAULT NULL,
    ADD COLUMN l1_info_root BYTEA DEFAULT NULL,
    ADD COLUMN l1_info_tree_index BIGINT DEFAULT NULL UNIQUE;
CREATE INDEX IF NOT EXISTS idx_exit_root_l1_info_tree_index ON state.exit_root (l1_info_tree_index);

ALTER TABLE state.transaction
    ADD COLUMN l2_hash VARCHAR UNIQUE,
    ADD COLUMN used_sha256_hashes INTEGER;

CREATE INDEX IF NOT EXISTS idx_transaction_l2_hash ON state.transaction (l2_hash);

ALTER TABLE state.batch
    ADD COLUMN wip BOOLEAN NOT NULL;

ALTER TABLE state.virtual_batch
    ADD COLUMN timestamp_batch_etrog TIMESTAMP WITH TIME ZONE NULL;

-- +migrate Down
ALTER TABLE state.exit_root
    DROP COLUMN prev_block_hash,
    DROP COLUMN l1_info_root,
    DROP COLUMN l1_info_tree_index;
DROP INDEX IF EXISTS state.idx_exit_root_l1_info_tree_index;

ALTER TABLE state.transaction
    DROP COLUMN l2_hash,
    DROP COLUMN used_sha256_hashes;

DROP INDEX IF EXISTS state.idx_transaction_l2_hash;

ALTER TABLE state.batch
    DROP COLUMN IF EXISTS wip;

ALTER TABLE state.virtual_batch
    DROP COLUMN IF EXISTS timestamp_batch_etrog;

