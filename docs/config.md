# Config

This document specifies the use and default values for all the configuration parameters.
Each config parameter can be specified through configuration file, such as [this one](../config/environments/public/public.node.config.toml), or environment variable and has a default value.
Note that not all the configuration parameters are used, this depends on the component(s) that are being run.

## Log

Used by: all the components

| config      | env                        | default       | description                                                                                                                      |
| ----------- | -------------------------- | ------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| Environment | ZKEVM_NODE_LOG_ENVIRONMENT | "development" | Can be "development" or "production". This will change the format of the logs, for productions logs are displayed in JOSN format |
| Level       | ZKEVM_NODE_LOG_LEVEL       | "debug"       | Can be "debug", "info", "warn" or "error". This will filter the logs according to it's level                                     |
| Outputs     | ZKEVM_NODE_LOG_OUTPUTS     | ["stderr"]    |                                                                                                                                  |

## StateDB

Used by: all the components

| config    | env                            | default          | description                                                 |
| --------- | ------------------------------ | ---------------- | ----------------------------------------------------------- |
| User      | ZKEVM_NODE_STATE_DB_USER       | "state_user"     | User name of the DB                                         |
| Password  | ZKEVM_NODE_STATE_DB_PASSWORD   | "state_password" | Password for the DB user                                    |
| Name      | ZKEVM_NODE_STATE_DB_NAME       | "state_db"       | Name of the DB                                              |
| Host      | ZKEVM_NODE_STATE_DB_HOST       | "localhost"      | Host of the DB                                              |
| Port      | ZKEVM_NODE_STATE_DB_PORT       | "5432"           | Port of the DB                                              |
| EnableLog | ZKEVM_NODE_STATE_DB_ENABLE_LOG | false            | If set to true the node will log the queries done to the DB |
| MaxConns  | ZKEVM_NODE_STATE_DB_MAX_CONNS  | 200              | Maximum connections that the node will use towards the DB   |

## Pool

Used by: `sequencer`, `rpc`, `synchronizer`, `l2gaspricer`

| config                         | env                                                 | default    | description                                                                                    |
| ------------------------------ | --------------------------------------------------- | ---------- | ---------------------------------------------------------------------------------------------- |
| FreeClaimGasLimit              | ZKEVM_NODE_POOL_FREE_CLAIM_GAS_LIMIT                | 150000     | If txs exceed this value in its gas limit, they won't be allowed to use 0 gas price            |
| MaxTxBytesSize                 | ZKEVM_NODE_POOL_MAX_TX_BYTES_SIZE                   | 30132      | Maximum byte size allowed per tx                                                               |
| MaxTxDataBytesSize             | ZKEVM_NODE_POOL_MAX_TX_DATA_BYTES_SIZE              | 30000      | Maximum byte size data allowed per tx                                                          |
| DefaultMinGasPriceAllowed      | ZKEVM_NODE_POOL_DEFAULT_MIN_GAS_PRICE_ALLOWED       | 1000000000 | If there are no updates on the suggested gas price, txs bellow this gas price will be rejected |
| MinAllowedGasPriceInterval     | ZKEVM_NODE_POOL_MIN_ALLOWED_GAS_PRICE_INTERVAL      | "5m"       | Txs bellow the minimum gas price suggested in the previous indicated interval will be rejected |
| PollMinAllowedGasPriceInterval | ZKEVM_NODE_POOL_POLL_MIN_ALLOWED_GAS_PRICE_INTERVAL | "15s"      | Frequency to update the minimum gas price allowed                                              |

### Pool.DB

