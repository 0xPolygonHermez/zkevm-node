# Schema Docs

**Type:** : `object`
**Description:** Config represents the configuration of the entire Hermez Node The file is TOML format You could find some examples:

[TOML format]: https://en.wikipedia.org/wiki/TOML

| Property                                             | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| ---------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [IsTrustedSequencer](#IsTrustedSequencer )         | No      | boolean | No         | -          | This define is a trusted node (\`true\`) or a permission less (\`false\`). If you don't known<br />set to \`false\`                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| - [ForkUpgradeBatchNumber](#ForkUpgradeBatchNumber ) | No      | integer | No         | -          | Last batch number before  a forkid change (fork upgrade). That implies that<br />greater batch numbers are going to be trusted but no virtualized neither verified.<br />So after the batch number \`ForkUpgradeBatchNumber\` is virtualized and verified you could update<br />the system (SC,...) to new forkId and remove this value to allow the system to keep<br />Virtualizing and verifying the new batchs.<br />Check issue [#2236](https://github.com/0xPolygonHermez/zkevm-node/issues/2236) to known more<br />This value overwrite \`SequenceSender.ForkUpgradeBatchNumber\` |
| - [ForkUpgradeNewForkId](#ForkUpgradeNewForkId )     | No      | integer | No         | -          | Which is the new forkId                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| - [Log](#Log )                                       | No      | object  | No         | -          | Configure Log level for all the services, allow also to store the logs in a file                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| - [Etherman](#Etherman )                             | No      | object  | No         | -          | Configuration of the etherman (client for access L1)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [EthTxManager](#EthTxManager )                     | No      | object  | No         | -          | Configuration for ethereum transaction manager                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| - [Pool](#Pool )                                     | No      | object  | No         | -          | Pool service configuration                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| - [RPC](#RPC )                                       | No      | object  | No         | -          | Configuration for RPC service. THis one offers a extended Ethereum JSON-RPC API interface to interact with the node                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| - [Synchronizer](#Synchronizer )                     | No      | object  | No         | -          | Configuration of service \`Syncrhonizer\`. For this service is also really important the value of \`IsTrustedSequencer\`<br />because depending of this values is going to ask to a trusted node for trusted transactions or not                                                                                                                                                                                                                                                                                                                                                          |
| - [Sequencer](#Sequencer )                           | No      | object  | No         | -          | Configuration of the sequencer service                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| - [SequenceSender](#SequenceSender )                 | No      | object  | No         | -          | Configuration of the sequence sender service                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| - [Aggregator](#Aggregator )                         | No      | object  | No         | -          | Configuration of the aggregator service                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| - [NetworkConfig](#NetworkConfig )                   | No      | object  | No         | -          | Configuration of the genesis of the network. This is used to known the initial state of the network                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| - [L2GasPriceSuggester](#L2GasPriceSuggester )       | No      | object  | No         | -          | Configuration of the gas price suggester service                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| - [Executor](#Executor )                             | No      | object  | No         | -          | Configuration of the executor service                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| - [MTClient](#MTClient )                             | No      | object  | No         | -          | Configuration of the merkle tree client service. Not use in the node, only for testing                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| - [Metrics](#Metrics )                               | No      | object  | No         | -          | Configuration of the metrics service, basically is where is going to publish the metrics                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| - [EventLog](#EventLog )                             | No      | object  | No         | -          | Configuration of the event database connection                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| - [HashDB](#HashDB )                                 | No      | object  | No         | -          | Configuration of the hash database connection                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| - [State](#State )                                   | No      | object  | No         | -          | State service configuration                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |

## <a name="IsTrustedSequencer"></a>1. `IsTrustedSequencer`

**Type:** : `boolean`

**Default:** `false`

**Description:** This define is a trusted node (`true`) or a permission less (`false`). If you don't known
set to `false`

**Example setting the default value** (false):
```
IsTrustedSequencer=false
```

## <a name="ForkUpgradeBatchNumber"></a>2. `ForkUpgradeBatchNumber`

**Type:** : `integer`

**Default:** `0`

**Description:** Last batch number before  a forkid change (fork upgrade). That implies that
greater batch numbers are going to be trusted but no virtualized neither verified.
So after the batch number `ForkUpgradeBatchNumber` is virtualized and verified you could update
the system (SC,...) to new forkId and remove this value to allow the system to keep
Virtualizing and verifying the new batchs.
Check issue [#2236](https://github.com/0xPolygonHermez/zkevm-node/issues/2236) to known more
This value overwrite `SequenceSender.ForkUpgradeBatchNumber`

**Example setting the default value** (0):
```
ForkUpgradeBatchNumber=0
```

## <a name="ForkUpgradeNewForkId"></a>3. `ForkUpgradeNewForkId`

**Type:** : `integer`

**Default:** `0`

**Description:** Which is the new forkId

**Example setting the default value** (0):
```
ForkUpgradeNewForkId=0
```

## <a name="Log"></a>4. `[Log]`

**Type:** : `object`
**Description:** Configure Log level for all the services, allow also to store the logs in a file

| Property                           | Pattern | Type             | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                                                                                               |
| ---------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [Environment](#Log_Environment ) | No      | enum (of string) | No         | -          | Environment defining the log format ("production" or "development").<br />In development mode enables development mode (which makes DPanicLevel logs panic), uses a console encoder, writes to standard error, and disables sampling. Stacktraces are automatically included on logs of WarnLevel and above.<br />Check [here](https://pkg.go.dev/go.uber.org/zap@v1.24.0#NewDevelopmentConfig) |
| - [Level](#Log_Level )             | No      | enum (of string) | No         | -          | Level of log. As lower value more logs are going to be generated                                                                                                                                                                                                                                                                                                                                |
| - [Outputs](#Log_Outputs )         | No      | array of string  | No         | -          | Outputs                                                                                                                                                                                                                                                                                                                                                                                         |

### <a name="Log_Environment"></a>4.1. `Log.Environment`

**Type:** : `enum (of string)`

**Default:** `"development"`

**Description:** Environment defining the log format ("production" or "development").
In development mode enables development mode (which makes DPanicLevel logs panic), uses a console encoder, writes to standard error, and disables sampling. Stacktraces are automatically included on logs of WarnLevel and above.
Check [here](https://pkg.go.dev/go.uber.org/zap@v1.24.0#NewDevelopmentConfig)

**Example setting the default value** ("development"):
```
[Log]
Environment="development"
```

Must be one of:
* "production"
* "development"

### <a name="Log_Level"></a>4.2. `Log.Level`

**Type:** : `enum (of string)`

**Default:** `"info"`

**Description:** Level of log. As lower value more logs are going to be generated

**Example setting the default value** ("info"):
```
[Log]
Level="info"
```

Must be one of:
* "debug"
* "info"
* "warn"
* "error"
* "dpanic"
* "panic"
* "fatal"

### <a name="Log_Outputs"></a>4.3. `Log.Outputs`

**Type:** : `array of string`

**Default:** `["stderr"]`

**Description:** Outputs

**Example setting the default value** (["stderr"]):
```
[Log]
Outputs=["stderr"]
```

## <a name="Etherman"></a>5. `[Etherman]`

**Type:** : `object`
**Description:** Configuration of the etherman (client for access L1)

| Property                                          | Pattern | Type    | Deprecated | Definition | Title/Description                                                                       |
| ------------------------------------------------- | ------- | ------- | ---------- | ---------- | --------------------------------------------------------------------------------------- |
| - [URL](#Etherman_URL )                           | No      | string  | No         | -          | URL is the URL of the Ethereum node for L1                                              |
| - [ForkIDChunkSize](#Etherman_ForkIDChunkSize )   | No      | integer | No         | -          | ForkIDChunkSize is the max interval for each call to L1 provider to get the forkIDs     |
| - [MultiGasProvider](#Etherman_MultiGasProvider ) | No      | boolean | No         | -          | allow that L1 gas price calculation use multiples sources                               |
| - [Etherscan](#Etherman_Etherscan )               | No      | object  | No         | -          | Configuration for use Etherscan as used as gas provider, basically it needs the API-KEY |

### <a name="Etherman_URL"></a>5.1. `Etherman.URL`

**Type:** : `string`

**Default:** `"http://localhost:8545"`

**Description:** URL is the URL of the Ethereum node for L1

**Example setting the default value** ("http://localhost:8545"):
```
[Etherman]
URL="http://localhost:8545"
```

### <a name="Etherman_ForkIDChunkSize"></a>5.2. `Etherman.ForkIDChunkSize`

**Type:** : `integer`

**Default:** `20000`

**Description:** ForkIDChunkSize is the max interval for each call to L1 provider to get the forkIDs

**Example setting the default value** (20000):
```
[Etherman]
ForkIDChunkSize=20000
```

### <a name="Etherman_MultiGasProvider"></a>5.3. `Etherman.MultiGasProvider`

**Type:** : `boolean`

**Default:** `false`

**Description:** allow that L1 gas price calculation use multiples sources

**Example setting the default value** (false):
```
[Etherman]
MultiGasProvider=false
```

### <a name="Etherman_Etherscan"></a>5.4. `[Etherman.Etherscan]`

**Type:** : `object`
**Description:** Configuration for use Etherscan as used as gas provider, basically it needs the API-KEY

| Property                                | Pattern | Type   | Deprecated | Definition | Title/Description                                                                                                                     |
| --------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| - [ApiKey](#Etherman_Etherscan_ApiKey ) | No      | string | No         | -          | Need API key to use etherscan, if it's empty etherscan is not used                                                                    |
| - [Url](#Etherman_Etherscan_Url )       | No      | string | No         | -          | URL of the etherscan API. Overwritten with a hardcoded URL: "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=" |

#### <a name="Etherman_Etherscan_ApiKey"></a>5.4.1. `Etherman.Etherscan.ApiKey`

**Type:** : `string`

**Default:** `""`

**Description:** Need API key to use etherscan, if it's empty etherscan is not used

**Example setting the default value** (""):
```
[Etherman.Etherscan]
ApiKey=""
```

#### <a name="Etherman_Etherscan_Url"></a>5.4.2. `Etherman.Etherscan.Url`

**Type:** : `string`

**Default:** `""`

**Description:** URL of the etherscan API. Overwritten with a hardcoded URL: "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey="

**Example setting the default value** (""):
```
[Etherman.Etherscan]
Url=""
```

## <a name="EthTxManager"></a>6. `[EthTxManager]`

**Type:** : `object`
**Description:** Configuration for ethereum transaction manager

| Property                                                        | Pattern | Type            | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| --------------------------------------------------------------- | ------- | --------------- | ---------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [FrequencyToMonitorTxs](#EthTxManager_FrequencyToMonitorTxs ) | No      | string          | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| - [WaitTxToBeMined](#EthTxManager_WaitTxToBeMined )             | No      | string          | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| - [PrivateKeys](#EthTxManager_PrivateKeys )                     | No      | array of object | No         | -          | PrivateKeys defines all the key store files that are going<br />to be read in order to provide the private keys to sign the L1 txs                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| - [ForcedGas](#EthTxManager_ForcedGas )                         | No      | integer         | No         | -          | ForcedGas is the amount of gas to be forced in case of gas estimation error                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| - [GasPriceMarginFactor](#EthTxManager_GasPriceMarginFactor )   | No      | number          | No         | -          | GasPriceMarginFactor is used to multiply the suggested gas price provided by the network<br />in order to allow a different gas price to be set for all the transactions and making it<br />easier to have the txs prioritized in the pool, default value is 1.<br /><br />ex:<br />suggested gas price: 100<br />GasPriceMarginFactor: 1<br />gas price = 100<br /><br />suggested gas price: 100<br />GasPriceMarginFactor: 1.1<br />gas price = 110                                                                                                                                                                                              |
| - [MaxGasPriceLimit](#EthTxManager_MaxGasPriceLimit )           | No      | integer         | No         | -          | MaxGasPriceLimit helps avoiding transactions to be sent over an specified<br />gas price amount, default value is 0, which means no limit.<br />If the gas price provided by the network and adjusted by the GasPriceMarginFactor<br />is greater than this configuration, transaction will have its gas price set to<br />the value configured in this config as the limit.<br /><br />ex:<br /><br />suggested gas price: 100<br />gas price margin factor: 20%<br />max gas price limit: 150<br />tx gas price = 120<br /><br />suggested gas price: 100<br />gas price margin factor: 20%<br />max gas price limit: 110<br />tx gas price = 110 |

### <a name="EthTxManager_FrequencyToMonitorTxs"></a>6.1. `EthTxManager.FrequencyToMonitorTxs`

**Title:** Duration

**Type:** : `string`

**Default:** `"1s"`

**Description:** FrequencyToMonitorTxs frequency of the resending failed txs

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1s"):
```
[EthTxManager]
FrequencyToMonitorTxs="1s"
```

### <a name="EthTxManager_WaitTxToBeMined"></a>6.2. `EthTxManager.WaitTxToBeMined`

**Title:** Duration

**Type:** : `string`

**Default:** `"2m0s"`

**Description:** WaitTxToBeMined time to wait after transaction was sent to the ethereum

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("2m0s"):
```
[EthTxManager]
WaitTxToBeMined="2m0s"
```

### <a name="EthTxManager_PrivateKeys"></a>6.3. `EthTxManager.PrivateKeys`

**Type:** : `array of object`
**Description:** PrivateKeys defines all the key store files that are going
to be read in order to provide the private keys to sign the L1 txs

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                      | Description                                                                                   |
| ---------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [PrivateKeys items](#EthTxManager_PrivateKeys_items) | KeystoreFileConfig has all the information needed to load a private key from a key store file |

#### <a name="autogenerated_heading_2"></a>6.3.1. [EthTxManager.PrivateKeys.PrivateKeys items]

**Type:** : `object`
**Description:** KeystoreFileConfig has all the information needed to load a private key from a key store file

| Property                                                | Pattern | Type   | Deprecated | Definition | Title/Description                                      |
| ------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------ |
| - [Path](#EthTxManager_PrivateKeys_items_Path )         | No      | string | No         | -          | Path is the file path for the key store file           |
| - [Password](#EthTxManager_PrivateKeys_items_Password ) | No      | string | No         | -          | Password is the password to decrypt the key store file |

##### <a name="EthTxManager_PrivateKeys_items_Path"></a>6.3.1.1. `EthTxManager.PrivateKeys.PrivateKeys items.Path`

**Type:** : `string`
**Description:** Path is the file path for the key store file

##### <a name="EthTxManager_PrivateKeys_items_Password"></a>6.3.1.2. `EthTxManager.PrivateKeys.PrivateKeys items.Password`

**Type:** : `string`
**Description:** Password is the password to decrypt the key store file

### <a name="EthTxManager_ForcedGas"></a>6.4. `EthTxManager.ForcedGas`

**Type:** : `integer`

**Default:** `0`

**Description:** ForcedGas is the amount of gas to be forced in case of gas estimation error

**Example setting the default value** (0):
```
[EthTxManager]
ForcedGas=0
```

### <a name="EthTxManager_GasPriceMarginFactor"></a>6.5. `EthTxManager.GasPriceMarginFactor`

**Type:** : `number`

**Default:** `1`

**Description:** GasPriceMarginFactor is used to multiply the suggested gas price provided by the network
in order to allow a different gas price to be set for all the transactions and making it
easier to have the txs prioritized in the pool, default value is 1.

ex:
suggested gas price: 100
GasPriceMarginFactor: 1
gas price = 100

suggested gas price: 100
GasPriceMarginFactor: 1.1
gas price = 110

**Example setting the default value** (1):
```
[EthTxManager]
GasPriceMarginFactor=1
```

### <a name="EthTxManager_MaxGasPriceLimit"></a>6.6. `EthTxManager.MaxGasPriceLimit`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxGasPriceLimit helps avoiding transactions to be sent over an specified
gas price amount, default value is 0, which means no limit.
If the gas price provided by the network and adjusted by the GasPriceMarginFactor
is greater than this configuration, transaction will have its gas price set to
the value configured in this config as the limit.

ex:

suggested gas price: 100
gas price margin factor: 20%
max gas price limit: 150
tx gas price = 120

suggested gas price: 100
gas price margin factor: 20%
max gas price limit: 110
tx gas price = 110

**Example setting the default value** (0):
```
[EthTxManager]
MaxGasPriceLimit=0
```

## <a name="Pool"></a>7. `[Pool]`

**Type:** : `object`
**Description:** Pool service configuration

| Property                                                                        | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                    |
| ------------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------------------------------------------------- |
| - [IntervalToRefreshBlockedAddresses](#Pool_IntervalToRefreshBlockedAddresses ) | No      | string  | No         | -          | Duration                                                                                             |
| - [IntervalToRefreshGasPrices](#Pool_IntervalToRefreshGasPrices )               | No      | string  | No         | -          | Duration                                                                                             |
| - [MaxTxBytesSize](#Pool_MaxTxBytesSize )                                       | No      | integer | No         | -          | MaxTxBytesSize is the max size of a transaction in bytes                                             |
| - [MaxTxDataBytesSize](#Pool_MaxTxDataBytesSize )                               | No      | integer | No         | -          | MaxTxDataBytesSize is the max size of the data field of a transaction in bytes                       |
| - [DB](#Pool_DB )                                                               | No      | object  | No         | -          | DB is the database configuration                                                                     |
| - [DefaultMinGasPriceAllowed](#Pool_DefaultMinGasPriceAllowed )                 | No      | integer | No         | -          | DefaultMinGasPriceAllowed is the default min gas price to suggest                                    |
| - [MinAllowedGasPriceInterval](#Pool_MinAllowedGasPriceInterval )               | No      | string  | No         | -          | Duration                                                                                             |
| - [PollMinAllowedGasPriceInterval](#Pool_PollMinAllowedGasPriceInterval )       | No      | string  | No         | -          | Duration                                                                                             |
| - [AccountQueue](#Pool_AccountQueue )                                           | No      | integer | No         | -          | AccountQueue represents the maximum number of non-executable transaction slots permitted per account |
| - [GlobalQueue](#Pool_GlobalQueue )                                             | No      | integer | No         | -          | GlobalQueue represents the maximum number of non-executable transaction slots for all accounts       |
| - [EffectiveGasPrice](#Pool_EffectiveGasPrice )                                 | No      | object  | No         | -          | EffectiveGasPrice is the config for the effective gas price calculation                              |

### <a name="Pool_IntervalToRefreshBlockedAddresses"></a>7.1. `Pool.IntervalToRefreshBlockedAddresses`

**Title:** Duration

**Type:** : `string`

**Default:** `"5m0s"`

**Description:** IntervalToRefreshBlockedAddresses is the time it takes to sync the
blocked address list from db to memory

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5m0s"):
```
[Pool]
IntervalToRefreshBlockedAddresses="5m0s"
```

### <a name="Pool_IntervalToRefreshGasPrices"></a>7.2. `Pool.IntervalToRefreshGasPrices`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** IntervalToRefreshGasPrices is the time to wait to refresh the gas prices

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Pool]
IntervalToRefreshGasPrices="5s"
```

### <a name="Pool_MaxTxBytesSize"></a>7.3. `Pool.MaxTxBytesSize`

**Type:** : `integer`

**Default:** `100132`

**Description:** MaxTxBytesSize is the max size of a transaction in bytes

**Example setting the default value** (100132):
```
[Pool]
MaxTxBytesSize=100132
```

### <a name="Pool_MaxTxDataBytesSize"></a>7.4. `Pool.MaxTxDataBytesSize`

**Type:** : `integer`

**Default:** `100000`

**Description:** MaxTxDataBytesSize is the max size of the data field of a transaction in bytes

**Example setting the default value** (100000):
```
[Pool]
MaxTxDataBytesSize=100000
```

### <a name="Pool_DB"></a>7.5. `[Pool.DB]`

**Type:** : `object`
**Description:** DB is the database configuration

| Property                           | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| ---------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#Pool_DB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#Pool_DB_User )           | No      | string  | No         | -          | Database User name                                         |
| - [Password](#Pool_DB_Password )   | No      | string  | No         | -          | Database Password of the user                              |
| - [Host](#Pool_DB_Host )           | No      | string  | No         | -          | Host address of database                                   |
| - [Port](#Pool_DB_Port )           | No      | string  | No         | -          | Port Number of database                                    |
| - [EnableLog](#Pool_DB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#Pool_DB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

#### <a name="Pool_DB_Name"></a>7.5.1. `Pool.DB.Name`

**Type:** : `string`

**Default:** `"pool_db"`

**Description:** Database name

**Example setting the default value** ("pool_db"):
```
[Pool.DB]
Name="pool_db"
```

#### <a name="Pool_DB_User"></a>7.5.2. `Pool.DB.User`

**Type:** : `string`

**Default:** `"pool_user"`

**Description:** Database User name

**Example setting the default value** ("pool_user"):
```
[Pool.DB]
User="pool_user"
```

#### <a name="Pool_DB_Password"></a>7.5.3. `Pool.DB.Password`

**Type:** : `string`

**Default:** `"pool_password"`

**Description:** Database Password of the user

**Example setting the default value** ("pool_password"):
```
[Pool.DB]
Password="pool_password"
```

#### <a name="Pool_DB_Host"></a>7.5.4. `Pool.DB.Host`

**Type:** : `string`

**Default:** `"zkevm-pool-db"`

**Description:** Host address of database

**Example setting the default value** ("zkevm-pool-db"):
```
[Pool.DB]
Host="zkevm-pool-db"
```

#### <a name="Pool_DB_Port"></a>7.5.5. `Pool.DB.Port`

**Type:** : `string`

**Default:** `"5432"`

**Description:** Port Number of database

**Example setting the default value** ("5432"):
```
[Pool.DB]
Port="5432"
```

#### <a name="Pool_DB_EnableLog"></a>7.5.6. `Pool.DB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[Pool.DB]
EnableLog=false
```

#### <a name="Pool_DB_MaxConns"></a>7.5.7. `Pool.DB.MaxConns`

**Type:** : `integer`

**Default:** `200`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (200):
```
[Pool.DB]
MaxConns=200
```

### <a name="Pool_DefaultMinGasPriceAllowed"></a>7.6. `Pool.DefaultMinGasPriceAllowed`

**Type:** : `integer`

**Default:** `1000000000`

**Description:** DefaultMinGasPriceAllowed is the default min gas price to suggest

**Example setting the default value** (1000000000):
```
[Pool]
DefaultMinGasPriceAllowed=1000000000
```

### <a name="Pool_MinAllowedGasPriceInterval"></a>7.7. `Pool.MinAllowedGasPriceInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"5m0s"`

**Description:** MinAllowedGasPriceInterval is the interval to look back of the suggested min gas price for a tx

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5m0s"):
```
[Pool]
MinAllowedGasPriceInterval="5m0s"
```

### <a name="Pool_PollMinAllowedGasPriceInterval"></a>7.8. `Pool.PollMinAllowedGasPriceInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"15s"`

**Description:** PollMinAllowedGasPriceInterval is the interval to poll the suggested min gas price for a tx

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("15s"):
```
[Pool]
PollMinAllowedGasPriceInterval="15s"
```

### <a name="Pool_AccountQueue"></a>7.9. `Pool.AccountQueue`

**Type:** : `integer`

**Default:** `64`

**Description:** AccountQueue represents the maximum number of non-executable transaction slots permitted per account

**Example setting the default value** (64):
```
[Pool]
AccountQueue=64
```

### <a name="Pool_GlobalQueue"></a>7.10. `Pool.GlobalQueue`

**Type:** : `integer`

**Default:** `1024`

**Description:** GlobalQueue represents the maximum number of non-executable transaction slots for all accounts

**Example setting the default value** (1024):
```
[Pool]
GlobalQueue=1024
```

### <a name="Pool_EffectiveGasPrice"></a>7.11. `[Pool.EffectiveGasPrice]`

**Type:** : `object`
**Description:** EffectiveGasPrice is the config for the effective gas price calculation

| Property                                                          | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                            |
| ----------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------------------------------------------------------------------------- |
| - [Enabled](#Pool_EffectiveGasPrice_Enabled )                     | No      | boolean | No         | -          | Enabled is a flag to enable/disable the effective gas price                                                                  |
| - [L1GasPriceFactor](#Pool_EffectiveGasPrice_L1GasPriceFactor )   | No      | number  | No         | -          | L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price                             |
| - [ByteGasCost](#Pool_EffectiveGasPrice_ByteGasCost )             | No      | integer | No         | -          | ByteGasCost is the gas cost per byte that is not 0                                                                           |
| - [ZeroByteGasCost](#Pool_EffectiveGasPrice_ZeroByteGasCost )     | No      | integer | No         | -          | ZeroByteGasCost is the gas cost per byte that is 0                                                                           |
| - [NetProfit](#Pool_EffectiveGasPrice_NetProfit )                 | No      | number  | No         | -          | NetProfit is the profit margin to apply to the calculated breakEvenGasPrice                                                  |
| - [BreakEvenFactor](#Pool_EffectiveGasPrice_BreakEvenFactor )     | No      | number  | No         | -          | BreakEvenFactor is the factor to apply to the calculated breakevenGasPrice when comparing it with the gasPriceSigned of a tx |
| - [FinalDeviationPct](#Pool_EffectiveGasPrice_FinalDeviationPct ) | No      | integer | No         | -          | FinalDeviationPct is the max allowed deviation percentage BreakEvenGasPrice on re-calculation                                |

#### <a name="Pool_EffectiveGasPrice_Enabled"></a>7.11.1. `Pool.EffectiveGasPrice.Enabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** Enabled is a flag to enable/disable the effective gas price

**Example setting the default value** (false):
```
[Pool.EffectiveGasPrice]
Enabled=false
```

#### <a name="Pool_EffectiveGasPrice_L1GasPriceFactor"></a>7.11.2. `Pool.EffectiveGasPrice.L1GasPriceFactor`

**Type:** : `number`

**Default:** `0.25`

**Description:** L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price

**Example setting the default value** (0.25):
```
[Pool.EffectiveGasPrice]
L1GasPriceFactor=0.25
```

#### <a name="Pool_EffectiveGasPrice_ByteGasCost"></a>7.11.3. `Pool.EffectiveGasPrice.ByteGasCost`

**Type:** : `integer`

**Default:** `16`

**Description:** ByteGasCost is the gas cost per byte that is not 0

**Example setting the default value** (16):
```
[Pool.EffectiveGasPrice]
ByteGasCost=16
```

#### <a name="Pool_EffectiveGasPrice_ZeroByteGasCost"></a>7.11.4. `Pool.EffectiveGasPrice.ZeroByteGasCost`

**Type:** : `integer`

**Default:** `4`

**Description:** ZeroByteGasCost is the gas cost per byte that is 0

**Example setting the default value** (4):
```
[Pool.EffectiveGasPrice]
ZeroByteGasCost=4
```

#### <a name="Pool_EffectiveGasPrice_NetProfit"></a>7.11.5. `Pool.EffectiveGasPrice.NetProfit`

**Type:** : `number`

**Default:** `1`

**Description:** NetProfit is the profit margin to apply to the calculated breakEvenGasPrice

**Example setting the default value** (1):
```
[Pool.EffectiveGasPrice]
NetProfit=1
```

#### <a name="Pool_EffectiveGasPrice_BreakEvenFactor"></a>7.11.6. `Pool.EffectiveGasPrice.BreakEvenFactor`

**Type:** : `number`

**Default:** `1.1`

**Description:** BreakEvenFactor is the factor to apply to the calculated breakevenGasPrice when comparing it with the gasPriceSigned of a tx

**Example setting the default value** (1.1):
```
[Pool.EffectiveGasPrice]
BreakEvenFactor=1.1
```

#### <a name="Pool_EffectiveGasPrice_FinalDeviationPct"></a>7.11.7. `Pool.EffectiveGasPrice.FinalDeviationPct`

**Type:** : `integer`

**Default:** `10`

**Description:** FinalDeviationPct is the max allowed deviation percentage BreakEvenGasPrice on re-calculation

**Example setting the default value** (10):
```
[Pool.EffectiveGasPrice]
FinalDeviationPct=10
```

## <a name="RPC"></a>8. `[RPC]`

**Type:** : `object`
**Description:** Configuration for RPC service. THis one offers a extended Ethereum JSON-RPC API interface to interact with the node

| Property                                                                     | Pattern | Type             | Deprecated | Definition | Title/Description                                                                                                                                                                     |
| ---------------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [Host](#RPC_Host )                                                         | No      | string           | No         | -          | Host defines the network adapter that will be used to serve the HTTP requests                                                                                                         |
| - [Port](#RPC_Port )                                                         | No      | integer          | No         | -          | Port defines the port to serve the endpoints via HTTP                                                                                                                                 |
| - [ReadTimeout](#RPC_ReadTimeout )                                           | No      | string           | No         | -          | Duration                                                                                                                                                                              |
| - [WriteTimeout](#RPC_WriteTimeout )                                         | No      | string           | No         | -          | Duration                                                                                                                                                                              |
| - [MaxRequestsPerIPAndSecond](#RPC_MaxRequestsPerIPAndSecond )               | No      | number           | No         | -          | MaxRequestsPerIPAndSecond defines how much requests a single IP can<br />send within a single second                                                                                  |
| - [SequencerNodeURI](#RPC_SequencerNodeURI )                                 | No      | string           | No         | -          | SequencerNodeURI is used allow Non-Sequencer nodes<br />to relay transactions to the Sequencer node                                                                                   |
| - [MaxCumulativeGasUsed](#RPC_MaxCumulativeGasUsed )                         | No      | integer          | No         | -          | MaxCumulativeGasUsed is the max gas allowed per batch                                                                                                                                 |
| - [WebSockets](#RPC_WebSockets )                                             | No      | object           | No         | -          | WebSockets configuration                                                                                                                                                              |
| - [EnableL2SuggestedGasPricePolling](#RPC_EnableL2SuggestedGasPricePolling ) | No      | boolean          | No         | -          | EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.                                                                     |
| - [BatchRequestsEnabled](#RPC_BatchRequestsEnabled )                         | No      | boolean          | No         | -          | BatchRequestsEnabled defines if the Batch requests are enabled or disabled                                                                                                            |
| - [BatchRequestsLimit](#RPC_BatchRequestsLimit )                             | No      | integer          | No         | -          | BatchRequestsLimit defines the limit of requests that can be incorporated into each batch request                                                                                     |
| - [L2Coinbase](#RPC_L2Coinbase )                                             | No      | array of integer | No         | -          | L2Coinbase defines which address is going to receive the fees                                                                                                                         |
| - [MaxLogsCount](#RPC_MaxLogsCount )                                         | No      | integer          | No         | -          | MaxLogsCount is a configuration to set the max number of logs that can be returned<br />in a single call to the state, if zero it means no limit                                      |
| - [MaxLogsBlockRange](#RPC_MaxLogsBlockRange )                               | No      | integer          | No         | -          | MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs<br />logs in a single call to the state, if zero it means no limit                       |
| - [MaxNativeBlockHashBlockRange](#RPC_MaxNativeBlockHashBlockRange )         | No      | integer          | No         | -          | MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying<br />native block hashes in a single call to the state, if zero it means no limit |
| - [EnableHttpLog](#RPC_EnableHttpLog )                                       | No      | boolean          | No         | -          | EnableHttpLog allows the user to enable or disable the logs related to the HTTP<br />requests to be captured by the server.                                                           |

### <a name="RPC_Host"></a>8.1. `RPC.Host`

**Type:** : `string`

**Default:** `"0.0.0.0"`

**Description:** Host defines the network adapter that will be used to serve the HTTP requests

**Example setting the default value** ("0.0.0.0"):
```
[RPC]
Host="0.0.0.0"
```

### <a name="RPC_Port"></a>8.2. `RPC.Port`

**Type:** : `integer`

**Default:** `8545`

**Description:** Port defines the port to serve the endpoints via HTTP

**Example setting the default value** (8545):
```
[RPC]
Port=8545
```

### <a name="RPC_ReadTimeout"></a>8.3. `RPC.ReadTimeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"1m0s"`

**Description:** ReadTimeout is the HTTP server read timeout
check net/http.server.ReadTimeout and net/http.server.ReadHeaderTimeout

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1m0s"):
```
[RPC]
ReadTimeout="1m0s"
```

### <a name="RPC_WriteTimeout"></a>8.4. `RPC.WriteTimeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"1m0s"`

**Description:** WriteTimeout is the HTTP server write timeout
check net/http.server.WriteTimeout

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1m0s"):
```
[RPC]
WriteTimeout="1m0s"
```

### <a name="RPC_MaxRequestsPerIPAndSecond"></a>8.5. `RPC.MaxRequestsPerIPAndSecond`

**Type:** : `number`

**Default:** `500`

**Description:** MaxRequestsPerIPAndSecond defines how much requests a single IP can
send within a single second

**Example setting the default value** (500):
```
[RPC]
MaxRequestsPerIPAndSecond=500
```

### <a name="RPC_SequencerNodeURI"></a>8.6. `RPC.SequencerNodeURI`

**Type:** : `string`

**Default:** `""`

**Description:** SequencerNodeURI is used allow Non-Sequencer nodes
to relay transactions to the Sequencer node

**Example setting the default value** (""):
```
[RPC]
SequencerNodeURI=""
```

### <a name="RPC_MaxCumulativeGasUsed"></a>8.7. `RPC.MaxCumulativeGasUsed`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxCumulativeGasUsed is the max gas allowed per batch

**Example setting the default value** (0):
```
[RPC]
MaxCumulativeGasUsed=0
```

### <a name="RPC_WebSockets"></a>8.8. `[RPC.WebSockets]`

**Type:** : `object`
**Description:** WebSockets configuration

| Property                                  | Pattern | Type    | Deprecated | Definition | Title/Description                                                               |
| ----------------------------------------- | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------------------- |
| - [Enabled](#RPC_WebSockets_Enabled )     | No      | boolean | No         | -          | Enabled defines if the WebSocket requests are enabled or disabled               |
| - [Host](#RPC_WebSockets_Host )           | No      | string  | No         | -          | Host defines the network adapter that will be used to serve the WS requests     |
| - [Port](#RPC_WebSockets_Port )           | No      | integer | No         | -          | Port defines the port to serve the endpoints via WS                             |
| - [ReadLimit](#RPC_WebSockets_ReadLimit ) | No      | integer | No         | -          | ReadLimit defines the maximum size of a message read from the client (in bytes) |

#### <a name="RPC_WebSockets_Enabled"></a>8.8.1. `RPC.WebSockets.Enabled`

**Type:** : `boolean`

**Default:** `true`

**Description:** Enabled defines if the WebSocket requests are enabled or disabled

**Example setting the default value** (true):
```
[RPC.WebSockets]
Enabled=true
```

#### <a name="RPC_WebSockets_Host"></a>8.8.2. `RPC.WebSockets.Host`

**Type:** : `string`

**Default:** `"0.0.0.0"`

**Description:** Host defines the network adapter that will be used to serve the WS requests

**Example setting the default value** ("0.0.0.0"):
```
[RPC.WebSockets]
Host="0.0.0.0"
```

#### <a name="RPC_WebSockets_Port"></a>8.8.3. `RPC.WebSockets.Port`

**Type:** : `integer`

**Default:** `8546`

**Description:** Port defines the port to serve the endpoints via WS

**Example setting the default value** (8546):
```
[RPC.WebSockets]
Port=8546
```

#### <a name="RPC_WebSockets_ReadLimit"></a>8.8.4. `RPC.WebSockets.ReadLimit`

**Type:** : `integer`

**Default:** `104857600`

**Description:** ReadLimit defines the maximum size of a message read from the client (in bytes)

**Example setting the default value** (104857600):
```
[RPC.WebSockets]
ReadLimit=104857600
```

### <a name="RPC_EnableL2SuggestedGasPricePolling"></a>8.9. `RPC.EnableL2SuggestedGasPricePolling`

**Type:** : `boolean`

**Default:** `true`

**Description:** EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.

**Example setting the default value** (true):
```
[RPC]
EnableL2SuggestedGasPricePolling=true
```

### <a name="RPC_BatchRequestsEnabled"></a>8.10. `RPC.BatchRequestsEnabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** BatchRequestsEnabled defines if the Batch requests are enabled or disabled

**Example setting the default value** (false):
```
[RPC]
BatchRequestsEnabled=false
```

### <a name="RPC_BatchRequestsLimit"></a>8.11. `RPC.BatchRequestsLimit`

**Type:** : `integer`

**Default:** `20`

**Description:** BatchRequestsLimit defines the limit of requests that can be incorporated into each batch request

**Example setting the default value** (20):
```
[RPC]
BatchRequestsLimit=20
```

### <a name="RPC_L2Coinbase"></a>8.12. `RPC.L2Coinbase`

**Type:** : `array of integer`
**Description:** L2Coinbase defines which address is going to receive the fees

### <a name="RPC_MaxLogsCount"></a>8.13. `RPC.MaxLogsCount`

**Type:** : `integer`

**Default:** `10000`

**Description:** MaxLogsCount is a configuration to set the max number of logs that can be returned
in a single call to the state, if zero it means no limit

**Example setting the default value** (10000):
```
[RPC]
MaxLogsCount=10000
```

### <a name="RPC_MaxLogsBlockRange"></a>8.14. `RPC.MaxLogsBlockRange`

**Type:** : `integer`

**Default:** `10000`

**Description:** MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs
logs in a single call to the state, if zero it means no limit

**Example setting the default value** (10000):
```
[RPC]
MaxLogsBlockRange=10000
```

### <a name="RPC_MaxNativeBlockHashBlockRange"></a>8.15. `RPC.MaxNativeBlockHashBlockRange`

**Type:** : `integer`

**Default:** `60000`

**Description:** MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying
native block hashes in a single call to the state, if zero it means no limit

**Example setting the default value** (60000):
```
[RPC]
MaxNativeBlockHashBlockRange=60000
```

### <a name="RPC_EnableHttpLog"></a>8.16. `RPC.EnableHttpLog`

**Type:** : `boolean`

**Default:** `true`

**Description:** EnableHttpLog allows the user to enable or disable the logs related to the HTTP
requests to be captured by the server.

**Example setting the default value** (true):
```
[RPC]
EnableHttpLog=true
```

## <a name="Synchronizer"></a>9. `[Synchronizer]`

**Type:** : `object`
**Description:** Configuration of service `Syncrhonizer`. For this service is also really important the value of `IsTrustedSequencer`
because depending of this values is going to ask to a trusted node for trusted transactions or not

| Property                                                                | Pattern | Type             | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                       |
| ----------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [SyncInterval](#Synchronizer_SyncInterval )                           | No      | string           | No         | -          | Duration                                                                                                                                                                                                                                                |
| - [SyncChunkSize](#Synchronizer_SyncChunkSize )                         | No      | integer          | No         | -          | SyncChunkSize is the number of blocks to sync on each chunk                                                                                                                                                                                             |
| - [TrustedSequencerURL](#Synchronizer_TrustedSequencerURL )             | No      | string           | No         | -          | TrustedSequencerURL is the rpc url to connect and sync the trusted state                                                                                                                                                                                |
| - [L1SynchronizationMode](#Synchronizer_L1SynchronizationMode )         | No      | enum (of string) | No         | -          | L1SynchronizationMode define how to synchronize with L1:<br />- parallel: Request data to L1 in parallel, and process sequentially. The advantage is that executor is not blocked waiting for L1 data<br />- sequential: Request data to L1 and execute |
| - [L1ParallelSynchronization](#Synchronizer_L1ParallelSynchronization ) | No      | object           | No         | -          | L1ParallelSynchronization Configuration for parallel mode (if L1SynchronizationMode equal to 'parallel')                                                                                                                                                |

### <a name="Synchronizer_SyncInterval"></a>9.1. `Synchronizer.SyncInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"1s"`

**Description:** SyncInterval is the delay interval between reading new rollup information

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1s"):
```
[Synchronizer]
SyncInterval="1s"
```

### <a name="Synchronizer_SyncChunkSize"></a>9.2. `Synchronizer.SyncChunkSize`

**Type:** : `integer`

**Default:** `100`

**Description:** SyncChunkSize is the number of blocks to sync on each chunk

**Example setting the default value** (100):
```
[Synchronizer]
SyncChunkSize=100
```

### <a name="Synchronizer_TrustedSequencerURL"></a>9.3. `Synchronizer.TrustedSequencerURL`

**Type:** : `string`

**Default:** `""`

**Description:** TrustedSequencerURL is the rpc url to connect and sync the trusted state

**Example setting the default value** (""):
```
[Synchronizer]
TrustedSequencerURL=""
```

### <a name="Synchronizer_L1SynchronizationMode"></a>9.4. `Synchronizer.L1SynchronizationMode`

**Type:** : `enum (of string)`

**Default:** `"parallel"`

**Description:** L1SynchronizationMode define how to synchronize with L1:
- parallel: Request data to L1 in parallel, and process sequentially. The advantage is that executor is not blocked waiting for L1 data
- sequential: Request data to L1 and execute

**Example setting the default value** ("parallel"):
```
[Synchronizer]
L1SynchronizationMode="parallel"
```

Must be one of:
* "sequential"
* "parallel"

### <a name="Synchronizer_L1ParallelSynchronization"></a>9.5. `[Synchronizer.L1ParallelSynchronization]`

**Type:** : `object`
**Description:** L1ParallelSynchronization Configuration for parallel mode (if L1SynchronizationMode equal to 'parallel')

| Property                                                                                                                    | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                                                             |
| --------------------------------------------------------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [MaxClients](#Synchronizer_L1ParallelSynchronization_MaxClients )                                                         | No      | integer | No         | -          | MaxClients Number of clients used to synchronize with L1                                                                                                                                      |
| - [MaxPendingNoProcessedBlocks](#Synchronizer_L1ParallelSynchronization_MaxPendingNoProcessedBlocks )                       | No      | integer | No         | -          | MaxPendingNoProcessedBlocks Size of the buffer used to store rollup information from L1, must be >= to NumberOfEthereumClientsToSync<br />sugested twice of NumberOfParallelOfEthereumClients |
| - [RequestLastBlockPeriod](#Synchronizer_L1ParallelSynchronization_RequestLastBlockPeriod )                                 | No      | string  | No         | -          | Duration                                                                                                                                                                                      |
| - [PerformanceWarning](#Synchronizer_L1ParallelSynchronization_PerformanceWarning )                                         | No      | object  | No         | -          | Consumer Configuration for the consumer of rollup information from L1                                                                                                                         |
| - [RequestLastBlockTimeout](#Synchronizer_L1ParallelSynchronization_RequestLastBlockTimeout )                               | No      | string  | No         | -          | Duration                                                                                                                                                                                      |
| - [RequestLastBlockMaxRetries](#Synchronizer_L1ParallelSynchronization_RequestLastBlockMaxRetries )                         | No      | integer | No         | -          | RequestLastBlockMaxRetries Max number of retries to request LastBlock On L1                                                                                                                   |
| - [StatisticsPeriod](#Synchronizer_L1ParallelSynchronization_StatisticsPeriod )                                             | No      | string  | No         | -          | Duration                                                                                                                                                                                      |
| - [TimeOutMainLoop](#Synchronizer_L1ParallelSynchronization_TimeOutMainLoop )                                               | No      | string  | No         | -          | Duration                                                                                                                                                                                      |
| - [RollupInfoRetriesSpacing](#Synchronizer_L1ParallelSynchronization_RollupInfoRetriesSpacing )                             | No      | string  | No         | -          | Duration                                                                                                                                                                                      |
| - [FallbackToSequentialModeOnSynchronized](#Synchronizer_L1ParallelSynchronization_FallbackToSequentialModeOnSynchronized ) | No      | boolean | No         | -          | FallbackToSequentialModeOnSynchronized if true switch to sequential mode if the system is synchronized                                                                                        |

#### <a name="Synchronizer_L1ParallelSynchronization_MaxClients"></a>9.5.1. `Synchronizer.L1ParallelSynchronization.MaxClients`

**Type:** : `integer`

**Default:** `10`

**Description:** MaxClients Number of clients used to synchronize with L1

**Example setting the default value** (10):
```
[Synchronizer.L1ParallelSynchronization]
MaxClients=10
```

#### <a name="Synchronizer_L1ParallelSynchronization_MaxPendingNoProcessedBlocks"></a>9.5.2. `Synchronizer.L1ParallelSynchronization.MaxPendingNoProcessedBlocks`

**Type:** : `integer`

**Default:** `25`

**Description:** MaxPendingNoProcessedBlocks Size of the buffer used to store rollup information from L1, must be >= to NumberOfEthereumClientsToSync
sugested twice of NumberOfParallelOfEthereumClients

**Example setting the default value** (25):
```
[Synchronizer.L1ParallelSynchronization]
MaxPendingNoProcessedBlocks=25
```

#### <a name="Synchronizer_L1ParallelSynchronization_RequestLastBlockPeriod"></a>9.5.3. `Synchronizer.L1ParallelSynchronization.RequestLastBlockPeriod`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** RequestLastBlockPeriod is the time to wait to request the
last block to L1 to known if we need to retrieve more data.
This value only apply when the system is synchronized

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Synchronizer.L1ParallelSynchronization]
RequestLastBlockPeriod="5s"
```

#### <a name="Synchronizer_L1ParallelSynchronization_PerformanceWarning"></a>9.5.4. `[Synchronizer.L1ParallelSynchronization.PerformanceWarning]`

**Type:** : `object`
**Description:** Consumer Configuration for the consumer of rollup information from L1

| Property                                                                                                                 | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                        |
| ------------------------------------------------------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------ |
| - [AceptableInacctivityTime](#Synchronizer_L1ParallelSynchronization_PerformanceWarning_AceptableInacctivityTime )       | No      | string  | No         | -          | Duration                                                                                                                 |
| - [ApplyAfterNumRollupReceived](#Synchronizer_L1ParallelSynchronization_PerformanceWarning_ApplyAfterNumRollupReceived ) | No      | integer | No         | -          | ApplyAfterNumRollupReceived is the number of iterations to<br />start checking the time waiting for new rollup info data |

##### <a name="Synchronizer_L1ParallelSynchronization_PerformanceWarning_AceptableInacctivityTime"></a>9.5.4.1. `Synchronizer.L1ParallelSynchronization.PerformanceWarning.AceptableInacctivityTime`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** AceptableInacctivityTime is the expected maximum time that the consumer
could wait until new data is produced. If the time is greater it emmit a log to warn about
that. The idea is keep working the consumer as much as possible, so if the producer is not
fast enought then you could increse the number of parallel clients to sync with L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Synchronizer.L1ParallelSynchronization.PerformanceWarning]
AceptableInacctivityTime="5s"
```

##### <a name="Synchronizer_L1ParallelSynchronization_PerformanceWarning_ApplyAfterNumRollupReceived"></a>9.5.4.2. `Synchronizer.L1ParallelSynchronization.PerformanceWarning.ApplyAfterNumRollupReceived`

**Type:** : `integer`

**Default:** `10`

**Description:** ApplyAfterNumRollupReceived is the number of iterations to
start checking the time waiting for new rollup info data

**Example setting the default value** (10):
```
[Synchronizer.L1ParallelSynchronization.PerformanceWarning]
ApplyAfterNumRollupReceived=10
```

#### <a name="Synchronizer_L1ParallelSynchronization_RequestLastBlockTimeout"></a>9.5.5. `Synchronizer.L1ParallelSynchronization.RequestLastBlockTimeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** RequestLastBlockTimeout Timeout for request LastBlock On L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Synchronizer.L1ParallelSynchronization]
RequestLastBlockTimeout="5s"
```

#### <a name="Synchronizer_L1ParallelSynchronization_RequestLastBlockMaxRetries"></a>9.5.6. `Synchronizer.L1ParallelSynchronization.RequestLastBlockMaxRetries`

**Type:** : `integer`

**Default:** `3`

**Description:** RequestLastBlockMaxRetries Max number of retries to request LastBlock On L1

**Example setting the default value** (3):
```
[Synchronizer.L1ParallelSynchronization]
RequestLastBlockMaxRetries=3
```

#### <a name="Synchronizer_L1ParallelSynchronization_StatisticsPeriod"></a>9.5.7. `Synchronizer.L1ParallelSynchronization.StatisticsPeriod`

**Title:** Duration

**Type:** : `string`

**Default:** `"5m0s"`

**Description:** StatisticsPeriod how ofter show a log with statistics (0 is disabled)

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5m0s"):
```
[Synchronizer.L1ParallelSynchronization]
StatisticsPeriod="5m0s"
```

#### <a name="Synchronizer_L1ParallelSynchronization_TimeOutMainLoop"></a>9.5.8. `Synchronizer.L1ParallelSynchronization.TimeOutMainLoop`

**Title:** Duration

**Type:** : `string`

**Default:** `"5m0s"`

**Description:** TimeOutMainLoop is the timeout for the main loop of the L1 synchronizer when is not updated

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5m0s"):
```
[Synchronizer.L1ParallelSynchronization]
TimeOutMainLoop="5m0s"
```

#### <a name="Synchronizer_L1ParallelSynchronization_RollupInfoRetriesSpacing"></a>9.5.9. `Synchronizer.L1ParallelSynchronization.RollupInfoRetriesSpacing`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** RollupInfoRetriesSpacing is the minimum time between retries to request rollup info (it will sleep for fulfill this time) to avoid spamming L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Synchronizer.L1ParallelSynchronization]
RollupInfoRetriesSpacing="5s"
```

#### <a name="Synchronizer_L1ParallelSynchronization_FallbackToSequentialModeOnSynchronized"></a>9.5.10. `Synchronizer.L1ParallelSynchronization.FallbackToSequentialModeOnSynchronized`

**Type:** : `boolean`

**Default:** `false`

**Description:** FallbackToSequentialModeOnSynchronized if true switch to sequential mode if the system is synchronized

**Example setting the default value** (false):
```
[Synchronizer.L1ParallelSynchronization]
FallbackToSequentialModeOnSynchronized=false
```

## <a name="Sequencer"></a>10. `[Sequencer]`

**Type:** : `object`
**Description:** Configuration of the sequencer service

| Property                                                                     | Pattern | Type    | Deprecated | Definition | Title/Description                                                                            |
| ---------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | -------------------------------------------------------------------------------------------- |
| - [WaitPeriodPoolIsEmpty](#Sequencer_WaitPeriodPoolIsEmpty )                 | No      | string  | No         | -          | Duration                                                                                     |
| - [BlocksAmountForTxsToBeDeleted](#Sequencer_BlocksAmountForTxsToBeDeleted ) | No      | integer | No         | -          | BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool |
| - [FrequencyToCheckTxsForDelete](#Sequencer_FrequencyToCheckTxsForDelete )   | No      | string  | No         | -          | Duration                                                                                     |
| - [TxLifetimeCheckTimeout](#Sequencer_TxLifetimeCheckTimeout )               | No      | string  | No         | -          | Duration                                                                                     |
| - [MaxTxLifetime](#Sequencer_MaxTxLifetime )                                 | No      | string  | No         | -          | Duration                                                                                     |
| - [Finalizer](#Sequencer_Finalizer )                                         | No      | object  | No         | -          | Finalizer's specific config properties                                                       |
| - [DBManager](#Sequencer_DBManager )                                         | No      | object  | No         | -          | DBManager's specific config properties                                                       |
| - [StreamServer](#Sequencer_StreamServer )                                   | No      | object  | No         | -          | StreamServerCfg is the config for the stream server                                          |

### <a name="Sequencer_WaitPeriodPoolIsEmpty"></a>10.1. `Sequencer.WaitPeriodPoolIsEmpty`

**Title:** Duration

**Type:** : `string`

**Default:** `"1s"`

**Description:** WaitPeriodPoolIsEmpty is the time the sequencer waits until
trying to add new txs to the state

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1s"):
```
[Sequencer]
WaitPeriodPoolIsEmpty="1s"
```

### <a name="Sequencer_BlocksAmountForTxsToBeDeleted"></a>10.2. `Sequencer.BlocksAmountForTxsToBeDeleted`

**Type:** : `integer`

**Default:** `100`

**Description:** BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool

**Example setting the default value** (100):
```
[Sequencer]
BlocksAmountForTxsToBeDeleted=100
```

### <a name="Sequencer_FrequencyToCheckTxsForDelete"></a>10.3. `Sequencer.FrequencyToCheckTxsForDelete`

**Title:** Duration

**Type:** : `string`

**Default:** `"12h0m0s"`

**Description:** FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("12h0m0s"):
```
[Sequencer]
FrequencyToCheckTxsForDelete="12h0m0s"
```

### <a name="Sequencer_TxLifetimeCheckTimeout"></a>10.4. `Sequencer.TxLifetimeCheckTimeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"10m0s"`

**Description:** TxLifetimeCheckTimeout is the time the sequencer waits to check txs lifetime

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("10m0s"):
```
[Sequencer]
TxLifetimeCheckTimeout="10m0s"
```

### <a name="Sequencer_MaxTxLifetime"></a>10.5. `Sequencer.MaxTxLifetime`

**Title:** Duration

**Type:** : `string`

**Default:** `"3h0m0s"`

**Description:** MaxTxLifetime is the time a tx can be in the sequencer/worker memory

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("3h0m0s"):
```
[Sequencer]
MaxTxLifetime="3h0m0s"
```

### <a name="Sequencer_Finalizer"></a>10.6. `[Sequencer.Finalizer]`

**Type:** : `object`
**Description:** Finalizer's specific config properties

| Property                                                                                                                       | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                                                                              |
| ------------------------------------------------------------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [GERDeadlineTimeout](#Sequencer_Finalizer_GERDeadlineTimeout )                                                               | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [ForcedBatchDeadlineTimeout](#Sequencer_Finalizer_ForcedBatchDeadlineTimeout )                                               | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [SleepDuration](#Sequencer_Finalizer_SleepDuration )                                                                         | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [ResourcePercentageToCloseBatch](#Sequencer_Finalizer_ResourcePercentageToCloseBatch )                                       | No      | integer | No         | -          | ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed                                                                                                    |
| - [GERFinalityNumberOfBlocks](#Sequencer_Finalizer_GERFinalityNumberOfBlocks )                                                 | No      | integer | No         | -          | GERFinalityNumberOfBlocks is number of blocks to consider GER final                                                                                                                                            |
| - [ClosingSignalsManagerWaitForCheckingL1Timeout](#Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingL1Timeout )         | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [ClosingSignalsManagerWaitForCheckingGER](#Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingGER )                     | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [ClosingSignalsManagerWaitForCheckingForcedBatches](#Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingForcedBatches ) | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [ForcedBatchesFinalityNumberOfBlocks](#Sequencer_Finalizer_ForcedBatchesFinalityNumberOfBlocks )                             | No      | integer | No         | -          | ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final                                                                                                                                  |
| - [TimestampResolution](#Sequencer_Finalizer_TimestampResolution )                                                             | No      | string  | No         | -          | Duration                                                                                                                                                                                                       |
| - [StopSequencerOnBatchNum](#Sequencer_Finalizer_StopSequencerOnBatchNum )                                                     | No      | integer | No         | -          | StopSequencerOnBatchNum specifies the batch number where the Sequencer will stop to process more transactions and generate new batches. The Sequencer will halt after it closes the batch equal to this number |
| - [SequentialReprocessFullBatch](#Sequencer_Finalizer_SequentialReprocessFullBatch )                                           | No      | boolean | No         | -          | SequentialReprocessFullBatch indicates if the reprocess of a closed batch (sanity check) must be done in a<br />sequential way (instead than in parallel)                                                      |

#### <a name="Sequencer_Finalizer_GERDeadlineTimeout"></a>10.6.1. `Sequencer.Finalizer.GERDeadlineTimeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** GERDeadlineTimeout is the time the finalizer waits after receiving closing signal to update Global Exit Root

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Sequencer.Finalizer]
GERDeadlineTimeout="5s"
```

#### <a name="Sequencer_Finalizer_ForcedBatchDeadlineTimeout"></a>10.6.2. `Sequencer.Finalizer.ForcedBatchDeadlineTimeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"1m0s"`

**Description:** ForcedBatchDeadlineTimeout is the time the finalizer waits after receiving closing signal to process Forced Batches

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1m0s"):
```
[Sequencer.Finalizer]
ForcedBatchDeadlineTimeout="1m0s"
```

#### <a name="Sequencer_Finalizer_SleepDuration"></a>10.6.3. `Sequencer.Finalizer.SleepDuration`

**Title:** Duration

**Type:** : `string`

**Default:** `"100ms"`

**Description:** SleepDuration is the time the finalizer sleeps between each iteration, if there are no transactions to be processed

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("100ms"):
```
[Sequencer.Finalizer]
SleepDuration="100ms"
```

#### <a name="Sequencer_Finalizer_ResourcePercentageToCloseBatch"></a>10.6.4. `Sequencer.Finalizer.ResourcePercentageToCloseBatch`

**Type:** : `integer`

**Default:** `10`

**Description:** ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed

**Example setting the default value** (10):
```
[Sequencer.Finalizer]
ResourcePercentageToCloseBatch=10
```

#### <a name="Sequencer_Finalizer_GERFinalityNumberOfBlocks"></a>10.6.5. `Sequencer.Finalizer.GERFinalityNumberOfBlocks`

**Type:** : `integer`

**Default:** `64`

**Description:** GERFinalityNumberOfBlocks is number of blocks to consider GER final

**Example setting the default value** (64):
```
[Sequencer.Finalizer]
GERFinalityNumberOfBlocks=64
```

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingL1Timeout"></a>10.6.6. `Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingL1Timeout`

**Title:** Duration

**Type:** : `string`

**Default:** `"10s"`

**Description:** ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("10s"):
```
[Sequencer.Finalizer]
ClosingSignalsManagerWaitForCheckingL1Timeout="10s"
```

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingGER"></a>10.6.7. `Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingGER`

**Title:** Duration

**Type:** : `string`

**Default:** `"10s"`

**Description:** ClosingSignalsManagerWaitForCheckingGER is used by the closing signals manager to wait for its operation

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("10s"):
```
[Sequencer.Finalizer]
ClosingSignalsManagerWaitForCheckingGER="10s"
```

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingForcedBatches"></a>10.6.8. `Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingForcedBatches`

**Title:** Duration

**Type:** : `string`

**Default:** `"10s"`

**Description:** ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("10s"):
```
[Sequencer.Finalizer]
ClosingSignalsManagerWaitForCheckingForcedBatches="10s"
```

#### <a name="Sequencer_Finalizer_ForcedBatchesFinalityNumberOfBlocks"></a>10.6.9. `Sequencer.Finalizer.ForcedBatchesFinalityNumberOfBlocks`

**Type:** : `integer`

**Default:** `64`

**Description:** ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final

**Example setting the default value** (64):
```
[Sequencer.Finalizer]
ForcedBatchesFinalityNumberOfBlocks=64
```

#### <a name="Sequencer_Finalizer_TimestampResolution"></a>10.6.10. `Sequencer.Finalizer.TimestampResolution`

**Title:** Duration

**Type:** : `string`

**Default:** `"10s"`

**Description:** TimestampResolution is the resolution of the timestamp used to close a batch

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("10s"):
```
[Sequencer.Finalizer]
TimestampResolution="10s"
```

#### <a name="Sequencer_Finalizer_StopSequencerOnBatchNum"></a>10.6.11. `Sequencer.Finalizer.StopSequencerOnBatchNum`

**Type:** : `integer`

**Default:** `0`

**Description:** StopSequencerOnBatchNum specifies the batch number where the Sequencer will stop to process more transactions and generate new batches. The Sequencer will halt after it closes the batch equal to this number

**Example setting the default value** (0):
```
[Sequencer.Finalizer]
StopSequencerOnBatchNum=0
```

#### <a name="Sequencer_Finalizer_SequentialReprocessFullBatch"></a>10.6.12. `Sequencer.Finalizer.SequentialReprocessFullBatch`

**Type:** : `boolean`

**Default:** `false`

**Description:** SequentialReprocessFullBatch indicates if the reprocess of a closed batch (sanity check) must be done in a
sequential way (instead than in parallel)

**Example setting the default value** (false):
```
[Sequencer.Finalizer]
SequentialReprocessFullBatch=false
```

### <a name="Sequencer_DBManager"></a>10.7. `[Sequencer.DBManager]`

**Type:** : `object`
**Description:** DBManager's specific config properties

| Property                                                                     | Pattern | Type   | Deprecated | Definition | Title/Description |
| ---------------------------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| - [PoolRetrievalInterval](#Sequencer_DBManager_PoolRetrievalInterval )       | No      | string | No         | -          | Duration          |
| - [L2ReorgRetrievalInterval](#Sequencer_DBManager_L2ReorgRetrievalInterval ) | No      | string | No         | -          | Duration          |

#### <a name="Sequencer_DBManager_PoolRetrievalInterval"></a>10.7.1. `Sequencer.DBManager.PoolRetrievalInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"500ms"`

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("500ms"):
```
[Sequencer.DBManager]
PoolRetrievalInterval="500ms"
```

#### <a name="Sequencer_DBManager_L2ReorgRetrievalInterval"></a>10.7.2. `Sequencer.DBManager.L2ReorgRetrievalInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Sequencer.DBManager]
L2ReorgRetrievalInterval="5s"
```

### <a name="Sequencer_StreamServer"></a>10.8. `[Sequencer.StreamServer]`

**Type:** : `object`
**Description:** StreamServerCfg is the config for the stream server

| Property                                        | Pattern | Type    | Deprecated | Definition | Title/Description                                     |
| ----------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------------------------------------------- |
| - [Port](#Sequencer_StreamServer_Port )         | No      | integer | No         | -          | Port to listen on                                     |
| - [Filename](#Sequencer_StreamServer_Filename ) | No      | string  | No         | -          | Filename of the binary data file                      |
| - [Enabled](#Sequencer_StreamServer_Enabled )   | No      | boolean | No         | -          | Enabled is a flag to enable/disable the data streamer |
| - [Log](#Sequencer_StreamServer_Log )           | No      | object  | No         | -          | Log is the log configuration                          |

#### <a name="Sequencer_StreamServer_Port"></a>10.8.1. `Sequencer.StreamServer.Port`

**Type:** : `integer`

**Default:** `0`

**Description:** Port to listen on

**Example setting the default value** (0):
```
[Sequencer.StreamServer]
Port=0
```

#### <a name="Sequencer_StreamServer_Filename"></a>10.8.2. `Sequencer.StreamServer.Filename`

**Type:** : `string`

**Default:** `""`

**Description:** Filename of the binary data file

**Example setting the default value** (""):
```
[Sequencer.StreamServer]
Filename=""
```

#### <a name="Sequencer_StreamServer_Enabled"></a>10.8.3. `Sequencer.StreamServer.Enabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** Enabled is a flag to enable/disable the data streamer

**Example setting the default value** (false):
```
[Sequencer.StreamServer]
Enabled=false
```

#### <a name="Sequencer_StreamServer_Log"></a>10.8.4. `[Sequencer.StreamServer.Log]`

**Type:** : `object`
**Description:** Log is the log configuration

| Property                                                  | Pattern | Type             | Deprecated | Definition | Title/Description |
| --------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------- |
| - [Environment](#Sequencer_StreamServer_Log_Environment ) | No      | enum (of string) | No         | -          | -                 |
| - [Level](#Sequencer_StreamServer_Log_Level )             | No      | enum (of string) | No         | -          | -                 |
| - [Outputs](#Sequencer_StreamServer_Log_Outputs )         | No      | array of string  | No         | -          | -                 |

##### <a name="Sequencer_StreamServer_Log_Environment"></a>10.8.4.1. `Sequencer.StreamServer.Log.Environment`

**Type:** : `enum (of string)`

**Default:** `""`

**Example setting the default value** (""):
```
[Sequencer.StreamServer.Log]
Environment=""
```

Must be one of:
* "production"
* "development"

##### <a name="Sequencer_StreamServer_Log_Level"></a>10.8.4.2. `Sequencer.StreamServer.Log.Level`

**Type:** : `enum (of string)`

**Default:** `""`

**Example setting the default value** (""):
```
[Sequencer.StreamServer.Log]
Level=""
```

Must be one of:
* "debug"
* "info"
* "warn"
* "error"
* "dpanic"
* "panic"
* "fatal"

##### <a name="Sequencer_StreamServer_Log_Outputs"></a>10.8.4.3. `Sequencer.StreamServer.Log.Outputs`

**Type:** : `array of string`

## <a name="SequenceSender"></a>11. `[SequenceSender]`

**Type:** : `object`
**Description:** Configuration of the sequence sender service

| Property                                                                                                | Pattern | Type             | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                                                                                                                             |
| ------------------------------------------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [WaitPeriodSendSequence](#SequenceSender_WaitPeriodSendSequence )                                     | No      | string           | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [LastBatchVirtualizationTimeMaxWaitPeriod](#SequenceSender_LastBatchVirtualizationTimeMaxWaitPeriod ) | No      | string           | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [MaxTxSizeForL1](#SequenceSender_MaxTxSizeForL1 )                                                     | No      | integer          | No         | -          | MaxTxSizeForL1 is the maximum size a single transaction can have. This field has<br />non-trivial consequences: larger transactions than 128KB are significantly harder and<br />more expensive to propagate; larger transactions also take more resources<br />to validate whether they fit into the pool or not.                                                                                                            |
| - [SenderAddress](#SequenceSender_SenderAddress )                                                       | No      | array of integer | No         | -          | SenderAddress defines which private key the eth tx manager needs to use<br />to sign the L1 txs                                                                                                                                                                                                                                                                                                                               |
| - [L2Coinbase](#SequenceSender_L2Coinbase )                                                             | No      | array of integer | No         | -          | L2Coinbase defines which address is going to receive the fees                                                                                                                                                                                                                                                                                                                                                                 |
| - [PrivateKey](#SequenceSender_PrivateKey )                                                             | No      | object           | No         | -          | PrivateKey defines all the key store files that are going<br />to be read in order to provide the private keys to sign the L1 txs                                                                                                                                                                                                                                                                                             |
| - [ForkUpgradeBatchNumber](#SequenceSender_ForkUpgradeBatchNumber )                                     | No      | integer          | No         | -          | Batch number where there is a forkid change (fork upgrade)                                                                                                                                                                                                                                                                                                                                                                    |
| - [GasOffset](#SequenceSender_GasOffset )                                                               | No      | integer          | No         | -          | GasOffset is the amount of gas to be added to the gas estimation in order<br />to provide an amount that is higher than the estimated one. This is used<br />to avoid the TX getting reverted in case something has changed in the network<br />state after the estimation which can cause the TX to require more gas to be<br />executed.<br /><br />ex:<br />gas estimation: 1000<br />gas offset: 100<br />final gas: 1100 |

### <a name="SequenceSender_WaitPeriodSendSequence"></a>11.1. `SequenceSender.WaitPeriodSendSequence`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** WaitPeriodSendSequence is the time the sequencer waits until
trying to send a sequence to L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[SequenceSender]
WaitPeriodSendSequence="5s"
```

### <a name="SequenceSender_LastBatchVirtualizationTimeMaxWaitPeriod"></a>11.2. `SequenceSender.LastBatchVirtualizationTimeMaxWaitPeriod`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** LastBatchVirtualizationTimeMaxWaitPeriod is time since sequences should be sent

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[SequenceSender]
LastBatchVirtualizationTimeMaxWaitPeriod="5s"
```

### <a name="SequenceSender_MaxTxSizeForL1"></a>11.3. `SequenceSender.MaxTxSizeForL1`

**Type:** : `integer`

**Default:** `131072`

**Description:** MaxTxSizeForL1 is the maximum size a single transaction can have. This field has
non-trivial consequences: larger transactions than 128KB are significantly harder and
more expensive to propagate; larger transactions also take more resources
to validate whether they fit into the pool or not.

**Example setting the default value** (131072):
```
[SequenceSender]
MaxTxSizeForL1=131072
```

### <a name="SequenceSender_SenderAddress"></a>11.4. `SequenceSender.SenderAddress`

**Type:** : `array of integer`
**Description:** SenderAddress defines which private key the eth tx manager needs to use
to sign the L1 txs

### <a name="SequenceSender_L2Coinbase"></a>11.5. `SequenceSender.L2Coinbase`

**Type:** : `array of integer`

**Default:** `"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"`

**Description:** L2Coinbase defines which address is going to receive the fees

**Example setting the default value** ("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"):
```
[SequenceSender]
L2Coinbase="0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
```

### <a name="SequenceSender_PrivateKey"></a>11.6. `[SequenceSender.PrivateKey]`

**Type:** : `object`
**Description:** PrivateKey defines all the key store files that are going
to be read in order to provide the private keys to sign the L1 txs

| Property                                           | Pattern | Type   | Deprecated | Definition | Title/Description                                      |
| -------------------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------ |
| - [Path](#SequenceSender_PrivateKey_Path )         | No      | string | No         | -          | Path is the file path for the key store file           |
| - [Password](#SequenceSender_PrivateKey_Password ) | No      | string | No         | -          | Password is the password to decrypt the key store file |

#### <a name="SequenceSender_PrivateKey_Path"></a>11.6.1. `SequenceSender.PrivateKey.Path`

**Type:** : `string`

**Default:** `"/pk/sequencer.keystore"`

**Description:** Path is the file path for the key store file

**Example setting the default value** ("/pk/sequencer.keystore"):
```
[SequenceSender.PrivateKey]
Path="/pk/sequencer.keystore"
```

#### <a name="SequenceSender_PrivateKey_Password"></a>11.6.2. `SequenceSender.PrivateKey.Password`

**Type:** : `string`

**Default:** `"testonly"`

**Description:** Password is the password to decrypt the key store file

**Example setting the default value** ("testonly"):
```
[SequenceSender.PrivateKey]
Password="testonly"
```

### <a name="SequenceSender_ForkUpgradeBatchNumber"></a>11.7. `SequenceSender.ForkUpgradeBatchNumber`

**Type:** : `integer`

**Default:** `0`

**Description:** Batch number where there is a forkid change (fork upgrade)

**Example setting the default value** (0):
```
[SequenceSender]
ForkUpgradeBatchNumber=0
```

### <a name="SequenceSender_GasOffset"></a>11.8. `SequenceSender.GasOffset`

**Type:** : `integer`

**Default:** `80000`

**Description:** GasOffset is the amount of gas to be added to the gas estimation in order
to provide an amount that is higher than the estimated one. This is used
to avoid the TX getting reverted in case something has changed in the network
state after the estimation which can cause the TX to require more gas to be
executed.

ex:
gas estimation: 1000
gas offset: 100
final gas: 1100

**Example setting the default value** (80000):
```
[SequenceSender]
GasOffset=80000
```

## <a name="Aggregator"></a>12. `[Aggregator]`

**Type:** : `object`
**Description:** Configuration of the aggregator service

| Property                                                                                            | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                                                                                                                             |
| --------------------------------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [Host](#Aggregator_Host )                                                                         | No      | string  | No         | -          | Host for the grpc server                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [Port](#Aggregator_Port )                                                                         | No      | integer | No         | -          | Port for the grpc server                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [RetryTime](#Aggregator_RetryTime )                                                               | No      | string  | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [VerifyProofInterval](#Aggregator_VerifyProofInterval )                                           | No      | string  | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [ProofStatePollingInterval](#Aggregator_ProofStatePollingInterval )                               | No      | string  | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [TxProfitabilityCheckerType](#Aggregator_TxProfitabilityCheckerType )                             | No      | string  | No         | -          | TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch<br />possible values: base/acceptall                                                                                                                                                                                                                                                                                           |
| - [TxProfitabilityMinReward](#Aggregator_TxProfitabilityMinReward )                                 | No      | object  | No         | -          | TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch<br />this parameter is used for the base tx profitability checker                                                                                                                                                                                                                                                   |
| - [IntervalAfterWhichBatchConsolidateAnyway](#Aggregator_IntervalAfterWhichBatchConsolidateAnyway ) | No      | string  | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [ChainID](#Aggregator_ChainID )                                                                   | No      | integer | No         | -          | ChainID is the L2 ChainID provided by the Network Config                                                                                                                                                                                                                                                                                                                                                                      |
| - [ForkId](#Aggregator_ForkId )                                                                     | No      | integer | No         | -          | ForkID is the L2 ForkID provided by the Network Config                                                                                                                                                                                                                                                                                                                                                                        |
| - [SenderAddress](#Aggregator_SenderAddress )                                                       | No      | string  | No         | -          | SenderAddress defines which private key the eth tx manager needs to use<br />to sign the L1 txs                                                                                                                                                                                                                                                                                                                               |
| - [CleanupLockedProofsInterval](#Aggregator_CleanupLockedProofsInterval )                           | No      | string  | No         | -          | Duration                                                                                                                                                                                                                                                                                                                                                                                                                      |
| - [GeneratingProofCleanupThreshold](#Aggregator_GeneratingProofCleanupThreshold )                   | No      | string  | No         | -          | GeneratingProofCleanupThreshold represents the time interval after<br />which a proof in generating state is considered to be stuck and<br />allowed to be cleared.                                                                                                                                                                                                                                                           |
| - [GasOffset](#Aggregator_GasOffset )                                                               | No      | integer | No         | -          | GasOffset is the amount of gas to be added to the gas estimation in order<br />to provide an amount that is higher than the estimated one. This is used<br />to avoid the TX getting reverted in case something has changed in the network<br />state after the estimation which can cause the TX to require more gas to be<br />executed.<br /><br />ex:<br />gas estimation: 1000<br />gas offset: 100<br />final gas: 1100 |

### <a name="Aggregator_Host"></a>12.1. `Aggregator.Host`

**Type:** : `string`

**Default:** `"0.0.0.0"`

**Description:** Host for the grpc server

**Example setting the default value** ("0.0.0.0"):
```
[Aggregator]
Host="0.0.0.0"
```

### <a name="Aggregator_Port"></a>12.2. `Aggregator.Port`

**Type:** : `integer`

**Default:** `50081`

**Description:** Port for the grpc server

**Example setting the default value** (50081):
```
[Aggregator]
Port=50081
```

### <a name="Aggregator_RetryTime"></a>12.3. `Aggregator.RetryTime`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** RetryTime is the time the aggregator main loop sleeps if there are no proofs to aggregate
or batches to generate proofs. It is also used in the isSynced loop

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Aggregator]
RetryTime="5s"
```

### <a name="Aggregator_VerifyProofInterval"></a>12.4. `Aggregator.VerifyProofInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"1m30s"`

**Description:** VerifyProofInterval is the interval of time to verify/send an proof in L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1m30s"):
```
[Aggregator]
VerifyProofInterval="1m30s"
```

### <a name="Aggregator_ProofStatePollingInterval"></a>12.5. `Aggregator.ProofStatePollingInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"5s"`

**Description:** ProofStatePollingInterval is the interval time to polling the prover about the generation state of a proof

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5s"):
```
[Aggregator]
ProofStatePollingInterval="5s"
```

### <a name="Aggregator_TxProfitabilityCheckerType"></a>12.6. `Aggregator.TxProfitabilityCheckerType`

**Type:** : `string`

**Default:** `"acceptall"`

**Description:** TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
possible values: base/acceptall

**Example setting the default value** ("acceptall"):
```
[Aggregator]
TxProfitabilityCheckerType="acceptall"
```

### <a name="Aggregator_TxProfitabilityMinReward"></a>12.7. `[Aggregator.TxProfitabilityMinReward]`

**Type:** : `object`
**Description:** TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
this parameter is used for the base tx profitability checker

### <a name="Aggregator_IntervalAfterWhichBatchConsolidateAnyway"></a>12.8. `Aggregator.IntervalAfterWhichBatchConsolidateAnyway`

**Title:** Duration

**Type:** : `string`

**Default:** `"0s"`

**Description:** IntervalAfterWhichBatchConsolidateAnyway this is interval for the main sequencer, that will check if there is no transactions

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("0s"):
```
[Aggregator]
IntervalAfterWhichBatchConsolidateAnyway="0s"
```

### <a name="Aggregator_ChainID"></a>12.9. `Aggregator.ChainID`

**Type:** : `integer`

**Default:** `0`

**Description:** ChainID is the L2 ChainID provided by the Network Config

**Example setting the default value** (0):
```
[Aggregator]
ChainID=0
```

### <a name="Aggregator_ForkId"></a>12.10. `Aggregator.ForkId`

**Type:** : `integer`

**Default:** `0`

**Description:** ForkID is the L2 ForkID provided by the Network Config

**Example setting the default value** (0):
```
[Aggregator]
ForkId=0
```

### <a name="Aggregator_SenderAddress"></a>12.11. `Aggregator.SenderAddress`

**Type:** : `string`

**Default:** `""`

**Description:** SenderAddress defines which private key the eth tx manager needs to use
to sign the L1 txs

**Example setting the default value** (""):
```
[Aggregator]
SenderAddress=""
```

### <a name="Aggregator_CleanupLockedProofsInterval"></a>12.12. `Aggregator.CleanupLockedProofsInterval`

**Title:** Duration

**Type:** : `string`

**Default:** `"2m0s"`

**Description:** CleanupLockedProofsInterval is the interval of time to clean up locked proofs.

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("2m0s"):
```
[Aggregator]
CleanupLockedProofsInterval="2m0s"
```

### <a name="Aggregator_GeneratingProofCleanupThreshold"></a>12.13. `Aggregator.GeneratingProofCleanupThreshold`

**Type:** : `string`

**Default:** `"10m"`

**Description:** GeneratingProofCleanupThreshold represents the time interval after
which a proof in generating state is considered to be stuck and
allowed to be cleared.

**Example setting the default value** ("10m"):
```
[Aggregator]
GeneratingProofCleanupThreshold="10m"
```

### <a name="Aggregator_GasOffset"></a>12.14. `Aggregator.GasOffset`

**Type:** : `integer`

**Default:** `0`

**Description:** GasOffset is the amount of gas to be added to the gas estimation in order
to provide an amount that is higher than the estimated one. This is used
to avoid the TX getting reverted in case something has changed in the network
state after the estimation which can cause the TX to require more gas to be
executed.

ex:
gas estimation: 1000
gas offset: 100
final gas: 1100

**Example setting the default value** (0):
```
[Aggregator]
GasOffset=0
```

## <a name="NetworkConfig"></a>13. `[NetworkConfig]`

**Type:** : `object`
**Description:** Configuration of the genesis of the network. This is used to known the initial state of the network

| Property                                                                     | Pattern | Type             | Deprecated | Definition | Title/Description                                                                   |
| ---------------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------------------------------------------------------------------------- |
| - [l1Config](#NetworkConfig_l1Config )                                       | No      | object           | No         | -          | L1: Configuration related to L1                                                     |
| - [L2GlobalExitRootManagerAddr](#NetworkConfig_L2GlobalExitRootManagerAddr ) | No      | array of integer | No         | -          | DEPRECATED L2: address of the \`PolygonZkEVMGlobalExitRootL2 proxy\` smart contract |
| - [L2BridgeAddr](#NetworkConfig_L2BridgeAddr )                               | No      | array of integer | No         | -          | L2: address of the \`PolygonZkEVMBridge proxy\` smart contract                      |
| - [Genesis](#NetworkConfig_Genesis )                                         | No      | object           | No         | -          | L1: Genesis of the rollup, first block number and root                              |

### <a name="NetworkConfig_l1Config"></a>13.1. `[NetworkConfig.l1Config]`

**Type:** : `object`
**Description:** L1: Configuration related to L1

| Property                                                                                          | Pattern | Type             | Deprecated | Definition | Title/Description                                |
| ------------------------------------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ------------------------------------------------ |
| - [chainId](#NetworkConfig_l1Config_chainId )                                                     | No      | integer          | No         | -          | Chain ID of the L1 network                       |
| - [polygonZkEVMAddress](#NetworkConfig_l1Config_polygonZkEVMAddress )                             | No      | array of integer | No         | -          | Address of the L1 contract                       |
| - [maticTokenAddress](#NetworkConfig_l1Config_maticTokenAddress )                                 | No      | array of integer | No         | -          | Address of the L1 Matic token Contract           |
| - [polygonZkEVMGlobalExitRootAddress](#NetworkConfig_l1Config_polygonZkEVMGlobalExitRootAddress ) | No      | array of integer | No         | -          | Address of the L1 GlobalExitRootManager contract |

#### <a name="NetworkConfig_l1Config_chainId"></a>13.1.1. `NetworkConfig.l1Config.chainId`

**Type:** : `integer`

**Default:** `0`

**Description:** Chain ID of the L1 network

**Example setting the default value** (0):
```
[NetworkConfig.l1Config]
chainId=0
```

#### <a name="NetworkConfig_l1Config_polygonZkEVMAddress"></a>13.1.2. `NetworkConfig.l1Config.polygonZkEVMAddress`

**Type:** : `array of integer`
**Description:** Address of the L1 contract

#### <a name="NetworkConfig_l1Config_maticTokenAddress"></a>13.1.3. `NetworkConfig.l1Config.maticTokenAddress`

**Type:** : `array of integer`
**Description:** Address of the L1 Matic token Contract

#### <a name="NetworkConfig_l1Config_polygonZkEVMGlobalExitRootAddress"></a>13.1.4. `NetworkConfig.l1Config.polygonZkEVMGlobalExitRootAddress`

**Type:** : `array of integer`
**Description:** Address of the L1 GlobalExitRootManager contract

### <a name="NetworkConfig_L2GlobalExitRootManagerAddr"></a>13.2. `NetworkConfig.L2GlobalExitRootManagerAddr`

**Type:** : `array of integer`
**Description:** DEPRECATED L2: address of the `PolygonZkEVMGlobalExitRootL2 proxy` smart contract

### <a name="NetworkConfig_L2BridgeAddr"></a>13.3. `NetworkConfig.L2BridgeAddr`

**Type:** : `array of integer`
**Description:** L2: address of the `PolygonZkEVMBridge proxy` smart contract

### <a name="NetworkConfig_Genesis"></a>13.4. `[NetworkConfig.Genesis]`

**Type:** : `object`
**Description:** L1: Genesis of the rollup, first block number and root

| Property                                                     | Pattern | Type             | Deprecated | Definition | Title/Description                                                                 |
| ------------------------------------------------------------ | ------- | ---------------- | ---------- | ---------- | --------------------------------------------------------------------------------- |
| - [GenesisBlockNum](#NetworkConfig_Genesis_GenesisBlockNum ) | No      | integer          | No         | -          | GenesisBlockNum is the block number where the polygonZKEVM smc was deployed on L1 |
| - [Root](#NetworkConfig_Genesis_Root )                       | No      | array of integer | No         | -          | Root hash of the genesis block                                                    |
| - [GenesisActions](#NetworkConfig_Genesis_GenesisActions )   | No      | array of object  | No         | -          | Contracts to be deployed to L2                                                    |

#### <a name="NetworkConfig_Genesis_GenesisBlockNum"></a>13.4.1. `NetworkConfig.Genesis.GenesisBlockNum`

**Type:** : `integer`

**Default:** `0`

**Description:** GenesisBlockNum is the block number where the polygonZKEVM smc was deployed on L1

**Example setting the default value** (0):
```
[NetworkConfig.Genesis]
GenesisBlockNum=0
```

#### <a name="NetworkConfig_Genesis_Root"></a>13.4.2. `NetworkConfig.Genesis.Root`

**Type:** : `array of integer`
**Description:** Root hash of the genesis block

#### <a name="NetworkConfig_Genesis_GenesisActions"></a>13.4.3. `NetworkConfig.Genesis.GenesisActions`

**Type:** : `array of object`
**Description:** Contracts to be deployed to L2

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                                     | Description                                                               |
| ------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| [GenesisActions items](#NetworkConfig_Genesis_GenesisActions_items) | GenesisAction represents one of the values set on the SMT during genesis. |

##### <a name="autogenerated_heading_3"></a>13.4.3.1. [NetworkConfig.Genesis.GenesisActions.GenesisActions items]

**Type:** : `object`
**Description:** GenesisAction represents one of the values set on the SMT during genesis.

| Property                                                                          | Pattern | Type    | Deprecated | Definition | Title/Description |
| --------------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------- |
| - [address](#NetworkConfig_Genesis_GenesisActions_items_address )                 | No      | string  | No         | -          | -                 |
| - [type](#NetworkConfig_Genesis_GenesisActions_items_type )                       | No      | integer | No         | -          | -                 |
| - [storagePosition](#NetworkConfig_Genesis_GenesisActions_items_storagePosition ) | No      | string  | No         | -          | -                 |
| - [bytecode](#NetworkConfig_Genesis_GenesisActions_items_bytecode )               | No      | string  | No         | -          | -                 |
| - [key](#NetworkConfig_Genesis_GenesisActions_items_key )                         | No      | string  | No         | -          | -                 |
| - [value](#NetworkConfig_Genesis_GenesisActions_items_value )                     | No      | string  | No         | -          | -                 |
| - [root](#NetworkConfig_Genesis_GenesisActions_items_root )                       | No      | string  | No         | -          | -                 |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_address"></a>13.4.3.1.1. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.address`

**Type:** : `string`

##### <a name="NetworkConfig_Genesis_GenesisActions_items_type"></a>13.4.3.1.2. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.type`

**Type:** : `integer`

##### <a name="NetworkConfig_Genesis_GenesisActions_items_storagePosition"></a>13.4.3.1.3. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.storagePosition`

**Type:** : `string`

##### <a name="NetworkConfig_Genesis_GenesisActions_items_bytecode"></a>13.4.3.1.4. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.bytecode`

**Type:** : `string`

##### <a name="NetworkConfig_Genesis_GenesisActions_items_key"></a>13.4.3.1.5. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.key`

**Type:** : `string`

##### <a name="NetworkConfig_Genesis_GenesisActions_items_value"></a>13.4.3.1.6. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.value`

**Type:** : `string`

##### <a name="NetworkConfig_Genesis_GenesisActions_items_root"></a>13.4.3.1.7. `NetworkConfig.Genesis.GenesisActions.GenesisActions items.root`

**Type:** : `string`

## <a name="L2GasPriceSuggester"></a>14. `[L2GasPriceSuggester]`

**Type:** : `object`
**Description:** Configuration of the gas price suggester service

| Property                                                                       | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                        |
| ------------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| - [Type](#L2GasPriceSuggester_Type )                                           | No      | string  | No         | -          | -                                                                                                                                        |
| - [DefaultGasPriceWei](#L2GasPriceSuggester_DefaultGasPriceWei )               | No      | integer | No         | -          | DefaultGasPriceWei is used to set the gas price to be used by the default gas pricer or as minimim gas price by the follower gas pricer. |
| - [MaxGasPriceWei](#L2GasPriceSuggester_MaxGasPriceWei )                       | No      | integer | No         | -          | MaxGasPriceWei is used to limit the gas price returned by the follower gas pricer to a maximum value. It is ignored if 0.                |
| - [MaxPrice](#L2GasPriceSuggester_MaxPrice )                                   | No      | object  | No         | -          | -                                                                                                                                        |
| - [IgnorePrice](#L2GasPriceSuggester_IgnorePrice )                             | No      | object  | No         | -          | -                                                                                                                                        |
| - [CheckBlocks](#L2GasPriceSuggester_CheckBlocks )                             | No      | integer | No         | -          | -                                                                                                                                        |
| - [Percentile](#L2GasPriceSuggester_Percentile )                               | No      | integer | No         | -          | -                                                                                                                                        |
| - [UpdatePeriod](#L2GasPriceSuggester_UpdatePeriod )                           | No      | string  | No         | -          | Duration                                                                                                                                 |
| - [CleanHistoryPeriod](#L2GasPriceSuggester_CleanHistoryPeriod )               | No      | string  | No         | -          | Duration                                                                                                                                 |
| - [CleanHistoryTimeRetention](#L2GasPriceSuggester_CleanHistoryTimeRetention ) | No      | string  | No         | -          | Duration                                                                                                                                 |
| - [Factor](#L2GasPriceSuggester_Factor )                                       | No      | number  | No         | -          | -                                                                                                                                        |

### <a name="L2GasPriceSuggester_Type"></a>14.1. `L2GasPriceSuggester.Type`

**Type:** : `string`

**Default:** `"follower"`

**Example setting the default value** ("follower"):
```
[L2GasPriceSuggester]
Type="follower"
```

### <a name="L2GasPriceSuggester_DefaultGasPriceWei"></a>14.2. `L2GasPriceSuggester.DefaultGasPriceWei`

**Type:** : `integer`

**Default:** `2000000000`

**Description:** DefaultGasPriceWei is used to set the gas price to be used by the default gas pricer or as minimim gas price by the follower gas pricer.

**Example setting the default value** (2000000000):
```
[L2GasPriceSuggester]
DefaultGasPriceWei=2000000000
```

### <a name="L2GasPriceSuggester_MaxGasPriceWei"></a>14.3. `L2GasPriceSuggester.MaxGasPriceWei`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxGasPriceWei is used to limit the gas price returned by the follower gas pricer to a maximum value. It is ignored if 0.

**Example setting the default value** (0):
```
[L2GasPriceSuggester]
MaxGasPriceWei=0
```

### <a name="L2GasPriceSuggester_MaxPrice"></a>14.4. `[L2GasPriceSuggester.MaxPrice]`

**Type:** : `object`

### <a name="L2GasPriceSuggester_IgnorePrice"></a>14.5. `[L2GasPriceSuggester.IgnorePrice]`

**Type:** : `object`

### <a name="L2GasPriceSuggester_CheckBlocks"></a>14.6. `L2GasPriceSuggester.CheckBlocks`

**Type:** : `integer`

**Default:** `0`

**Example setting the default value** (0):
```
[L2GasPriceSuggester]
CheckBlocks=0
```

### <a name="L2GasPriceSuggester_Percentile"></a>14.7. `L2GasPriceSuggester.Percentile`

**Type:** : `integer`

**Default:** `0`

**Example setting the default value** (0):
```
[L2GasPriceSuggester]
Percentile=0
```

### <a name="L2GasPriceSuggester_UpdatePeriod"></a>14.8. `L2GasPriceSuggester.UpdatePeriod`

**Title:** Duration

**Type:** : `string`

**Default:** `"10s"`

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("10s"):
```
[L2GasPriceSuggester]
UpdatePeriod="10s"
```

### <a name="L2GasPriceSuggester_CleanHistoryPeriod"></a>14.9. `L2GasPriceSuggester.CleanHistoryPeriod`

**Title:** Duration

**Type:** : `string`

**Default:** `"1h0m0s"`

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1h0m0s"):
```
[L2GasPriceSuggester]
CleanHistoryPeriod="1h0m0s"
```

### <a name="L2GasPriceSuggester_CleanHistoryTimeRetention"></a>14.10. `L2GasPriceSuggester.CleanHistoryTimeRetention`

**Title:** Duration

**Type:** : `string`

**Default:** `"5m0s"`

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("5m0s"):
```
[L2GasPriceSuggester]
CleanHistoryTimeRetention="5m0s"
```

### <a name="L2GasPriceSuggester_Factor"></a>14.11. `L2GasPriceSuggester.Factor`

**Type:** : `number`

**Default:** `0.15`

**Example setting the default value** (0.15):
```
[L2GasPriceSuggester]
Factor=0.15
```

## <a name="Executor"></a>15. `[Executor]`

**Type:** : `object`
**Description:** Configuration of the executor service

| Property                                                                  | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                       |
| ------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------- |
| - [URI](#Executor_URI )                                                   | No      | string  | No         | -          | -                                                                                                                       |
| - [MaxResourceExhaustedAttempts](#Executor_MaxResourceExhaustedAttempts ) | No      | integer | No         | -          | MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion |
| - [WaitOnResourceExhaustion](#Executor_WaitOnResourceExhaustion )         | No      | string  | No         | -          | Duration                                                                                                                |
| - [MaxGRPCMessageSize](#Executor_MaxGRPCMessageSize )                     | No      | integer | No         | -          | -                                                                                                                       |

### <a name="Executor_URI"></a>15.1. `Executor.URI`

**Type:** : `string`

**Default:** `"zkevm-prover:50071"`

**Example setting the default value** ("zkevm-prover:50071"):
```
[Executor]
URI="zkevm-prover:50071"
```

### <a name="Executor_MaxResourceExhaustedAttempts"></a>15.2. `Executor.MaxResourceExhaustedAttempts`

**Type:** : `integer`

**Default:** `3`

**Description:** MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion

**Example setting the default value** (3):
```
[Executor]
MaxResourceExhaustedAttempts=3
```

### <a name="Executor_WaitOnResourceExhaustion"></a>15.3. `Executor.WaitOnResourceExhaustion`

**Title:** Duration

**Type:** : `string`

**Default:** `"1s"`

**Description:** WaitOnResourceExhaustion is the time to wait before retrying a transaction because of resource exhaustion

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("1s"):
```
[Executor]
WaitOnResourceExhaustion="1s"
```

### <a name="Executor_MaxGRPCMessageSize"></a>15.4. `Executor.MaxGRPCMessageSize`

**Type:** : `integer`

**Default:** `100000000`

**Example setting the default value** (100000000):
```
[Executor]
MaxGRPCMessageSize=100000000
```

## <a name="MTClient"></a>16. `[MTClient]`

**Type:** : `object`
**Description:** Configuration of the merkle tree client service. Not use in the node, only for testing

| Property                | Pattern | Type   | Deprecated | Definition | Title/Description      |
| ----------------------- | ------- | ------ | ---------- | ---------- | ---------------------- |
| - [URI](#MTClient_URI ) | No      | string | No         | -          | URI is the server URI. |

### <a name="MTClient_URI"></a>16.1. `MTClient.URI`

**Type:** : `string`

**Default:** `"zkevm-prover:50061"`

**Description:** URI is the server URI.

**Example setting the default value** ("zkevm-prover:50061"):
```
[MTClient]
URI="zkevm-prover:50061"
```

## <a name="Metrics"></a>17. `[Metrics]`

**Type:** : `object`
**Description:** Configuration of the metrics service, basically is where is going to publish the metrics

| Property                                         | Pattern | Type    | Deprecated | Definition | Title/Description                                                   |
| ------------------------------------------------ | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------- |
| - [Host](#Metrics_Host )                         | No      | string  | No         | -          | Host is the address to bind the metrics server                      |
| - [Port](#Metrics_Port )                         | No      | integer | No         | -          | Port is the port to bind the metrics server                         |
| - [Enabled](#Metrics_Enabled )                   | No      | boolean | No         | -          | Enabled is the flag to enable/disable the metrics server            |
| - [ProfilingHost](#Metrics_ProfilingHost )       | No      | string  | No         | -          | ProfilingHost is the address to bind the profiling server           |
| - [ProfilingPort](#Metrics_ProfilingPort )       | No      | integer | No         | -          | ProfilingPort is the port to bind the profiling server              |
| - [ProfilingEnabled](#Metrics_ProfilingEnabled ) | No      | boolean | No         | -          | ProfilingEnabled is the flag to enable/disable the profiling server |

### <a name="Metrics_Host"></a>17.1. `Metrics.Host`

**Type:** : `string`

**Default:** `"0.0.0.0"`

**Description:** Host is the address to bind the metrics server

**Example setting the default value** ("0.0.0.0"):
```
[Metrics]
Host="0.0.0.0"
```

### <a name="Metrics_Port"></a>17.2. `Metrics.Port`

**Type:** : `integer`

**Default:** `9091`

**Description:** Port is the port to bind the metrics server

**Example setting the default value** (9091):
```
[Metrics]
Port=9091
```

### <a name="Metrics_Enabled"></a>17.3. `Metrics.Enabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** Enabled is the flag to enable/disable the metrics server

**Example setting the default value** (false):
```
[Metrics]
Enabled=false
```

### <a name="Metrics_ProfilingHost"></a>17.4. `Metrics.ProfilingHost`

**Type:** : `string`

**Default:** `""`

**Description:** ProfilingHost is the address to bind the profiling server

**Example setting the default value** (""):
```
[Metrics]
ProfilingHost=""
```

### <a name="Metrics_ProfilingPort"></a>17.5. `Metrics.ProfilingPort`

**Type:** : `integer`

**Default:** `0`

**Description:** ProfilingPort is the port to bind the profiling server

**Example setting the default value** (0):
```
[Metrics]
ProfilingPort=0
```

### <a name="Metrics_ProfilingEnabled"></a>17.6. `Metrics.ProfilingEnabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** ProfilingEnabled is the flag to enable/disable the profiling server

**Example setting the default value** (false):
```
[Metrics]
ProfilingEnabled=false
```

## <a name="EventLog"></a>18. `[EventLog]`

**Type:** : `object`
**Description:** Configuration of the event database connection

| Property              | Pattern | Type   | Deprecated | Definition | Title/Description                |
| --------------------- | ------- | ------ | ---------- | ---------- | -------------------------------- |
| - [DB](#EventLog_DB ) | No      | object | No         | -          | DB is the database configuration |

### <a name="EventLog_DB"></a>18.1. `[EventLog.DB]`

**Type:** : `object`
**Description:** DB is the database configuration

| Property                               | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| -------------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#EventLog_DB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#EventLog_DB_User )           | No      | string  | No         | -          | Database User name                                         |
| - [Password](#EventLog_DB_Password )   | No      | string  | No         | -          | Database Password of the user                              |
| - [Host](#EventLog_DB_Host )           | No      | string  | No         | -          | Host address of database                                   |
| - [Port](#EventLog_DB_Port )           | No      | string  | No         | -          | Port Number of database                                    |
| - [EnableLog](#EventLog_DB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#EventLog_DB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

#### <a name="EventLog_DB_Name"></a>18.1.1. `EventLog.DB.Name`

**Type:** : `string`

**Default:** `""`

**Description:** Database name

**Example setting the default value** (""):
```
[EventLog.DB]
Name=""
```

#### <a name="EventLog_DB_User"></a>18.1.2. `EventLog.DB.User`

**Type:** : `string`

**Default:** `""`

**Description:** Database User name

**Example setting the default value** (""):
```
[EventLog.DB]
User=""
```

#### <a name="EventLog_DB_Password"></a>18.1.3. `EventLog.DB.Password`

**Type:** : `string`

**Default:** `""`

**Description:** Database Password of the user

**Example setting the default value** (""):
```
[EventLog.DB]
Password=""
```

#### <a name="EventLog_DB_Host"></a>18.1.4. `EventLog.DB.Host`

**Type:** : `string`

**Default:** `""`

**Description:** Host address of database

**Example setting the default value** (""):
```
[EventLog.DB]
Host=""
```

#### <a name="EventLog_DB_Port"></a>18.1.5. `EventLog.DB.Port`

**Type:** : `string`

**Default:** `""`

**Description:** Port Number of database

**Example setting the default value** (""):
```
[EventLog.DB]
Port=""
```

#### <a name="EventLog_DB_EnableLog"></a>18.1.6. `EventLog.DB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[EventLog.DB]
EnableLog=false
```

#### <a name="EventLog_DB_MaxConns"></a>18.1.7. `EventLog.DB.MaxConns`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (0):
```
[EventLog.DB]
MaxConns=0
```

## <a name="HashDB"></a>19. `[HashDB]`

**Type:** : `object`
**Description:** Configuration of the hash database connection

| Property                          | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| --------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#HashDB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#HashDB_User )           | No      | string  | No         | -          | Database User name                                         |
| - [Password](#HashDB_Password )   | No      | string  | No         | -          | Database Password of the user                              |
| - [Host](#HashDB_Host )           | No      | string  | No         | -          | Host address of database                                   |
| - [Port](#HashDB_Port )           | No      | string  | No         | -          | Port Number of database                                    |
| - [EnableLog](#HashDB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#HashDB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

### <a name="HashDB_Name"></a>19.1. `HashDB.Name`

**Type:** : `string`

**Default:** `"prover_db"`

**Description:** Database name

**Example setting the default value** ("prover_db"):
```
[HashDB]
Name="prover_db"
```

### <a name="HashDB_User"></a>19.2. `HashDB.User`

**Type:** : `string`

**Default:** `"prover_user"`

**Description:** Database User name

**Example setting the default value** ("prover_user"):
```
[HashDB]
User="prover_user"
```

### <a name="HashDB_Password"></a>19.3. `HashDB.Password`

**Type:** : `string`

**Default:** `"prover_pass"`

**Description:** Database Password of the user

**Example setting the default value** ("prover_pass"):
```
[HashDB]
Password="prover_pass"
```

### <a name="HashDB_Host"></a>19.4. `HashDB.Host`

**Type:** : `string`

**Default:** `"zkevm-state-db"`

**Description:** Host address of database

**Example setting the default value** ("zkevm-state-db"):
```
[HashDB]
Host="zkevm-state-db"
```

### <a name="HashDB_Port"></a>19.5. `HashDB.Port`

**Type:** : `string`

**Default:** `"5432"`

**Description:** Port Number of database

**Example setting the default value** ("5432"):
```
[HashDB]
Port="5432"
```

### <a name="HashDB_EnableLog"></a>19.6. `HashDB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[HashDB]
EnableLog=false
```

### <a name="HashDB_MaxConns"></a>19.7. `HashDB.MaxConns`

**Type:** : `integer`

**Default:** `200`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (200):
```
[HashDB]
MaxConns=200
```

## <a name="State"></a>20. `[State]`

**Type:** : `object`
**Description:** State service configuration

| Property                                                               | Pattern | Type            | Deprecated | Definition | Title/Description                                                                                                                                                                     |
| ---------------------------------------------------------------------- | ------- | --------------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [MaxCumulativeGasUsed](#State_MaxCumulativeGasUsed )                 | No      | integer         | No         | -          | MaxCumulativeGasUsed is the max gas allowed per batch                                                                                                                                 |
| - [ChainID](#State_ChainID )                                           | No      | integer         | No         | -          | ChainID is the L2 ChainID provided by the Network Config                                                                                                                              |
| - [ForkIDIntervals](#State_ForkIDIntervals )                           | No      | array of object | No         | -          | ForkIdIntervals is the list of fork id intervals                                                                                                                                      |
| - [MaxResourceExhaustedAttempts](#State_MaxResourceExhaustedAttempts ) | No      | integer         | No         | -          | MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion                                                               |
| - [WaitOnResourceExhaustion](#State_WaitOnResourceExhaustion )         | No      | string          | No         | -          | Duration                                                                                                                                                                              |
| - [ForkUpgradeBatchNumber](#State_ForkUpgradeBatchNumber )             | No      | integer         | No         | -          | Batch number from which there is a forkid change (fork upgrade)                                                                                                                       |
| - [ForkUpgradeNewForkId](#State_ForkUpgradeNewForkId )                 | No      | integer         | No         | -          | New fork id to be used for batches greaters than ForkUpgradeBatchNumber (fork upgrade)                                                                                                |
| - [DB](#State_DB )                                                     | No      | object          | No         | -          | DB is the database configuration                                                                                                                                                      |
| - [Batch](#State_Batch )                                               | No      | object          | No         | -          | Configuration for the batch constraints                                                                                                                                               |
| - [MaxLogsCount](#State_MaxLogsCount )                                 | No      | integer         | No         | -          | MaxLogsCount is a configuration to set the max number of logs that can be returned<br />in a single call to the state, if zero it means no limit                                      |
| - [MaxLogsBlockRange](#State_MaxLogsBlockRange )                       | No      | integer         | No         | -          | MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs<br />logs in a single call to the state, if zero it means no limit                       |
| - [MaxNativeBlockHashBlockRange](#State_MaxNativeBlockHashBlockRange ) | No      | integer         | No         | -          | MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying<br />native block hashes in a single call to the state, if zero it means no limit |

### <a name="State_MaxCumulativeGasUsed"></a>20.1. `State.MaxCumulativeGasUsed`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxCumulativeGasUsed is the max gas allowed per batch

**Example setting the default value** (0):
```
[State]
MaxCumulativeGasUsed=0
```

### <a name="State_ChainID"></a>20.2. `State.ChainID`

**Type:** : `integer`

**Default:** `0`

**Description:** ChainID is the L2 ChainID provided by the Network Config

**Example setting the default value** (0):
```
[State]
ChainID=0
```

### <a name="State_ForkIDIntervals"></a>20.3. `State.ForkIDIntervals`

**Type:** : `array of object`
**Description:** ForkIdIntervals is the list of fork id intervals

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                       | Description                          |
| ----------------------------------------------------- | ------------------------------------ |
| [ForkIDIntervals items](#State_ForkIDIntervals_items) | ForkIDInterval is a fork id interval |

#### <a name="autogenerated_heading_4"></a>20.3.1. [State.ForkIDIntervals.ForkIDIntervals items]

**Type:** : `object`
**Description:** ForkIDInterval is a fork id interval

| Property                                                           | Pattern | Type    | Deprecated | Definition | Title/Description |
| ------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ----------------- |
| - [FromBatchNumber](#State_ForkIDIntervals_items_FromBatchNumber ) | No      | integer | No         | -          | -                 |
| - [ToBatchNumber](#State_ForkIDIntervals_items_ToBatchNumber )     | No      | integer | No         | -          | -                 |
| - [ForkId](#State_ForkIDIntervals_items_ForkId )                   | No      | integer | No         | -          | -                 |
| - [Version](#State_ForkIDIntervals_items_Version )                 | No      | string  | No         | -          | -                 |
| - [BlockNumber](#State_ForkIDIntervals_items_BlockNumber )         | No      | integer | No         | -          | -                 |

##### <a name="State_ForkIDIntervals_items_FromBatchNumber"></a>20.3.1.1. `State.ForkIDIntervals.ForkIDIntervals items.FromBatchNumber`

**Type:** : `integer`

##### <a name="State_ForkIDIntervals_items_ToBatchNumber"></a>20.3.1.2. `State.ForkIDIntervals.ForkIDIntervals items.ToBatchNumber`

**Type:** : `integer`

##### <a name="State_ForkIDIntervals_items_ForkId"></a>20.3.1.3. `State.ForkIDIntervals.ForkIDIntervals items.ForkId`

**Type:** : `integer`

##### <a name="State_ForkIDIntervals_items_Version"></a>20.3.1.4. `State.ForkIDIntervals.ForkIDIntervals items.Version`

**Type:** : `string`

##### <a name="State_ForkIDIntervals_items_BlockNumber"></a>20.3.1.5. `State.ForkIDIntervals.ForkIDIntervals items.BlockNumber`

**Type:** : `integer`

### <a name="State_MaxResourceExhaustedAttempts"></a>20.4. `State.MaxResourceExhaustedAttempts`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion

**Example setting the default value** (0):
```
[State]
MaxResourceExhaustedAttempts=0
```

### <a name="State_WaitOnResourceExhaustion"></a>20.5. `State.WaitOnResourceExhaustion`

**Title:** Duration

**Type:** : `string`

**Default:** `"0s"`

**Description:** WaitOnResourceExhaustion is the time to wait before retrying a transaction because of resource exhaustion

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

**Example setting the default value** ("0s"):
```
[State]
WaitOnResourceExhaustion="0s"
```

### <a name="State_ForkUpgradeBatchNumber"></a>20.6. `State.ForkUpgradeBatchNumber`

**Type:** : `integer`

**Default:** `0`

**Description:** Batch number from which there is a forkid change (fork upgrade)

**Example setting the default value** (0):
```
[State]
ForkUpgradeBatchNumber=0
```

### <a name="State_ForkUpgradeNewForkId"></a>20.7. `State.ForkUpgradeNewForkId`

**Type:** : `integer`

**Default:** `0`

**Description:** New fork id to be used for batches greaters than ForkUpgradeBatchNumber (fork upgrade)

**Example setting the default value** (0):
```
[State]
ForkUpgradeNewForkId=0
```

### <a name="State_DB"></a>20.8. `[State.DB]`

**Type:** : `object`
**Description:** DB is the database configuration

| Property                            | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| ----------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#State_DB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#State_DB_User )           | No      | string  | No         | -          | Database User name                                         |
| - [Password](#State_DB_Password )   | No      | string  | No         | -          | Database Password of the user                              |
| - [Host](#State_DB_Host )           | No      | string  | No         | -          | Host address of database                                   |
| - [Port](#State_DB_Port )           | No      | string  | No         | -          | Port Number of database                                    |
| - [EnableLog](#State_DB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#State_DB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

#### <a name="State_DB_Name"></a>20.8.1. `State.DB.Name`

**Type:** : `string`

**Default:** `"state_db"`

**Description:** Database name

**Example setting the default value** ("state_db"):
```
[State.DB]
Name="state_db"
```

#### <a name="State_DB_User"></a>20.8.2. `State.DB.User`

**Type:** : `string`

**Default:** `"state_user"`

**Description:** Database User name

**Example setting the default value** ("state_user"):
```
[State.DB]
User="state_user"
```

#### <a name="State_DB_Password"></a>20.8.3. `State.DB.Password`

**Type:** : `string`

**Default:** `"state_password"`

**Description:** Database Password of the user

**Example setting the default value** ("state_password"):
```
[State.DB]
Password="state_password"
```

#### <a name="State_DB_Host"></a>20.8.4. `State.DB.Host`

**Type:** : `string`

**Default:** `"zkevm-state-db"`

**Description:** Host address of database

**Example setting the default value** ("zkevm-state-db"):
```
[State.DB]
Host="zkevm-state-db"
```

#### <a name="State_DB_Port"></a>20.8.5. `State.DB.Port`

**Type:** : `string`

**Default:** `"5432"`

**Description:** Port Number of database

**Example setting the default value** ("5432"):
```
[State.DB]
Port="5432"
```

#### <a name="State_DB_EnableLog"></a>20.8.6. `State.DB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[State.DB]
EnableLog=false
```

#### <a name="State_DB_MaxConns"></a>20.8.7. `State.DB.MaxConns`

**Type:** : `integer`

**Default:** `200`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (200):
```
[State.DB]
MaxConns=200
```

### <a name="State_Batch"></a>20.9. `[State.Batch]`

**Type:** : `object`
**Description:** Configuration for the batch constraints

| Property                                   | Pattern | Type   | Deprecated | Definition | Title/Description |
| ------------------------------------------ | ------- | ------ | ---------- | ---------- | ----------------- |
| - [Constraints](#State_Batch_Constraints ) | No      | object | No         | -          | -                 |

#### <a name="State_Batch_Constraints"></a>20.9.1. `[State.Batch.Constraints]`

**Type:** : `object`

| Property                                                                 | Pattern | Type    | Deprecated | Definition | Title/Description |
| ------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ----------------- |
| - [MaxTxsPerBatch](#State_Batch_Constraints_MaxTxsPerBatch )             | No      | integer | No         | -          | -                 |
| - [MaxBatchBytesSize](#State_Batch_Constraints_MaxBatchBytesSize )       | No      | integer | No         | -          | -                 |
| - [MaxCumulativeGasUsed](#State_Batch_Constraints_MaxCumulativeGasUsed ) | No      | integer | No         | -          | -                 |
| - [MaxKeccakHashes](#State_Batch_Constraints_MaxKeccakHashes )           | No      | integer | No         | -          | -                 |
| - [MaxPoseidonHashes](#State_Batch_Constraints_MaxPoseidonHashes )       | No      | integer | No         | -          | -                 |
| - [MaxPoseidonPaddings](#State_Batch_Constraints_MaxPoseidonPaddings )   | No      | integer | No         | -          | -                 |
| - [MaxMemAligns](#State_Batch_Constraints_MaxMemAligns )                 | No      | integer | No         | -          | -                 |
| - [MaxArithmetics](#State_Batch_Constraints_MaxArithmetics )             | No      | integer | No         | -          | -                 |
| - [MaxBinaries](#State_Batch_Constraints_MaxBinaries )                   | No      | integer | No         | -          | -                 |
| - [MaxSteps](#State_Batch_Constraints_MaxSteps )                         | No      | integer | No         | -          | -                 |

##### <a name="State_Batch_Constraints_MaxTxsPerBatch"></a>20.9.1.1. `State.Batch.Constraints.MaxTxsPerBatch`

**Type:** : `integer`

**Default:** `300`

**Example setting the default value** (300):
```
[State.Batch.Constraints]
MaxTxsPerBatch=300
```

##### <a name="State_Batch_Constraints_MaxBatchBytesSize"></a>20.9.1.2. `State.Batch.Constraints.MaxBatchBytesSize`

**Type:** : `integer`

**Default:** `120000`

**Example setting the default value** (120000):
```
[State.Batch.Constraints]
MaxBatchBytesSize=120000
```

##### <a name="State_Batch_Constraints_MaxCumulativeGasUsed"></a>20.9.1.3. `State.Batch.Constraints.MaxCumulativeGasUsed`

**Type:** : `integer`

**Default:** `30000000`

**Example setting the default value** (30000000):
```
[State.Batch.Constraints]
MaxCumulativeGasUsed=30000000
```

##### <a name="State_Batch_Constraints_MaxKeccakHashes"></a>20.9.1.4. `State.Batch.Constraints.MaxKeccakHashes`

**Type:** : `integer`

**Default:** `2145`

**Example setting the default value** (2145):
```
[State.Batch.Constraints]
MaxKeccakHashes=2145
```

##### <a name="State_Batch_Constraints_MaxPoseidonHashes"></a>20.9.1.5. `State.Batch.Constraints.MaxPoseidonHashes`

**Type:** : `integer`

**Default:** `252357`

**Example setting the default value** (252357):
```
[State.Batch.Constraints]
MaxPoseidonHashes=252357
```

##### <a name="State_Batch_Constraints_MaxPoseidonPaddings"></a>20.9.1.6. `State.Batch.Constraints.MaxPoseidonPaddings`

**Type:** : `integer`

**Default:** `135191`

**Example setting the default value** (135191):
```
[State.Batch.Constraints]
MaxPoseidonPaddings=135191
```

##### <a name="State_Batch_Constraints_MaxMemAligns"></a>20.9.1.7. `State.Batch.Constraints.MaxMemAligns`

**Type:** : `integer`

**Default:** `236585`

**Example setting the default value** (236585):
```
[State.Batch.Constraints]
MaxMemAligns=236585
```

##### <a name="State_Batch_Constraints_MaxArithmetics"></a>20.9.1.8. `State.Batch.Constraints.MaxArithmetics`

**Type:** : `integer`

**Default:** `236585`

**Example setting the default value** (236585):
```
[State.Batch.Constraints]
MaxArithmetics=236585
```

##### <a name="State_Batch_Constraints_MaxBinaries"></a>20.9.1.9. `State.Batch.Constraints.MaxBinaries`

**Type:** : `integer`

**Default:** `473170`

**Example setting the default value** (473170):
```
[State.Batch.Constraints]
MaxBinaries=473170
```

##### <a name="State_Batch_Constraints_MaxSteps"></a>20.9.1.10. `State.Batch.Constraints.MaxSteps`

**Type:** : `integer`

**Default:** `7570538`

**Example setting the default value** (7570538):
```
[State.Batch.Constraints]
MaxSteps=7570538
```

### <a name="State_MaxLogsCount"></a>20.10. `State.MaxLogsCount`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxLogsCount is a configuration to set the max number of logs that can be returned
in a single call to the state, if zero it means no limit

**Example setting the default value** (0):
```
[State]
MaxLogsCount=0
```

### <a name="State_MaxLogsBlockRange"></a>20.11. `State.MaxLogsBlockRange`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs
logs in a single call to the state, if zero it means no limit

**Example setting the default value** (0):
```
[State]
MaxLogsBlockRange=0
```

### <a name="State_MaxNativeBlockHashBlockRange"></a>20.12. `State.MaxNativeBlockHashBlockRange`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying
native block hashes in a single call to the state, if zero it means no limit

**Example setting the default value** (0):
```
[State]
MaxNativeBlockHashBlockRange=0
```

----------------------------------------------------------------------------------------------------------------------------
Generated using [json-schema-for-humans](https://github.com/coveooss/json-schema-for-humans)
