CREATE DATABASE prover_db;
\connect prover_db;

CREATE SCHEMA state;

CREATE TABLE state.merkletree (
    hash BYTEA PRIMARY KEY,
    data BYTEA NOT NULL
);

CREATE USER prover_user with password 'prover_pass';
GRANT CONNECT ON DATABASE prover_db TO prover_user;
GRANT USAGE ON SCHEMA state TO prover_user;
GRANT ALL PRIVILEGES ON TABLE state.merkletree TO prover_user;
