# Polygon CDK Validium Node

For a full overview of the CDK-Validium please reference the [CDK documentation](https://wiki.polygon.technology/docs/cdk/).

The CDK-Validium solution is made up of several components, start with the [CDK Validium Node](https://github.com/0xPolygon/cdk-validium-node). However, for quick reference, the complete list of components are outlined below:

| Component                                                                     | Description                                                          |
| ----------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| [CDK Validium Node](https://github.com/0xPolygon/cdk-validium-node)           | Node implementation for the CDK networks in Validium mode            |
| [CDK Validium Contracts](https://github.com/0xPolygon/cdk-validium-contracts) | Smart contract implementation for the CDK networks in Validium mode |
| [CDK Data Availability](https://github.com/0xPolygon/cdk-data-availability)   | Data availability implementation for the CDK networks          |
| [Prover / Executor](https://github.com/0xPolygonHermez/zkevm-prover)          | zkEVM engine and prover implementation                               |
| [Bridge Service](https://github.com/0xPolygonHermez/zkevm-bridge-service)     | Bridge service implementation for CDK networks                       |
| [Bridge UI](https://github.com/0xPolygonHermez/zkevm-bridge-ui)               | UI for the CDK networks bridge                                       |

Understanding the underlying protocol is crucial when working with an implementation. This project is based on the Polygon zkEVM network, which is designed to bring scalability to Ethereum-compatible blockchains.

For an in-depth understanding of the protocol’s specifications, please refer to the [zkEVM Protocol Overview](https://wiki.polygon.technology/docs/zkevm/)

## Run a CDK Validium

> This repo is a fork of the [zkevm-node](https://github.com/0xPolygonHermez/zkevm-node), more information and code diff explained [here](./docs/diff/diff.md)

### Development

> ARM devices (such as Apple M1 and M2) are not supported

For a streamlined development experience, it’s highly recommended to utilize the make utility for tasks such as building and testing the code. To view a comprehensive list of available commands, simply execute `make help` in your terminal.

This step by step guide will result in a local environment that has everything needed to test and develop on a CDK Validium, but note that:

- everything will be run on a ephemeral and local L1 network, once the environment is shutdown, all progress will be lost
- ZK Proofs are mocked
- Bridge service and UI is not included as part of this setup, instead there is a pre-funded account

#### Steps

1. Clone this GitHub repository to your local machine:

```
git clone https://github.com/0xPolygon/cdk-validium-node.git
```

2. Navigate to the cloned directory:

```
cd cdk-validium-node
```

3. Build the Docker image using the provided Dockerfile:

```
make build-docker
```

4. Navigate to the test directory:

```
cd test
```

5. Run all needed components:

```
make run
```

#### Usage

- L2 RPC endpoint: `http://localhost:8123`
- L2 Chain ID: 1001
- L1 RPC endpoint: `http:localhost:8545`
- L1 Chain ID: 1337
- Pre funded account private key: `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`

#### Troubleshooting

Everything is run using docker, so if anything is not working, first thing is to identify what containers are running or not:

```
docker compose ps
```

Then check the logs:

```
docker logs <problematic container, example: cdk-validium-sync>
```

Aditionaly, it can be worth checking the DBs:

- StateDB: `psql -h localhost -p 5432 -U state_user state_db`, password: `state_password`
- PoolDB: `psql -h localhost -p 5433 -U pool_user pool_db`, password: `pool_password`

#### Advanced config

In order to go beyond the default configuration, you can edit the config files:

- `./test/config/test.node.config.toml`: configuration of the node, documented [here](./docs/config-file/node-config-doc.md)
- `./test/config/test.genesis.config.toml`: configuration of the network, documented [here](./docs/config-file/custom_network-config-doc.md)
- `./test/config/test.prover.config.json`: configuration of the prover/executor

## Key Components

### Aggregator

The Aggregator is responsible for submitting validity proofs of the L2 state to L1. To do so, it fetches the batches sequencced by the sequencer, and interacts with the provers to generate the ZeroKnowledge Proofs (ZKPs).

To do so in a efficient way, the Aggregator will:

- Orchestrate communication with one or multiple provers
- Aggregate proofs from many batches, a single proof can verify multiple batches
- Send the aggregated proof to L1 using the EthTxManager

### Prover

The Prover is tasked with generating proofs for the batched transactions. These proofs are essential for the subsequent validation of the transactions on the Ethereum mainnet. In general:

- ZKP Generation: Creates cryptographic proofs for each batch of transactions or for a combination of batches (proof aggregation).
- Optimization: Utilizes parallel computing to speed up the proof generation process.
- Ethereum Mainnet Preparation: Formats the proofs for validation on the Ethereum mainnet.

Note that this software is not implemented in this repo, but in [this one](https://github.com/0xPolygonHermez/zkevm-prover)

### Sequencer

The Sequencer is responsible for ordering transactions, in other words, making the state move forward:

- Transaction Ordering: Get transactions from the pool and adds them into the state.
- State Transition: Collaborates with the Executor to process transactions and update the state.
- Trusted finality: Once the sequencer has added a transaction into the state, it shares this information with other nodes, making the transaction final. Other nodes will need to trust that this transaction is added into the state until they get data availability (DA) and validity (ZKPs) confirmations

### SequenceSender

The SequenceSender’s role is to send the ordered list of transactions, known as a sequence, to the Ethereum mainnet. It also collaborates with the Data Availability layer, ensuring that all transaction data is accessible off-chain. It plays a pivotal role in finalizing the rollup:

- Sequence Transmission: Sends a fingerprint of the ordered transaction batches to the Ethereum mainnet.
- Data Availability: Works in tandem with the Data Availability layer to ensure off-chain data is accessible.
- L1 Interaction: Utilizes the EthTxManager to handle L1 transaction nuances like nonce management and gas price adjustments.

### Synchronizer

The Synchronizer keeps the node’s local state in sync with the Ethereum mainnet. It listens for events emitted by the smart contract on the mainnet and updates the local state to match. The Synchronizer acts as the bridge between the Ethereum mainnet and the node:

- Event Listening: Monitors events emitted by the smart contract on the Ethereum mainnet.
- Data Availability: downloads data from the Data Availability layer based on L1 events
- State Updating: Aligns the local state with the mainnet, ensuring consistency.
- Reorg Handling: Detects and manages blockchain reorganizations to maintain data integrity.

### Data Availability Configuration

The Data Availability (DA) layer is a crucial component that ensures all transaction data is available when needed. This off-chain storage solution is configurable, allowing operators to set parameters that best suit their needs. The DA layer is essential for the Validium system, where data availability is maintained off-chain but can be made available for verification when required. In general:

- Off-Chain Storage: Stores all transaction data off-chain but ensures it’s readily available for verification.
- Configurability: Allows chain operators to customize data storage parameters.
- Data Verification: Provides mechanisms for data integrity checks, crucial for the Validium model.

### Executor

The Executor is the state transition implementation, in this case a EVM implementation:

- Batch execution: receives requests to execute batch of transactions.
- EVM Implementation: Provides an EVM-compatible implementation for transaction processing.
- Metadata Retrieval: Retrieves necessary metadata like state root, transaction receipts, and logs from the execution.

Note that this software is not implemented in this repo, but in [this one](https://github.com/0xPolygonHermez/zkevm-prover)

### EthTxManager

The EthTxManager is crucial for interacting with the Ethereum mainnet:

- L1 Transaction Handling: Manages requests from the SequenceSender and Aggregator to send transactions to L1.
- Nonce Management: Takes care of the nonce for each account involved in a transaction.
- Gas Price Adjustment: Dynamically adjusts the gas price to ensure that transactions are mined in a timely manner.

### State

The State component is the backbone of the node’s data management:

- State Management: Handles all state-related data, including batches, blocks, and transactions.
- Executor Integration: Communicates with the Executor to process transactions and update the state.
- StateDB: used for persistance

### Pool

The Pool serves as a temporary storage for transactions:

- Transaction Storage: Holds transactions submitted via the RPC.
- Sequencer Interaction: Provides transactions to the Sequencer for ordering and batch creation.

### JSON RPC

The JSON RPC serves as the HTTP interface for user interaction:

- User Interface: Allows users and dApps to interact with the node, following the Ethereum standard:
    - [Endpoint compatibility](./docs/json-rpc-endpoints.md)
    - [Custom endpoints](./docs/zkEVM-custom-endpoints.md)
- State Interaction: Retrieves data from the state and processes transactions.
- Pool Interaction: Stores transactions in the pool.

### L2GasPricer

The L2GasPricer is responsible for calculating the gas price on L2 based on the L1 gas price:

- L1 Gas Price Fetching: Retrieves the current L1 gas price.
- Gas Price Calculation: Applies a formula to calculate the suggested L2 gas price.
- Pool Storage: Stores the calculated L2 gas price in the pool for consumption by the rpc.

## Contribute

Before opening a pull request, please read [this guide](./CONTRIBUTING.md).

## License

The cdk-validium-node project is licensed under the GNU Affero General Public License free software license.
