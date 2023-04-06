# Configure the Node: Different modes of execution

## JSON RPC Service:

This service will sync transactions from L2 to L1, it does not require a Prover to be working, just an RPC and Synchronizer. 

### Services needed:

*Please perform each of these steps (downloading and running) before continuing!*

- [RPCDB and StateDB Database](./components/databases.md)
- [Synchronizer](./components/synchronizer.md)
- [RPC](./components/rpc.md)
- [MT and Executor](./components/prover.md)

By default the config files found in the repository will spin up the Node in JSON RPC Service mode, which will not require a Prover (but will require a MT and Executor service).

**This is considered to be the base, all modes require a Synchronizer Node container to be spun up*

This will syncronize with the Trusted Sequencer (run by Polygon).

Use the default [public config file](https://github.com/0xPolygonHermez/zkevm-node/blob/develop/config/environments/public/public.node.config.toml), and make sure the following values are set to:

```toml
[RPC]
...
SequencerNodeURI = "https://public.zkevm-test.net:2083"
```

Same goes for the Prover Config ([prover-config.json](https://github.com/0xPolygonHermez/zkevm-node/blob/develop/config/environments/public/public.prover.config.json)):

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

Additionally, the [`production-setup.md`](./production-setup.md) goes through the setup of both a synchronizer and RPC components of the node.

### Docker services:

- `zkevm-sync`
- `zkevm-prover` (`Merkle Tree`, `Executor`)
- `zkevm-rpc` 
- Databases

## If you want to create Proofs:

This mode is a tad more complicated, as it will require more services and more machines:

Requirements for the Prover service (sans MT/Executor): 1TB RAM 128 cores

### Services needed: 

*Please perform each of these steps (downloading and running) before continuing!*

- [StateDB Database](./components/databases.md)
- [Synchronizer](./components/synchronizer.md)
- [Aggregator](./components/aggregator.md)
- [Prover, MT and Executor](./components/prover.md)

Machine 0:

- Synchronizer
- Aggregator
- MT and Executor
- Databases

Machine 1:

- Prover only

#### Machine 1

Use default [prover config](https://github.com/0xPolygonHermez/zkevm-node/blob/develop/config/environments/public/public.prover.config.json) but change the following values (`runProverServer` set to true, rest false):

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

### Docker services:

- `zkevm-sync`
- `zkevm-prover` (`Prover`, `Merkle Tree`, `Executor`)
- `zkevm-aggregator` 
- Databases