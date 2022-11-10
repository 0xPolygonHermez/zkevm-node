-- +migrate Down
DROP TABLE state.sequence_group;

ALTER TABLE state.proof DROP COLUMN tx_hash;
ALTER TABLE state.proof DROP COLUMN tx_nonce;
ALTER TABLE state.proof DROP COLUMN status;
ALTER TABLE state.proof DROP COLUMN created_at;
ALTER TABLE state.proof DROP COLUMN updated_at;

-- +migrate Up
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

UPDATE state.proof SET status = 'confirmed', created_at = NOW(), updated_at = NOW();

ALTER TABLE state.proof ALTER COLUMN status SET NOT NULL;
ALTER TABLE state.proof ALTER COLUMN created_at SET NOT NULL;
