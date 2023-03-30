-- +migrate Up
CREATE TABLE pool.blocked
(
    addr varchar NOT NULL PRIMARY KEY
);

-- +migrate Down
DROP TABLE pool.blocked;
