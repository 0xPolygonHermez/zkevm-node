-- +migrate Up
CREATE INDEX IF NOT EXISTS idx_receipt_tx_index ON state.receipt (block_num, tx_index);

-- +migrate Down
DROP INDEX IF EXISTS state.idx_receipt_tx_index;
    
