-- +migrate Up
ALTER TABLE state.log
ALTER COLUMN topic0 DROP NOT NULL;

-- +migrate Down
DELETE FROM state.log WHERE topic0 IS NULL;
ALTER TABLE state.log
ALTER COLUMN topic0 SET NOT NULL;