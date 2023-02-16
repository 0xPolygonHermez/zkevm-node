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

CREATE TABLE state.forced_batch
(
    forced_batch_num BIGINT PRIMARY KEY,
    global_exit_root VARCHAR,
    timestamp        TIMESTAMP WITH TIME ZONE NOT NULL,
    raw_txs_data     VARCHAR,
    coinbase         VARCHAR,
    block_num        BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.batch
( --batch abstraction: will be created through trusted state
    batch_num        BIGINT PRIMARY KEY,
    global_exit_root VARCHAR,
    local_exit_root  VARCHAR,
    state_root       VARCHAR,
    acc_input_hash   VARCHAR,
    timestamp        TIMESTAMP WITH TIME ZONE,
    coinbase         VARCHAR,
    raw_txs_data     BYTEA,
    forced_batch_num BIGINT REFERENCES state.forced_batch(forced_batch_num)
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
    state_root VARCHAR,
    block_num  BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    is_trusted BOOLEAN DEFAULT true
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
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
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
    id                   SERIAL PRIMARY KEY,
    block_num            BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    timestamp            TIMESTAMP WITH TIME ZONE NOT NULL,
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
    batch_num  BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    batch_num_final BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    proof VARCHAR,
    proof_id VARCHAR,
    input_prover VARCHAR,
    prover VARCHAR,
    generating BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (batch_num, batch_num_final)    
);

CREATE TABLE IF NOT EXISTS state.sequences
( --Allowed Verifications
    from_batch_num BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    to_batch_num   BIGINT REFERENCES state.batch (batch_num) ON DELETE CASCADE
);

CREATE TABLE state.monitored_txs
(
    owner      VARCHAR NOT NULL,
    id         VARCHAR NOT NULL,
    from_addr  VARCHAR NOT NULL,
    to_addr    VARCHAR,
    nonce      DECIMAL(78, 0) NOT NULL,
    value      DECIMAL(78, 0),
    data       VARCHAR,
    gas        DECIMAL(78, 0) NOT NULL,
    gas_price  DECIMAL(78, 0) NOT NULL,
    status     VARCHAR NOT NULL,
    history    VARCHAR[],
    block_num  BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (owner, id)
);

CREATE TABLE state.debug
(
    error_type  VARCHAR,
    timestamp timestamp,
    payload VARCHAR  
);