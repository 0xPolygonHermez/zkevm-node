-- +migrate Up
ALTER TABLE state.transaction
    ADD COLUMN effective_percentage SMALLINT;

-- +migrate Down
ALTER TABLE state.transaction
    DROP COLUMN effective_percentage;