-- +migrate Up
CREATE TABLE state.l1_tx_queue
(
    hash          VARCHAR PRIMARY KEY,
    raw_tx        VARCHAR,
    data          jsonb,
    type          VARCHAR,
    status        VARCHAR,
    wait_duration BIGINT
);

-- +migrate Down
DROP TABLE state.l1_tx_queue;
