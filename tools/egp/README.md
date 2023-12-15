# EGP TOOL
## Introduction
A Go tool to analyze and simulate the use of Effective Gas Price (EGP) feature. This tool has 2 main functionalities:

- Calculate real statistics based on the EGP logs stored in the State database.

- Simulate results that would have been obtained using different parameters for the EGP.

## Running the tool
### Help
Executing the following command line will display the help with the available parameters in the tool.
```sh
go run main.go help
```
```
NAME:
   main - Analyze stats for EGP

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --from value           stats from L2 block onwards (default: 18446744073709551615)
   --to value             stats until L2 block (optional) (default: 18446744073709551615)
   --showerror            show transactions with EGP errors (default: false)
   --showloss             show transactions with losses (default: false)
   --showdetail           show full detail record when show loss/error (default: false)
   --cfg value, -c value  configuration file
   --db value             DB connection string: "host=xxx port=xxx user=xxx dbname=xxx password=xxx"
   --help, -h             show help
```

### Statistics
Running the tool without specifying a configuration file will only calculate the real EGP statistics from the `state`.`transaction` table.

> The `--db` parameter specifying the DB connection string is required

```sh
go run main.go --db "host=X port=X user=X dbname=X password=X"
```
```
EGP REAL STATS:
Total Tx.........: [100000]
Error Tx.........: [0] (0.00%)
Total No EGP info: [0] (0.00%)
Total Tx EGP info: [100000] (100.00%)
    EGP enable.......: [0] (0.00%)
    Reprocessed Tx...: [82] (0.08%)
        Suspicious Tx....: [2] (2.44%)
    Final gas:
        Used EGP1........: [99890] (99.89%)
        Used EGP2........: [50] (0.05%)
        Used User Gas....: [60] (0.06%)
        Used Weird Gas...: [0] (0.00%)
    Gas average..........: [11413043906] (11413 MWei) (0.000000011 ETH)
    Loss count.......: [58] (0.06%)
    Loss total.......: [194697356796] (194697 MWei) (0.000000195 ETH)
    Loss average.....: [3356850979] (3356 MWei) (0.000000003 ETH)
```

### Simulation
Specifying the parameter `--cfg` with a configuration file, the tool will in addition to calculate real statistics, perform a simulation of the results with that config file parameters.

#### Config file parameters (e.g. `config.toml`)

```toml
# gas cost of 1 byte
ByteGasCost = 16

# gas cost of 1 byte zero
ZeroGasCost = 4

# L2 network profit factor
NetProfitFactor = 1.0

# L1 gas price factor
L1GasPriceFactor = 0.25

# L2 gas price suggester factor
L2GasPriceSugFactor = 0.5

# Max final deviation percentage
FinalDeviationPct = 10

# Min gas price allowed
MinGasPriceAllowed = 1000000000

```

```sh
go run main.go --cfg cfg/egp1.config.toml --db "host=X port=X user=X dbname=X password=X"
```
```
EGP REAL STATS:
Total Tx.........: [100000]
Error Tx.........: [0] (0.00%)
Total No EGP info: [0] (0.00%)
Total Tx EGP info: [100000] (100.00%)
    EGP enable.......: [0] (0.00%)
    Reprocessed Tx...: [82] (0.08%)
        Suspicious Tx....: [2] (2.44%)
    Final gas:
        Used EGP1........: [99890] (99.89%)
        Used EGP2........: [50] (0.05%)
        Used User Gas....: [60] (0.06%)
        Used Weird Gas...: [0] (0.00%)
    Gas average..........: [11413043906] (11413 MWei) (0.000000011 ETH)
    Loss count.......: [58] (0.06%)
    Loss total.......: [194697356796] (194697 MWei) (0.000000195 ETH)
    Loss average.....: [3356850979] (3356 MWei) (0.000000003 ETH)

EGP SIMULATION STATS:
Total Tx.........: [100000]
Error Tx.........: [0] (0.00%)
Total No EGP info: [0] (0.00%)
Total Tx EGP info: [100000] (100.00%)
    EGP enable.......: [0] (0.00%)
    Reprocessed Tx...: [110] (0.11%)
        Suspicious Tx....: [2] (1.82%)
    Final gas:
        Used EGP1........: [99867] (99.87%)
        Used EGP2........: [78] (0.08%)
        Used User Gas....: [55] (0.06%)
        Used Weird Gas...: [0] (0.00%)
    Gas average..........: [9594986075] (9594 MWei) (0.000000010 ETH)
    Loss count.......: [53] (0.05%)
    Loss total.......: [90503106670] (90503 MWei) (0.000000091 ETH)
    Loss average.....: [1707605786] (1707 MWei) (0.000000002 ETH)
PARAMS: byte[16] zero[4] netFactor[1.00] L1factor[0.20] L2sugFactor[0.50] devPct[10] minGas[1000000000]
```
