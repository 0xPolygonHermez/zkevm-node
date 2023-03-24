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
| ReadTimeoutInSec   | ZKEVM_NODE_RPC_READ_TIMEOUT_IN_SEC   | 60   | Maximum time allowed to read requests  |
| WriteTimeoutInSec   | ZKEVM_NODE_RPC_WRITE_TIMEOUT_IN_SEC   | 60   | Maximum time allowed to write requests  |
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

NotSyncedWait = "1s"
WaitPeriodSendSequence = "5s"
LastBatchVirtualizationTimeMaxWaitPeriod = "5s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
MaxTxsPerBatch = 150
MaxBatchBytesSize = 129848
MaxCumulativeGasUsed = 30000000
MaxKeccakHashes = 468
MaxPoseidonHashes = 279620
MaxPoseidonPaddings = 149796
MaxMemAligns = 262144
MaxArithmetics = 262144
MaxBinaries = 262144
MaxSteps = 8388608
WeightBatchBytesSize = 1
WeightCumulativeGasUsed = 1
WeightKeccakHashes = 1
WeightPoseidonHashes = 1
WeightPoseidonPaddings = 1
WeightMemAligns = 1
WeightArithmetics = 1
WeightBinaries = 1
WeightSteps = 1
TxLifetimeCheckTimeout = "10m"
MaxTxLifetime = "3h"
MaxTxSizeForL1 = 131072

| config  | env  | default  | description  |
|---|---|---|---|
| NotSyncedWait   | ZKEVM_NODE_SEQUENCER_NOT_SYNCED_WAIT   | "1s"   | Time too wait when the sequencer is waiting for th synchronizer to catch up with L1   |
| WaitPeriodSendSequence   | ZKEVM_NODE_SEQUENCER_   | "5s"   |   |
| LastBatchVirtualizationTimeMaxWaitPeriod   | ZKEVM_NODE_SEQUENCER_   | "5s"   |   |
| BlocksAmountForTxsToBeDeleted   | ZKEVM_NODE_SEQUENCER_   | 100   |   |
| FrequencyToCheckTxsForDelete   | ZKEVM_NODE_SEQUENCER_   | "12h"   |   |
| MaxTxsPerBatch   | ZKEVM_NODE_SEQUENCER_   | 150   |   |
| MaxBatchBytesSize   | ZKEVM_NODE_SEQUENCER_   | 129848   |   |
| MaxCumulativeGasUsed   | ZKEVM_NODE_SEQUENCER_   | 30000000   |   |
| MaxKeccakHashes   | ZKEVM_NODE_SEQUENCER_   | 468   |   |
| MaxPoseidonHashes   | ZKEVM_NODE_SEQUENCER_   | 279620   |   |
| MaxPoseidonPaddings   | ZKEVM_NODE_SEQUENCER_   | 149796   |   |
| MaxMemAligns   | ZKEVM_NODE_SEQUENCER_   | 262144   |   |
| MaxArithmetics   | ZKEVM_NODE_SEQUENCER_   | 262144   |   |
| MaxBinaries   | ZKEVM_NODE_SEQUENCER_   | 262144   |   |
| MaxSteps   | ZKEVM_NODE_SEQUENCER_   | 8388608   |   |
| WeightBatchBytesSize   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightCumulativeGasUsed   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightKeccakHashes   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightPoseidonHashes   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightPoseidonPaddings   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightMemAligns   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightArithmetics   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightBinaries   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| WeightSteps   | ZKEVM_NODE_SEQUENCER_   | 1   |   |
| TxLifetimeCheckTimeout   | ZKEVM_NODE_SEQUENCER_   | "10m"   |   |
| MaxTxLifetime   | ZKEVM_NODE_SEQUENCER_   | "3h"   |   |
| MaxTxSizeForL1   | ZKEVM_NODE_SEQUENCER_   | 131072   |   |

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