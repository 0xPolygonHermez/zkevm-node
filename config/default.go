package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false

[Log]
Environment = "development" # "production" or "development"
Level = "debug"
Outputs = ["stderr"]

[StateDB]
User = "state_user"
Password = "state_password"
Name = "state_db"
Host = "localhost"
Port = "5432"
EnableLog = false
MaxConns = 200

[Pool]
FreeClaimGasLimit = 150000
MaxTxBytesSize=30132
MaxTxDataBytesSize=30000
DefaultMinGasPriceAllowed = 1000000000
MinAllowedGasPriceInterval = "5m"
PollMinAllowedGasPriceInterval = "15s"
	[Pool.DB]
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
PoEAddr = "0x610178dA211FEF7D417bC0e6FeD39F05609AD788"
MaticAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
GlobalExitRootManagerAddr = "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"
MultiGasProvider = true
	[Etherman.Etherscan]
		ApiKey = ""

[EthTxManager]
FrequencyToMonitorTxs = "1s"
WaitTxToBeMined = "2m"
ForcedGas = 0

[RPC]
Host = "0.0.0.0"
Port = 8123
ReadTimeoutInSec = 60
WriteTimeoutInSec = 60
MaxRequestsPerIPAndSecond = 50
SequencerNodeURI = ""
BroadcastURI = "127.0.0.1:61090"
DefaultSenderAddress = "0x1111111111111111111111111111111111111111"
EnableL2SuggestedGasPricePolling = true
	[RPC.WebSockets]
		Enabled = false
		Port = 8133

[Synchronizer]
SyncInterval = "0s"
SyncChunkSize = 100
GenBlockNumber = 67
trustedSequencerURL = ""

[Sequencer]
WaitPeriodPoolIsEmpty = "1s"
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
	[Sequencer.Finalizer]
		GERDeadlineTimeoutInSec = "5s"
		ForcedBatchDeadlineTimeoutInSec = "60s"
		SendingToL1DeadlineTimeoutInSec = "20s"
		SleepDurationInMs = "100ms"
		ResourcePercentageToCloseBatch = 10
		GERFinalityNumberOfBlocks = 64
		ClosingSignalsManagerWaitForCheckingL1Timeout = "10s"
		ClosingSignalsManagerWaitForCheckingGER = "10s"
		ClosingSignalsManagerWaitForCheckingForcedBatches = "10s"
		ForcedBatchesFinalityNumberOfBlocks = 64
		SenderAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		PrivateKeys = [{Path = "/pk/sequencer.keystore", Password = "testonly"}]
	[Sequencer.DBManager]
		PoolRetrievalInterval = "500ms"

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
Type = "default"
DefaultGasPriceWei = 1000000000

[Prover]
ProverURI = "0.0.0.0:50051"

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
