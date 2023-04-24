-- +migrate Up
ALTER TABLE state.batch
    ADD COLUMN batch_resources JSONB,
    ADD COLUMN closing_reason VARCHAR;

-- +migrate Down
ALTER TABLE state.batch
    DROP COLUMN batch_resources,
    DROP COLUMN closing_reason;