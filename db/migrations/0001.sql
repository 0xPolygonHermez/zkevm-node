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
    batch_num    BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE
);
    
CREATE TABLE state.sequencer
(
    address     BYTEA NOT NULL,
    url         VARCHAR NOT NULL,
    chain_id    BIGINT PRIMARY KEY,
    block_num   BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE
);

CREATE TABLE state.misc
(
    last_batch_num_seen BIGINT
);

CREATE OR REPLACE FUNCTION populate_misc() RETURNS TRIGGER 
AS 'BEGIN INSERT INTO state.misc (last_batch_num_seen) VALUES(0); RETURN NULL; END;'
LANGUAGE plpgsql;

CREATE TRIGGER delete_misc
    AFTER DELETE ON state.misc
    EXECUTE PROCEDURE populate_misc();

DELETE FROM state.misc;