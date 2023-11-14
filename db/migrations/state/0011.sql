-- +migrate Up
CREATE INDEX IF NOT EXISTS l2block_created_at_idx ON state.l2block (created_at);

CREATE INDEX IF NOT EXISTS log_log_index_idx ON state.log (log_index);
CREATE INDEX IF NOT EXISTS log_topic0_idx ON state.log (topic0);
CREATE INDEX IF NOT EXISTS log_topic1_idx ON state.log (topic1);
CREATE INDEX IF NOT EXISTS log_topic2_idx ON state.log (topic2);
CREATE INDEX IF NOT EXISTS log_topic3_idx ON state.log (topic3);

ALTER TABLE state.transaction ADD COLUMN egp_log JSONB;

-- +migrate Down
DROP INDEX IF EXISTS state.l2block_created_at_idx;

DROP INDEX IF EXISTS state.log_log_index_idx;
DROP INDEX IF EXISTS state.log_topic0_idx;
DROP INDEX IF EXISTS state.log_topic1_idx;
DROP INDEX IF EXISTS state.log_topic2_idx;
DROP INDEX IF EXISTS state.log_topic3_idx;

ALTER TABLE state.transaction DROP COLUMN egp_log;