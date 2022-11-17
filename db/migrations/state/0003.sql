-- +migrate Up
ALTER TABLE state.batch ADD COLUMN acc_input_hash VARCHAR;
CREATE TABLE IF NOT EXISTS state.sequences
( --Allowed Verifications
    from_batch_num BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    to_batch_num   BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE
);

-- +migrate Down
ALTER TABLE state.batch DROP COLUMN acc_input_hash;
DROP TABLE state.sequences;
