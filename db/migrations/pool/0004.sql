-- +migrate Up
ALTER TABLE pool.transaction
ALTER COLUMN ip SET NOT NULL;
ALTER TABLE pool.transaction
ALTER COLUMN ip SET DEFAULT '';

-- +migrate Down
ALTER TABLE pool.transaction
ALTER COLUMN ip DROP NOT NULL;
ALTER TABLE pool.transaction
ALTER COLUMN ip DROP DEFAULT;
