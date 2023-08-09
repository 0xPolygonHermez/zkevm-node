-- +migrate Up
ALTER TABLE state.transaction
    ADD COLUMN effective_percentage SMALLINT DEFAULT 255;

ALTER TABLE state.receipt
    ADD COLUMN effective_gas_price BIGINT;

-- +migrate Down
ALTER TABLE state.transaction
    DROP COLUMN effective_percentage;

ALTER TABLE state.receipt
    DROP COLUMN effective_gas_price;
