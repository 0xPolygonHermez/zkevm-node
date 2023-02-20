-- +migrate Down
ALTER TABLE pool.transaction ADD COLUMN preprocessed_state_root varchar;

-- +migrate Up
ALTER TABLE pool.transaction DROP COLUMN preprocessed_state_root;