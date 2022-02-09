# Production Setup

This document will guide you through all the steps needed to setup your own Hermez Node for production.

---

## Network Components

Required:

- Ethereum Node - Interface to L1
- Hermez zk-EVM-Node - Interface to L2
  - Sequencer - Responsible to select transactions from the pool and propose new batches
  - Aggregator - Responsible to consolidate the changes in the state proposed by the Sequencers
- Hermez zk-Prover - Zero knowledge proof generator

Optional:

- Metamask - Wallet to manage your accounts
- Block Scout Explorer - Web UI to interact with the network information

---

## Requirements

- docker: All components have images available to be used, check this <https://www.docker.com/get-started>

---

## Recommendations

- We strongly recommend using docker-compose to setup the services, check this <https://docs.docker.com/compose/install/>
- It's also recommended that you create a directory to add the files we are going to create during this document, we are going to refer to this directory as `hermez` directory. To create this directory, run the following command:

```bash
mkdir -p /$HOME/hermez
```
---

## Ethereum Node Setup

Let's go!

The first component we are going to setup is the Ethereum Node, it is the first because this is going to take a lot of time to synchronize the Ethereum network, so we will keep it synchronizing while we setup the others components to take advantage of this required time.

Before we start:

> There are may ways to setup an Ethereum L1 environment, we are going to use Geth for this.

We recommend you to use a dedicated machine to this component, this can be shared by multiple Hermez zk-EVM-Node if you want to have more than one in your infrastructure.

First of all, we need to create a folder to store the Ethereum node data outside of the container, in order to not lose all the data if the container is restarted.

```bash
mkdir -p /$HOME/hermez/.ethereum
```

In order to run the Ethereum node instance, create a file called `docker-compose.yml` inside of the directory `hermez`

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
        - /$HOME/hermez/.ethereum:/$HOME/geth/.ethereum
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

To run the Ethereum node instance, go to the `hermez` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

If you want to follow the logs of the synchronization, run the following command:

```bash
docker logs -f eth-node
```

---

## Postgres Setup

Before we start:

> It's important to say that running the instance of Postgres in a docker container is just one way of running it. We strongly recommend you to have a specialized infrastructure to the DB like AWS RDS, a On-site server or any other Postgres DB dedicated infrastructure.

Also:

> It's not required to have a backup, since all the data is available on L1 to be resynchronized if it was lost, but it's strongly recommended to have a backup in order to avoid resynchronizing the whole network in case of a problem with the db, because the synchronization is a process that can take a lot of time and this time is going to ever increase as the network continues to roll.

With that said, we must setup a Postgres instance to be shared between the Node and the Prover.

- Node requires a full access user to run the migrations and control the data.
- Prover only needs a readonly user to access the historical data and compute the proofs.

In order to run the Postgres instance, create a file called `docker-compose.yml` inside of the directory `hermez`

> We recommend you to customize the ENVIRONMENT variables values in the file below to your preference:

```dockercompose
version: '3'

services:

  hez-postgres:
    container_name: hez-postgres
      image: postgres
      ports:
        - 5432:5432
      environment:
        - POSTGRES_USER=test_user
        - POSTGRES_PASSWORD=test_password
        - POSTGRES_DB=test_db
```

To run the postgres instance, go to the `hermez` directory in your terminal and run the following command:

```bash
docker-compose up -d
```

Congratulations, your postgres instance is ready!

---

## Prove Setup

Before we start:

> It's very important to say that the Prover is a software that requires a lot of technology power to be executed properly, with that said, we recommend you to have a dedicated machine with the following configuration to run the prover:

- TBD
- TBD
- TBD
- TBD

Also: 

> The prover depends on the Postgres instance we created before, so make sure it has network access to this.

- TDB how to setup de prover, docker, downloads, dependencies, etc

---

## Hermez zk-EVM-Node Setup

Very well, we already have the Postgres, Prover and Ethereum Node instances running, now it's time so setup the Hermez zk-EVM-Node.

> The node depends on the Postgres and Prover instances, so make sure it has network access to them.
> We also expect you to run the node in a dedicated machine

Before we start:

> The node requires the L1 account credentials in order to send L1 transactions to propose new batches and consolidate the state.

The node expected to read a `keystore` file, which is an encrypted file containing your credentials.
To create this file, go to the `hermez` directory and run the following command:

> Remember to replace the `--pk` and `--pw` parameter values by the L1 account private key and the password you want to use to encrypt the file, the password will be required in the future to configure the node, so make sure you will remember it.

```bash
docker run --rm hermeznetwork/hermez-node-zkevm:latest sh -c "./hezcore encryptKey --pk=<account private key> --pw=<password to encrypt> --output=./keystore; cat ./keystore/*" > acc.keystore
```

The command above will create the file `acc.keystore` inside of the `hermez` directory.

After it we need to create a configuration file to provide the configurations to the node, to achieve this create a file called `config.toml` inside of the `hermez` directory with this:

> - Remember to replace the database information if you set it differently while setting up the Postgres instance
> - Remember to set the Database Host with the Postgres instance IP
> - Remember to set the Etherman URL with the JSON RPC URL of the Ethereum node, which is "http://\<Ethereum Node Instance IP>:\<PORT>"
> - Remember to set the Etherman Password to allow the node to decrypt the keystore file
> - Remember to set the Prover URI the IP and port of the Prover Instance like this "\<IP>:\<PORT>"

```toml
[Log]
Level = "info"
Outputs = ["stdout"]

[Database]
User = "test_user"
Password = "test_password"
Name = "test_db"
Host = 
Port = "5432"

[Etherman]
URL = 
PrivateKeyPath = "/pk/keystore"
PrivateKeyPassword = 

[RPC]
Host = "0.0.0.0"
Port = 8123

[Synchronizer]
SyncInterval = "5s"
SyncChunkSize = 100

[Sequencer]
AllowNonRegistered = "false"
IntervalToProposeBatch = "15s"
SyncedBlockDif = 1
    [Sequencer.Strategy]
        [Sequencer.Strategy.TxSelector]
            Type = "acceptall"
            TxSorterType = "bycostandnonce"
        [Sequencer.Strategy.TxProfitabilityChecker]
            Type = "acceptall"
            MinReward = "1.1"

[Aggregator]
IntervalToConsolidateState = "10s"
TxProfitabilityCheckerType = "acceptall"
TxProfitabilityMinReward = "1.1"

[Prover]
ProverURI = 

```

Now we are going to put everything together in order to run the Hermez zk-EVM-Node instance.

Create a file called `docker-compose.yml` inside of the directory `hermez`.

```docker-compose
version: '3'

services:
  
  hez-core:
    container_name: hez-core
    image: hezcore
    ports:
        - 8123:8123
    volumes:
      - /$HOME/hermez/acc/keystore:/pk/keystore
      - /$HOME/hermez/config.toml:/app/config.toml
    command: ["./hezcore", "run", "--network", "internaltestnet", "--cfg", "/app/config.toml"]
```

To run the Hermez zk-EVM-Node instance, go to the `hermez` directory in your terminal and run the following command:

```bash
docker-compose up -d
```