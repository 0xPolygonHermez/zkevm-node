-- +migrate Up
CREATE INDEX IF NOT EXISTS l2block_block_hash_idx ON state.l2block (block_hash);

DELETE FROM state.sequences a USING (
    SELECT MIN(ctid) as ctid, from_batch_num
    FROM state.sequences 
    GROUP BY from_batch_num HAVING COUNT(*) > 1
) b
WHERE a.from_batch_num = b.from_batch_num 
AND a.ctid <> b.ctid;

ALTER TABLE state.sequences ADD PRIMARY KEY(from_batch_num);
ALTER TABLE state.trusted_reorg  ADD PRIMARY KEY(timestamp);
ALTER TABLE state.sync_info ADD PRIMARY KEY(last_batch_num_seen, last_batch_num_consolidated, init_sync_batch);

-- +migrate Down
DROP INDEX IF EXISTS state.l2block_block_hash_idx;

ALTER TABLE state.sequences DROP CONSTRAINT sequences_pkey;
ALTER TABLE state.trusted_reorg DROP CONSTRAINT trusted_reorg_pkey;
ALTER TABLE state.sync_info DROP CONSTRAINT sync_info_pkey;
