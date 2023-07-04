## Configuration

To configure a node you need 3 files: 
- Node configuration
- Genesis configuration
- Prover configuration

### Node Config
This file is a [TOML](https://en.wikipedia.org/wiki/TOML#) formatted file. 
You could find some examples here: 
 - `config/environments/local/local.node.config.toml`: running a permisionless node
  - `config/environments/mainnet/public.node.config.toml`
  - `config/environments/public/public.node.config.toml`
  - `test/config/test.node.config.toml`: configuration for a trusted node used in CI

  For details about the contents you can read specifications [here](config-file/node-config-doc.md)

This file is used for trusted and for permisionless nodes. In the case of permissionless node you only need to setup next sections: 



### Network Genesis Config
This contain all the info information relating to the genesis of the network (e.g. contracts, etc..)

You could find some examples here: 
- `config/environments/local/local.genesis.config.json`:

For details about the contents you can read specifications [here](config-file/custom_network-config-doc.md)


### Prover Config

Please check [prover repository](https://github.com/0xPolygonHermez/zkevm-prover)  for further information

Examples: 
 - `config/environments/mainnet/public.prover.config.json`
 - `config/environments/testnet/testnet.prover.config.json`
