-- +migrate Up
ALTER TABLE state.exit_root
    ADD COLUMN prev_block_hash BYTEA DEFAULT NULL,
    ADD COLUMN l1_info_root BYTEA DEFAULT NULL,
    ADD COLUMN l1_info_tree_index BIGINT DEFAULT NULL UNIQUE;
CREATE INDEX IF NOT EXISTS exit_root_l1_info_tree_index ON state.exit_root (l1_info_tree_index);

-- +migrate Down
ALTER TABLE state.exit_root
    DROP COLUMN prev_block_hash,
    DROP COLUMN l1_info_root,
    DROP COLUMN l1_info_tree_index;
DROP INDEX IF EXISTS state.exit_root_l1_info_tree_index;