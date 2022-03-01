-- +migrate Up
CREATE SCHEMA state

-- History
CREATE TABLE state.block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  BYTEA                       NOT NULL,
    parent_hash BYTEA,

    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE state.batch
(
    batch_num            BIGINT PRIMARY KEY,
    batch_hash           BYTEA,
    block_num            BIGINT                      NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    sequencer            BYTEA,
    aggregator           BYTEA,
    consolidated_tx_hash BYTEA,
    header               jsonb,
    uncles               jsonb,
    raw_txs_data         BYTEA,
    matic_collateral     NUMERIC(78,0),
    chain_id             BIGINT,
    global_exit_root     BYTEA,

    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    consolidated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE state.transaction
(
    hash         BYTEA PRIMARY KEY,
    from_address BYTEA,
    encoded      VARCHAR,
    decoded      jsonb,
    batch_num    BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    tx_index     integer
);

CREATE TABLE state.sequencer
(
    address     BYTEA NOT NULL,
    url         VARCHAR NOT NULL,
    chain_id    BIGINT PRIMARY KEY,
    block_num   BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.receipt
(
    type integer,
    post_state BYTEA,
    status BIGINT,
    cumulative_gas_used BIGINT,
    gas_used BIGINT,
    batch_num BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    batch_hash BYTEA NOT NULL,
    tx_hash BYTEA NOT NULL REFERENCES state.transaction (hash) ON DELETE CASCADE,
    tx_index integer,
    tx_from BYTEA,
    tx_to BYTEA,
    contract_address BYTEA
);

CREATE TABLE state.misc
(
    last_batch_num_seen BIGINT,
    last_batch_num_consolidated BIGINT
);

CREATE TABLE state.exit_root
(
    global_exit_root_num BIGINT,
    mainnet_exit_root    BYTEA,
    rollup_exit_root     BYTEA
);

-- Insert default values into misc table
INSERT INTO state.misc (last_batch_num_seen, last_batch_num_consolidated) VALUES(0, 0);
