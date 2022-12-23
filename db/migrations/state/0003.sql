-- +migrate Down
DROP SCHEMA IF EXISTS txman CASCADE;

-- +migrate Up
CREATE SCHEMA txman;

CREATE TABLE txman.monitored_txs
(
    id         VARCHAR NOT NULL PRIMARY KEY,
    from       VARCHAR NOT NULL,
    to         VARCHAR,
    nonce      DECIMAL(78, 0) NOT NULL,
    value      DECIMAL(78, 0),
    data       VARCHAR,
    gas        BIGINT NOT NULL
    gas_price  DECIMAL(78, 0) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    status     VARCHAR,
    history    VARCHAR[]
);
