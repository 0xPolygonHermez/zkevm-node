> WARNING: This documentation is outdated, it will be updated soon

# Production Setup

This document will guide you through all the steps needed to setup your own `zkEVM-Node` for production.

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

- All components have docker images available on docker hub, so it's important that you have an account to download and use them, please check the links below for more details:
  - [docker: Get-started](https://www.docker.com/get-started)
  - [docker hub](https://hub.docker.com)

Some of the images are still private, so make sure to login and check if you have access to the [Hermez organization](https://hub.docker.com/orgs/hermeznetwork) before trying to download them. Once you have docker installed on your machine, run the following command to login:

```bash
docker login
```

- The examples on this document assume you have `docker-compose` installed, if you need help with the installation, please check the link below:
  - [docker-compose: Install](https://docs.docker.com/compose/install/)

## Recommendations

- It's recommended that you create a directory to add the files we are going to create during this document, we are going to refer to this directory as `zkevm-node` directory. To create this directory, run the following command:

```bash
mkdir -p /$HOME/zkevm-node
```

## Ethereum Node Setup

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

```dockercompose
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

## Postgres Setup

Before we start:

> It's important to say that running the instance of Postgres in a docker container is just one way of running it. We strongly recommend you to have a specialized infrastructure to the DB like AWS RDS, a On-site server or any other Postgres DB dedicated infrastructure.

Also:

> It's not required to have a backup, since all the data is available on L1 to be resynchronized if it was lost, but it's strongly recommended to have a backup in order to avoid resynchronizing the whole network in case of a problem with the db, because the synchronization is a process that can take a lot of time and this time is going to ever increase as the network continues to roll.

With that said, we must setup a Postgres instance to be shared between the Node and the Prover.

- Node requires a full access user to run the migrations and control the data.
- Prover only needs a readonly user to access the Merkletree data and compute the proofs.

We need to create a folder to store the Postgres data outside of the container, in order to not lose all the data if the container is restarted.

```bash
mkdir -p /$HOME/zkevm-node/.postgres
```

In order to run the Postgres instance, create a file called `docker-compose.yml` inside of the directory `zkevm-node`

> We recommend you to customize the ENVIRONMENT variables values in the file below to your preference:

```docker-compose
version: '3'

services:

  zkevm-db:
    container_name: zkevm-db
    image: postgres
    ports:
      - 5432:5432
    environment:
      - POSTGRES_USER=test_user
      - POSTGRES_PASSWORD=test_password
      - POSTGRES_DB=test_db
    volumes:
      - /$HOME/zkevm-node/.postgres:./postgres-data
```

To run the postgres instance, go to the `zkevm-node` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

Congratulations, your postgres instance is ready!

## Prover Setup

Before we start:

> It's very important to say that the Prover is a software that requires a lot of technology power to be executed properly, with that said, we recommend you to have a dedicated machine with the following configuration to run the prover:

- TBD
- TBD
- TBD
- TBD

Also: 

> The prover depends on the Postgres instance we created before, so make sure it has network access to this.

- TDB how to setup de prover, docker, downloads, dependencies, etc

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

After it we need to create a configuration file to provide the configurations to the node, to achieve this create a file called `config.toml` inside of the `zkevm-node` directory, then go to the example [config file](config/config.debug.toml) and `copy/paste` the content into the `config.toml` you'll actually use.

Remember to:

- replace the database information if you set it differently while setting up the Postgres instance
- set the `Database Host` with the `Postgres instance IP`
- set the `Etherman URL` with the `JSON RPC URL` of the `Ethereum node`
- set the `Etherman Password` to allow the node to decrypt the `keystore file`
- set the `Prover URI` the `IP and port` of the `Prover Instance`



In order to be able to propose batches we are going to register our Ethereum account as a Sequencer,
to do this execute this command:

```bash
docker run --rm -v /$HOME/zkevm-node/config.toml:/app/config.toml hermeznetwork/zkevm-node:latest sh -c "./zkevm-node register --cfg=/app/config.toml --network=internaltestnet --y <public IP or URL for users to access the sequencer> "
```

In order to propose new batches, you must approve the Tokens to be used by the Roll-up on your behalf, to do this execute this command:
> remember to set the value of the parameter amount before executing

```bash
docker run --rm -v /$HOME/zkevm-node/config.toml:/app/config.toml hermeznetwork/zkevm-node:latest sh -c "./zkevm-node approve --cfg=/app/config.toml --network=internaltestnet --address=poe --amount=0 --y"
```

Now we are going to put everything together in order to run the `zkEVM-Node` instance.

Create a file called `docker-compose.yml` inside of the directory `zkevm-node`.

```docker-compose
version: '3'

services:
  
  zkevm-node:
    container_name: zkevm-node
    image: zkevm-node
    ports:
        - 8545:8545
    volumes:
      - /$HOME/zkevm-node/acc/keystore:/pk/keystore
      - /$HOME/zkevm-node/config.toml:/app/config.toml
    command: ["./zkevm-node", "run", "--network", "internaltestnet", "--cfg", "/app/config.toml"]
```

To run the `zkEVM-Node` instance, go to the `zkevm-node` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

## Setup Explorer

To have a visual access to the network we are going to setup a Block Scout instance.

For more details about Block Scout, check it here: <https://docs.blockscout.com/>

Block Scout requires access to the `zkEVM-Node` instance to have access to the network
via the JSON RPC Server and a dedicated Postgres Instance in order to save its own data.

We recommend you use a dedicated machine for the Explorer.

Create a file called `docker-compose.yml` inside of the `zkevm-node` directory.

> Feel free to customize the environment variables to set the user, password and
> database for the Explore Postgres instance, but make sure to also update the url to connect
> to the DB in the Explorer environment variable called DATABASE_URL
> Remember to set the environment variable ETHEREUM_JSONRPC_HTTP_URL with the `zkEVM-Node` IP and PORT

```docker-compose
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
            - ETHEREUM_JSONRPC_HTTP_URL=http://:8545 # Set the IP and PORT of the zkEVM-Node
            - DATABASE_URL=postgres://test_user:test_password@zkevm-explorer-db:5432/explorer
            - ECTO_USE_SSL=false
            - MIX_ENV=prod
            - LOGO=/images/blockscout_logo.svg
            - LOGO_FOOTER=/images/blockscout_logo.svg
        command: ["/bin/sh", "-c", "mix do ecto.create, ecto.migrate; mix phx.server"]
```

To run the Explorer, execute the following command:

```bash
docker-compose up -d zkevm-explorer-db
sleep5
docker-compose up -d zkevm-explorer
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
    3. Chain ID: TBD
    4. Currency Symbol: ETH
    5. Block Explorer URL: <http://IP-And-Port-of-Explorer-Instance>
6. Click on Save
7. Click on the X in the right top corner to close the Settings
8. Click in the list of networks on the top right corner
9. Select Polygon Hermez - Goerli
