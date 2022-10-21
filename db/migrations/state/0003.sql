-- +migrate Down
DROP TABLE state.sequences;

-- +migrate Up
CREATE TABLE state.sequences
(
    batch_num        BIGINT NOT NULL PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    state_root       VARCHAR NOT NULL,
    global_exit_root VARCHAR NOT NULL,
    local_exit_root  VARCHAR NOT NULL,
    timestamp        TIMESTAMP NOT NULL,
    txs              VARCHAR[] NOT NULL,

    status        VARCHAR(15) NOT NULL,
    l1_tx_hash    VARCHAR,
    l1_tx_encoded jsonb[],

    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);