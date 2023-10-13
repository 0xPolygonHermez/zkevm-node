# Component: Synchronizer

## XAGON Synchronizer:

The XAGON Synchronizer is the **base** component for which all others will depend on. You can *mix and match* different components to achieve a different outcome, be it sending transactions or computing proofs, but the Sync module will need to be up and running.

This module syncs data between the Layer 1 Ethereum network and XAGON L2 network.

## Running:

The preferred way to run the XAGON Synchronizer component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/xagon-node
```

To orchestrate multiple deployments of the different XAGON Node components, a `docker-compose.yaml` file for Docker Compose can be used:

**THIS STEP IS MANDATORY ON ALL DEPLOYMENT MODES**

```yaml
  xagon-sync:
    container_name: xagon-sync
    image: xagon-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/xagon-node run --genesis /app/genesis.json --cfg /app/config.toml --components synchronizer"
```

The container alone needs some parameters configured, access to certain configuration files and the appropriate ports exposed.

- environment: Env variables that supersede the config file
    - `ZKEVM_NODE_STATEDB_HOST`: Name of StateDB Database Host
- volumes:
    - `your config.toml file`: /app/config.toml
    - `your genesis.json file`: /app/genesis.json
