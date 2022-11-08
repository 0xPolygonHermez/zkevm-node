-- +migrate Up
ALTER TABLE state.batch ADD COLUMN acc_input_hash VARCHAR;

-- +migrate Down
ALTER TABLE state.batch DROP COLUMN acc_input_hash;

