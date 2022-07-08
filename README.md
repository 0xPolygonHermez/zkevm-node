# zkEVM Node

zkEVM Node is a Go implementation of a node that operates the Polygon zkEVM Network.

## About the Polygon zkEVM network

Since this is an implementation of a protocol it's fundamental to understand it, [here TODO]() you can find the specification of the protocol.

Glossary:

- L1: Base blockchain where the rollup smart contracts are deployed. It's Ethereum or a testnet of Ethereum, but it could be any EVM compatible blockchain.
- L2: the rollup network aka the Polygon zkEVM network.
- Batch: group of transactions that are executed/proved, using the [zkEVM prover TODO]() and sent to / synchronized from L1
- Sequencer: actor that is responsible for selecting transactions, put them in a specific order and send them in batches to L1
- Trusted sequencer: sequencer that has special privileges, there can only be one trusted sequencer. The privileges granted to the trusted sequencer make it available to forecast the batches that will be applied to L1. This way it can commit to a specific sequence before interacting with L1. This is done in order to achieve fast finality and reduce costs associated to using the network (lower gas fees)
- Permissionless sequencer: sequencer role that can be performed by anyone. It has competitive disadvantadges compared to the trusted sequencer (slow finality, MEV attacks). The main purpose of it is to provide censorship resistance and unstoppability features to the network.
- Sequence: Group of batches and other metadata that the trusted sequencer sends to L1 in order to update the state
- Forced batch: batch that is sent by permissionless sequencers to L1 in order to update the state
- L2 Block: Same as a L1 block, but for L2. This is mostly used by the JSON RPC interface. Currently, all the L2 Blocks are set to only include one transaction, this is done to achieve instant finality: it's not necessary to close a batch to allow the JSON RPC to expose results of already processed transactions
- Trusted state: state reached through processing transactions that have been shared by the trusted sequencer. This state is considered trusted as the trusted sequencer could commit to a certain sequence, and then send a different one to L1
- Virtual state: state reached through processing transactions that have already been submitted to L1. This transactions are sent in batches buy either trusted or permissionless sequencers. Those batches are also called virtual batches. Note that this state is trustless as it relies on L1 security assumptions
- Consolidated state: state that is proven on-chain by submitting a ZKP (Zero Knowledge Proof) that proofs the execution of a sequence of the last virtual batch.
- Invalid transaction: transaction that can't be processed and doesn't affect the state. Note that such a transaction could be included in a virtual batch. The reason for a transaction to not be valid could be related to the Ethereum protocol (invalid nonce, not enough balance, ...) or due to limitations introduced by the zkEVM (each batch can make use of a limited ammount of resources such as the total amount of keccak hashes that can be computed)
- Reverted transaction: transaction that is executed, but is reverted (because of smart contract logic). Main difference with *invalid transaction* is that this transaction modifies the state, at least to increment nonce off the sender.
- Proof of Efficiency (PoE): name of the protocol used by the network, it's enforced by the [smart contracts TODO]()

## Architecture

<p align="center">
  <img src="./docs/architecture.drawio.png"/>
</p>

The diagram represents the main components of the software and how they interact between them. Note that this reflects a single entity running a node, in particular a node that acts as trusted sequencer. But there are many entities running nodes into the network, and each of this entities can perform different roles. More on this later.

- (JSON) RPC: interface that allow users (metamask, etherscan, ...) to interact with the node. Fully compatible with Ethereum RPC + some extra endpoints specifics of the network. It interacts with the `state` to get data and process transactions and with the `pool` to store transactions
- Pool: DB that stores transactions by the `RPC` to be selected/discarded by the `sequencer` later on
- Trusted Sequencer: get transactions from the `pool`, check if they are valid by processing them using the `state` and create sequences. Once transactions are added into the state, they are immediatley available to the `broadcast` service. Sequences are sent to L1 using the `etherman`
- Broadcast: API used by the `synchronizer` of nodes that are not the `trusted sequencer` to synchornize the trusted state
- Permissionless Sequenser: *comming soon*
- Etherman: abstraction that implements the needed methods to interact with the Ethereum network and the relevant smart contracts.
- Synchronizer: Updates the `state` by fetching data from Ethereum through the `etherman`. If the node is not a `trusted sequencer` it also updates the state with the data fetched from the `broadcast` of the `trusted sequencer`. It also detect and handles reorgs that can happen if the `trusted sequencer` sends different data in the broadcast vs the sequences sent to L1 (trusted vs virtual state)
- State: Responsible for mannaging the state data (batches, blocks, transactions, ...) that is stored on the `state SB`. It also handles the integration with the `executor` and the `Merkletree` service
- State DB: persistance layer for the state data (except the Merkletree that is handled by the `Merkletree` service)
- Aggregator: consolidates batches by generating ZKPs (Zero Knowledge proofs). To do so it gathers the necessary data that the `prover` needs as input thorugh the `state` and sends a request to it. Once the prove is generated it's sent to Ethereum through the `etherman`
- Prover/Executor: service that generates ZK proofs. Note that this component is not implemented in this repository, and it's treated as a "black box" from the perspective of the node. The prover/executor has two implementations: [JS reference implementation TODO](https://github.com/hermeznetwork/zkproverjs) and [C production ready implementation TODO](https://github.com/hermeznetwork/zkproverc). Although it's the same software/service, it has two very different purposes:
  - Provide an EVM implementation that allows to process transactions and get all needed results metadata (state root, receipts, logs, ...)
  - Generate ZKPs
- Merkletree: serivce that stores the Merkletree, containing all the account information (balances, nonces, smart contract code and smart contract storage). This component is also not implemented in this repo and is consumed as an external service by teh node. The implementation can be found [here TODO]()

## Roles of the network



### Trusted sequencer

Explained on the diagram above. It requires all the comopnents except the aggregator, which is optional. This role can only be performed by a single entitiy. This is enforced in the smart contract, as the related methods of the trusted sequencer can only be performed by the owner of a particular private key.

## Development

It's recommended to use `make` for building, testing the code, ... Run `make help` to get a list of the available commands.

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

Before opening a pull request, please read this [guide](docs/contribute-guie.md)
