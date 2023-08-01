-- +migrate Up
ALTER TABLE pool.transaction ADD COLUMN break_even_gas_price DECIMAL(78, 0);
ALTER TABLE pool.gas_price ADD COLUMN l1_price DECIMAL(78, 0);


-- +migrate Down
ALTER TABLE pool.transaction DROP COLUMN break_even_gas_price;
ALTER TABLE pool.gas_price DROP COLUMN l1_price;
