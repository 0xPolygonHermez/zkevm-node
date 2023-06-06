package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false

[Log]
Environment = "development" # "production" or "development"
Level = "info"
Outputs = ["stderr"]

[StateDB]
User = "state_user"
Password = "state_password"
Name = "state_db"
Host = "zkevm-state-db"
Port = "5432"
EnableLog = false
MaxConns = 200

[Pool]
IntervalToRefreshBlockedAddresses = "5m"
MaxTxBytesSize=100132
MaxTxDataBytesSize=100000
DefaultMinGasPriceAllowed = 1000000000
MinAllowedGasPriceInterval = "5m"
PollMinAllowedGasPriceInterval = "15s"
	[Pool.DB]
	User = "pool_user"
	Password = "pool_password"
	Name = "pool_db"
	Host = "zkevm-pool-db"
	Port = "5432"
	EnableLog = false
	MaxConns = 200

[Etherman]
URL = "http://localhost:8545"
MultiGasProvider = false
	[Etherman.Etherscan]
		ApiKey = ""

[EthTxManager]
FrequencyToMonitorTxs = "1s"
WaitTxToBeMined = "2m"
ForcedGas = 0

[RPC]
Host = "0.0.0.0"
Port = 8545
ReadTimeout = "60s"
WriteTimeout = "60s"
MaxRequestsPerIPAndSecond = 500
SequencerNodeURI = ""
EnableL2SuggestedGasPricePolling = true
	[RPC.WebSockets]
		Enabled = true
		Host = "0.0.0.0"
		Port = 8546

[Synchronizer]
SyncInterval = "1s"
SyncChunkSize = 100
TrustedSequencerURL = "" # If it is empty or not specified, then the value is read from the smc

[Sequencer]
WaitPeriodPoolIsEmpty = "1s"
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
	[Sequencer.Finalizer]
		GERDeadlineTimeout = "5s"
		ForcedBatchDeadlineTimeout = "60s"
		SleepDuration = "100ms"
		ResourcePercentageToCloseBatch = 10
		GERFinalityNumberOfBlocks = 64
		ClosingSignalsManagerWaitForCheckingL1Timeout = "10s"
		ClosingSignalsManagerWaitForCheckingGER = "10s"
		ClosingSignalsManagerWaitForCheckingForcedBatches = "10s"
		ForcedBatchesFinalityNumberOfBlocks = 64
		TimestampResolution = "10s"	
	[Sequencer.DBManager]
		PoolRetrievalInterval = "500ms"
		L2ReorgRetrievalInterval = "5s"
	[Sequencer.Worker]
		ResourceCostMultiplier = 1000

[SequenceSender]
WaitPeriodSendSequence = "5s"
LastBatchVirtualizationTimeMaxWaitPeriod = "5s"
MaxTxSizeForL1 = 131072
SenderAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
PrivateKeys = [{Path = "/pk/sequencer.keystore", Password = "testonly"}]

[PriceGetter]
Type = "default"
DefaultPrice = "2000"

[Aggregator]
Host = "0.0.0.0"
Port = 50081
ForkId = 2
RetryTime = "5s"
VerifyProofInterval = "90s"
TxProfitabilityCheckerType = "acceptall"
TxProfitabilityMinReward = "1.1"
ProofStatePollingInterval = "5s"
CleanupLockedProofsInterval = "2m"
GeneratingProofCleanupThreshold = "10m"

[L2GasPriceSuggester]
Type = "follower"
UpdatePeriod = "10s"
Factor = 0.15
DefaultGasPriceWei = 2000000000

[MTClient]
URI = "zkevm-prover:50061"

[Executor]
URI = "zkevm-prover:50071"
MaxGRPCMessageSize = 100000000

[Metrics]
Host = "0.0.0.0"
Port = 9091
Enabled = false
`
