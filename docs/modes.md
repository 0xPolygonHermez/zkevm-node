## Configure the Node: Different modes of execution

### Sync-only:

By default the config files found in the repository will spin up the Node in `sync-only` mode, which will not require a Prover (but will require a MT and Executor service).

It will syncronize with the Trusted Sequencer (run by Polygon Hermez).

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

### With Prover:

This will act as a Trusted Sequencer:

```toml
[RPC]
...
SequencerNodeURI = ""
BroadcastURI = "zkevm-broadcast:61090"
```

You will need to spin up the `zkevm-broadcast` service, which is a subcommand of the Node (`broadcast-trusted-state`). Find how to do it via the `test/docker-compose.yaml` file.

Prover Config:

```json
{
	...
    "runProverServer": true,
    "runProverServerMock": false,
    "runProverClient": false,
    "runExecutorServer": true,
    "runExecutorClient": false,
    "runStateDBServer": true
}
```