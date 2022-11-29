## Configure the Node: Different modes of execution

### Sync-only (*read-only*):

By default the config files found in the repository will spin up the Node in `sync-only` mode, which will not require a Prover (but will require a MT and Executor service).

**This is considered to be the base, all modes require a Synchronizer Node container to be spun up*

This will syncronize with the Trusted Sequencer (run by Polygon Hermez).

Config:

```toml
[RPC]
...
SequencerNodeURI = "https://public.zkevm-test.net:2083"
BroadcastURI = "public-grpc.zkevm-test.net:61090"
```

Prover Config:

```json
{
	...
    "runProverServer": false,
    "runProverServerMock": false,
    "runProverClient": false,
    "runExecutorServer": true,
    "runExecutorClient": false,
    "runStateDBServer": true
}
```

ZKEVM RPC component is also needed.

The `zkevm-rpc` component will act as a relay between the Trusted Sequencer and the Synchronizer (`zkevm-sync`). 

The [`production-setup.md`](./production-setup.md) goes through the setup of both a synchronizer and RPC components of the node.

##### Docker services:

- `zkevm-sync`
- `zkevm-prover` (`Merkle Tree`, `Executor`)
- `zkevm-rpc` 
- Databases

### Create Proofs: Base + RPC + Aggregator + Prover:

Use an Aggregator with a Prover to create proofs.

On a separate machine from the *Merkle Tree/Executor* `zkevm-prover` container, spin up a Prover:

Use stock Prover config for Merkle Tree and Executor `zkevm-prover` image on the other machine.

For *only* Prover Config (`only-prover-config.json`):

```json
{
	...
    "runProverServer": true,
    "runProverServerMock": false,
    "runProverClient": false,
    "runExecutorServer": false,
    "runExecutorClient": false,
    "runStateDBServer": false
}
```

*docker-compose.yaml*:

```yaml
  zkevm-only-prover:
    container_name: zkevm-prover
    image: hermeznetwork/zkevm-prover:develop
    ports:
      - 50051:50051 # Prover
    volumes:
      - ./only-prover-config.json:/usr/src/app/config.json
    command: >
      zkProver -c /usr/src/app/config.json
```

For aggregator, here's how to spin it up using docker compose:

```yaml
  zkevm-aggregator:
    container_name: zkevm-aggregator
    image: zkevm-node
    environment:
      - ZKEVM_NODE_STATEDB_HOST=zkevm-state-db
    volumes:
      - ./config.toml:/app/config.toml
      - ./genesis.json:/app/genesis.json
    command:
      - "/bin/sh"
      - "-c"
      - "/app/zkevm-node run --genesis /app/genesis.json --cfg /app/config.toml --components aggregator"
```