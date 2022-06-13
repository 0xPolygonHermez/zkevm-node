-- +migrate Up
CREATE SCHEMA statev2 CREATE TABLE statev2.block (
    block_num BIGINT PRIMARY KEY,
    block_hash VARCHAR NOT NULL,
    parent_hash VARCHAR,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE statev2.batch (  --batch abstraction 
    id SERIAL PRIMARY KEY,
    global_exit_root VARCHAR,
    batch_hash VARCHAR,
    tx_hash VARCHAR,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.sequenced_batch (  --virtual state
    batch_num BIGINT PRIMARY KEY,
    force_batches_num INTEGER,  --if empty, it means this batch is from sequenceForceBatches
    tx_hash VARCHAR,
    sequencer VARCHAR,
    sequenced_at TIMESTAMP WITH TIME ZONE NOT NULL,
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id)
);

CREATE TABLE statev2.verified_batch (  --consolidated state
    batch_num BIGINT PRIMARY KEY,
    tx_hash VARCHAR,
    aggregator VARCHAR,
    consolidated_at TIMESTAMP WITH TIME ZONE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id),
    block_num BIGINT NOT NULL REFERENCES statev2.block (block_num) ON DELETE CASCADE
);

CREATE TABLE statev2.forced_batch (
    forced_batch_num BIGINT PRIMARY KEY,
    sequencer VARCHAR,
    forced_at TIMESTAMP WITH TIME ZONE NOT NULL,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

-- All txs must be stored here. We only will manage the relationship between transactions and batches regarding the reorg.
CREATE TABLE statev2.transaction (  --transaction abstraction
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR,
    from_address VARCHAR,
    encoded VARCHAR,
    decoded jsonb,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE statev2.tx_batch_relationship (
    tx_id BIGINT NOT NULL REFERENCES statev2.transaction (id) ON DELETE CASCADE,
    batch_id BIGINT NOT NULL REFERENCES statev2.batch (id) ON DELETE CASCADE
);

CREATE TABLE statev2.ts_transaction (  --transaction in the trusted sequencer
    tx_id BIGINT NOT NULL REFERENCES statev2.transaction (id) ON DELETE CASCADE,
    tx_index BIGINT PRIMARY KEY,  --l2Block identifier
    header jsonb,
    uncles jsonb
);