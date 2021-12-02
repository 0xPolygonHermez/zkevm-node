-- +migrate Up

-- Table that stores all MerkleTree nodes
CREATE TABLE merkletree
(
    key  BYTEA PRIMARY KEY,
    data BYTEA NOT NULL
);
