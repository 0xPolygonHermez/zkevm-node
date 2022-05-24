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

[RPC]
Host = "0.0.0.0"
Port = 8123
MaxRequestsPerIPAndSecond = 50

[Synchronizer]
SyncInterval = "0s"
SyncChunkSize = 100

[Sequencer]
AllowNonRegistered = "false"
IntervalToProposeBatch = "15s"
SyncedBlockDif = 1
InitBatchProcessorIfDiffType = "synced"
    [Sequencer.Strategy]
        [Sequencer.Strategy.TxSelector]
            TxSelectorType = "acceptall"
            TxSorterType = "bycostandnonce"
        [Sequencer.Strategy.TxProfitabilityChecker]
            TxProfitabilityCheckerType = "acceptall"
            MinReward = "1.1"
			RewardPercentageToAggregator = 50
	[Sequencer.PriceGetter]
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

[ExecutorServer]
Host = "0.0.0.0"
Port = 0
`
