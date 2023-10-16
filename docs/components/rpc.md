# Component: RPC

## XGON RPC:

The XGON RPC relays transactions to the Trusted sequencer.

## Hard dependencies:

- [Synchronizer](./synchronizer.md)
- [StateDB Database](./databases.md)
- [RPCDB Database](./databases.md)
- [Merkle Tree and Executor](./prover.md)

## Running:

The preferred way to run the XGON RPC component is via Docker and Docker Compose.

```bash
docker pull okx/xgon-node
```

To orchestrate multiple deployments of the different XGON Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  xgon-rpc:
    container_name: xgon-rpc
    image: xgon-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/xgon-node run --genesis /app/genesis.json --cfg /app/config.toml --components rpc"
```

The container alone needs some parameters configured, access to certain configuration files and the appropriate ports exposed.

- ports:
    - `8545:8545`: RPC Port
    - `9091:9091`: Needed if Prometheus metrics are enabled
- environment: Env variables that supersede the config file
    - `XGON_NODE_STATEDB_HOST`: Name of StateDB Database Host
    - `XGON_NODE_POOL_HOST`: Name of PoolDB Database Host 
    - `XGON_NODE_RPC_DB_HOST`: Name of RPCDB Database Host
- volumes:
    - `your config.toml file`: /app/config.toml
    - `your genesis file`: /app/genesis.json
