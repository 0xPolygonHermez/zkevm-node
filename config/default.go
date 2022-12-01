package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false

[Log]
Level = "debug"
Outputs = ["stdout"]

[StateDB]
User = "state_user"
Password = "state_password"
Name = "state_db"
Host = "localhost"
Port = "5432"
EnableLog = false
MaxConns = 200

[PoolDB]
User = "pool_user"
Password = "pool_password"
Name = "pool_db"
Host = "localhost"
Port = "5432"
EnableLog = false
MaxConns = 200

[Etherman]
URL = "http://localhost:8545"
L1ChainID = 1337
PoEAddr = "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"
MaticAddr = "0x0165878A594ca255338adfa4d48449f69242Eb8F"
GlobalExitRootManagerAddr = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
MultiGasProvider = true
	[Etherscan]
		ApiKey = ""

[EthTxManager]
MaxSendBatchTxRetries = 10
MaxVerifyBatchTxRetries = 10
FrequencyForResendingFailedSendBatches = "1s"
FrequencyForResendingFailedVerifyBatch = "1s"
WaitTxToBeMined = "2m"
PercentageToIncreaseGasPrice = 10
PercentageToIncreaseGasLimit = 10

[RPC]
Host = "0.0.0.0"
Port = 8123
ReadTimeoutInSec = 60
WriteTimeoutInSec = 60
MaxRequestsPerIPAndSecond = 50
SequencerNodeURI = ""
BroadcastURI = "127.0.0.1:61090"
DefaultSenderAddress = "0x1111111111111111111111111111111111111111"
	[RPC.DB]
		User = "rpc_user"
		Password = "rpc_password"
		Name = "rpc_db"
		Host = "localhost"
		Port = "5432"
		EnableLog = false
		MaxConns = 200

[Synchronizer]
SyncInterval = "0s"
SyncChunkSize = 100
TrustedSequencerURI = ""
GenBlockNumber = 1

[Sequencer]
MaxSequenceSize = "2000000"
WaitPeriodPoolIsEmpty = "1s"
WaitPeriodSendSequence = "15s"
LastBatchVirtualizationTimeMaxWaitPeriod = "300s"
WaitBlocksToUpdateGER = 10
WaitBlocksToConsiderGerFinal = 10
ElapsedTimeToCloseBatchWithoutTxsDueToNewGER = "60s"
MaxTimeForBatchToBeOpen = "15s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
MaxTxsPerBatch = 150
MaxBatchBytesSize = 30000
MaxCumulativeGasUsed = 30000000
MaxKeccakHashes = 468
MaxPoseidonHashes = 279620
MaxPoseidonPaddings = 149796
MaxMemAligns = 262144
MaxArithmetics = 262144
MaxBinaries = 262144
MaxSteps = 8388608
MaxAllowedFailedCounter = 50
	[Sequencer.ProfitabilityChecker]
		SendBatchesEvenWhenNotProfitable = "true"

[PriceGetter]
Type = "default"
DefaultPrice = "2000"

[Aggregator]
IntervalFrequencyToGetProofGenerationState = "5s"
IntervalToConsolidateState = "3s"
TxProfitabilityCheckerType = "acceptall"
TxProfitabilityMinReward = "1.1"

[GasPriceEstimator]
Type = "default"
DefaultGasPriceWei = 1000000000

[Prover]
ProverURI = "0.0.0.0:50051"

[MTServer]
Host = "0.0.0.0"
Port = 50060
StoreBackend = "PostgreSQL"

[MTClient]
URI = "127.0.0.1:50061"

[Executor]
URI = "127.0.0.1:50071"

[BroadcastServer]
Host = "0.0.0.0"
Port = 61090

[Metrics]
Host = "0.0.0.0"
Port = 9091
Enabled = false
`
