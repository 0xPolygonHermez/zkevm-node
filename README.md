# hermez-core

Hermez Core is a Go implementation of a node that operates the zkEVM Hermez Network.

## About the zkEVM Hermez Network

Since this is an implementation of a protocol it's fundamental to understand it, [here](https://hackmd.io/tEny6MhSQaqPpu4ltUyC_w) you can find the specification of the protocol.

Glossary:

- Batch: equivalent of a block, but in the L2 realm. Note that the `jsonrpc` uses "block" instead of "batch" as it has to be compatible with the API of L1, so for this particular package block means batch.
- Virtual state: state reached by virtual batches, which are sequences of transactions sent to L1 by sequencers. In order to understand the outcome of executing this transactions, the nodes have to synchronize those batches and execute the sequences.
- Consolidated state: state that is proven on-chain by submitting a ZKP (Zero Knowledge Proof) that proofs the execution of a sequence of the last unconsolidated batch (virtual batch that has not been consolidated yet).
- Invalid transaction: transaction that can't be processed and doesn't affect the state. Note that such a transaction could be included in a virtual batch.
- Reverted transaction: transaction that is executed, but is reverted (because of smart contract logic). Main difference with *invalid transaction* is that this transaction modifies the state, at least to increment nonce off the sender.
- v1.0: refers to the payment rollup already in production, and has nothing to do with this project (zkEVM).
- v1.3: first iteration of the zkEVM project, includes processing Ether transfers with ZKPs
- v1.5: the Node implemented in this repo is able to provide most of the expected functionalities by a EVM compatible network, as an example, executing Uniswap smart contracts. However ZKPs will not be included in this version
- v2.0: Implements all the RPC features, and uses ZKPs to process any kind of transaction

## Architecture

*It's important to note that currently this implementation is a monolith monorepo, but this is very likely to change in the near future.*

<p align="center">
  <img src="./docs/architecture.drawio.png"/>
</p>

The diagram represents the main components of the software and how they interact between them:

- RPC: interface that allow users (metamask, etherscan, ...) to interact with the node. Fully compatible with Ethereum RPC + some extra endpoints specifics of the network.
- Pool: DB that stores txs to be selected/discarded by the `sequencer`
- Sequencer: get txs from the `pool`, check if they are valid by processing them using the `state` and group a sequence of txs into a batch that then sends to Ethereum using the `etherman`. The goal of the sequencer is to maximize profit extracted of the txs gas fees, as described in the PoE protocol.
- Etherman: abstraction that implements the needed methods to interact with the Ethereum network and the relevant smart contracts.
- Synchronizer: Updates the `state` by fetching data from Ethereum through the `etherman`.
- State: process txs using the EVM and the Merkletree, and manage all the metadata (batches, events, receipts, logs, ...).
- Aggregator: consolidates batches by generating ZK proofs. To do so it gathers the necessary data that the `prover` needs as input. Once the prove is generated it's sent to Ethereum through the `etherman`.
- Prover: service that generates ZK proofs. Note that this component is not implemented in this repository, and it's treated as a "black box" from the perspective of the node. The prover has two implementations: [JS reference implementation](https://github.com/hermeznetwork/zkproverjs) and [C production ready implementation](https://github.com/hermeznetwork/zkproverc).

## Development

It's recommended to use `make` for building, testing the code, ... Run `make help` to get a list of the available commands.

### Requirements

- Go 1.16
- Docker
- Docker Compose
- Make
- GCC

## Contribute

Before creating a PR, please read this [guide](docs/contribute-guie.md)