| config    | env                            | default          | description                                                 |
| --------- | ------------------------------ | ---------------- | ----------------------------------------------------------- |
| User      | ZKEVM_NODE_POOL_DB_USER        | "state_user"     | User name of the DB                                         |
| Password  | ZKEVM_NODE_POOL_DB_PASSWORD    | "state_password" | Password for the DB user                                    |
| Name      | ZKEVM_NODE_POOL_DB_NAME        | "state_db"       | Name of the DB                                              |
| Host      | ZKEVM_NODE_POOL_DB_HOST        | "localhost"      | Host of the DB                                              |
| Port      | ZKEVM_NODE_POOL_DB_PORT        | "5432"           | Port of the DB                                              |
| EnableLog | ZKEVM_NODE_POOL_DB_ENABLE_LOG  | false            | If set to true the node will log the queries done to the DB |
| MaxConns  | ZKEVM_NODE_POOL_DB_MAX_CONNS   | 200              | Maximum connections that the node will use towards the DB   |

## Etherman

Used by: `aggregator`, `sequencer`, `synchronizer`, `eth-tx-manager`

| config                    | env                                               | default                                      | description                                           |
| ------------------------- | ------------------------------------------------- | -------------------------------------------- | ----------------------------------------------------- |
| URL                       | ZKEVM_NODE_ETHERMAN_URL                           | "http://localhost:8545"                    | URL of the L1 JSON RPC node                           |
| L1ChainID                 | ZKEVM_NODE_ETHERMAN_L1_CHAIN_ID                   | 1337                                         | Chain ID of the L1 network                            |
| PoEAddr                   | ZKEVM_NODE_ETHERMAN_POE_ADDR                      | "0x610178dA211FEF7D417bC0e6FeD39F05609AD788" | Address of the zkEVM smart contract                   |
| MaticAddr                 | ZKEVM_NODE_ETHERMAN_MATIC_ADDR                    | "0x5FbDB2315678afecb367f032d93F642f64180aa3" | Address of the MATIC token                            |
| GlobalExitRootManagerAddr | ZKEVM_NODE_ETHERMAN_GLOBAL_EXIT_ROOT_MANAGER_ADDR | "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6" | Address of the Global Exit Root Manager contract      |
| MultiGasProvider          | ZKEVM_NODE_ETHERMAN_MULTI_GAS_PROVIDER            | true                                         | Indicate if more than one gas provider should be used |

### Etherman.Etherscan

| config | env                           | default | description                                                                                                               |
| ------ | ----------------------------- | ------- | ------------------------------------------------------------------------------------------------------------------------- |
| ApiKey | ZKEVM_NODE_ETHERMAN_ETHERSCAN_API_KEY | ""      | If MultiGasProvider is set to true this API key will be used to cnnect to Etherscan in order to get gas price information |

## EthTxManager

Used by: `eth-tx-manager`

| config                | env                                                | default | description                                                                                                                        |
| --------------------- | -------------------------------------------------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| FrequencyToMonitorTxs | ZKEVM_NODE_ETH_TX_MANAGER_FREQUENCY_TO_MONITOR_TXS | "1s"    | Frequency in which transactions are checked for status changes                                                                     |
| WaitTxToBeMined       | ZKEVM_NODE_ETH_TX_MANAGER_WAIT_TX_TO_BEMINED       | "2m"    | Period of time used to wait for a tx to be mined. If tx is not mined in this period of time the manager will try to send a new one |
| ForcedGas             | ZKEVM_NODE_ETH_TX_MANAGER_FORCED_GAS               | 0       | If this value is greater than 0, txs being sent to L1 will use this value for gas limit **when gas estimation fails**              |

## RPC

Used by: `rpc`



