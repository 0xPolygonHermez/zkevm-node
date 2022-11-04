-- +migrate Up
CREATE TABLE IF NOT EXISTS state.sequences
( --Allowed Verifications
    last_verified_batch_num BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    new_verified_batch_num   BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE state.sequences;
