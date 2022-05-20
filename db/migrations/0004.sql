-- +migrate Up

CREATE SCHEMA rpc

CREATE TABLE rpc.filters (
    id          SERIAL PRIMARY KEY,
    filter_type VARCHAR(15),
    parameters  JSONB,
    last_poll   TIMESTAMP NOT NULL
);