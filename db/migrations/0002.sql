-- +migrate Up

CREATE SCHEMA pool

CREATE TABLE pool.txs (
    hash      VARCHAR PRIMARY KEY,
    encoded   VARCHAR,
    decoded   jsonb,
    state     varchar(15),
    received_at timestamp
);

-- create json indexes to query ordered by nonce and by tx state

CREATE TABLE pool.gas_price (
    item_id SERIAL PRIMARY KEY,
    price DECIMAL(78,0),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);