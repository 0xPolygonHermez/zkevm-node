# Production Setup for an RPC node:

This document will guide you through all the steps needed to setup your own `zkEVM-Node` for production.

# Warning:

>Currently the Executor/Prover does not run on ARM-powered Macs. For Windows users, WSL/WSL2 use is not recommended. 
> - Recommended specs: 
>    - Node: 16G RAM 4 cores
>    - Prover: 1TB RAM 128 cores
> - Unfortunately, M1 chips are not supported - for now since some optimizations on the prover require specific Intel instructions, this means some non-M1 computers won't work regardless of the OS, eg: AMD

## Network Components

Required:

- `Ethereum Node` - L1 Network
- `zkEVM-Node` - L2 Network
  - `JSON RPC Server` - Interface to L2 network
  - `Synchronizer` - Responsible to synchronize data between L1 and L2
  - `Sequencer` - Responsible to select transactions from the pool and propose new batches
  - `Aggregator`  - Responsible to consolidate the changes in the state proposed by the `Sequencers`
- `zk-Prover` - Zero knowledge proof generator

Optional:

- `Metamask` - Wallet to manage blockchain accounts
- `Block Scout Explorer` - Web UI to interact with the network information

## Requirements

- The examples on this document assume you have `docker-compose` installed, if you need help with the installation, please check the link below:
  - [docker-compose: Install](https://docs.docker.com/compose/install/)

## Recommendations

- It's recommended that you create a directory to add the files we are going to create during this document, we are going to refer to this directory as `zkevm-node` directory. To create this directory, run the following command:

```bash
mkdir -p /$HOME/zkevm-node
```

## Ethereum Node Setup

<details>
  <summary>Running your own Ethereum L1 network Geth node:</summary>
Let's go!

The first component we are going to setup is the Ethereum Node, it is the first because this is going to take a lot of time to synchronize the Ethereum network, so we will keep it synchronizing while we setup the others components to take advantage of this required time.

Before we start:

> There are many ways to setup an Ethereum L1 environment, we are going to use Geth for this.

We recommend you to use a dedicated machine to this component, this can be shared by multiple zkEVM-Node if you want to have more than one in your infrastructure.

First of all, we need to create a folder to store the Ethereum node data outside of the container, in order to not lose all the data if the container is restarted.

```bash
mkdir -p /$HOME/zkevm-node/.ethereum
```

In order to run the Ethereum node instance, create a file called `docker-compose.yml` inside of the directory `zkevm-node`

```yaml
version: '3'

services:

  eth-node:
    container_name: eth-node
    image: ethereum/client-go:stable
    ports:
        - 8545:8545
        - 8546:8546
        - 30303:30303
    volumes:
        - /$HOME/zkevm-node/.ethereum:/$HOME/geth/.ethereum
    command: [
        "--goerli",
        "--http",
        "--http.addr=0.0.0.0",
        "--http.corsdomain=*",
        "--http.vhosts=*",
        "--http.api=admin,eth,debug,miner,net,txpool,personal,web3",
        "--ws",
        "--ws.addr=0.0.0.0",
        "--ws.origins=*", 
        "--graphql", 
        "--graphql.corsdomain=*", 
        "--graphql.vhosts=*", 
        "--vmdebug", 
        "--metrics",
        "--datadir=/$HOME/geth/.ethereum"
    ]
```

To run the Ethereum node instance, go to the `zkevm-node` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

If you want to follow the logs of the synchronization, run the following command:

```bash
docker logs -f eth-node
```

</details>

---

We suggest using geth, but any Goerli node should work.

## Postgres Setup

Before we start:

> It's important to say that running the instances of Postgres in a docker container is just one way of running it. We strongly recommend you to have a specialized infrastructure to the DB like AWS RDS, a On-site server or any other Postgres DB dedicated infrastructure.

Also:

> It's not required to have a backup, since all the data is available on L1 to be resynchronized if it was lost, but it's strongly recommended to have a backup in order to avoid resynchronizing the whole network in case of a problem with the db, because the synchronization is a process that can take a lot of time and this time is going to ever increase as the network continues to roll.

With that said, we must setup several Postgres instances to be shared between the Node and the Prover/Executor.

- Node requires a full access user to run the migrations and control the data.
- Prover only needs a readonly user to access the Merkletree data and compute the proofs. Executor will need read/write access. Migration file `init_prover_db.sql` will create the merkle tree table in state DB.

We need to create several directories to store the Postgres data outside of the container, in order to not lose all the data if the container is restarted.

```bash
mkdir -p /$HOME/zkevm-node/.postgres-state
mkdir -p /$HOME/zkevm-node/.postgres-pool
mkdir -p /$HOME/zkevm-node/.postgres-rpc
```

Download the init schema for the prover DB: [./db/scripts/init_prover_db.sql](https://github.com/0xPolygonHermez/zkevm-node/blob/develop/db/scripts/init_prover_db.sql) to the directory `zkevm-node`.

In order to run the Postgres instance, create a file called `docker-compose.yml` inside of the directory `zkevm-node`

> We recommend you to customize the ENVIRONMENT variables values in the file below to your preference:

```yaml
version: '3'

services:
  zkevm-state-db:
    container_name: zkevm-state-db
    image: postgres
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
    ports:
      - 5432:5432
    volumes:
      - ./init_prover_db.sql:/docker-entrypoint-initdb.d/init.sql
      - /$HOME/zkevm-node/.postgres-state:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=state_user
      - POSTGRES_PASSWORD=state_password
      - POSTGRES_DB=state_db
    command: ["postgres", "-N", "500"]

  zkevm-pool-db:
    container_name: zkevm-pool-db
    image: postgres
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
    ports:
      - 5433:5432
    volumes:
      - /$HOME/zkevm-node/.postgres-pool:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=pool_user
      - POSTGRES_PASSWORD=pool_password
      - POSTGRES_DB=pool_db
    command: ["postgres", "-N", "500"]
```

To run the postgres instance, go to the `zkevm-node` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

Congratulations, your postgres instances are ready!

## Executor Setup

Before we start:

> It's very important to say that the Prover is a software that requires a lot of technology power to be executed properly, with that said, we recommend you to have a dedicated machine with the following configuration to run the prover:

- 128 CPU cores
- 1TB RAM

Also: 

> The prover depends on the Postgres instance we created before, so make sure it has network access to this.

The prover is available on Docker Registry, start by pulling the image:

```bash
docker pull hermeznetwork/zkevm-prover
```

Then download [the sample Prover config file](../config/environments/public/public.prover.config.json) (`./config/environments/public/public.prover.config.json`) and store it as `prover-config.json` inside the `zkevm-node` directory.

Finally, add the following entry to the `docker-compose.yml` file:

```yaml
  zkevm-prover:
    container_name: zkevm-prover
    image: hermeznetwork/zkevm-prover:develop
    ports:
      - 50061:50061 # MT
      - 50071:50071 # Executor
    volumes:
      - ./prover-config.json:/usr/src/app/config.json
    command: >
      zkProver -c /usr/src/app/config.json
```

This will spin up the Executor and MT, for a prover setup, exposing the `50051` port is needed:

```yaml
  zkevm-prover:
    container_name: zkevm-prover
    image: hermeznetwork/zkevm-prover:develop
    ports:
      - 50051:50051 # Prover
      - 50061:50061 # MT
      - 50071:50071 # Executor
    volumes:
      - ./prover-config.json:/usr/src/app/config.json
    command: >
      zkProver -c /usr/src/app/config.json
```

For more information visit [the Prover repository](https://github.com/0xPolygonHermez/zkevm-prover)

## zkEVM-Node Setup

Very well, we already have the Postgres, Prover and Ethereum Node instances running, now it's time so setup the zkEVM-Node.

> The node depends on the Postgres, Prover and Ethereum Node instances, so make sure it has network access to them. We also expect the node to have its own dedicated machine

Before we start, the node requires an Ethereum account with:

- Funds on L1 in order to propose new batches and consolidate the state
- Tokens to pay the collateral for batch proposal
- Approval of these tokens to be used by the roll-up SC on behalf of the Ethereum account owner
- Register this account as a sequencer

The node expected to read a `keystore` file, which is an encrypted file containing your credentials.
To create this file, go to the `zkevm-node` directory and run the following command:

> Remember to replace the `--pk` and `--pw` parameter values by the L1 account private key and the password you want to use to encrypt the file, the password will be required in the future to configure the node, so make sure you will remember it.

```bash
docker run --rm hermeznetwork/zkevm-node:latest sh -c "/app/zkevm-node encryptKey --pk=<account private key> --pw=<password to encrypt> --output=./keystore; cat ./keystore/*" > acc.keystore
```

The command above will create the file `acc.keystore` inside of the `zkevm-node` directory.

After it we need to create a configuration file to provide the configurations to the node, to achieve this create a file called `config.toml` inside of the `zkevm-node` directory, then go to the example [config file](../config/environments/public/public.node.config.toml) (`./config/environments/public/public.node.config.toml`) and `copy/paste` the content into the `config.toml` you'll actually use.

Do the same for the `genesis` file: [genesis file](../config/environments/public/public.genesis.config.json) (`./config/environments/public/public.genesis.config.json`)

Remember to:

- replace the database information if you set it differently while setting up the Postgres instance
- set the `Database Host` with the `Postgres instance IP`
- set the `Etherman URL` with the `JSON RPC URL` of the `Ethereum node` you created earlier *or* use any L1 Goerli service
- set the `Etherman Password` (`config.json` => `PrivateKeyPassword` field, defaults to `testonly`) to allow the node to decrypt the `keystore file`
- set the `MT / Executor URIs` the `IP and port` of the `MT/Executor Instances` and change the array of provers if a prover was spun up

Now we are going to put everything together in order to run the `zkEVM-Node` instance.

Add the following entries to the `docker-compose.yml` file

```yaml 
  zkevm-rpc:
    container_name: zkevm-rpc
    image: zkevm-node
    ports:
      - 8545:8545
    environment:
      - ZKEVM_NODE_STATEDB_HOST=zkevm-state-db
      - ZKEVM_NODE_POOL_HOST=zkevm-pool-db
      - ZKEVM_NODE_RPC_BROADCASTURI=public-grpc.zkevm-test.net:61090
    volumes:
      - ./acc.keystore:/pk/keystore
      - ./config.toml:/app/config.toml
      - ./genesis.json:/app/genesis.json
    command:
      - "/bin/sh"
      - "-c"
      - "/app/zkevm-node run --genesis /app/genesis.json --cfg /app/config.toml --components rpc"

  zkevm-sync:
    container_name: zkevm-sync
    image: zkevm-node
    environment:
      - ZKEVM_NODE_STATEDB_HOST=zkevm-state-db
    volumes:
      - ./acc.keystore:/pk/keystore
      - ./config.toml:/app/config.toml
      - ./genesis.json:/app/genesis.json
    command:
      - "/bin/sh"
      - "-c"
      - "/app/zkevm-node run --genesis /app/genesis.json --cfg /app/config.toml --components synchronizer"

```

To run the `zkEVM-Node` instance, go to the `zkevm-node` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

## Setup Explorer

To have a visual access to the network we are going to setup a Block Scout instance.

For more details about Block Scout, check it here: <https://docs.blockscout.com/>

Block Scout requires access to its own `zkEVM-Node` RPC-only instance to have access to the network
via the JSON RPC Server and a dedicated Postgres Instance in order to save its own data.

> Feel free to customize the environment variables to set the user, password and
> database for the Explore Postgres instance, but make sure to also update the url to connect
> to the DB in the Explorer environment variable called DATABASE_URL

```yaml 
version: '3'

services:

    zkevm-explorer-db:
        container_name: zkevm-explorer-db
        image: postgres
        ports:
            - 5432:5432
        environment:
            - POSTGRES_USER=test_user
            - POSTGRES_PASSWORD=test_password
            - POSTGRES_DB=explorer

    zkevm-explorer:
        container_name: zkevm-explorer
        image: hermeznetwork/hermez-node-blockscout:latest
        ports:
            - 4000:4000
        environment:
            - NETWORK=POE
            - SUBNETWORK=Polygon Hermez
            - COIN=ETH
            - ETHEREUM_JSONRPC_VARIANT=geth
            - ETHEREUM_JSONRPC_HTTP_URL=http://zkevm-explorer-zknode:8124
            - DATABASE_URL=postgres://test_user:test_password@zkevm-explorer-db:5432/explorer
            - ECTO_USE_SSL=false
            - MIX_ENV=prod
            - LOGO=/images/blockscout_logo.svg
            - LOGO_FOOTER=/images/blockscout_logo.svg
        command: ["/bin/sh", "-c", "mix do ecto.create, ecto.migrate; mix phx.server"]


    zkevm-explorer-zknode:
      container_name: zkevm-explorer-zknode
      image: zkevm-node
      ports:
        - 8124:8124
      environment:
        - ZKEVM_NODE_STATEDB_HOST=zkevm-state-db
        - ZKEVM_NODE_POOL_HOST=zkevm-pool-db
        - ZKEVM_NODE_RPC_PORT=8124
      volumes:
        - ./config/test.node.config.toml:/app/config.toml
        - ./config/test.genesis.config.json:/app/genesis.json
      command:
        - "/bin/sh"
        - "-c"
        - "/app/zkevm-node run --genesis /app/genesis.json --cfg /app/config.toml --components rpc --http.api eth,net,debug,zkevm,txpool,web3"
```

To run the Explorer, execute the following command:

```bash
docker-compose up -d zkevm-explorer-db
sleep 5
docker-compose up -d zkevm-explorer
sleep 5
docker-compose up -d zkevm-explorer-zknode
```

## Setup Metamask

To be able to use the Network via Metamask, a custom network must be configured.

> IMPORTANT: Metamask only allows custom networks to be added if the network is
> up and running, so make sure the whole environment is up and running before
> trying to add it as a custom network

To configure a custom network follow these steps:

1. Login to you Metamask account
2. Click in the circle with a picture on the top right side to open the Menu
3. Click on Settings
4. On the Left menu click com Networks
5. Fill up the following fields:
    1. Network Name: Polygon Hermez - Goerli
    2. New RPC URL: <http://IP-And-Port-of-zkEVM-Node-Instance>
    3. Chain ID: `1402`
    4. Currency Symbol: ETH
    5. Block Explorer URL: <http://IP-And-Port-of-Explorer-Instance>
6. Click on Save
7. Click on the X in the right top corner to close the Settings
8. Click in the list of networks on the top right corner
9. Select Polygon Hermez - Goerli
