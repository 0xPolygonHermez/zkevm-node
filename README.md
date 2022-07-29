# zkEVM Node

zkEVM Node is a Go implementation of a node that operates the Polygon zkEVM Network.

## About the Polygon zkEVM network

Since this is an implementation of a protocol it's fundamental to understand it, [here](https://docs.hermez.io/) you can find the specification of the protocol.

Glossary:

- L1: Base blockchain where the rollup smart contracts are deployed. It's Ethereum or a testnet of Ethereum, but it could be any EVM compatible blockchain.
- L2: the rollup network aka the Polygon zkEVM network.
- Batch: a group of transactions that are executed/proved, using the [zkEVM prover](https://github.com/0xPolygonHermez/zkevm-prover) and sent to / synchronized from L1
- Sequencer: the actor that is responsible for selecting transactions, putting them in a specific order, and sending them in batches to L1
- Trusted sequencer: sequencer that has special privileges, there can only be one trusted sequencer. The privileges granted to the trusted sequencer allow it to forecast the batches that will be applied to L1. This way it can commit to a specific sequence before interacting with L1. This is done to achieve fast finality and reduce costs associated with using the network (lower gas fees)
- Permissionless sequencer: sequencer role that can be performed by anyone. It has competitive disadvantages compared to the trusted sequencer (slow finality, MEV attacks). Its main purpose is to provide censorship resistance and unstoppability features to the network.
- Sequence: Group of batches and other metadata that the trusted sequencer sends to L1 to update the state
- Forced batch: batch that is sent by permissionless sequencers to L1 to update the state
- L2 Block: Same as an L1 block, but for L2. This is mostly used by the JSON RPC interface. Currently, all the L2 Blocks are set to only include one transaction, this is done to achieve instant finality: it's not necessary to close a batch to allow the JSON RPC to expose results of already processed transactions
- Trusted state: state reached through processing transactions that have been shared by the trusted sequencer. This state is considered trusted as the trusted sequencer could commit to a certain sequence, and then send a different one to L1
- Virtual state: state reached through processing transactions that have already been submitted to L1. These transactions are sent in batches by either trusted or permissionless sequencers. Those batches are also called virtual batches. Note that this state is trustless as it relies on L1 security assumptions
- Consolidated state: state that is proven on-chain by submitting a ZKP (Zero Knowledge Proof) that proves the execution of a sequence of the last virtual batch.
- Invalid transaction: a transaction that can't be processed and doesn't affect the state. Note that such a transaction could be included in a virtual batch. The reason for a transaction to be invalid could be related to the Ethereum protocol (invalid nonce, not enough balance, ...) or due to limitations introduced by the zkEVM (each batch can make use of a limited amount of resources such as the total amount of keccak hashes that can be computed)
- Reverted transaction: a transaction that is executed, but is reverted (because of smart contract logic). The main difference with *invalid transaction* is that this transaction modifies the state, at least to increment nonce of the sender.
- Proof of Efficiency (PoE): name of the protocol used by the network, it's enforced by the [smart contracts](https://github.com/0xPolygonHermez/zkevm-contracts)

## Architecture

<p align="center">
  <img src="./docs/architecture.drawio.png"/>
</p>

The diagram represents the main components of the software and how they interact between them. Note that this reflects a single entity running a node, in particular a node that acts as the trusted sequencer. But there are many entities running nodes in the network, and each of these entities can perform different roles. More on this later.

- (JSON) RPC: an interface that allows users (metamask, etherscan, ...) to interact with the node. Fully compatible with Ethereum RPC + some extra endpoints specifics of the network. It interacts with the `state` to get data and process transactions and with the `pool` to store transactions
- Pool: DB that stores transactions by the `RPC` to be selected/discarded by the `sequencer` later on
- Trusted Sequencer: get transactions from the `pool`, check if they are valid by processing them using the `state`, and create sequences. Once transactions are added into the state, they are immediately available to the `broadcast` service. Sequences are sent to L1 using the `etherman`
- Broadcast: API used by the `synchronizer` of nodes that are not the `trusted sequencer` to synchronize the trusted state
- Permissionless Sequencer: *coming soon*
- Etherman: abstraction that implements the needed methods to interact with the Ethereum network and the relevant smart contracts.
- Synchronizer: Updates the `state` by fetching data from Ethereum through the `etherman`. If the node is not a `trusted sequencer` it also updates the state with the data fetched from the `broadcast` of the `trusted sequencer`. It also detects and handles reorgs that can happen if the `trusted sequencer` sends different data in the broadcast vs the sequences sent to L1 (trusted vs virtual state)
- State: Responsible for managing the state data (batches, blocks, transactions, ...) that is stored on the `state SB`. It also handles the integration with the `executor` and the `Merkletree` service
- State DB: persistence layer for the state data (except the Merkletree that is handled by the `Merkletree` service)
- Aggregator: consolidates batches by generating ZKPs (Zero Knowledge proofs). To do so it gathers the necessary data that the `prover` needs as input through the `state` and sends a request to it. Once the proof is generated it's sent to Ethereum through the `etherman`
- Prover/Executor: service that generates ZK proofs. Note that this component is not implemented in this repository, and it's treated as a "black box" from the perspective of the node. The prover/executor has two implementations: [JS reference implementation](https://github.com/0xPolygonHermez/zkevm-proverjs) and [C production-ready implementation](https://github.com/0xPolygonHermez/zkevm-prover). Although it's the same software/service, it has two very different purposes:
  - Provide an EVM implementation that allows processing transactions and getting all needed results metadata (state root, receipts, logs, ...)
  - Generate ZKPs
- Merkletree: service that stores the Merkletree, containing all the account information (balances, nonces, smart contract code, and smart contract storage). This component is also not implemented in this repo and is consumed as an external service by the node. The implementation can be found [here](https://github.com/0xPolygonHermez/zkevm-prover)

## Roles of the network

The node software is designed to support the execution of multiple roles. Each role requires different services to work. Most of the services can run in different instances, and the JSON RPC can run in many instances (all the other services must have a single instance)

### RPC

This role can be performed by anyone.

Required services and components:

- JSON RPC: can run in a separated instance, and can have multiple instances
- Synchronizer: single instance that can run on a separate instance
- Executor & Merkletree: service that can run on a separate instance
- State DB: Postgres SQL that can be run in a separate instance

There must be only one synchronizer, and it's recommended that it has exclusive access to an executor instance, although it's not necessary. This role can perfectly be run in a single instance, however, the JSON RPC and executor services can benefit from running in multiple instances, if the performance decreases due to the number of requests received

### Trusted sequencer

This role can only be performed by a single entity. This is enforced in the smart contract, as the related methods of the trusted sequencer can only be performed by the owner of a particular private key.

Required services and components:

- JSON RPC: can run in a separated instance, and can have multiple instances
- Sequencer & Synchronizer: single instance that needs to run together
- Executor & Merkletree: service that can run on a separate instance
- Broadcast: can run on a separate instance
- Pool DB: Postgres SQL that can be run in a separate instance
- State DB: Postgres SQL that can be run in a separate instance

Note that the JSON RPC is required to receive transactions. It's recommended that the JSON RPC runs on separated instances, and potentially more than one (depending on the load of the network). It's also recommended that the JSON RPC and the Sequencer don't share the same executor instance, to make sure that the sequencer has exclusive access to an executor

### Permissionless sequencer

TBD

### Aggregator

This role can be performed by anyone.

Required services and components:

- Synchronizer: single instance that can run on a separated instance
- Executor & Merkletree: service that can run on a separate instance
- State DB: Postgres SQL that can be run in a separate instance
- Aggregator: single instance that can run on a separated instance
- Prover: single instance that can run on a separated instance
- Executor: single instance that can run on a separated instance

It's recommended that the prover is run on a separate instance, as it has important hardware requirements. On the other hand, all the other components can run on a single instance,

## Development

It's recommended to use `make` for building, and testing the code, ... Run `make help` to get a list of the available commands.

## Running the node

- [Running localy](docs/running_local.md)
- [Running on production](docs/production-setup.md)

### Requirements

- Go 1.17
- Docker
- Docker Compose
- Make
- GCC

## Contribute

Before opening a pull request, please read this [guide](CONTRIBUTING.md)

## Disclaimer

This code has not yet been audited, and should not be used in any production systems.
