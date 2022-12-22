-- +migrate Down
DROP SCHEMA IF EXISTS txman CASCADE;

-- +migrate Up
CREATE SCHEMA txman;

CREATE TABLE txman.monitored_txs
(
    hash          VARCHAR PRIMARY KEY,
    raw_tx        VARCHAR,
    status        VARCHAR,
);
