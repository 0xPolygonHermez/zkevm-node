-- +migrate Up
CREATE TABLE state.blob_inner
(
    blob_inner_num BIGINT PRIMARY KEY,
    data           BYTEA,
    block_num      BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE    
);
ALTER TABLE state.virtual_batch
    ADD COLUMN IF NOT EXISTS blob_inner_num BIGINT REFERENCES state.blob_inner(blob_inner_num),
    ADD COLUMN IF NOT EXISTS prev_l1_it_root VARCHAR,
    ADD COLUMN IF NOT EXISTS prev_l1_it_index BIGINT;

-- +migrate Down
ALTER TABLE state.virtual_batch
    DROP COLUMN IF EXISTS blob_inner_num,
    DROP COLUMN IF EXISTS prev_l1_it_root,
    DROP COLUMN IF EXISTS prev_l1_it_index;

DROP TABLE state.blob_inner;