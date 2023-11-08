-- +migrate Up
ALTER TABLE state.monitored_txs
    ADD COLUMN gas_offset DECIMAL(78, 0) NOT NULL DEFAULT 0;
ALTER TABLE state.monitored_txs ALTER COLUMN gas_offset DROP DEFAULT;

-- +migrate Down
ALTER TABLE state.monitored_txs
    DROP COLUMN gas_offset;