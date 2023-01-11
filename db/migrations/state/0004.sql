-- +migrate Up
CREATE TABLE state.debug
(
    error_type  VARCHAR,
    timestamp timestamp,
    payload VARCHAR  
);

-- +migrate Down
DROP TABLE state.debug;