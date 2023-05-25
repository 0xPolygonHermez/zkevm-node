create user prover with encrypted password '${PROVER_PASSWORD:-default_prover_password}';
grant usage on schema state to prover;
grant select on state.merkletree to prover;
