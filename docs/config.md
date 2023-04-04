# Config

This document specifies the use and default values for all the configuration parameters.
Each config parameter can be specified through configuration file, such as [this one](../config/environments/public/public.node.config.toml), or environment variable and has a default value.
Note that some config parameters will be used and some others won't, depending on the components that are being run.

## Log

> Used by: all the components

#### Environment

- Description: Can be "development" or "production". This will change the format of the logs, for productions logs are displayed in JOSN format
- Default: "development" 
- env: ZKEVM_NODE_LOG_ENVIRONMENT

#### Level

- Description: Can be "debug", "info", "warn" or "error". This will filter the logs according to it's level
- Default: "debug"       
- env: ZKEVM_NODE_LOG_LEVEL

#### Outputs

- Description: ~
- Default: ["stderr"]    
- env: ZKEVM_NODE_LOG_OUTPUTS

## StateDB

> Used by: all the components

#### User

- Description: User name of the DB
- Default: "state_user"     
- env: ZKEVM_NODE_STATE_DB_USER

#### Password

- Description: Password for the DB user
- Default: "state_password" 
- env: ZKEVM_NODE_STATE_DB_PASSWORD

#### Name

- Description: Name of the DB
- Default: "state_db"       
- env: ZKEVM_NODE_STATE_DB_NAME

#### Host

- Description: Host of the DB
- Default: "localhost"      
- env: ZKEVM_NODE_STATE_DB_HOST

#### Port

- Description: Port of the DB
- Default: "5432"           
- env: ZKEVM_NODE_STATE_DB_PORT

#### EnableLog

- Description: If set to true the node will log the queries done to the DB
- Default: false            
- env: ZKEVM_NODE_STATE_DB_ENABLE_LOG

#### MaxConns

- Description: Maximum connections that the node will use towards the DB
- Default: 200              
- env: ZKEVM_NODE_STATE_DB_MAX_CONNS


## Pool

> Used by: `sequencer`, `rpc`, `synchronizer`, `l2gaspricer`

#### FreeClaimGasLimit

- Description: If txs exceed this value in its gas limit, they won't be allowed to use 0 gas price            
- Default: 150000     
- env: ZKEVM_NODE_POOL_FREE_CLAIM_GAS_LIMIT

#### MaxTxBytesSize

- Description: Maximum byte size allowed per tx                                                               
- Default: 30132      
- env: ZKEVM_NODE_POOL_MAX_TX_BYTES_SIZE

#### MaxTxDataBytesSize

- Description: Maximum byte size data allowed per tx                                                          
- Default: 30000      
- env: ZKEVM_NODE_POOL_MAX_TX_DATA_BYTES_SIZE

#### DefaultMinGasPriceAllowed

- Description: If there are no updates on the suggested gas price, txs bellow this gas price will be rejected
- Default: 1000000000 
- env: ZKEVM_NODE_POOL_DEFAULT_MIN_GAS_PRICE_ALLOWED

#### MinAllowedGasPriceInterval

- Description: Txs bellow the minimum gas price suggested in the previous indicated interval will be rejected
- Default: "5m"       
- env: ZKEVM_NODE_POOL_MIN_ALLOWED_GAS_PRICE_INTERVAL

#### PollMinAllowedGasPriceInterval

- Description: Frequency to update the minimum gas price allowed
- Default: "15s"      
- env: ZKEVM_NODE_POOL_POLL_MIN_ALLOWED_GAS_PRICE_INTERVAL


### Pool.DB

#### User

- Description: User name of the DB
- Default: "state_user"     
- env: ZKEVM_NODE_POOL_DB_USER

#### Password

- Description: Password for the DB user
- Default: "state_password" 
- env: ZKEVM_NODE_POOL_DB_PASSWORD

#### Name

- Description: Name of the DB
- Default: "state_db"       
- env: ZKEVM_NODE_POOL_DB_NAME

#### Host

- Description: Host of the DB
- Default: "localhost"      
- env: ZKEVM_NODE_POOL_DB_HOST

#### Port

- Description: Port of the DB
- Default: "5432"           
- env: ZKEVM_NODE_POOL_DB_PORT

#### EnableLog

- Description: If set to true the node will log the queries done to the DB
- Default: false            
- env: ZKEVM_NODE_POOL_DB_ENABLE_LOG

#### MaxConns

- Description: Maximum connections that the node will use towards the DB
- Default: 200              
- env: ZKEVM_NODE_POOL_DB_MAX_CONNS


## Etherman

> Used by: `aggregator`, `sequencer`, `synchronizer`, `eth-tx-manager`

#### URL

- Description: URL of the L1 JSON RPC node
- Default: "http://localhost:8545"                    
- env: ZKEVM_NODE_ETHERMAN_URL

#### L1ChainID

