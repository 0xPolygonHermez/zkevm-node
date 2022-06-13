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
    aggregator BYTEA,
    global_exit_root BYTEA,
    transactions BYTEA,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.sequenced_batch (
    batch_num BIGINT PRIMARY KEY,
    force_batch_num INTEGER,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.verified_batch (
    batch_num BIGINT PRIMARY KEY,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.forced_batch (
    batch_num BIGINT PRIMARY KEY,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.force_sequenced_batch (
    batch_num BIGINT PRIMARY KEY,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);