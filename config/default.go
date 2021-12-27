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

[Etherman]
URL = "http://localhost"
PrivateKeyPath = "./test/test.keystore"
PrivateKeyPassword = "testonly"

[RPC]
Host = "0.0.0.0"
Port = 8123

[Synchronizer]
SyncInterval = "15s"

[Sequencer]
IntervalToProposeBatch = "15s"
SyncedBlockDif = 1
    [Sequencer.Strategy]
    Type = "acceptall"
    TxSorterType = "bycostandnonce"
    TxProfitabilityCheckerType = "base"
    MinReward = "1.1"
    PossibleTimeToSendTx = "60s"

[Aggregator]
IntervalToConsolidateState = "3s"
TxProfitabilityCheckerType = "acceptall"
TxProfitabilityMinReward = "1.1"

[Prover]
ProverURI = "0.0.0.0:50051"
`
