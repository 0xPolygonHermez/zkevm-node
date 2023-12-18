
# Benchmark Sequencer Scripts

This repository contains scripts to benchmark a sequencer. The main script is written in Go and can be used to run a series of commands and perform various operations.

## Usage

### 1. Clone the repository:
   ```
   git clone git@github.com:0xPolygonHermez/zkevm-node.git
   cd zkevm-node/test/benchmarks/sequencer/scripts
   ```

### 2. Setup Environment Variables:
   Copy the `.env.example` file to `.env` and populate it with the appropriate values. 
   
   #### Required environment variables are:
   - `BASTION_HOST`: The IP address or domain name of the bastion host. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_DBNAME`: Database name for the pool. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_EP`: Endpoint for the pool database. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_PASS`: Password for the pool database. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_USER`: User for the pool database. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `SEQUENCER_IP`: The IP address of the sequencer. (`sequencer.zkevm-public.aws` for `public testnet`, `sequencer.zkevm-internal.aws` for `internal testnet`, `sequencer.zkevm-dev.aws` for `dev testnet`)
   - `RPC_URL`: The URL for the Remote Procedure Call (RPC) server. (From `Deployments.doc` under `Public URLs` section as a bullet point to `RPC` for the specific `Environment`)
   - `CHAIN_ID`: The ID of the blockchain network. (From `Deployments.doc` under `Public URLs` section as a bullet point to `RPC` for the specific `Environment`)
   - `PRIVATE_KEY`: Your private key.

   #### Optional environment variables:
   - `POOLDB_PORT`: Port for the pool database. (Default is `5433`)

   Example:
   ```
   cp .env.example .env
   nano .env
   ```
### 3. Run the Benchmark Script:
   Run the `main.go` script with the following command-line flags:
   - `--type`: The type of transactions to test. Accepted values are `eth`, `erc20` or `uniswap`.
   - `--num-ops` (optional): The number of operations to run. Default is 200.
   - `--help` (optional): Display the help message.

   Example:
   ```
   go run main.go --type erc20 --sequencer-ip <Private IP>
   ```

## Notes

- Ensure that the `.env` file exists and contains all the required environment variables before running the script.
- The script will perform various operations based on the provided command-line flags and environment variables.
- Ensure that Go is installed on your system to run the script.