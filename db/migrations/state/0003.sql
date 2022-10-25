-- +migrate Down
DROP TRIGGER state.tr_del_sequence_group_cascade_batch;
DROP FUNCTION state.fn_del_sequence_group_cascade_batch;
DROP TABLE state.sequence_group;
DROP TABLE state.sequence;

-- +migrate Up
CREATE TABLE state.sequence
(
    batch_num        BIGINT NOT NULL PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    state_root       VARCHAR NOT NULL,
    global_exit_root VARCHAR NOT NULL,
    local_exit_root  VARCHAR NOT NULL,
    timestamp        TIMESTAMP NOT NULL,
    txs              VARCHAR[] NOT NULL
);

CREATE TABLE state.sequence_group
(
    tx_hash      VARCHAR,
    tx_nonce     DECIMAL(78, 0),
    batch_nums   BIGINT[],
    status       VARCHAR(15) NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at   TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY  (tx_hash)
);

-- -- +migrate StatementBegin
-- CREATE FUNCTION state.fn_del_sequence_group_cascade_batch() RETURNS trigger AS $fn_del_sequence_group_cascade_batch$
--     BEGIN
--         DELETE FROM state.sequence_group
--          WHERE OLD.batch_num = ANY(batch_nums);
--         RETURN OLD;
--     END;
-- $fn_del_sequence_group_cascade_batch$ LANGUAGE plpgsql;
-- -- +migrate StatementEnd

-- CREATE TRIGGER state.tr_del_sequence_group_cascade_batch BEFORE DELETE ON state.batch
--    FOR EACH ROW EXECUTE FUNCTION state.fn_del_sequence_group_cascade_batch();
