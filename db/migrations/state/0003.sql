-- +migrate Up
ALTER TABLE state.log
ALTER COLUMN topic0 DROP NOT NULL;

-- +migrate Down
ALTER TABLE state.log
ALTER COLUMN topic0 SET NOT NULL;