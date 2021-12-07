-- +migrate Down
DROP TABLE IF EXISTS state.merkletree;
DROP TABLE IF EXISTS state.merkletree_roots;

-- +migrate Up

-- Table that stores all MerkleTree nodes
CREATE TABLE state.merkletree
(
    hash BYTEA PRIMARY KEY,
    data BYTEA NOT NULL
);

-- Table that stores MerkleTree root for each batch number
CREATE TABLE state.merkletree_roots
(
    batch_num BIGINT PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    hash BYTEA NOT NULL
);
