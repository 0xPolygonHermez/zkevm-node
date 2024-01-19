-- +migrate Up
ALTER TABLE state.l2block
    ADD COLUMN IF NOT EXISTS ger VARCHAR UNIQUE;

CREATE INDEX IF NOT EXISTS idx_l2block_ger ON state.l2block (ger);

-- +migrate Down
DROP INDEX IF EXISTS state.idx_l2block_ger;

ALTER TABLE state.l2block
    DROP COLUMN IF EXISTS ger;
    
