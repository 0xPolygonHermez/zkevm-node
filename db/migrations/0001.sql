-- +migrate Up
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
    block_num BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    tx_hash BYTEA NOT NULL REFERENCES state.transaction (hash) ON DELETE CASCADE,
    tx_index integer
);

CREATE TABLE state.misc
(
    last_batch_num_seen BIGINT
);

CREATE OR REPLACE FUNCTION populate_misc() RETURNS TRIGGER 
AS 'BEGIN INSERT INTO state.misc (last_batch_num_seen) VALUES(0); RETURN NULL; END;'
LANGUAGE plpgsql;

-- This trigger ensures a register is always present in the table and Updates will always work
-- This situation should never happen, just by human mistakes, as there is no code
-- to delete on this table. It does not affect performance as is "after delete"

CREATE TRIGGER delete_misc
    AFTER DELETE ON state.misc
    EXECUTE PROCEDURE populate_misc();

DELETE FROM state.misc;
