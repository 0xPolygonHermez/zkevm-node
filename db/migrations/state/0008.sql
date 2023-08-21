-- +migrate Up
CREATE TABLE IF NOT EXISTS state.fork_id
(
    fork_id        BIGINT NOT NULL PRIMARY KEY,
    from_batch_num BIGINT NOT NULL,
    to_batch_num   BIGINT NOT NULL,
    version        VARCHAR,
    block_num      BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS state.fork_id;

