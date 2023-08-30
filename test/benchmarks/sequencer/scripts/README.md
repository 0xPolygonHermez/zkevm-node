
# Benchmark Sequencer Scripts

This repository contains scripts to benchmark a sequencer. The main script is written in Go and can be used to run a series of commands and perform various operations.

## Usage

1. **Clone the repository**:
   ```
   git clone git@github.com:0xPolygonHermez/zkevm-node.git
   cd zkevm-node/test/benchmarks/sequencer/scripts
   ```

2. **Setup Environment Variables**:
   Copy the `.env.example` file to `.env` and populate it with the appropriate values. The following environment variables are required:
   - `BASTION_HOST`: The IP address or domain name of the bastion host. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_DBNAME`: Database name for the pool. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_EP`: Endpoint for the pool database. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_PASS`: Password for the pool database. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_USER`: User for the pool database. (From `Deployments.doc` under `BASH VARIABLES` section for the specific `Environment`)
   - `POOLDB_LOCALPORT`: Local port for accessing the pool database. (From `Deployments.doc` under `Access to databases` section for the specific `Environment`)
   - `RPC_URL`: The URL for the Remote Procedure Call (RPC) server. (From `Deployments.doc` under `Public URLs` section as a bullet point to `RPC` for the specific `Environment`)
   - `CHAIN_ID`: The ID of the blockchain network. (From `Deployments.doc` under `Public URLs` section as a bullet point to `RPC` for the specific `Environment`)
   - `PRIVATE_KEY`: Your private key.

   Example:
   ```
   cp .env.example .env
   nano .env
   ```

3. **Get the Sequencer IP Address**:
   The value for the `sequencer-ip` argument needs to be obtained from the `Amazon Elastic Container Service`:
   - Go to https://polygon-technology.awsapps.com/start#/
   - Select the `account`/`environment` (`zkEVM-dev` for `dev` or `zkEVM-staging` for `internal` or `public` testnets).
   - Go to the `Amazon Elastic Container Service`.
   - Navigate to the `Service` section and select `sequencer`.
   - In the `Task` section, select the only `task` definition.
   - Find the `Private IP` field and use it as the value for the `sequencer-ip` argument.

4. **Run the Benchmark Script**:
   Run the `main.go` script with the following command-line flags:
   - `--type`: The type of transactions to test. Accepted values are `eth`, `erc20` or `uniswap`.
   - `--sequencer-ip`: The IP address of the sequencer (**obtained in step 3**).
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