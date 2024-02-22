-- +migrate Up
ALTER TABLE pool.transaction
    ADD COLUMN reserved_zkcounters jsonb;

-- +migrate Down
ALTER TABLE pool.transaction
    DROP COLUMN reserved_zkcounters;
