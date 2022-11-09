-- +migrate Up
ALTER TABLE state.batch ADD COLUMN acc_input_hash VARCHAR;
CREATE TABLE IF NOT EXISTS state.sequences
( --Allowed Verifications
    last_verified_batch_num BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    new_verified_batch_num   BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE
);

-- +migrate Down
ALTER TABLE state.batch DROP COLUMN acc_input_hash;
DROP TABLE state.sequences;
