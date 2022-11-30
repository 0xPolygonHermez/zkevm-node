## Component: Sequencer

### ZKEVM Sequencer:

The ZKEVM Sequencer is an optional but ancillary module that proposes new batches using transactions stored in the Pool Database.

### Running:

The preferred way to run the ZKEVM Sequencer component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/zkevm-node
```

To orchestrate multiple deployments of the different ZKEVM Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  zkevm-sequencer:
    container_name: zkevm-sequencer
    image: zkevm-node
	command:
		- "/bin/sh"
		- "-c"
		- "/app/zkevm-node run --genesis /app/genesis.json --cfg /app/config.toml --components sequencer"
```

The container alone needs some parameters configured, access to certain configuration files and the appropiate ports exposed.

- environment: Env variables that supersede the config file
	- `ZKEVM_NODE_POOLDB_HOST`: Name of PoolDB Database Host
	- `ZKEVM_NODE_STATEDB_HOST`: Name of StateDB Database Host
- volumes:
	- [your Account Keystore file]:/pk/keystore (note, this `/pk/keystore` value is the default path that's written in the Public Configuration files on this repo, meant to expedite deployments, it can be superseded via an env flag `ZKEVM_NODE_ETHERMAN_PRIVATEKEYPATH`.)
	- [your config.toml file]:/app/config.toml
	- [your genesis.json file]:/app/genesis.json

### Generating an Account Keystore file:

```bash
docker run --rm hermeznetwork/zkevm-node:latest sh -c "/app/zkevm-node encryptKey --pk=[your private key] --pw=[password to encrypt file] --output=./keystore; cat ./keystore/*" > account.keystore
```

**NOTE**:

- Replace `[your private key]` with your Ethereum L1 account private key
- Replace `[password to encrypt file]` with a password used for file encryption. This password must be passed to the Node later on via env variable (`ZKEVM_NODE_ETHERMAN_PRIVATEKEYPASSWORD`)
- The resulting pipe file must match in name to what's on the docker-compose file (which will bind the host file to `/pk/keystore`)