-- +migrate Up
ALTER TABLE state.batch ADD COLUMN acc_input_hash VARCHAR;

CREATE TABLE state.sequence_group
(
    tx_hash        VARCHAR,
    tx_nonce       DECIMAL(78, 0),
    from_batch_num BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    to_batch_num   BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    status         VARCHAR(15) NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY  (tx_hash)
);

ALTER TABLE state.proof ADD COLUMN tx_hash    VARCHAR;
ALTER TABLE state.proof ADD COLUMN tx_nonce   DECIMAL(78, 0);
ALTER TABLE state.proof ADD COLUMN status     VARCHAR(15);
ALTER TABLE state.proof ADD COLUMN created_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE state.proof ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE;

UPDATE state.proof SET created_at = NOW();
UPDATE state.proof SET status = 'pending' WHERE batch_num > (SELECT batch_num from state.verified_batch ORDER BY batch_num DESC LIMIT 1);
UPDATE state.proof SET status = 'confirmed', updated_at = NOW() WHERE batch_num <= (SELECT batch_num from state.verified_batch ORDER BY batch_num DESC LIMIT 1);

ALTER TABLE state.proof ALTER COLUMN status SET NOT NULL;
ALTER TABLE state.proof ALTER COLUMN created_at SET NOT NULL;

-- +migrate Down
ALTER TABLE state.batch DROP COLUMN acc_input_hash;

DROP TABLE state.sequence_group;
ALTER TABLE state.proof DROP COLUMN tx_hash;
ALTER TABLE state.proof DROP COLUMN tx_nonce;
ALTER TABLE state.proof DROP COLUMN status;
ALTER TABLE state.proof DROP COLUMN created_at;
ALTER TABLE state.proof DROP COLUMN updated_at;
