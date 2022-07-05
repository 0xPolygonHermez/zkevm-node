-- +migrate Down
DROP SCHEMA IF EXISTS state CASCADE;

-- +migrate Up
CREATE SCHEMA state;

-- table schema from https://github.com/hermeznetwork/zkproverc/blob/94961614779317436cd90f1a48d4965df76f2695/src/database.cpp
CREATE TABLE state.merkletree
(
    hash BYTEA PRIMARY KEY,
    data BYTEA NOT NULL
);
