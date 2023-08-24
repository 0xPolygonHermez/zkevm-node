-- +migrate Up
CREATE TABLE IF NOT EXISTS state.fork_id
(
    fork_id        BIGINT NOT NULL PRIMARY KEY,
    from_batch_num numeric NOT NULL,
    to_batch_num   numeric NOT NULL,
    version        VARCHAR,
    block_num      BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS receipt_block_num_idx ON state.receipt USING btree (block_num);

-- +migrate Down
DROP INDEX IF EXISTS state.receipt_block_num_idx;

DROP TABLE IF EXISTS state.fork_id;
