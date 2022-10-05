-- +migrate Down
DROP SCHEMA IF EXISTS state CASCADE;

-- +migrate Up
CREATE SCHEMA state;

CREATE TABLE state.block
( --L1 block
    block_num   BIGINT PRIMARY KEY,
    block_hash  VARCHAR                  NOT NULL,
    parent_hash VARCHAR,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE state.batch
( --batch abstraction: will be created through trusted state
    batch_num        BIGINT PRIMARY KEY,
    global_exit_root VARCHAR,
    local_exit_root  VARCHAR,
    state_root       VARCHAR,
    timestamp        TIMESTAMP,
    coinbase         VARCHAR,
    raw_txs_data     BYTEA
);

CREATE TABLE state.virtual_batch
( --virtual state
    batch_num BIGINT PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    tx_hash   VARCHAR,
    coinbase  VARCHAR,
    block_num BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.verified_batch
( --consolidated state
    batch_num  BIGINT PRIMARY KEY REFERENCES state.virtual_batch (batch_num) ON DELETE CASCADE,
    tx_hash    VARCHAR,
    aggregator VARCHAR,
    block_num  BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.forced_batch
(
    forced_batch_num BIGINT PRIMARY KEY,
    global_exit_root VARCHAR,
    timestamp        TIMESTAMP,
    raw_txs_data     VARCHAR,
    coinbase         VARCHAR,
    batch_num        BIGINT, -- It can be null if the batch state is not trusted
    block_num        BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.l2block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  VARCHAR                  NOT NULL,
    header      jsonb,
    uncles      jsonb,
    parent_hash VARCHAR,
    state_root  VARCHAR,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    batch_num   BIGINT                   NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE
);

CREATE TABLE state.transaction
(
    hash         VARCHAR PRIMARY KEY,
    encoded      VARCHAR NOT NULL,
    decoded      jsonb,
    l2_block_num BIGINT  NOT NULL REFERENCES state.l2block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.exit_root
(
    block_num            BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    global_exit_root_num BIGINT PRIMARY KEY,
    mainnet_exit_root    BYTEA,
    rollup_exit_root     BYTEA,
    global_exit_root     BYTEA
);

CREATE TABLE state.sync_info
(
    last_batch_num_seen         BIGINT,
    last_batch_num_consolidated BIGINT,
    init_sync_batch             BIGINT
);

-- Insert default values into sync_info table
INSERT INTO state.sync_info (last_batch_num_seen, last_batch_num_consolidated, init_sync_batch)
VALUES (0, 0, 0);

CREATE TABLE state.receipt
(
    tx_hash             VARCHAR NOT NULL PRIMARY KEY REFERENCES state.transaction (hash) ON DELETE CASCADE,
    type                integer,
    post_state          BYTEA,
    status              BIGINT,
    cumulative_gas_used BIGINT,
    gas_used            BIGINT,
    block_num           BIGINT  NOT NULL REFERENCES state.l2block (block_num) ON DELETE CASCADE,
    tx_index            integer,
    contract_address    VARCHAR
);

CREATE TABLE state.log
(
    tx_hash   VARCHAR NOT NULL REFERENCES state.transaction (hash) ON DELETE CASCADE,
    log_index integer,
    address   VARCHAR NOT NULL,
    data      VARCHAR,
    topic0    VARCHAR NOT NULL,
    topic1    VARCHAR,
    topic2    VARCHAR,
    topic3    VARCHAR,
    PRIMARY KEY (tx_hash, log_index)
);

CREATE TABLE state.proof
(
    batch_num  BIGINT NOT NULL PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    proof jsonb
);
