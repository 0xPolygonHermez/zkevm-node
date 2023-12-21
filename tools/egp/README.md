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
   --showreprocess        show transactions reprocessed (default: false)
   --showdetail           show full detail record when showing error/loss/reprocess (default: false)
   --showalways           show always full detailed record (default: false)
   --cfg value, -c value  simulation configuration file
   --onlycfg              show only simulation results (default: false)
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
Total Tx.........: [10000]
Error Tx.........: [0] (0.00%)
Total No EGP info: [0] (0.00%)
Total Tx EGP info: [10000] (100.00%)
    EGP enable.......: [0] (0.00%)
    Reprocessed Tx...: [7] (0.07%)
        Suspicious Tx....: [0] (0.00%)
    Final gas:
        Used EGP1........: [9989] (99.89%)
        Used EGP2........: [3] (0.03%)
        Used User Gas....: [8] (0.08%)
        Used Weird Gas...: [0] (0.00%)
    Gas price avg........: [18941296931] (18.941 GWei) (0.000000019 ETH)
    Tx fee avg...........: [1319258335287442] (1319258.335 GWei) (0.001319258 ETH)
    Gas pri.avg preEGP...: [5413503250] (5.414 GWei) (0.000000005 ETH)
    Tx fee avg preEGP....: [421947362151699] (421947.362 GWei) (0.000421947 ETH)
    Diff fee EGP-preEGP..: [8973109731357435904] (8973109731.357 Gwei) (8.973109731 ETH)
    Loss count.......: [8] (0.08%)
    Loss total.......: [43211133382] (43.211 GWei) (0.000000043 ETH)
    Loss average.....: [5401391673] (5 GWei) (0.000000005 ETH)
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
NetProfitFactor = 1.2

# L1 gas price factor
L1GasPriceFactor = 0.04

# L2 gas price suggester factor
L2GasPriceSugFactor = 0.30

# Max final deviation percentage
FinalDeviationPct = 10

# Min gas price allowed
MinGasPriceAllowed = 1000000000

# L2 gas price suggester factor pre EGP
L2GasPriceSugFactorPreEGP = 0.1
```

```sh
go run main.go --cfg cfg/config.toml --db "host=X port=X user=X dbname=X password=X"
```
```
EGP REAL STATS:
Total Tx.........: [10000]
Error Tx.........: [0] (0.00%)
Total No EGP info: [0] (0.00%)
Total Tx EGP info: [10000] (100.00%)
    EGP enable.......: [0] (0.00%)
    Reprocessed Tx...: [7] (0.07%)
        Suspicious Tx....: [0] (0.00%)
    Final gas:
        Used EGP1........: [9989] (99.89%)
        Used EGP2........: [3] (0.03%)
        Used User Gas....: [8] (0.08%)
        Used Weird Gas...: [0] (0.00%)
    Gas price avg........: [18941296931] (18.941 GWei) (0.000000019 ETH)
    Tx fee avg...........: [1319258335287442] (1319258.335 GWei) (0.001319258 ETH)
    Gas pri.avg preEGP...: [5413503250] (5.414 GWei) (0.000000005 ETH)
    Tx fee avg preEGP....: [421947362151699] (421947.362 GWei) (0.000421947 ETH)
    Diff fee EGP-preEGP..: [8973109731357425664] (8973109731.357 Gwei) (8.973109731 ETH)
    Loss count.......: [8] (0.08%)
    Loss total.......: [43211133382] (43.211 GWei) (0.000000043 ETH)
    Loss average.....: [5401391673] (5 GWei) (0.000000005 ETH)

EGP SIMULATION STATS:
Total Tx.........: [10000]
Error Tx.........: [0] (0.00%)
Total No EGP info: [0] (0.00%)
Total Tx EGP info: [10000] (100.00%)
    EGP enable.......: [0] (0.00%)
    Reprocessed Tx...: [16] (0.16%)
        Suspicious Tx....: [0] (0.00%)
    Final gas:
        Used EGP1........: [9867] (98.67%)
        Used EGP2........: [12] (0.12%)
        Used User Gas....: [121] (1.21%)
        Used Weird Gas...: [0] (0.00%)
    Gas price avg........: [9073552262] (9.074 GWei) (0.000000009 ETH)
    Tx fee avg...........: [519499850778700] (519499.851 GWei) (0.000519500 ETH)
    Gas pri.avg preEGP...: [5413503250] (5.414 GWei) (0.000000005 ETH)
    Tx fee avg preEGP....: [421947362151699] (421947.362 GWei) (0.000421947 ETH)
    Diff fee EGP-preEGP..: [975524886270010368] (975524886.270 Gwei) (0.975524886 ETH)
    Loss count.......: [121] (1.21%)
    Loss total.......: [194278383566] (194.278 GWei) (0.000000194 ETH)
    Loss average.....: [1605606476] (2 GWei) (0.000000002 ETH)
PARAMS: byte[16] zero[4] netFactor[1.20] L1factor[0.04] L2sugFactor[0.30] devPct[10] minGas[1000000000] L2sugPreEGP[0.10]
```
> To show only the result of the simulation, use the flag `--onlycfg`
