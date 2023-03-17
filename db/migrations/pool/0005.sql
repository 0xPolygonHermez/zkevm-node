-- +migrate Up
UPDATE TABLE pool.transaction SET ip = '' where ip is null;

-- +migrate Down
-- Nothing to do here