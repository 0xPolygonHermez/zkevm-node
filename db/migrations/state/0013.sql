-- +migrate Up
ALTER TABLE state.exit_root
    ADD COLUMN prev_block_hash BYTEA DEFAULT NULL,
    ADD COLUMN l1_info_root BYTEA DEFAULT NULL,
    ADD COLUMN l1_info_tree_index BIGINT DEFAULT NULL UNIQUE;
CREATE INDEX IF NOT EXISTS exit_root_l1_info_tree_index ON state.exit_root (l1_info_tree_index);

ALTER TABLE state.transaction
    ADD COLUMN l2_hash VARCHAR UNIQUE;

CREATE INDEX IF NOT EXISTS transaction_l2_hash ON state.transaction (l2_hash);

-- +migrate Down
ALTER TABLE state.exit_root
    DROP COLUMN prev_block_hash,
    DROP COLUMN l1_info_root,
    DROP COLUMN l1_info_tree_index;
DROP INDEX IF EXISTS state.exit_root_l1_info_tree_index;

ALTER TABLE state.transaction
    DROP COLUMN l2_hash;

DROP INDEX IF EXISTS state.transaction_l2_hash;