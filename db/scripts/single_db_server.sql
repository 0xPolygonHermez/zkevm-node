CREATE DATABASE state_db;
CREATE DATABASE pool_db;
CREATE DATABASE rpc_db;

CREATE DATABASE prover_db;
\connect prover_db;

CREATE SCHEMA state;

CREATE TABLE state.nodes (hash BYTEA PRIMARY KEY, data BYTEA NOT NULL);
CREATE TABLE state.program (hash BYTEA PRIMARY KEY, data BYTEA NOT NULL);

CREATE USER prover_user with password 'prover_pass';
ALTER DATABASE prover_db OWNER TO prover_user;
ALTER SCHEMA state OWNER TO prover_user;
ALTER SCHEMA public OWNER TO prover_user;
ALTER TABLE state.nodes OWNER TO prover_user;
ALTER TABLE state.program OWNER TO prover_user;
ALTER USER prover_user SET SEARCH_PATH=state;