- Description: Chain ID of the L1 network
- Default: 1337                                         
- env: ZKEVM_NODE_ETHERMAN_L1_CHAIN_ID

#### PoEAddr

- Description: Address of the zkEVM smart contract
- Default: "0x610178dA211FEF7D417bC0e6FeD39F05609AD788" 
- env: ZKEVM_NODE_ETHERMAN_POE_ADDR

#### MaticAddr

- Description: Address of the MATIC token
- Default: "0x5FbDB2315678afecb367f032d93F642f64180aa3" 
- env: ZKEVM_NODE_ETHERMAN_MATIC_ADDR

#### GlobalExitRootManagerAddr

- Description: Address of the Global Exit Root Manager contract
- Default: "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6" 
- env: ZKEVM_NODE_ETHERMAN_GLOBAL_EXIT_ROOT_MANAGER_ADDR

#### MultiGasProvider

- Description: Indicate if more than one gas provider should be used
- Default: true                                         
- env: ZKEVM_NODE_ETHERMAN_MULTI_GAS_PROVIDER


### Etherman.Etherscan

#### ApiKey

- Description: If MultiGasProvider is set to true this API key will be used to cnnect to Etherscan in order to get gas price information
- Default: ""
- env: ZKEVM_NODE_ETHERMAN_ETHERSCAN_API_KEY

## EthTxManager

> Used by: `eth-tx-manager`

#### FrequencyToMonitorTxs

- Description: Frequency in which transactions are checked for status changes                                                                     
- Default: "1s"    
- env: ZKEVM_NODE_ETH_TX_MANAGER_FREQUENCY_TO_MONITOR_TXS

#### WaitTxToBeMined

- Description: Period of time used to wait for a tx to be mined. If tx is not mined in this period of time the manager will try to send a new one 
- Default: "2m"    
- env: ZKEVM_NODE_ETH_TX_MANAGER_WAIT_TX_TO_BEMINED

#### ForcedGas

- Description: If this value is greater than 0, txs being sent to L1 will use this value for gas limit **when gas estimation fails**              
- Default: 0       
- env: ZKEVM_NODE_ETH_TX_MANAGER_FORCED_GAS


## RPC

> Used by: `rpc`



#### Host

- Description: Host of the server  
- Default: "0.0.0.0"   
- env: ZKEVM_NODE_RPC_HOST

#### Port

- Description: Port of the server  
- Default: 8123   
- env: ZKEVM_NODE_RPC_PORT

#### ReadTimeout

- Description: Maximum time allowed to read requests  
- Default: 60   
- env: ZKEVM_NODE_RPC_READ_TIMEOUT

#### WriteTimeout

- Description: Maximum time allowed to write requests  
- Default: 60   
- env: ZKEVM_NODE_RPC_WRITE_TIMEOUT

#### MaxRequestsPerIPAndSecond

- Description: Maximum amount of requests per second allowed per IP address  
- Default: 50   
- env: ZKEVM_NODE_RPC_MAX_REQUESTS_PER_IP_PER_SECOND

#### SequencerNodeURI

- Description: RPC URL of the trusted sequencer. Used to proxy some pool related endpoints (such as sending txs)  
- Default: ""   
- env: ZKEVM_NODE_RPC_SEQUENCER_NODE_URI

#### BroadcastURI

- Description: URL of the broadcast service (deprecated)  
- Default: "127.0.0.1:61090"   
- env: ZKEVM_NODE_RPC_BROADCAST_URI

#### DefaultSenderAddress

- Description: Used to set the "from" of a tx when not provided for unsigned tx methods  
- Default: "0x1111111111111111111111111111111111111111"   
- env: ZKEVM_NODE_RPC_DEFAULT_SENDER_ADDRESS

#### EnableL2SuggestedGasPricePolling

- Description: When true, gas price suggestions will be updated  
- Default: true   
- env: ZKEVM_NODE_RPC_ENABLE_L2_SUGGESTED_GAS_PRICE_POLLING


## Synchronizer

> Used by: `synchronizer`

#### SyncInterval

- Description: Amount of time waited between sync loops                   
- Default: "0s"    
- env: ZKEVM_NODE_SYNCHRONIZER_SYNC_INTERVAL

#### SyncChunkSize

- Description: Amount of L1 blocks fetched per sync loop                  
- Default: 100     
- env: ZKEVM_NODE_SYNCHRONIZER_CHUNK_SIZE

#### GenBlockNumber

- Description: L1 block in which the rollup smart contracts were deployed 
- Default: 67      
- env: ZKEVM_NODE_SYNCHRONIZER_GEN_BLOKC_NUMBER

## Sequencer

> Used by: `sequencer`

#### NotSyncedWait

- Description: Time too wait when the sequencer is waiting for the synchronizer to catch up with L1
- Default: "1s"
- env: ZKEVM_NODE_SEQUENCER_NOT_SYNCED_WAIT

