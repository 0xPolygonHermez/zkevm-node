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

[RPC]
Host = "0.0.0.0"
Port = 8123
ChainID = 2576980377

[Synchronizer]
GenesisBlock = 1

[Synchronizer.Etherman]
URL = "http://localhost"
PoEAddress = "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"
PrivateKeyPath = "./test/test.keystore"
PrivateKeyPassword = "testonly"

[Sequencer]
IntervalToProposeBatch = "15s"
SyncedBlockDif = 1

[Sequencer.Etherman]
URL = "http://localhost"
PoEAddress = "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"
PrivateKeyPath = "../test/test.keystore"
PrivateKeyPassword = "testonly"

[Aggregator]
IntervalToConsolidateState = "3s"

[Aggregator.Etherman]
URL = "http://localhost"
PoEAddress = "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"
PrivateKeyPath = "../test/test.keystore"
PrivateKeyPassword = "testonly"

[Prover]
ProverURI = "0.0.0.0:50051"
`
