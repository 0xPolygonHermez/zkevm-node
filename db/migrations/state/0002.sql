-- +migrate Up
DROP TABLE state.proof
CREATE TABLE state.proof
(
    batch_num  BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    batch_num_final BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    proof VARCHAR,
    proof_id VARCHAR,
    input_prover VARCHAR,
    prover VARCHAR,
    generating BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (batch_num, batch_num_final)    
);

-- +migrate Down
DROP TABLE state.proof
CREATE TABLE state.proof
(
    batch_num  BIGINT NOT NULL PRIMARY KEY REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    proof jsonb,
    proof_id VARCHAR,
    input_prover jsonb,
    prover VARCHAR
);
