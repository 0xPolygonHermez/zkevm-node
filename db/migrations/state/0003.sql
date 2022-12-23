-- +migrate Down
DROP SCHEMA IF EXISTS txman CASCADE;

-- +migrate Up
CREATE SCHEMA txman;

CREATE TABLE txman.monitored_txs
(
    id         VARCHAR NOT NULL PRIMARY KEY,
    from_addr  VARCHAR NOT NULL,
    to_addr    VARCHAR,
    nonce      DECIMAL(78, 0) NOT NULL,
    value      DECIMAL(78, 0),
    data       VARCHAR,
    gas        DECIMAL(78, 0) NOT NULL,
    gas_price  DECIMAL(78, 0) NOT NULL,
    status     VARCHAR,
    history    VARCHAR[],
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
