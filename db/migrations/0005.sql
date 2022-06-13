-- +migrate Up
CREATE SCHEMA statev2;

CREATE TABLE statev2.block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  BYTEA                       NOT NULL,
    parent_hash BYTEA,

    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE statev2.batch (
    id SERIAL PRIMARY KEY,
    global_exit_root BYTEA,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.sequenced_batch (
    batch_num BIGINT PRIMARY KEY,
    force_batch_num BIGINT,
    sequencer BYTEA,
    sequenced_at TIMESTAMP WITH TIME ZONE NOT NULL,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.verified_batch (
    batch_num BIGINT PRIMARY KEY,
    aggregator BYTEA,
    consolidated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.forced_batch (
    forced_batch_num BIGINT PRIMARY KEY,
    sequencer BYTEA,
    forced_at TIMESTAMP WITH TIME ZONE NOT NULL,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.transaction (
    tx_hash BYTEA PRIMARY KEY,
    tx_index BIGINT,
    from_address BYTEA,
    encoded VARCHAR,
    decoded jsonb,
    header jsonb,
    uncles jsonb,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);