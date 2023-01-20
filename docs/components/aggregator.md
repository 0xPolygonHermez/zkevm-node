# Component: Aggregator

## ZKEVM Aggregator:

The ZKEVM Aggregator is an optional module responsible for receiving connections from Prover(s) in order to generate the proofs for the batches not proven yet.

## Hard dependencies:

- [Synchronizer](./synchronizer.md)
- [StateDB Database](./databases.md)
- [Prover, Merkle Tree and Executor](./prover.md)

## Running:

The preferred way to run the ZKEVM Aggregator component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/zkevm-node
```

To orchestrate multiple deployments of the different ZKEVM Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  zkevm-aggregator:
    container_name: zkevm-aggregator
    image: zkevm-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/zkevm-node run --genesis /app/genesis.json --cfg /app/config.toml --components aggregator"
```

The container alone needs some parameters configured, access to certain configuration files and the appropriate ports exposed.

- volumes:
    - `your Account Keystore file`: /pk/keystore (note, this `/pk/keystore` value is the default path that's written in the Public Configuration files on this repo, meant to expedite deployments, it can be superseded via an env flag `ZKEVM_NODE_ETHERMAN_PRIVATEKEYPATH`.)
    - `your config.toml file`: /app/config.toml
    - `your genesis.json file`: /app/genesis.json

- environment: Env variables that supersede the config file
    - `ZKEVM_NODE_STATEDB_HOST`: Name of StateDB Database Host

### The Account Keystore file:

Since the Aggregator will send transactions to L1 you'll need to generate an account keystore:

[Generate an Account Keystore file](./account_keystore.md)
