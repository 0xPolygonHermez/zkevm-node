-- +migrate Down
DROP SCHEMA IF EXISTS pool CASCADE;

-- +migrate Up
CREATE SCHEMA pool;

CREATE TABLE pool.transaction
(
    hash                   VARCHAR PRIMARY KEY,
    encoded                VARCHAR,
    decoded                jsonb,
    status                 varchar(15),
    gas_price              DECIMAL(78, 0),
    nonce                  DECIMAL(78, 0),
    is_claims              BOOLEAN,
    cumulative_gas_used    BIGINT,
    used_keccak_hashes     INTEGER,
    used_poseidon_hashes   INTEGER,
    used_poseidon_paddings INTEGER,
    used_mem_aligns        INTEGER,
    used_arithmetics       INTEGER,
    used_binaries          INTEGER,
    used_steps             INTEGER,
    failed_counter         DECIMAL(78, 0) DEFAULT 0,
    received_at            TIMESTAMP WITH TIME ZONE NOT NULL,
    from_address           varchar                  NOT NULL
);

CREATE INDEX idx_state_gas_price_nonce ON pool.transaction (status, gas_price, nonce);
CREATE INDEX idx_failed_counter ON pool.transaction (failed_counter);

CREATE TABLE pool.gas_price
(
    item_id   SERIAL PRIMARY KEY,
    price     DECIMAL(78, 0),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);
