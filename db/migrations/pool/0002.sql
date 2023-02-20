-- +migrate Down
ALTER TABLE pool.transactions ADD COLUMN preprocessed_state_root varchar;

-- +migrate Up
ALTER TABLE pool.transactions DROP COLUMN preprocessed_state_root;