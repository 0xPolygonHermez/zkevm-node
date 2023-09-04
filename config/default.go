package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false
ForkUpgradeBatchNumber = 0
ForkUpgradeNewForkId = 0

[Log]
Environment = "development" # "production" or "development"
Level = "info"
Outputs = ["stderr"]

[State]
AccountQueue = 64
	[State.DB]
	User = "state_user"
	Password = "state_password"
	Name = "state_db"
	Host = "zkevm-state-db"
	Port = "5432"
	EnableLog = false	
	MaxConns = 200
	[State.Batch]
		[State.Batch.Constraints]
		MaxTxsPerBatch = 300
		MaxBatchBytesSize = 120000
		MaxCumulativeGasUsed = 30000000
		MaxKeccakHashes = 2145
		MaxPoseidonHashes = 252357
		MaxPoseidonPaddings = 135191
		MaxMemAligns = 236585
		MaxArithmetics = 236585
		MaxBinaries = 473170
		MaxSteps = 7570538

[Pool]
IntervalToRefreshBlockedAddresses = "5m"
IntervalToRefreshGasPrices = "5s"
MaxTxBytesSize=100132
MaxTxDataBytesSize=100000
DefaultMinGasPriceAllowed = 1000000000
MinAllowedGasPriceInterval = "5m"
PollMinAllowedGasPriceInterval = "15s"
AccountQueue = 64
GlobalQueue = 1024
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
ForkIDChunkSize = 20000
MultiGasProvider = false
	[Etherman.Etherscan]
		ApiKey = ""

[EthTxManager]
FrequencyToMonitorTxs = "1s"
WaitTxToBeMined = "2m"
ForcedGas = 0
GasPriceMarginFactor = 1
MaxGasPriceLimit = 0

[RPC]
Host = "0.0.0.0"
Port = 8545
ReadTimeout = "60s"
WriteTimeout = "60s"
MaxRequestsPerIPAndSecond = 500
SequencerNodeURI = ""
EnableL2SuggestedGasPricePolling = true
TraceBatchUseHTTPS = true
BatchRequestsEnabled = false
BatchRequestsLimit = 20
	[RPC.WebSockets]
		Enabled = true
		Host = "0.0.0.0"
		Port = 8546
		ReadLimit = 104857600

[Synchronizer]
SyncInterval = "1s"
SyncChunkSize = 100
TrustedSequencerURL = "" # If it is empty or not specified, then the value is read from the smc
UseParallelModeForL1Synchronization = true
	[Synchronizer.L1ParallelSynchronization]
		NumberOfParallelOfEthereumClients = 2
		CapacityOfBufferingRollupInfoFromL1 = 10

[Sequencer]
WaitPeriodPoolIsEmpty = "1s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
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
		StopSequencerOnBatchNum = 0
		SequentialReprocessFullBatch = false
	[Sequencer.DBManager]
		PoolRetrievalInterval = "500ms"
		L2ReorgRetrievalInterval = "5s"
	[Sequencer.EffectiveGasPrice]
		MaxBreakEvenGasPriceDeviationPercentage = 10
		L1GasPriceFactor = 0.25
		ByteGasCost = 16
		MarginFactor = 1
		Enabled = false

[SequenceSender]
WaitPeriodSendSequence = "5s"
LastBatchVirtualizationTimeMaxWaitPeriod = "5s"
MaxTxSizeForL1 = 131072
L2Coinbase = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
PrivateKey = {Path = "/pk/sequencer.keystore", Password = "testonly"}

[Aggregator]
Host = "0.0.0.0"
Port = 50081
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
MaxGasPriceWei = 0
CleanHistoryPeriod = "1h"
CleanHistoryTimeRetention = "5m"

[MTClient]
URI = "zkevm-prover:50061"

[Executor]
URI = "zkevm-prover:50071"
MaxResourceExhaustedAttempts = 3
WaitOnResourceExhaustion = "1s"
MaxGRPCMessageSize = 100000000

[Metrics]
Host = "0.0.0.0"
Port = 9091
Enabled = false

[HashDB]
User = "prover_user"
Password = "prover_pass"
Name = "prover_db"
Host = "zkevm-state-db"
Port = "5432"
EnableLog = false
MaxConns = 200
`
