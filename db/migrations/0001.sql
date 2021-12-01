-- +migrate Up

-- NOTE: We use "DECIMAL(78,0)" to encode go *big.Int types.  All the *big.Int
-- that we deal with represent a value in the SNARK field, which is an integer
-- of 256 bits.  `log(2**256, 10) = 77.06`: that is, a 256 bit number can have
-- at most 78 digits, so we use this value to specify the precision in the
-- PostgreSQL DECIMAL guaranteeing that we will never lose precision.

CREATE SCHEMA state

-- History
CREATE TABLE state.block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  BYTEA                       NOT NULL,
    parent_hash BYTEA,

    received_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
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
    raw_txs_data         BYTEA
);

CREATE TABLE state.transaction
(
    hash         BYTEA PRIMARY KEY,
    from_address BYTEA,
    encoded      VARCHAR,
    decoded      jsonb,
    batch_num    BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE
);
    
CREATE TABLE state.sequencer
(
    address     BYTEA NOT NULL,
    url         VARCHAR NOT NULL,
    chain_id    BIGINT PRIMARY KEY,
    block_num   BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);
