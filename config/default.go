package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false
DefaultForkID = 1

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
PoEAddr = "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318"
MaticAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
GlobalExitRootManagerAddr = "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853"
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
	[RPC.WebSockets]
		Enabled = false
		Port = 8133

[Synchronizer]
SyncInterval = "0s"
SyncChunkSize = 100
GenBlockNumber = 66

[Sequencer]
MaxSequenceSize = "2000000"
WaitPeriodPoolIsEmpty = "1s"
WaitPeriodSendSequence = "5s"
LastBatchVirtualizationTimeMaxWaitPeriod = "5s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
MaxTxsPerBatch = 150
MaxBatchBytesSize = 150000
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
MaxAllowedFailedCounter = 50
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

[PriceGetter]
Type = "default"
DefaultPrice = "2000"

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
