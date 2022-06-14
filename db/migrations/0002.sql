-- +migrate Up
CREATE SCHEMA statev2

-- History
CREATE TABLE statev2.block ( --L1 block
    block_num BIGINT PRIMARY KEY,
    block_hash VARCHAR NOT NULL,
    parent_hash VARCHAR,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE statev2.batch (  --batch abstraction: will be created through trusted state
    batch_num BIGINT PRIMARY KEY,
    global_exit_root VARCHAR,
    timestamp TIMESTAMP,
    raw_txs_data VARCHAR
);

CREATE TABLE statev2.virtual_batch (  --virtual state
    batch_num BIGINT PRIMARY KEY REFERENCES statev2.batch (batch_num) ON DELETE CASCADE,
    tx_hash VARCHAR,
    sequencer VARCHAR,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.verified_batch (  --consolidated state
    batch_num BIGINT PRIMARY KEY REFERENCES statev2.batch (batch_num) ON DELETE CASCADE,
    tx_hash VARCHAR,
    aggregator VARCHAR,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.forced_batch (
    forced_batch_num BIGINT PRIMARY KEY,
    global_exit_root VARCHAR,
    timestamp TIMESTAMP,
    raw_txs_data VARCHAR,
    sequencer VARCHAR,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.transaction (  --transaction abstraction. transaction == L2 block
    hash VARCHAR PRIMARY KEY,
    from_address VARCHAR,
    encoded VARCHAR,
    decoded jsonb,
    header jsonb,
    uncles jsonb,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    batch_num BIGINT NOT NULL REFERENCES statev2.batch (batch_num) ON DELETE CASCADE
);

CREATE TABLE statev2.exit_root
(
    block_num               BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    global_exit_root_num    BIGINT PRIMARY KEY,
    mainnet_exit_root       BYTEA,
    rollup_exit_root        BYTEA
);