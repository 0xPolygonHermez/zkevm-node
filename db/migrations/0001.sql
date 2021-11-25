-- +migrate Up

-- NOTE: We use "DECIMAL(78,0)" to encode go *big.Int types.  All the *big.Int
-- that we deal with represent a value in the SNARK field, which is an integer
-- of 256 bits.  `log(2**256, 10) = 77.06`: that is, a 256 bit number can have
-- at most 78 digits, so we use this value to specify the precision in the
-- PostgreSQL DECIMAL guaranteeing that we will never lose precision.

-- History
CREATE TABLE block
(
    eth_block_num BIGINT PRIMARY KEY,
    received_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    hash          BYTEA                       NOT NULL
);

CREATE TABLE transaction
(
    hash      VARCHAR PRIMARY KEY,
    encoded   VARCHAR,
    decoded   jsonb,
    batch_num BIGINT NOT NULL REFERENCES batch (batch_num) ON DELETE CASCADE
);

CREATE TABLE batch
(
    batch_num     BIGINT PRIMARY KEY,
    hash          varchar,
    eth_block_num BIGINT         NOT NULL REFERENCES block (eth_block_num) ON DELETE CASCADE,
    state_root    DECIMAL(78, 0) NOT NULL,
    is_virtual    bool
)