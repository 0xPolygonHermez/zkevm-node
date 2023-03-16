-- +migrate Up
CREATE TABLE state.event
(
    event_type  VARCHAR NOT NULL,
    timestamp   TIMESTAMP WITH TIME ZONE NOT NULL,
    ip          VARCHAR,
    tx_hash     VARCHAR,
    payload     VARCHAR
);

-- +migrate Down
DROP table state.event;
