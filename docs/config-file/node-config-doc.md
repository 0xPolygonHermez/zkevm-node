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
| - [StateDB](#StateDB )                               | No      | object  | No         | -          | Configuration of the state database connection                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| - [Metrics](#Metrics )                               | No      | object  | No         | -          | Configuration of the metrics service, basically is where is going to publish the metrics                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| - [EventLog](#EventLog )                             | No      | object  | No         | -          | Configuration of the event database connection                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| - [HashDB](#HashDB )                                 | No      | object  | No         | -          | Configuration of the hash database connection                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |

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

### <a name="Etherman_MultiGasProvider"></a>5.2. `Etherman.MultiGasProvider`

**Type:** : `boolean`

**Default:** `false`

**Description:** allow that L1 gas price calculation use multiples sources

**Example setting the default value** (false):
```
[Etherman]
MultiGasProvider=false
```

### <a name="Etherman_Etherscan"></a>5.3. `[Etherman.Etherscan]`

**Type:** : `object`
**Description:** Configuration for use Etherscan as used as gas provider, basically it needs the API-KEY

| Property                                | Pattern | Type   | Deprecated | Definition | Title/Description                                                                                                                     |
| --------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| - [ApiKey](#Etherman_Etherscan_ApiKey ) | No      | string | No         | -          | Need API key to use etherscan, if it's empty etherscan is not used                                                                    |
| - [Url](#Etherman_Etherscan_Url )       | No      | string | No         | -          | URL of the etherscan API. Overwritten with a hardcoded URL: "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=" |

#### <a name="Etherman_Etherscan_ApiKey"></a>5.3.1. `Etherman.Etherscan.ApiKey`

**Type:** : `string`

**Default:** `""`

**Description:** Need API key to use etherscan, if it's empty etherscan is not used

**Example setting the default value** (""):
```
[Etherman.Etherscan]
ApiKey=""
```

#### <a name="Etherman_Etherscan_Url"></a>5.3.2. `Etherman.Etherscan.Url`

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

| Property                                                        | Pattern | Type            | Deprecated | Definition | Title/Description                                                                                                                  |
| --------------------------------------------------------------- | ------- | --------------- | ---------- | ---------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| - [FrequencyToMonitorTxs](#EthTxManager_FrequencyToMonitorTxs ) | No      | string          | No         | -          | Duration                                                                                                                           |
| - [WaitTxToBeMined](#EthTxManager_WaitTxToBeMined )             | No      | string          | No         | -          | Duration                                                                                                                           |
| - [PrivateKeys](#EthTxManager_PrivateKeys )                     | No      | array of object | No         | -          | PrivateKeys defines all the key store files that are going<br />to be read in order to provide the private keys to sign the L1 txs |
| - [ForcedGas](#EthTxManager_ForcedGas )                         | No      | integer         | No         | -          | ForcedGas is the amount of gas to be forced in case of gas estimation error                                                        |

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

| Each item of this array must be                      | Description                                                                          |
| ---------------------------------------------------- | ------------------------------------------------------------------------------------ |
| [PrivateKeys items](#EthTxManager_PrivateKeys_items) | KeystoreFileConfig has all the information needed to load a private key from a k ... |

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

## <a name="RPC"></a>8. `[RPC]`

**Type:** : `object`
**Description:** Configuration for RPC service. THis one offers a extended Ethereum JSON-RPC API interface to interact with the node

| Property                                                                     | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                                                          |
| ---------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| - [Host](#RPC_Host )                                                         | No      | string  | No         | -          | Host defines the network adapter that will be used to serve the HTTP requests                                                                                                              |
| - [Port](#RPC_Port )                                                         | No      | integer | No         | -          | Port defines the port to serve the endpoints via HTTP                                                                                                                                      |
| - [ReadTimeout](#RPC_ReadTimeout )                                           | No      | string  | No         | -          | Duration                                                                                                                                                                                   |
| - [WriteTimeout](#RPC_WriteTimeout )                                         | No      | string  | No         | -          | Duration                                                                                                                                                                                   |
| - [MaxRequestsPerIPAndSecond](#RPC_MaxRequestsPerIPAndSecond )               | No      | number  | No         | -          | MaxRequestsPerIPAndSecond defines how much requests a single IP can<br />send within a single second                                                                                       |
| - [SequencerNodeURI](#RPC_SequencerNodeURI )                                 | No      | string  | No         | -          | SequencerNodeURI is used allow Non-Sequencer nodes<br />to relay transactions to the Sequencer node                                                                                        |
| - [MaxCumulativeGasUsed](#RPC_MaxCumulativeGasUsed )                         | No      | integer | No         | -          | MaxCumulativeGasUsed is the max gas allowed per batch                                                                                                                                      |
| - [WebSockets](#RPC_WebSockets )                                             | No      | object  | No         | -          | WebSockets configuration                                                                                                                                                                   |
| - [EnableL2SuggestedGasPricePolling](#RPC_EnableL2SuggestedGasPricePolling ) | No      | boolean | No         | -          | EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.                                                                          |
| - [TraceBatchUseHTTPS](#RPC_TraceBatchUseHTTPS )                             | No      | boolean | No         | -          | TraceBatchUseHTTPS enables, in the debug_traceBatchByNum endpoint, the use of the HTTPS protocol (instead of HTTP)<br />to do the parallel requests to RPC.debug_traceTransaction endpoint |

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

| Property                              | Pattern | Type    | Deprecated | Definition | Title/Description                                                           |
| ------------------------------------- | ------- | ------- | ---------- | ---------- | --------------------------------------------------------------------------- |
| - [Enabled](#RPC_WebSockets_Enabled ) | No      | boolean | No         | -          | Enabled defines if the WebSocket requests are enabled or disabled           |
| - [Host](#RPC_WebSockets_Host )       | No      | string  | No         | -          | Host defines the network adapter that will be used to serve the WS requests |
| - [Port](#RPC_WebSockets_Port )       | No      | integer | No         | -          | Port defines the port to serve the endpoints via WS                         |

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

### <a name="RPC_EnableL2SuggestedGasPricePolling"></a>8.9. `RPC.EnableL2SuggestedGasPricePolling`

**Type:** : `boolean`

**Default:** `true`

**Description:** EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.

**Example setting the default value** (true):
```
[RPC]
EnableL2SuggestedGasPricePolling=true
```

### <a name="RPC_TraceBatchUseHTTPS"></a>8.10. `RPC.TraceBatchUseHTTPS`

**Type:** : `boolean`

**Default:** `true`

**Description:** TraceBatchUseHTTPS enables, in the debug_traceBatchByNum endpoint, the use of the HTTPS protocol (instead of HTTP)
to do the parallel requests to RPC.debug_traceTransaction endpoint

**Example setting the default value** (true):
```
[RPC]
TraceBatchUseHTTPS=true
```

## <a name="Synchronizer"></a>9. `[Synchronizer]`

**Type:** : `object`
**Description:** Configuration of service `Syncrhonizer`. For this service is also really important the value of `IsTrustedSequencer`
because depending of this values is going to ask to a trusted node for trusted transactions or not

| Property                                                    | Pattern | Type    | Deprecated | Definition | Title/Description                                                        |
| ----------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------------ |
| - [SyncInterval](#Synchronizer_SyncInterval )               | No      | string  | No         | -          | Duration                                                                 |
| - [SyncChunkSize](#Synchronizer_SyncChunkSize )             | No      | integer | No         | -          | SyncChunkSize is the number of blocks to sync on each chunk              |
| - [TrustedSequencerURL](#Synchronizer_TrustedSequencerURL ) | No      | string  | No         | -          | TrustedSequencerURL is the rpc url to connect and sync the trusted state |

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

## <a name="Sequencer"></a>10. `[Sequencer]`

**Type:** : `object`
**Description:** Configuration of the sequencer service

| Property                                                                     | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                  |
| ---------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [WaitPeriodPoolIsEmpty](#Sequencer_WaitPeriodPoolIsEmpty )                 | No      | string  | No         | -          | Duration                                                                                                                                           |
| - [BlocksAmountForTxsToBeDeleted](#Sequencer_BlocksAmountForTxsToBeDeleted ) | No      | integer | No         | -          | BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool                                                       |
| - [FrequencyToCheckTxsForDelete](#Sequencer_FrequencyToCheckTxsForDelete )   | No      | string  | No         | -          | Duration                                                                                                                                           |
| - [MaxTxsPerBatch](#Sequencer_MaxTxsPerBatch )                               | No      | integer | No         | -          | MaxTxsPerBatch is the maximum amount of transactions in the batch                                                                                  |
| - [MaxBatchBytesSize](#Sequencer_MaxBatchBytesSize )                         | No      | integer | No         | -          | MaxBatchBytesSize is the maximum batch size in bytes<br />(subtracted bits of all types.Sequence fields excluding BatchL2Data from MaxTxSizeForL1) |
| - [MaxCumulativeGasUsed](#Sequencer_MaxCumulativeGasUsed )                   | No      | integer | No         | -          | MaxCumulativeGasUsed is max gas amount used by batch                                                                                               |
| - [MaxKeccakHashes](#Sequencer_MaxKeccakHashes )                             | No      | integer | No         | -          | MaxKeccakHashes is max keccak hashes used by batch                                                                                                 |
| - [MaxPoseidonHashes](#Sequencer_MaxPoseidonHashes )                         | No      | integer | No         | -          | MaxPoseidonHashes is max poseidon hashes batch can handle                                                                                          |
| - [MaxPoseidonPaddings](#Sequencer_MaxPoseidonPaddings )                     | No      | integer | No         | -          | MaxPoseidonPaddings is max poseidon paddings batch can handle                                                                                      |
| - [MaxMemAligns](#Sequencer_MaxMemAligns )                                   | No      | integer | No         | -          | MaxMemAligns is max mem aligns batch can handle                                                                                                    |
| - [MaxArithmetics](#Sequencer_MaxArithmetics )                               | No      | integer | No         | -          | MaxArithmetics is max arithmetics batch can handle                                                                                                 |
| - [MaxBinaries](#Sequencer_MaxBinaries )                                     | No      | integer | No         | -          | MaxBinaries is max binaries batch can handle                                                                                                       |
| - [MaxSteps](#Sequencer_MaxSteps )                                           | No      | integer | No         | -          | MaxSteps is max steps batch can handle                                                                                                             |
| - [WeightBatchBytesSize](#Sequencer_WeightBatchBytesSize )                   | No      | integer | No         | -          | WeightBatchBytesSize is the cost weight for the BatchBytesSize batch resource                                                                      |
| - [WeightCumulativeGasUsed](#Sequencer_WeightCumulativeGasUsed )             | No      | integer | No         | -          | WeightCumulativeGasUsed is the cost weight for the CumulativeGasUsed batch resource                                                                |
| - [WeightKeccakHashes](#Sequencer_WeightKeccakHashes )                       | No      | integer | No         | -          | WeightKeccakHashes is the cost weight for the KeccakHashes batch resource                                                                          |
| - [WeightPoseidonHashes](#Sequencer_WeightPoseidonHashes )                   | No      | integer | No         | -          | WeightPoseidonHashes is the cost weight for the PoseidonHashes batch resource                                                                      |
| - [WeightPoseidonPaddings](#Sequencer_WeightPoseidonPaddings )               | No      | integer | No         | -          | WeightPoseidonPaddings is the cost weight for the PoseidonPaddings batch resource                                                                  |
| - [WeightMemAligns](#Sequencer_WeightMemAligns )                             | No      | integer | No         | -          | WeightMemAligns is the cost weight for the MemAligns batch resource                                                                                |
| - [WeightArithmetics](#Sequencer_WeightArithmetics )                         | No      | integer | No         | -          | WeightArithmetics is the cost weight for the Arithmetics batch resource                                                                            |
| - [WeightBinaries](#Sequencer_WeightBinaries )                               | No      | integer | No         | -          | WeightBinaries is the cost weight for the Binaries batch resource                                                                                  |
| - [WeightSteps](#Sequencer_WeightSteps )                                     | No      | integer | No         | -          | WeightSteps is the cost weight for the Steps batch resource                                                                                        |
| - [TxLifetimeCheckTimeout](#Sequencer_TxLifetimeCheckTimeout )               | No      | string  | No         | -          | Duration                                                                                                                                           |
| - [MaxTxLifetime](#Sequencer_MaxTxLifetime )                                 | No      | string  | No         | -          | Duration                                                                                                                                           |
| - [Finalizer](#Sequencer_Finalizer )                                         | No      | object  | No         | -          | Finalizer's specific config properties                                                                                                             |
| - [DBManager](#Sequencer_DBManager )                                         | No      | object  | No         | -          | DBManager's specific config properties                                                                                                             |
| - [Worker](#Sequencer_Worker )                                               | No      | object  | No         | -          | Worker's specific config properties                                                                                                                |
| - [EffectiveGasPrice](#Sequencer_EffectiveGasPrice )                         | No      | object  | No         | -          | EffectiveGasPrice is the config for the gas price                                                                                                  |

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

### <a name="Sequencer_MaxTxsPerBatch"></a>10.4. `Sequencer.MaxTxsPerBatch`

**Type:** : `integer`

**Default:** `300`

**Description:** MaxTxsPerBatch is the maximum amount of transactions in the batch

**Example setting the default value** (300):
```
[Sequencer]
MaxTxsPerBatch=300
```

### <a name="Sequencer_MaxBatchBytesSize"></a>10.5. `Sequencer.MaxBatchBytesSize`

**Type:** : `integer`

**Default:** `120000`

**Description:** MaxBatchBytesSize is the maximum batch size in bytes
(subtracted bits of all types.Sequence fields excluding BatchL2Data from MaxTxSizeForL1)

**Example setting the default value** (120000):
```
[Sequencer]
MaxBatchBytesSize=120000
```

### <a name="Sequencer_MaxCumulativeGasUsed"></a>10.6. `Sequencer.MaxCumulativeGasUsed`

**Type:** : `integer`

**Default:** `30000000`

**Description:** MaxCumulativeGasUsed is max gas amount used by batch

**Example setting the default value** (30000000):
```
[Sequencer]
MaxCumulativeGasUsed=30000000
```

### <a name="Sequencer_MaxKeccakHashes"></a>10.7. `Sequencer.MaxKeccakHashes`

**Type:** : `integer`

**Default:** `2145`

**Description:** MaxKeccakHashes is max keccak hashes used by batch

**Example setting the default value** (2145):
```
[Sequencer]
MaxKeccakHashes=2145
```

### <a name="Sequencer_MaxPoseidonHashes"></a>10.8. `Sequencer.MaxPoseidonHashes`

**Type:** : `integer`

**Default:** `252357`

**Description:** MaxPoseidonHashes is max poseidon hashes batch can handle

**Example setting the default value** (252357):
```
[Sequencer]
MaxPoseidonHashes=252357
```

### <a name="Sequencer_MaxPoseidonPaddings"></a>10.9. `Sequencer.MaxPoseidonPaddings`

**Type:** : `integer`

**Default:** `135191`

**Description:** MaxPoseidonPaddings is max poseidon paddings batch can handle

**Example setting the default value** (135191):
```
[Sequencer]
MaxPoseidonPaddings=135191
```

### <a name="Sequencer_MaxMemAligns"></a>10.10. `Sequencer.MaxMemAligns`

**Type:** : `integer`

**Default:** `236585`

**Description:** MaxMemAligns is max mem aligns batch can handle

**Example setting the default value** (236585):
```
[Sequencer]
MaxMemAligns=236585
```

### <a name="Sequencer_MaxArithmetics"></a>10.11. `Sequencer.MaxArithmetics`

**Type:** : `integer`

**Default:** `236585`

**Description:** MaxArithmetics is max arithmetics batch can handle

**Example setting the default value** (236585):
```
[Sequencer]
MaxArithmetics=236585
```

### <a name="Sequencer_MaxBinaries"></a>10.12. `Sequencer.MaxBinaries`

**Type:** : `integer`

**Default:** `473170`

**Description:** MaxBinaries is max binaries batch can handle

**Example setting the default value** (473170):
```
[Sequencer]
MaxBinaries=473170
```

### <a name="Sequencer_MaxSteps"></a>10.13. `Sequencer.MaxSteps`

**Type:** : `integer`

**Default:** `7570538`

**Description:** MaxSteps is max steps batch can handle

**Example setting the default value** (7570538):
```
[Sequencer]
MaxSteps=7570538
```

### <a name="Sequencer_WeightBatchBytesSize"></a>10.14. `Sequencer.WeightBatchBytesSize`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightBatchBytesSize is the cost weight for the BatchBytesSize batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightBatchBytesSize=1
```

### <a name="Sequencer_WeightCumulativeGasUsed"></a>10.15. `Sequencer.WeightCumulativeGasUsed`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightCumulativeGasUsed is the cost weight for the CumulativeGasUsed batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightCumulativeGasUsed=1
```

### <a name="Sequencer_WeightKeccakHashes"></a>10.16. `Sequencer.WeightKeccakHashes`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightKeccakHashes is the cost weight for the KeccakHashes batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightKeccakHashes=1
```

### <a name="Sequencer_WeightPoseidonHashes"></a>10.17. `Sequencer.WeightPoseidonHashes`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightPoseidonHashes is the cost weight for the PoseidonHashes batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightPoseidonHashes=1
```

### <a name="Sequencer_WeightPoseidonPaddings"></a>10.18. `Sequencer.WeightPoseidonPaddings`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightPoseidonPaddings is the cost weight for the PoseidonPaddings batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightPoseidonPaddings=1
```

### <a name="Sequencer_WeightMemAligns"></a>10.19. `Sequencer.WeightMemAligns`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightMemAligns is the cost weight for the MemAligns batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightMemAligns=1
```

### <a name="Sequencer_WeightArithmetics"></a>10.20. `Sequencer.WeightArithmetics`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightArithmetics is the cost weight for the Arithmetics batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightArithmetics=1
```

### <a name="Sequencer_WeightBinaries"></a>10.21. `Sequencer.WeightBinaries`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightBinaries is the cost weight for the Binaries batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightBinaries=1
```

### <a name="Sequencer_WeightSteps"></a>10.22. `Sequencer.WeightSteps`

**Type:** : `integer`

**Default:** `1`

**Description:** WeightSteps is the cost weight for the Steps batch resource

**Example setting the default value** (1):
```
[Sequencer]
WeightSteps=1
```

### <a name="Sequencer_TxLifetimeCheckTimeout"></a>10.23. `Sequencer.TxLifetimeCheckTimeout`

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

### <a name="Sequencer_MaxTxLifetime"></a>10.24. `Sequencer.MaxTxLifetime`

**Title:** Duration

**Type:** : `string`

**Default:** `"3h0m0s"`

**Description:** MaxTxLifetime is the time a tx can be in the sequencer memory

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

### <a name="Sequencer_Finalizer"></a>10.25. `[Sequencer.Finalizer]`

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

#### <a name="Sequencer_Finalizer_GERDeadlineTimeout"></a>10.25.1. `Sequencer.Finalizer.GERDeadlineTimeout`

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

#### <a name="Sequencer_Finalizer_ForcedBatchDeadlineTimeout"></a>10.25.2. `Sequencer.Finalizer.ForcedBatchDeadlineTimeout`

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

#### <a name="Sequencer_Finalizer_SleepDuration"></a>10.25.3. `Sequencer.Finalizer.SleepDuration`

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

#### <a name="Sequencer_Finalizer_ResourcePercentageToCloseBatch"></a>10.25.4. `Sequencer.Finalizer.ResourcePercentageToCloseBatch`

**Type:** : `integer`

**Default:** `10`

**Description:** ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed

**Example setting the default value** (10):
```
[Sequencer.Finalizer]
ResourcePercentageToCloseBatch=10
```

#### <a name="Sequencer_Finalizer_GERFinalityNumberOfBlocks"></a>10.25.5. `Sequencer.Finalizer.GERFinalityNumberOfBlocks`

**Type:** : `integer`

**Default:** `64`

**Description:** GERFinalityNumberOfBlocks is number of blocks to consider GER final

**Example setting the default value** (64):
```
[Sequencer.Finalizer]
GERFinalityNumberOfBlocks=64
```

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingL1Timeout"></a>10.25.6. `Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingL1Timeout`

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

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingGER"></a>10.25.7. `Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingGER`

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

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingForcedBatches"></a>10.25.8. `Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingForcedBatches`

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

#### <a name="Sequencer_Finalizer_ForcedBatchesFinalityNumberOfBlocks"></a>10.25.9. `Sequencer.Finalizer.ForcedBatchesFinalityNumberOfBlocks`

**Type:** : `integer`

**Default:** `64`

**Description:** ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final

**Example setting the default value** (64):
```
[Sequencer.Finalizer]
ForcedBatchesFinalityNumberOfBlocks=64
```

#### <a name="Sequencer_Finalizer_TimestampResolution"></a>10.25.10. `Sequencer.Finalizer.TimestampResolution`

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

#### <a name="Sequencer_Finalizer_StopSequencerOnBatchNum"></a>10.25.11. `Sequencer.Finalizer.StopSequencerOnBatchNum`

**Type:** : `integer`

**Default:** `0`

**Description:** StopSequencerOnBatchNum specifies the batch number where the Sequencer will stop to process more transactions and generate new batches. The Sequencer will halt after it closes the batch equal to this number

**Example setting the default value** (0):
```
[Sequencer.Finalizer]
StopSequencerOnBatchNum=0
```

### <a name="Sequencer_DBManager"></a>10.26. `[Sequencer.DBManager]`

**Type:** : `object`
**Description:** DBManager's specific config properties

| Property                                                                     | Pattern | Type   | Deprecated | Definition | Title/Description |
| ---------------------------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| - [PoolRetrievalInterval](#Sequencer_DBManager_PoolRetrievalInterval )       | No      | string | No         | -          | Duration          |
| - [L2ReorgRetrievalInterval](#Sequencer_DBManager_L2ReorgRetrievalInterval ) | No      | string | No         | -          | Duration          |

#### <a name="Sequencer_DBManager_PoolRetrievalInterval"></a>10.26.1. `Sequencer.DBManager.PoolRetrievalInterval`

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

#### <a name="Sequencer_DBManager_L2ReorgRetrievalInterval"></a>10.26.2. `Sequencer.DBManager.L2ReorgRetrievalInterval`

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

### <a name="Sequencer_Worker"></a>10.27. `[Sequencer.Worker]`

**Type:** : `object`
**Description:** Worker's specific config properties

| Property                                                              | Pattern | Type   | Deprecated | Definition | Title/Description                                              |
| --------------------------------------------------------------------- | ------- | ------ | ---------- | ---------- | -------------------------------------------------------------- |
| - [ResourceCostMultiplier](#Sequencer_Worker_ResourceCostMultiplier ) | No      | number | No         | -          | ResourceCostMultiplier is the multiplier for the resource cost |

#### <a name="Sequencer_Worker_ResourceCostMultiplier"></a>10.27.1. `Sequencer.Worker.ResourceCostMultiplier`

**Type:** : `number`

**Default:** `1000`

**Description:** ResourceCostMultiplier is the multiplier for the resource cost

**Example setting the default value** (1000):
```
[Sequencer.Worker]
ResourceCostMultiplier=1000
```

### <a name="Sequencer_EffectiveGasPrice"></a>10.28. `[Sequencer.EffectiveGasPrice]`

**Type:** : `object`
**Description:** EffectiveGasPrice is the config for the gas price

| Property                                                                                                           | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                   |
| ------------------------------------------------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| - [MaxBreakEvenGasPriceDeviationPercentage](#Sequencer_EffectiveGasPrice_MaxBreakEvenGasPriceDeviationPercentage ) | No      | integer | No         | -          | MaxBreakEvenGasPriceDeviationPercentage is the max allowed deviation percentage BreakEvenGasPrice on re-calculation                 |
| - [L1GasPriceFactor](#Sequencer_EffectiveGasPrice_L1GasPriceFactor )                                               | No      | number  | No         | -          | L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price                                    |
| - [ByteGasCost](#Sequencer_EffectiveGasPrice_ByteGasCost )                                                         | No      | integer | No         | -          | ByteGasCost is the gas cost per byte                                                                                                |
| - [MarginFactor](#Sequencer_EffectiveGasPrice_MarginFactor )                                                       | No      | number  | No         | -          | MarginFactor is the margin factor percentage to be added to the L2 min gas price                                                    |
| - [Enabled](#Sequencer_EffectiveGasPrice_Enabled )                                                                 | No      | boolean | No         | -          | Enabled is a flag to enable/disable the effective gas price                                                                         |
| - [DefaultMinGasPriceAllowed](#Sequencer_EffectiveGasPrice_DefaultMinGasPriceAllowed )                             | No      | integer | No         | -          | DefaultMinGasPriceAllowed is the default min gas price to suggest<br />This value is assigned from [Pool].DefaultMinGasPriceAllowed |

#### <a name="Sequencer_EffectiveGasPrice_MaxBreakEvenGasPriceDeviationPercentage"></a>10.28.1. `Sequencer.EffectiveGasPrice.MaxBreakEvenGasPriceDeviationPercentage`

**Type:** : `integer`

**Default:** `10`

**Description:** MaxBreakEvenGasPriceDeviationPercentage is the max allowed deviation percentage BreakEvenGasPrice on re-calculation

**Example setting the default value** (10):
```
[Sequencer.EffectiveGasPrice]
MaxBreakEvenGasPriceDeviationPercentage=10
```

#### <a name="Sequencer_EffectiveGasPrice_L1GasPriceFactor"></a>10.28.2. `Sequencer.EffectiveGasPrice.L1GasPriceFactor`

**Type:** : `number`

**Default:** `0.25`

**Description:** L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price

**Example setting the default value** (0.25):
```
[Sequencer.EffectiveGasPrice]
L1GasPriceFactor=0.25
```

#### <a name="Sequencer_EffectiveGasPrice_ByteGasCost"></a>10.28.3. `Sequencer.EffectiveGasPrice.ByteGasCost`

**Type:** : `integer`

**Default:** `16`

**Description:** ByteGasCost is the gas cost per byte

**Example setting the default value** (16):
```
[Sequencer.EffectiveGasPrice]
ByteGasCost=16
```

#### <a name="Sequencer_EffectiveGasPrice_MarginFactor"></a>10.28.4. `Sequencer.EffectiveGasPrice.MarginFactor`

**Type:** : `number`

**Default:** `1`

**Description:** MarginFactor is the margin factor percentage to be added to the L2 min gas price

**Example setting the default value** (1):
```
[Sequencer.EffectiveGasPrice]
MarginFactor=1
```

#### <a name="Sequencer_EffectiveGasPrice_Enabled"></a>10.28.5. `Sequencer.EffectiveGasPrice.Enabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** Enabled is a flag to enable/disable the effective gas price

**Example setting the default value** (false):
```
[Sequencer.EffectiveGasPrice]
Enabled=false
```

#### <a name="Sequencer_EffectiveGasPrice_DefaultMinGasPriceAllowed"></a>10.28.6. `Sequencer.EffectiveGasPrice.DefaultMinGasPriceAllowed`

**Type:** : `integer`

**Default:** `0`

**Description:** DefaultMinGasPriceAllowed is the default min gas price to suggest
This value is assigned from [Pool].DefaultMinGasPriceAllowed

**Example setting the default value** (0):
```
[Sequencer.EffectiveGasPrice]
DefaultMinGasPriceAllowed=0
```

## <a name="SequenceSender"></a>11. `[SequenceSender]`

**Type:** : `object`
**Description:** Configuration of the sequence sender service

| Property                                                                                                | Pattern | Type            | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                  |
| ------------------------------------------------------------------------------------------------------- | ------- | --------------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| - [WaitPeriodSendSequence](#SequenceSender_WaitPeriodSendSequence )                                     | No      | string          | No         | -          | Duration                                                                                                                                                                                                                                                                                                           |
| - [LastBatchVirtualizationTimeMaxWaitPeriod](#SequenceSender_LastBatchVirtualizationTimeMaxWaitPeriod ) | No      | string          | No         | -          | Duration                                                                                                                                                                                                                                                                                                           |
| - [MaxTxSizeForL1](#SequenceSender_MaxTxSizeForL1 )                                                     | No      | integer         | No         | -          | MaxTxSizeForL1 is the maximum size a single transaction can have. This field has<br />non-trivial consequences: larger transactions than 128KB are significantly harder and<br />more expensive to propagate; larger transactions also take more resources<br />to validate whether they fit into the pool or not. |
| - [SenderAddress](#SequenceSender_SenderAddress )                                                       | No      | string          | No         | -          | SenderAddress defines which private key the eth tx manager needs to use<br />to sign the L1 txs                                                                                                                                                                                                                    |
| - [PrivateKeys](#SequenceSender_PrivateKeys )                                                           | No      | array of object | No         | -          | PrivateKeys defines all the key store files that are going<br />to be read in order to provide the private keys to sign the L1 txs                                                                                                                                                                                 |
| - [ForkUpgradeBatchNumber](#SequenceSender_ForkUpgradeBatchNumber )                                     | No      | integer         | No         | -          | Batch number where there is a forkid change (fork upgrade)                                                                                                                                                                                                                                                         |

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

**Type:** : `string`

**Default:** `"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"`

**Description:** SenderAddress defines which private key the eth tx manager needs to use
to sign the L1 txs

**Example setting the default value** ("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"):
```
[SequenceSender]
SenderAddress="0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
```

### <a name="SequenceSender_PrivateKeys"></a>11.5. `SequenceSender.PrivateKeys`

**Type:** : `array of object`

**Default:** `[{"Path": "/pk/sequencer.keystore", "Password": "testonly"}]`

**Description:** PrivateKeys defines all the key store files that are going
to be read in order to provide the private keys to sign the L1 txs

**Example setting the default value** ([{"Path": "/pk/sequencer.keystore", "Password": "testonly"}]):
```
[SequenceSender]
PrivateKeys=[{"Path": "/pk/sequencer.keystore", "Password": "testonly"}]
```

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                        | Description                                                                          |
| ------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| [PrivateKeys items](#SequenceSender_PrivateKeys_items) | KeystoreFileConfig has all the information needed to load a private key from a k ... |

#### <a name="autogenerated_heading_3"></a>11.5.1. [SequenceSender.PrivateKeys.PrivateKeys items]

**Type:** : `object`
**Description:** KeystoreFileConfig has all the information needed to load a private key from a key store file

| Property                                                  | Pattern | Type   | Deprecated | Definition | Title/Description                                      |
| --------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------ |
| - [Path](#SequenceSender_PrivateKeys_items_Path )         | No      | string | No         | -          | Path is the file path for the key store file           |
| - [Password](#SequenceSender_PrivateKeys_items_Password ) | No      | string | No         | -          | Password is the password to decrypt the key store file |

##### <a name="SequenceSender_PrivateKeys_items_Path"></a>11.5.1.1. `SequenceSender.PrivateKeys.PrivateKeys items.Path`

**Type:** : `string`
**Description:** Path is the file path for the key store file

##### <a name="SequenceSender_PrivateKeys_items_Password"></a>11.5.1.2. `SequenceSender.PrivateKeys.PrivateKeys items.Password`

**Type:** : `string`
**Description:** Password is the password to decrypt the key store file

### <a name="SequenceSender_ForkUpgradeBatchNumber"></a>11.6. `SequenceSender.ForkUpgradeBatchNumber`

**Type:** : `integer`

**Default:** `0`

**Description:** Batch number where there is a forkid change (fork upgrade)

**Example setting the default value** (0):
```
[SequenceSender]
ForkUpgradeBatchNumber=0
```

## <a name="Aggregator"></a>12. `[Aggregator]`

**Type:** : `object`
**Description:** Configuration of the aggregator service

| Property                                                                                            | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                                                                           |
| --------------------------------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [Host](#Aggregator_Host )                                                                         | No      | string  | No         | -          | Host for the grpc server                                                                                                                                                    |
| - [Port](#Aggregator_Port )                                                                         | No      | integer | No         | -          | Port for the grpc server                                                                                                                                                    |
| - [RetryTime](#Aggregator_RetryTime )                                                               | No      | string  | No         | -          | Duration                                                                                                                                                                    |
| - [VerifyProofInterval](#Aggregator_VerifyProofInterval )                                           | No      | string  | No         | -          | Duration                                                                                                                                                                    |
| - [ProofStatePollingInterval](#Aggregator_ProofStatePollingInterval )                               | No      | string  | No         | -          | Duration                                                                                                                                                                    |
| - [TxProfitabilityCheckerType](#Aggregator_TxProfitabilityCheckerType )                             | No      | string  | No         | -          | TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch<br />possible values: base/acceptall                                         |
| - [TxProfitabilityMinReward](#Aggregator_TxProfitabilityMinReward )                                 | No      | object  | No         | -          | TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch<br />this parameter is used for the base tx profitability checker |
| - [IntervalAfterWhichBatchConsolidateAnyway](#Aggregator_IntervalAfterWhichBatchConsolidateAnyway ) | No      | string  | No         | -          | Duration                                                                                                                                                                    |
| - [ChainID](#Aggregator_ChainID )                                                                   | No      | integer | No         | -          | ChainID is the L2 ChainID provided by the Network Config                                                                                                                    |
| - [ForkId](#Aggregator_ForkId )                                                                     | No      | integer | No         | -          | ForkID is the L2 ForkID provided by the Network Config                                                                                                                      |
| - [SenderAddress](#Aggregator_SenderAddress )                                                       | No      | string  | No         | -          | SenderAddress defines which private key the eth tx manager needs to use<br />to sign the L1 txs                                                                             |
| - [CleanupLockedProofsInterval](#Aggregator_CleanupLockedProofsInterval )                           | No      | string  | No         | -          | Duration                                                                                                                                                                    |
| - [GeneratingProofCleanupThreshold](#Aggregator_GeneratingProofCleanupThreshold )                   | No      | string  | No         | -          | GeneratingProofCleanupThreshold represents the time interval after<br />which a proof in generating state is considered to be stuck and<br />allowed to be cleared.         |

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

##### <a name="autogenerated_heading_4"></a>13.4.3.1. [NetworkConfig.Genesis.GenesisActions.GenesisActions items]

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

## <a name="StateDB"></a>17. `[StateDB]`

**Type:** : `object`
**Description:** Configuration of the state database connection

| Property                           | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| ---------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#StateDB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#StateDB_User )           | No      | string  | No         | -          | Database User name                                         |
| - [Password](#StateDB_Password )   | No      | string  | No         | -          | Database Password of the user                              |
| - [Host](#StateDB_Host )           | No      | string  | No         | -          | Host address of database                                   |
| - [Port](#StateDB_Port )           | No      | string  | No         | -          | Port Number of database                                    |
| - [EnableLog](#StateDB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#StateDB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

### <a name="StateDB_Name"></a>17.1. `StateDB.Name`

**Type:** : `string`

**Default:** `"state_db"`

**Description:** Database name

**Example setting the default value** ("state_db"):
```
[StateDB]
Name="state_db"
```

### <a name="StateDB_User"></a>17.2. `StateDB.User`

**Type:** : `string`

**Default:** `"state_user"`

**Description:** Database User name

**Example setting the default value** ("state_user"):
```
[StateDB]
User="state_user"
```

### <a name="StateDB_Password"></a>17.3. `StateDB.Password`

**Type:** : `string`

**Default:** `"state_password"`

**Description:** Database Password of the user

**Example setting the default value** ("state_password"):
```
[StateDB]
Password="state_password"
```

### <a name="StateDB_Host"></a>17.4. `StateDB.Host`

**Type:** : `string`

**Default:** `"zkevm-state-db"`

**Description:** Host address of database

**Example setting the default value** ("zkevm-state-db"):
```
[StateDB]
Host="zkevm-state-db"
```

### <a name="StateDB_Port"></a>17.5. `StateDB.Port`

**Type:** : `string`

**Default:** `"5432"`

**Description:** Port Number of database

**Example setting the default value** ("5432"):
```
[StateDB]
Port="5432"
```

### <a name="StateDB_EnableLog"></a>17.6. `StateDB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[StateDB]
EnableLog=false
```

### <a name="StateDB_MaxConns"></a>17.7. `StateDB.MaxConns`

**Type:** : `integer`

**Default:** `200`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (200):
```
[StateDB]
MaxConns=200
```

## <a name="Metrics"></a>18. `[Metrics]`

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

### <a name="Metrics_Host"></a>18.1. `Metrics.Host`

**Type:** : `string`

**Default:** `"0.0.0.0"`

**Description:** Host is the address to bind the metrics server

**Example setting the default value** ("0.0.0.0"):
```
[Metrics]
Host="0.0.0.0"
```

### <a name="Metrics_Port"></a>18.2. `Metrics.Port`

**Type:** : `integer`

**Default:** `9091`

**Description:** Port is the port to bind the metrics server

**Example setting the default value** (9091):
```
[Metrics]
Port=9091
```

### <a name="Metrics_Enabled"></a>18.3. `Metrics.Enabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** Enabled is the flag to enable/disable the metrics server

**Example setting the default value** (false):
```
[Metrics]
Enabled=false
```

### <a name="Metrics_ProfilingHost"></a>18.4. `Metrics.ProfilingHost`

**Type:** : `string`

**Default:** `""`

**Description:** ProfilingHost is the address to bind the profiling server

**Example setting the default value** (""):
```
[Metrics]
ProfilingHost=""
```

### <a name="Metrics_ProfilingPort"></a>18.5. `Metrics.ProfilingPort`

**Type:** : `integer`

**Default:** `0`

**Description:** ProfilingPort is the port to bind the profiling server

**Example setting the default value** (0):
```
[Metrics]
ProfilingPort=0
```

### <a name="Metrics_ProfilingEnabled"></a>18.6. `Metrics.ProfilingEnabled`

**Type:** : `boolean`

**Default:** `false`

**Description:** ProfilingEnabled is the flag to enable/disable the profiling server

**Example setting the default value** (false):
```
[Metrics]
ProfilingEnabled=false
```

## <a name="EventLog"></a>19. `[EventLog]`

**Type:** : `object`
**Description:** Configuration of the event database connection

| Property              | Pattern | Type   | Deprecated | Definition | Title/Description                |
| --------------------- | ------- | ------ | ---------- | ---------- | -------------------------------- |
| - [DB](#EventLog_DB ) | No      | object | No         | -          | DB is the database configuration |

### <a name="EventLog_DB"></a>19.1. `[EventLog.DB]`

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

#### <a name="EventLog_DB_Name"></a>19.1.1. `EventLog.DB.Name`

**Type:** : `string`

**Default:** `""`

**Description:** Database name

**Example setting the default value** (""):
```
[EventLog.DB]
Name=""
```

#### <a name="EventLog_DB_User"></a>19.1.2. `EventLog.DB.User`

**Type:** : `string`

**Default:** `""`

**Description:** Database User name

**Example setting the default value** (""):
```
[EventLog.DB]
User=""
```

#### <a name="EventLog_DB_Password"></a>19.1.3. `EventLog.DB.Password`

**Type:** : `string`

**Default:** `""`

**Description:** Database Password of the user

**Example setting the default value** (""):
```
[EventLog.DB]
Password=""
```

#### <a name="EventLog_DB_Host"></a>19.1.4. `EventLog.DB.Host`

**Type:** : `string`

**Default:** `""`

**Description:** Host address of database

**Example setting the default value** (""):
```
[EventLog.DB]
Host=""
```

#### <a name="EventLog_DB_Port"></a>19.1.5. `EventLog.DB.Port`

**Type:** : `string`

**Default:** `""`

**Description:** Port Number of database

**Example setting the default value** (""):
```
[EventLog.DB]
Port=""
```

#### <a name="EventLog_DB_EnableLog"></a>19.1.6. `EventLog.DB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[EventLog.DB]
EnableLog=false
```

#### <a name="EventLog_DB_MaxConns"></a>19.1.7. `EventLog.DB.MaxConns`

**Type:** : `integer`

**Default:** `0`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (0):
```
[EventLog.DB]
MaxConns=0
```

## <a name="HashDB"></a>20. `[HashDB]`

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

### <a name="HashDB_Name"></a>20.1. `HashDB.Name`

**Type:** : `string`

**Default:** `"prover_db"`

**Description:** Database name

**Example setting the default value** ("prover_db"):
```
[HashDB]
Name="prover_db"
```

### <a name="HashDB_User"></a>20.2. `HashDB.User`

**Type:** : `string`

**Default:** `"prover_user"`

**Description:** Database User name

**Example setting the default value** ("prover_user"):
```
[HashDB]
User="prover_user"
```

### <a name="HashDB_Password"></a>20.3. `HashDB.Password`

**Type:** : `string`

**Default:** `"prover_pass"`

**Description:** Database Password of the user

**Example setting the default value** ("prover_pass"):
```
[HashDB]
Password="prover_pass"
```

### <a name="HashDB_Host"></a>20.4. `HashDB.Host`

**Type:** : `string`

**Default:** `"zkevm-state-db"`

**Description:** Host address of database

**Example setting the default value** ("zkevm-state-db"):
```
[HashDB]
Host="zkevm-state-db"
```

### <a name="HashDB_Port"></a>20.5. `HashDB.Port`

**Type:** : `string`

**Default:** `"5432"`

**Description:** Port Number of database

**Example setting the default value** ("5432"):
```
[HashDB]
Port="5432"
```

### <a name="HashDB_EnableLog"></a>20.6. `HashDB.EnableLog`

**Type:** : `boolean`

**Default:** `false`

**Description:** EnableLog

**Example setting the default value** (false):
```
[HashDB]
EnableLog=false
```

### <a name="HashDB_MaxConns"></a>20.7. `HashDB.MaxConns`

**Type:** : `integer`

**Default:** `200`

**Description:** MaxConns is the maximum number of connections in the pool.

**Example setting the default value** (200):
```
[HashDB]
MaxConns=200
```

----------------------------------------------------------------------------------------------------------------------------
Generated using [json-schema-for-humans](https://github.com/coveooss/json-schema-for-humans)
