-- +migrate Up
UPDATE pool.transaction 
SET ip = '' WHERE ip IS NULL;

-- +migrate Down
-- Nothing to do here