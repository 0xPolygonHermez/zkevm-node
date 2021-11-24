-- +migrate Up

CREATE SCHEMA pool

CREATE TABLE pool.txs (
    hash      VARCHAR PRIMARY KEY,
    encoded   VARCHAR,
    decoded   jsonb,
    state     varchar(15)
);