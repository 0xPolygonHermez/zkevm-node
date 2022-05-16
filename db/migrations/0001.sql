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
    rollup_exit_root     BYTEA,

    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    consolidated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX ON state.batch((header->>'stateRoot'));

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

CREATE TABLE state.log
(
    log_index integer,
    transaction_index integer,
    transaction_hash BYTEA NOT NULL REFERENCES state.transaction (hash) ON DELETE CASCADE,
    batch_hash BYTEA NOT NULL,
    batch_num BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    address BYTEA NOT NULL,
    data BYTEA,
    topic0 BYTEA NOT NULL,
    topic1 BYTEA,
    topic2 BYTEA,
    topic3 BYTEA
);

CREATE TABLE state.misc
(
    last_batch_num_seen BIGINT,
    last_batch_num_consolidated BIGINT,
    init_sync_block BIGINT
);

-- Insert default values into misc table
INSERT INTO state.misc (last_batch_num_seen, last_batch_num_consolidated, init_sync_block) VALUES(0, 0, 0);
