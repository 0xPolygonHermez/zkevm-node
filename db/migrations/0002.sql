-- +migrate Up

CREATE SCHEMA pool

CREATE TABLE pool.txs (
    hash      VARCHAR PRIMARY KEY,
    encoded   VARCHAR,
    decoded   jsonb,
    state     varchar(15),
    gas_price DECIMAL(78,0),
    nonce     DECIMAL(78,0),
    received_at timestamp
);

CREATE INDEX idx_state_gas_price_nonce ON pool.txs(state, gas_price, nonce);

CREATE TABLE pool.gas_price (
    item_id SERIAL PRIMARY KEY,
    price DECIMAL(78,0),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);