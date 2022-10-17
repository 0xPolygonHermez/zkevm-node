> WARNING: This documentation is outdated, it will be updated soon

# Steps to run environment locally

## Overview

This documentation will help you running the following components:

- zkEVM Node Database
- Explorer Database
- L1 Network
- Prover
- zkEVM Node
- Explorer

## Requirements

The current version of the environment requires `go`, `docker` and `docker-compose` to be previously installed, check the links bellow to understand how to install them:

- <https://go.dev/doc/install>
- <https://www.docker.com/get-started>
- <https://docs.docker.com/compose/install/>

The `zkevm-node` docker image must be built at least once and every time a change is made to the code.
If you haven't build the `zkevm-node` image yet, you must run:

```bash
make build-docker
```

## Controlling the environment

> All the data is stored inside of each docker container, this means once you remove the container, the data will be lost.

To run the environment:

```bash
make run
```

To stop the environment:

```bash
make stop
```

To restart the environment:

```bash
make restart
```

## Sample data

The `make run` will execute the containers needed to run the environment but this will not execute anything else, so the L2 will be basically empty.

If you need sample data already deployed to the network, we have the following scripts:

First initialize the network for the L2 node:

```bash
make init-network
```

To add some examples of transactions and smart contracts:

```bash
make deploy-sc
```

To deploy a full a uniswap environment:

```bash
make deploy-uniswap
```

## Accessing the environment

- zkEVM Node Database 
  - `Type:` Postgres DB
  - `User:` test_user
  - `Password:` test_password
  - `Database:` test_db
  - `Host:` localhost
  - `Port:` 5432
  - `Url:` <postgres://test_user:test_password@localhost:5432/test_db>
- Explorer Database
  - `Type:` Postgres DB
  - `User:` test_user
  - `Password:` test_password
  - `Database:` explorer
  - `Host:` localhost
  - `Port:` 5433
  - `Url:` <postgres://test_user:test_password@localhost:5433/explorer>
- L1 Network
  - `Type:` Geth
  - `Host:` localhost
  - `Port:` 8545
  - `Url:` <http://localhost:8545>
- Prover
  - `Type:` Mock
  - `Host:` localhost
  - `Port:` 50001
  - `Url:` <http://localhost:50001>
- zkEVM Node
  - `Type:` JSON RPC
  - `Host:` localhost
  - `Port:` 8123
  - `Url:` <http://localhost:8123>
- Explorer
  - `Type:` Web
  - `Host:` localhost
  - `Port:` 4000
  - `Url:` <http://localhost:4000>

## Metamask

> Metamask requires the network to be running while configuring it, so make sure your network is running before starting.

To configure your Metamask to use your local environment, follow these steps:

1. Log in to your Metamask wallet
2. Click on your account picture and then on Settings
3. On the left menu, click on Networks
4. Click on `Add Network` button
5. Fill up the L2 network information
    1. `Network Name:` Polygon Hermez - Local
    2. `New RPC URL:` <http://localhost:8123>
    3. `ChainID:` 1000
    4. `Currency Symbol:` ETH
    5. `Block Explorer URL:` <http://localhost:4000>
6. Click on Save
7. Click on `Add Network` button
8. Fill up the L1 network information
    1. `Network Name:` Geth - Local
    2. `New RPC URL:` <http://localhost:8545>
    3. `ChainID:` 1337
    4. `Currency Symbol:` ETH
9. Click on Save

## L1 Addresses

| Address | Description |
|---|---|
| 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9 | Proof of Efficiency |
| 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 | Bridge |
| 0x5FbDB2315678afecb367f032d93F642f64180aa3 | Matic token |
| 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 | GlobalExitRootManager |

## Deployer Account

| Address | Private Key |
|---|---|
| 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 | 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 |

## Sequencer Account

| Address | Private Key |
|---|---|
| 0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D | 0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e |

## Accounts

| Address | Private Key |
|---|---|
| 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 | 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d |
| 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC | 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a |
| 0x90F79bf6EB2c4f870365E785982E1f101E93b906 | 0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6 |
| 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65 | 0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a |
| 0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc | 0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba |
| 0x976EA74026E726554dB657fA54763abd0C3a0aa9 | 0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e |
| 0x14dC79964da2C08b23698B3D3cc7Ca32193d9955 | 0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356 |
| 0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f | 0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97 |
| 0xa0Ee7A142d267C1f36714E4a8F75612F20a79720 | 0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6 |
| 0xBcd4042DE499D14e55001CcbB24a551F3b954096 | 0xf214f2b2cd398c806f84e317254e0f0b801d0643303237d97a22a48e01628897 |
| 0x71bE63f3384f5fb98995898A86B02Fb2426c5788 | 0x701b615bbdfb9de65240bc28bd21bbc0d996645a3dd57e7b12bc2bdf6f192c82 |
| 0xFABB0ac9d68B0B445fB7357272Ff202C5651694a | 0xa267530f49f8280200edf313ee7af6b827f2a8bce2897751d06a843f644967b1 |
| 0x1CBd3b2770909D4e10f157cABC84C7264073C9Ec | 0x47c99abed3324a2707c28affff1267e45918ec8c3f20b8aa892e8b065d2942dd |
| 0xdF3e18d64BC6A983f673Ab319CCaE4f1a57C7097 | 0xc526ee95bf44d8fc405a158bb884d9d1238d99f0612e9f33d006bb0789009aaa |
| 0xcd3B766CCDd6AE721141F452C550Ca635964ce71 | 0x8166f546bab6da521a8369cab06c5d2b9e46670292d85c875ee9ec20e84ffb61 |
| 0x2546BcD3c84621e976D8185a91A922aE77ECEc30 | 0xea6c44ac03bff858b476bba40716402b03e41b8e97e276d1baec7c37d42484a0 |
| 0xbDA5747bFD65F08deb54cb465eB87D40e51B197E | 0x689af8efa8c651a91ad287602527f3af2fe9f6501a7ac4b061667b5a93e037fd |
| 0xdD2FD4581271e230360230F9337D5c0430Bf44C0 | 0xde9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0 |
| 0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199 | 0xdf57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e |
