package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false

[Log]
Level = "debug"
Outputs = ["stdout"]

[Database]
User = "test_user"
Password = "test_password"
Name = "test_db"
Host = "localhost"
Port = "5432"
EnableLog = false
MaxConns = 200

[Etherman]
URL = "http://localhost:8545"
L1ChainID = 1337
PrivateKeyPath = "./test/test.keystore"
PrivateKeyPassword = "testonly"
PoEAddr = "0x0000000000000000000000000000000000000000"
MaticAddr = "0x0000000000000000000000000000000000000000"
GlobalExitRootManagerAddr = "0x0000000000000000000000000000000000000000"

[EthTxManager]
MaxSendBatchTxRetries = 10
MaxVerifyBatchTxRetries = 10
FrequencyForResendingFailedSendBatches = "1s"
FrequencyForResendingFailedVerifyBatch = "1s"

[RPC]
Host = "0.0.0.0"
Port = 8123
MaxRequestsPerIPAndSecond = 50
SequencerNodeURI = ""
BroadcastURI = "127.0.0.1:61090"
DefaultSenderAddress = "0x1111111111111111111111111111111111111111"
L2ChainID = 1000

[Synchronizer]
SyncInterval = "0s"
SyncChunkSize = 100
TrustedSequencerURI = ""
GenBlockNumber = 1

[Sequencer]
WaitPeriodPoolIsEmpty = "1s"
WaitPeriodSendSequence = "15s"
LastBatchVirtualizationTimeMaxWaitPeriod = "300s"
WaitBlocksToUpdateGER = 10
MaxTimeForBatchToBeOpen = "15s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
MaxCumulativeGasUsed = 30000000
MaxKeccakHashes = 468
MaxPoseidonHashes = 279620
MaxPoseidonPaddings = 149796
MaxMemAligns = 262144
MaxArithmetics = 262144
MaxBinaries = 262144
MaxSteps = 8388608
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
`
