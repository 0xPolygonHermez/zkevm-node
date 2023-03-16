-- +migrate Up
ALTER TABLE pool.transaction
DROP COLUMN failed_counter;

ALTER TABLE pool.transaction
ADD COLUMN ip VARCHAR;

-- +migrate Down
ALTER TABLE pool.transaction
ADD COLUMN failed_counter DECIMAL(78, 0) DEFAULT 0;

ALTER TABLE pool.transaction
DROP COLUMN ip;
