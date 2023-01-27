-- +migrate Up
<<<<<<< HEAD
ALTER TABLE state.forced_batch
DROP COLUMN IF EXISTS batch_num;

ALTER TABLE state.batch
ADD COLUMN forced_batch_num BIGINT;
ALTER TABLE state.batch
ADD FOREIGN KEY (forced_batch_num) REFERENCES state.forced_batch(forced_batch_num);

=======
>>>>>>> develop
CREATE TABLE state.debug
(
    error_type  VARCHAR,
    timestamp timestamp,
    payload VARCHAR  
);

-- +migrate Down
<<<<<<< HEAD
ALTER TABLE state.batch
DROP COLUMN IF EXISTS forced_batch_num;

ALTER TABLE state.forced_batch
ADD COLUMN batch_num BIGINT;

DROP TABLE state.debug;
=======
DROP TABLE state.debug;
>>>>>>> develop