| config  | env  | default  | description  |
|---|---|---|---|
| Host   | ZKEVM_NODE_RPC_HOST   | "0.0.0.0"   | Host of the server  |
| Port   | ZKEVM_NODE_RPC_PORT   | 8123   | Port of the server  |
| ReadTimeout   | ZKEVM_NODE_RPC_READ_TIMEOUT   | 60   | Maximum time allowed to read requests  |
| WriteTimeout   | ZKEVM_NODE_RPC_WRITE_TIMEOUT   | 60   | Maximum time allowed to write requests  |
| MaxRequestsPerIPAndSecond   | ZKEVM_NODE_RPC_MAX_REQUESTS_PER_IP_PER_SECOND   | 50   | Maximum amount of requests per second allowed per IP address  |
| SequencerNodeURI   | ZKEVM_NODE_RPC_SEQUENCER_NODE_URI   | ""   | RPC URL of the trusted sequencer. Used to proxy some pool related endpoints (such as sending txs)  |
| BroadcastURI   | ZKEVM_NODE_RPC_BROADCAST_URI   | "127.0.0.1:61090"   | URL of the broadcast service (deprecated)  |
| DefaultSenderAddress   | ZKEVM_NODE_RPC_DEFAULT_SENDER_ADDRESS   | "0x1111111111111111111111111111111111111111"   | Used to set the "from" of a tx when not provided for unsigned tx methods  |
| EnableL2SuggestedGasPricePolling   | ZKEVM_NODE_RPC_ENABLE_L2_SUGGESTED_GAS_PRICE_POLLING   | true   | When true, gas price suggestions will be updated  |

## Synchronizer

Used by: `synchronizer`

SyncInterval = "0s"
SyncChunkSize = 100
GenBlockNumber = 67

| config         | env                                      | default | description                                                |
| -------------- | ---------------------------------------- | ------- | ---------------------------------------------------------- |
| SyncInterval   | ZKEVM_NODE_SYNCHRONIZER_SYNC_INTERVAL    | "0s"    | Amount of time waited between sync loops                   |
| SyncChunkSize  | ZKEVM_NODE_SYNCHRONIZER_CHUNK_SIZE       | 100     | Amount of L1 blocks fetched per sync loop                  |
| GenBlockNumber | ZKEVM_NODE_SYNCHRONIZER_GEN_BLOKC_NUMBER | 67      | L1 block in which the rollup smart contracts were deployed |

## Sequencer

