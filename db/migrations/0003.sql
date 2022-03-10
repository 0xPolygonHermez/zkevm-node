-- +migrate Down
DROP TABLE IF EXISTS state.merkletree;
DROP TABLE IF EXISTS state.sc_code;

-- +migrate Up

-- Table that stores all MerkleTree nodes
CREATE TABLE state.merkletree
(
    hash BYTEA PRIMARY KEY,
    data BYTEA NOT NULL
);

-- Table that stores all smart contract code
CREATE TABLE state.sc_code
(
    hash BYTEA PRIMARY KEY,
    data BYTEA
);
