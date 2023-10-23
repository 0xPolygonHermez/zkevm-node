-- +migrate Up
CREATE TABLE pool.whitelisted (
	addr VARCHAR PRIMARY KEY
);

-- +migrate Down
DROP TABLE pool.whitelisted;
