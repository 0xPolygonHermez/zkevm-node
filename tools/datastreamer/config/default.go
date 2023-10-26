package config

// DefaultValues is the default configuration
const DefaultValues = `
ChainID = 1440

[Online]
URI = "zkevm-sequencer:6900"
StreamType = 1

[Offline]
Port = 6901
Filename = "datastreamer.bin"

[StateDB]
User = "state_user"
Password = "state_password"
Name = "state_db"
Host = "localhost"
Port = "5432"
EnableLog = false	
MaxConns = 200

[Executor]
URI = "zkevm-prover:50071"
MaxGRPCMessageSize = 100000000

[MerkeTree]
URI = "zkevm-prover:50061"

[Log]
Environment = "development" # "production" or "development"
Level = "info"
Outputs = ["stderr"]
`
