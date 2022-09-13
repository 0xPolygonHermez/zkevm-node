-- +migrate Down
DROP SCHEMA IF EXISTS rpc CASCADE;

-- +migrate Up
CREATE SCHEMA rpc;

CREATE TABLE rpc.filters
(
    id          SERIAL PRIMARY KEY,
    filter_type VARCHAR(15) NOT NULL,
    parameters  JSONB       NOT NULL,
    last_poll   TIMESTAMP   NOT NULL
);
