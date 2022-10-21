-- +migrate Down
DROP TABLE state.sequences_group;
DROP TABLE state.sequences;

-- +migrate Up
CREATE TABLE state.pending_sequences
(
    batch_num        BIGINT NOT NULL PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    state_root       VARCHAR NOT NULL,
    global_exit_root VARCHAR NOT NULL,
    local_exit_root  VARCHAR NOT NULL,
    timestamp        TIMESTAMP NOT NULL,
    txs              VARCHAR[] NOT NULL,
);

CREATE TABLE state.pending_sequences_groups
(
    tx_hash     VARCHAR,
    tx_encoded  jsonb,
    batch_nums  BIGINT[],
    status      VARCHAR(15) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (tx_hash)
);

CREATE FUNCTION state.del_sequences_group_cascade_batch() RETURNS trigger AS $del_sequences_group_cascade_batch$
    BEGIN
        DELETE state.sequences_group
         WHERE OLD.batch_num = ANY(batch_nums);
        RETURN OLD;
    END;
$del_sequences_group_cascade_batch$ LANGUAGE plpgsql;

CREATE TRIGGER state.del_seq_group BEFORE DELETE ON state.batch
   FOR EACH ROW EXECUTE FUNCTION state.del_sequences_group_cascade_batch();
