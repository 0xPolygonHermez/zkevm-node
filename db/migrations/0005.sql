-- +migrate Down
DROP TABLE IF EXISTS state.l2block;

-- +migrate Up

-- Table that stores L2 nodes
CREATE TABLE state.l2block
(
    block_num BIGSERIAL PRIMARY KEY,
    tx_hash BYTEA NOT NULL,
    parent_tx_hash BYTEA NOT NULL,

    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);
