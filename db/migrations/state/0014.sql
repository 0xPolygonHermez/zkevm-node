-- +migrate Up
CREATE INDEX IF NOT EXISTS idx_batch_global_exit_root ON state.batch (global_exit_root, batch_num);

-- +migrate Down
DROP INDEX IF EXISTS state.idx_batch_global_exit_root;
    
