-- +migrate Up
ALTER TABLE state.proof ADD COLUMN proof_id VARCHAR;
ALTER TABLE state.proof ADD COLUMN input_prover jsonb;
ALTER TABLE state.proof ADD COLUMN prover VARCHAR;

-- +migrate Down
ALTER TABLE state.proof DROP COLUMN prover;
ALTER TABLE state.proof DROP COLUMN input_prover;
ALTER TABLE state.proof DROP COLUMN proof_id;


