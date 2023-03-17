-- +migrate Up
UPDATE TABLE pool.transaction 
SET ip = '' WHERE ip IS NULL;

-- +migrate Down
-- Nothing to do here