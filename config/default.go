package config

// DefaultValues is the default configuration
const DefaultValues = `
[Log]
Level = "debug"
Outputs = ["stdout"]

[Database]
Name = "polygon-hermez"
User = "hermez"
Password = "polygon"
Host = "localhost"
Port = "5432"

[Etherman]
URL = "http://localhost"
PoEAddress = "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"
PrivateKeyPath = "../test/test.keystore"
PrivateKeyPassword = "testonly"

[Synchronizer]
GenesisBlock = 1

[RPC]
Host = "0.0.0.0"
Port = 8123

[Sequencer]
IntervalToProposeBatch = "15s"
SyncedBlockDif = 1

[Aggregator]
IntervalToConsolidateState = "3s"

[Prover]
ProverURI = "0.0.0.0:50051"

[State]
Arity = 16
`
