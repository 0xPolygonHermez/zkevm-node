-- +migrate Up
ALTER TABLE state.forced_batch
DROP COLUMN IF EXISTS batch_num;

ALTER TABLE state.batch
ADD COLUMN forced_batch_num BIGINT;
ALTER TABLE state.batch
ADD FOREIGN KEY (forced_batch_num) REFERENCES state.forced_batch(forced_batch_num);

CREATE TABLE state.closing_signals
(
    sent_forced_batch_timestamp  TIMESTAMP WITH TIME ZONE NOT NULL,
    sent_to_l1_timestamp  TIMESTAMP WITH TIME ZONE NOT NULL,
    last_ger VARCHAR
);

-- +migrate Down
ALTER TABLE state.batch
DROP COLUMN IF EXISTS forced_batch_num;

ALTER TABLE state.forced_batch
ADD COLUMN batch_num BIGINT;

DROP TABLE state.closing_signals