#### WaitPeriodSendSequence

- Description: Frequency in which the send sequences to L1 loop is run (doesn't mean that sequences will be sent with this frequency)
- Default: "5s"     
- env: ZKEVM_NODE_SEQUENCER_WAIT_PERIOD_SEND_SEQUENCE

#### LastBatchVirtualizationTimeMaxWaitPeriod

- Description: If there has not been batches sent to L1 for longer than the specified duration, a batch will be sent even if the L1 tx could be more packed
- Default: "1h"     
- env: ZKEVM_NODE_SEQUENCER_LAST_BATCH_VIRTUALIZATION_TIME_MAX_WAIT_PERIOD

#### BlocksAmountForTxsToBeDeleted

- Description: Will delete txs from the pool that have been mined on L1 (virtual state) for more than the indicated amount of L1 blocks
- Default: 100
- env: ZKEVM_NODE_SEQUENCER_BLOCKS_AMOUNT_FORTXS_TO_BE_DELETED
      
#### FrequencyToCheckTxsForDelete

- Description: Frequency in which the deletion process of already mined txs will happen
- Default: "12h"    
- env: ZKEVM_NODE_SEQUENCER_FREQUENCY_TO_CHECK_TXS_FOR_DELETE

#### MaxTxsPerBatch

- Description: Maximum amount of txs that the sequencer will include in a batch
- Default: 150      
- env: ZKEVM_NODE_SEQUENCER_MAX_TXS_PER_BATCH

#### MaxBatchBytesSize

- Description: Sequencer will close a batch before it reaches the indicated size
- Default: 129848   
- env: ZKEVM_NODE_SEQUENCER_MAX_BATCH_BYTES_SIZE

#### MaxCumulativeGasUsed

- Description: Sequencer will close a batch before it consumes the indicated amount of gas
- Default: 30000000 
- env: ZKEVM_NODE_SEQUENCER_MAX_CUMULATIVE_GAS_USED

#### MaxKeccakHashes

- Description: Sequencer will close a batch before it consumes the indicated amount of Keccack hashes
- Default: 468      
- env: ZKEVM_NODE_SEQUENCER_MAX_KECCAK_HASHES

#### MaxPoseidonHashes

- Description: Sequencer will close a batch before it consumes the indicated amount of Poseidon hashes
- Default: 279620   
- env: ZKEVM_NODE_SEQUENCER_MAX_POSEIDON_HASHES

#### MaxPoseidonPaddings

- Description: Sequencer will close a batch before it consumes the indicated amount of Poseidon paddings
- Default: 149796   
- env: ZKEVM_NODE_SEQUENCER_MAX_POSEIDON_PADDINGS

#### MaxMemAligns

- Description: Sequencer will close a batch before it consumes the indicated amount of memory alignments
- Default: 262144   
- env: ZKEVM_NODE_SEQUENCER_MAX_MEM_ALIGNS

#### MaxArithmetics

- Description: Sequencer will close a batch before it consumes the indicated amount of arithmetic operations
- Default: 262144   
- env: ZKEVM_NODE_SEQUENCER_MAX_ARITHEMTICS

#### MaxBinaries

- Description: Sequencer will close a batch before it consumes the indicated amount of binary operations
- Default: 262144   
- env: ZKEVM_NODE_SEQUENCER_MAX_BINARIES

#### MaxSteps

- Description: Sequencer will close a batch before it consumes the indicated amount of steps
- Default: 8388608  
- env: ZKEVM_NODE_SEQUENCER_MAX_STEPS

#### WeightBatchBytesSize

- Description: Factor multiplied on the size by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_BATCH_BYTE_SIZE

#### WeightCumulativeGasUsed

- Description: Factor multiplied on the gas used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_CUMULATIVE_GAS_USED

#### WeightKeccakHashes

- Description: Factor multiplied on the keccack hashes used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_KECCACK_HASHES

#### WeightPoseidonHashes

- Description: Factor multiplied on the Poseidon hashes used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_POSEIDON_HASHES

#### WeightPoseidonPaddings

- Description: Factor multiplied on the Poseidon paddings used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_POSEIDON_PADDINGS

#### WeightMemAligns

- Description: Factor multiplied on the mem aligns used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_MEM_ALIGNS

#### WeightArithmetics

- Description: Factor multiplied on the arithemtics used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_ARITHMETICS

#### WeightBinaries

- Description: Factor multiplied on the binnaries used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_BINNARIES

#### WeightSteps

- Description: Factor multiplied on the steps used by a tx to set the efficiency score
- Default: 1        
- env: ZKEVM_NODE_SEQUENCER_WEIGHT_STEPS

#### TxLifetimeCheckTimeout

- Description: Frequency in which the expire logic for pending txs will be triggered
- Default: "10m"    
- env: ZKEVM_NODE_SEQUENCER_TX_LIFETIME_CHECK_TIMEOUT

#### MaxTxLifetime

- Description: Maximum time for txs to be in pending state
- Default: "3h"     
- env: ZKEVM_NODE_SEQUENCER_MAX_TX_LIFETIME

#### MaxTxSizeForL1

- Description: Max size of an L1 tx
- Default: 131072   
- env: ZKEVM_NODE_SEQUENCER_MAX_TX_SIZE_FOR_L1

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

- Description: Waiting period between iterations for checking L1 timeout
- Default: 10s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_CLOSING_SIGNALS_MANAGER_WAIT_FOR_CHECKING_L1_TIMEOUT

#### ClosingSignalsManagerWaitForCheckingGER

- Description: Waiting period between iterations for checking new GERs
- Default: 10s
- env: ZKEVM_NODE_SEQUENCER_FINALIZER_CLOSING_SIGNALS_MANAGER_WAIT_FOR_CHECKING_GER

#### ClosingSignalsManagerWaitForCheckingForcedBatches

- Description: Waiting period between iterations for checking new forced batches
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

#### Host

- Description: Host for the gRPC interface in which provers will conect to
- Default: "0.0.0.0"
- env: ZKEVM_NODE_AGGREGATOR_HOST

#### Port

- Description: Port for the gRPC interface in which provers will conect to
- Default: 50081
- env: ZKEVM_NODE_AGGREGATOR_PORT

#### RetryTime

- Description: RetryTime is the time the aggregator main loop sleeps if there are no proofs to aggregate or batches to generate proofs. It is also used in the isSynced loop
- Default: "5s"
- env: ZKEVM_NODE_AGGREGATOR_RETRY_TIME

#### VerifyProofInterval

- Description: VerifyProofInterval is the interval of time to verify/send an proof in L1
- Default: "90s"
- env: ZKEVM_NODE_AGGREGATOR_VERIFY_PROOF_INTERVAL

#### TxProfitabilityCheckerType

- Description: type for checking is it profitable for aggregator to validate batch
- Default: "acceptall"
- env: ZKEVM_NODE_AGGREGATOR_TX_PROFITABILITY_CHECKER_TYPE

#### TxProfitabilityMinReward

- Description: min reward for base tx profitability checker when aggregator will validate batch this parameter is used for the base tx profitability checker
- Default: "1.1"
- env: ZKEVM_NODE_AGGREGATOR_TX_PROFITABILITY_MIN_REWARD

#### ProofStatePollingInterval

- Description: is the interval time to polling the prover about the generation state of a proof
- Default: "5s"
- env: ZKEVM_NODE_AGGREGATOR_PROOF_STATE_POLLING_INTERVAL

#### CleanupLockedProofsInterval

- Description: is the interval of time to clean up locked proofs.
- Default: "2m"
- env: ZKEVM_NODE_AGGREGATOR_CLEANUP_LOCKED_PROOFS_INTERVAL

#### GeneratingProofCleanupThreshold

- Description: represents the time interval after which a proof in generating state is considered to be stuck and allowed to be cleared.
- Default: "10m"
- env: ZKEVM_NODE_AGGREGATOR_GENERATING_PROOF_CLEANUP_TRESHOLD


## L2GasPriceSuggester

> Used by: `l2gaspricer`

#### Type

- Description: type of L2 gas suggestion strategy. Available values:
  - "default": default gas price from config is set
  - "lastnbatches": calculate average gas tip from last n batches.
  - "follower": calculate the gas price basing on the L1 gasPrice.
- Default: "default"
- env: ZKEVM_NODE_L2_GAS_PRICE_SUGGESTER_TYPE

#### DefaultGasPriceWei

- Description: suggested gas price when using the `Type: default`
- Default: 1000000000
- env: ZKEVM_NODE_L2_GAS_PRICE_SUGGESTER_DEFAULT_GAS_PRICE_WEI

## MTClient

> Used by: all components

#### URI

- Description: URI of the MT service aka HashDB, which is part of the prover binary
- Default: "127.0.0.1:50061"
- env: ZKEVM_NODE_MT_CLIENT

## Executor

> Used by: all components

#### URI

- Description: URI of the executor, which is part of the prover binary
- Default: "127.0.0.1:50071"
- env: ZKEVM_NODE_MT_CLIENT

## BroadcastServer

> Deprecated

## Metrics

> Used by: all components

#### Host

- Description: host of the Prometheus server
- Default: "0.0.0.0"
- env: ZKEVM_NODE_METRICS_HOST

#### Port

- Description: port of the Prometheus server
- Default: 9091
- env: ZKEVM_NODE_METRICS_ 

#### Enabled

- Description: enable/disable Prometheus server
- Default: false
- env: ZKEVM_NODE_METRICS_ 