Used by: `sequencer`
| config                                   | env                                                                 | default  | description                                                                                                                                  |
| ---------------------------------------- | ------------------------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| NotSyncedWait                            | ZKEVM_NODE_SEQUENCER_NOT_SYNCED_WAIT                                | "1s"     | Time too wait when the sequencer is waiting for th synchronizer to catch up with L1                                                          |
| WaitPeriodSendSequence                   | ZKEVM_NODE_SEQUENCER_WAIT_PERIOD_SEND_SEQUENCE                      | "5s"     | Frequency in which the send sequences to L1 loop is run (doesn't mean that sequences will be sent with this frequency)                       |
| LastBatchVirtualizationTimeMaxWaitPeriod | ZKEVM_NODE_SEQUENCER_LAST_BATCH_VIRTUALIZATION_TIME_MAX_WAIT_PERIOD | "1h"     | If there has not been batches sent to L1 for longer than the specified duration, a batch will be sent even if the L1 tx could be more packed |
| BlocksAmountForTxsToBeDeleted            | ZKEVM_NODE_SEQUENCER_BLOCKS_AMOUNT_FORTXS_TO_BE_DELETED             | 100      | Will delete txs from the pool that have been mined on L1 (virtual state) for more than the indicated amount of L1 blocks                     |
| FrequencyToCheckTxsForDelete             | ZKEVM_NODE_SEQUENCER_FREQUENCY_TO_CHECK_TXS_FOR_DELETE              | "12h"    | Frequency in which the deletion process of already mined txs will happen                                                                     |
| MaxTxsPerBatch                           | ZKEVM_NODE_SEQUENCER_MAX_TXS_PER_BATCH                              | 150      | Maximum amount of txs that the sequencer will include in a batch                                                                             |
| MaxBatchBytesSize                        | ZKEVM_NODE_SEQUENCER_MAX_BATCH_BYTES_SIZE                           | 129848   | Sequencer will close a batch before it reaches the indicated size                                                                            |
| MaxCumulativeGasUsed                     | ZKEVM_NODE_SEQUENCER_MAX_CUMULATIVE_GAS_USED                        | 30000000 | Sequencer will close a batch before it consumes the indicated amount of gas                                                                  |
| MaxKeccakHashes                          | ZKEVM_NODE_SEQUENCER_MAX_KECCAK_HASHES                              | 468      | Sequencer will close a batch before it consumes the indicated amount of Keccack hashes                                                       |
| MaxPoseidonHashes                        | ZKEVM_NODE_SEQUENCER_MAX_POSEIDON_HASHES                            | 279620   | Sequencer will close a batch before it consumes the indicated amount of Poseidon hashes                                                      |
| MaxPoseidonPaddings                      | ZKEVM_NODE_SEQUENCER_MAX_POSEIDON_PADDINGS                          | 149796   | Sequencer will close a batch before it consumes the indicated amount of Poseidon paddings                                                    |
| MaxMemAligns                             | ZKEVM_NODE_SEQUENCER_MAX_MEM_ALIGNS                                 | 262144   | Sequencer will close a batch before it consumes the indicated amount of memory alignments                                                    |
| MaxArithmetics                           | ZKEVM_NODE_SEQUENCER_MAX_ARITHEMTICS                                | 262144   | Sequencer will close a batch before it consumes the indicated amount of arithmetic operations                                                |
| MaxBinaries                              | ZKEVM_NODE_SEQUENCER_MAX_BINARIES                                   | 262144   | Sequencer will close a batch before it consumes the indicated amount of binary operations                                                    |
| MaxSteps                                 | ZKEVM_NODE_SEQUENCER_MAX_STEPS                                      | 8388608  | Sequencer will close a batch before it consumes the indicated amount of steps                                                                |
| WeightBatchBytesSize                     | ZKEVM_NODE_SEQUENCER_WEIGHT_BATCH_BYTE_SIZE                         | 1        | Factor multiplied on the size by a tx to set the efficiency score                                                                            |
| WeightCumulativeGasUsed                  | ZKEVM_NODE_SEQUENCER_WEIGHT_CUMULATIVE_GAS_USED                     | 1        | Factor multiplied on the gas used by a tx to set the efficiency score                                                                        |
| WeightKeccakHashes                       | ZKEVM_NODE_SEQUENCER_WEIGHT_KECCACK_HASHES                          | 1        | Factor multiplied on the keccack hashes used by a tx to set the efficiency score                                                             |
| WeightPoseidonHashes                     | ZKEVM_NODE_SEQUENCER_WEIGHT_POSEIDON_HASHES                         | 1        | Factor multiplied on the Poseidon hashes used by a tx to set the efficiency score                                                            |
| WeightPoseidonPaddings                   | ZKEVM_NODE_SEQUENCER_WEIGHT_POSEIDON_PADDINGS                       | 1        | Factor multiplied on the Poseidon paddings used by a tx to set the efficiency score                                                          |
| WeightMemAligns                          | ZKEVM_NODE_SEQUENCER_WEIGHT_MEM_ALIGNS                              | 1        | Factor multiplied on the mem aligns used by a tx to set the efficiency score                                                                 |
| WeightArithmetics                        | ZKEVM_NODE_SEQUENCER_WEIGHT_ARITHMETICS                             | 1        | Factor multiplied on the arithemtics used by a tx to set the efficiency score                                                                |
| WeightBinaries                           | ZKEVM_NODE_SEQUENCER_WEIGHT_BINNARIES                               | 1        | Factor multiplied on the binnaries used by a tx to set the efficiency score                                                                  |
| WeightSteps                              | ZKEVM_NODE_SEQUENCER_WEIGHT_STEPS                                   | 1        | Factor multiplied on the steps used by a tx to set the efficiency score                                                                      |
| TxLifetimeCheckTimeout                   | ZKEVM_NODE_SEQUENCER_TX_LIFETIME_CHECK_TIMEOUT                      | "10m"    | Frequency in which the expire logic for pending txs will be triggered                                                                        |
| MaxTxLifetime                            | ZKEVM_NODE_SEQUENCER_MAX_TX_LIFETIME                                | "3h"     | Maximum time for txs to be in pending state                                                                                                  |
| MaxTxSizeForL1                           | ZKEVM_NODE_SEQUENCER_MAX_TX_SIZE_FOR_L1                             | 131072   | Max size of an L1 tx                                                                                                                         |

### Sequencer.Finalizer

#### GERDeadlineTimeout

- Description: Grace period for the finalizer to update the GER in a new batch without having to close the current open one in a forced way
- Default: 5s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_GER_DEADLINE_TIMEOUT

#### ForcedBatchDeadlineTimeout

- Description: Grace period for the finalizer to update include forced batches in a new batch without having to close the current open one in a forced way
- Default: 60s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_FORCED_BATCH_DEADLINE_TIMEOUT

#### SendingToL1DeadlineTimeout

- Description: Grace period for the finalizer to close the current open batch
- Default: 20s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_SENDING_TO_L1_DEADLINE_TIMEOUT

#### SleepDuration

- Description: Sleep time for the finalizer if there are no new processable txs
- Default: 100ms
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_SLEEP_DURATION

#### ResourcePercentageToCloseBatch

- Description: If a given resource of a batch exceeds this percentadge and there are no fitting txs in the worker, close batch
- Default: 10
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_RESOURCE_PERCENTAGE_TO_CLOSE_BATCH

#### GERFinalityNumberOfBlocks

- Description: Amount of (L1) blocks to consider that a tx that updates the GER is final
- Default: 64
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_GER_FINALITY_NUMBER_OF_BLOCKS

#### ClosingSignalsManagerWaitForCheckingL1Timeout

- Description: 
- Default: 10s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_CLOSING_SIGNALS_MANAGER_WAIT_FOR_CHECKING_L1_TIMEOUT

#### ClosingSignalsManagerWaitForCheckingGER

- Description: 
- Default: 10s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_CLOSING_SIGNALS_MANAGER_WAIT_FOR_CHECKING_GER

#### ClosingSignalsManagerWaitForCheckingForcedBatches

- Description: 
- Default: 10s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_CLOSING_SIGNALS_MANAGER_WAIT_FOR_CHECKING_FORCED_BATCHES

#### ForcedBatchesFinalityNumberOfBlocks

- Description: Amount of (L1) blocks to consider that a tx that forces a batch is final
- Default: 64
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_FORCED_BATCHES_FINALITY_NUMBER_OF_BLOCKS

#### SenderAddress

- Description: Address of the trusted sequencer
- Default: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_SENDER_ADDRESS

## Aggregator

Used by: `aggregator`

| config  | env  | default  | description  |
|---|---|---|---|
|   |   |   |   |
|   |   |   |   |
|   |   |   |   |

## L2GasPriceSuggester

Used by: `l2gaspricer`

| config  | env  | default  | description  |
|---|---|---|---|
|   |   |   |   |
|   |   |   |   |
|   |   |   |   |

## MTClient

Used by: all components

| config  | env  | default  | description  |
|---|---|---|---|
|   |   |   |   |
|   |   |   |   |
|   |   |   |   |

## Executor

Used by: all components

| config  | env  | default  | description  |
|---|---|---|---|
|   |   |   |   |
|   |   |   |   |
|   |   |   |   |

## BroadcastServer

Used by: `broadcast-trusted-state`

| config  | env  | default  | description  |
|---|---|---|---|
|   |   |   |   |
|   |   |   |   |
|   |   |   |   |

## Metrics

Used by: all components

| config  | env  | default  | description  |
|---|---|---|---|
|   |   |   |   |
|   |   |   |   |
|   |   |   |   |