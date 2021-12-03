-- +migrate Up

-- Table that stores all MerkleTree nodes
CREATE TABLE state.merkletree
(
    hash BYTEA PRIMARY KEY,
    data BYTEA NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS state.merkletree;