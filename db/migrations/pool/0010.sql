-- +migrate Up
ALTER TABLE pool.transaction ADD COLUMN break_even_gas_price DECIMAL(78, 0);

-- +migrate Down
ALTER TABLE pool.transaction DROP COLUMN break_even_gas_price;
