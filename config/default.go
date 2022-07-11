package config

// DefaultValues is the default configuration
const DefaultValues = `
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
URL = "http://localhost"
PrivateKeyPath = "./test/test.keystore"
PrivateKeyPassword = "testonly"

[EthTxManager]
MaxSendBatchTxRetries = 10
FrequencyForResendingFailedSendBatchesInMilliseconds = 1000

[RPC]
Host = "0.0.0.0"
Port = 8123
MaxRequestsPerIPAndSecond = 50
ChainID = 1001

[Synchronizer]
SyncInterval = "0s"
SyncChunkSize = 100

[Sequencer]
WaitPeriodPoolIsEmpty = "15s"
LastBatchVirtualizationTimeMaxWaitPeriod = "15s"
WaitBlocksToUpdateGER = 10
LastTimeBatchMaxWaitPeriod = "15s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
	[Sequencer.ProfitabilityChecker]
		SendBatchesEvenWhenNotProfitable = "true"

[PriceGetter]
Type = "default"
DefaultPrice = "2000"

[Aggregator]
IntervalFrequencyToGetProofGenerationStateInSeconds = "5s"
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
URI = "127.0.0.1:50060"

[Executor]
URI = "51.210.116.237:50071"

[BroadcastServer]
Host = "0.0.0.0"
Port = 61090

[BroadcastClient]
URI = "127.0.0.1:61090"
`
