## Configure the Node: Different modes of execution

### Sync-only:

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

### RPC:

The `zkevm-rpc` component will act as a relay between the Trusted Sequencer and the Synchronizer (`zkevm-sync`). 

The [`production-setup.md`](./production-setup.md) goes through the setup of both a synchronizer and RPC components of the node.

### With Prover: Aggregator mode:

Node Config:

```toml
[RPC]
...
SequencerNodeURI = ""
BroadcastURI = "zkevm-broadcast:61090"
```

You will need to spin up the `zkevm-broadcast` service, which is a subcommand of the Node (`broadcast-trusted-state`). Find how to do it via the `test/docker-compose.yaml` file.

Use stock Prover config for Merkle Tree and Executor `zkevm-prover` image.

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

Run the Prover separately from the MerkleTree/Executor.

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
