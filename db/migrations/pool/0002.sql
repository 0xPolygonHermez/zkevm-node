-- +migrate Up
ALTER TABLE pool.transaction
ADD COLUMN is_wip BOOLEAN;

-- +migrate Down
ALTER TABLE pool.transaction
DROP COLUMN is_wip;
