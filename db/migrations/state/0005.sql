-- +migrate Up
DROP table state.event;
DROP table state.debug;

-- +migrate Down
CREATE TABLE state.event
(
    event_type  VARCHAR NOT NULL,
    timestamp   TIMESTAMP WITH TIME ZONE NOT NULL,
    ip          VARCHAR,
    tx_hash     VARCHAR,
    payload     VARCHAR
);

CREATE TABLE state.debug
(
    error_type  VARCHAR,
    timestamp timestamp,
    payload VARCHAR  
);
