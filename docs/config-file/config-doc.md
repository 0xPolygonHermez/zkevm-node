# Schema Docs

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** Config represents the configuration of the entire Hermez Node The file is TOML format You could find some examples:

[TOML format]: https://en.wikipedia.org/wiki/TOML

| Property                                       | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                   |
| ---------------------------------------------- | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------- |
| - [IsTrustedSequencer](#IsTrustedSequencer )   | No      | boolean | No         | -          | This define is a trusted node (\`true\`) or a permission less (\`false\`). If you don't known<br />set to \`false\` |
| - [Log](#Log )                                 | No      | object  | No         | -          | Configure Log level for all the services, allow also to store the logs in a file                                    |
| - [Etherman](#Etherman )                       | No      | object  | No         | -          | Configure service \`Etherman\` responsible to interact with L1.                                                     |
| - [EthTxManager](#EthTxManager )               | No      | object  | No         | -          | -                                                                                                                   |
| - [Pool](#Pool )                               | No      | object  | No         | -          | -                                                                                                                   |
| - [RPC](#RPC )                                 | No      | object  | No         | -          | -                                                                                                                   |
| - [Synchronizer](#Synchronizer )               | No      | object  | No         | -          | Configuration of service \`Syncrhonizer\`. For this service is also important the value of \`IsTrustedSequencer\`   |
| - [Sequencer](#Sequencer )                     | No      | object  | No         | -          | -                                                                                                                   |
| - [SequenceSender](#SequenceSender )           | No      | object  | No         | -          | -                                                                                                                   |
| - [Aggregator](#Aggregator )                   | No      | object  | No         | -          | -                                                                                                                   |
| - [NetworkConfig](#NetworkConfig )             | No      | object  | No         | -          | -                                                                                                                   |
| - [L2GasPriceSuggester](#L2GasPriceSuggester ) | No      | object  | No         | -          | -                                                                                                                   |
| - [Executor](#Executor )                       | No      | object  | No         | -          | -                                                                                                                   |
| - [MTClient](#MTClient )                       | No      | object  | No         | -          | -                                                                                                                   |
| - [StateDB](#StateDB )                         | No      | object  | No         | -          | -                                                                                                                   |
| - [Metrics](#Metrics )                         | No      | object  | No         | -          | -                                                                                                                   |
| - [EventLog](#EventLog )                       | No      | object  | No         | -          | -                                                                                                                   |
| - [HashDB](#HashDB )                           | No      | object  | No         | -          | -                                                                                                                   |

## <a name="IsTrustedSequencer"></a>1. Property `root > IsTrustedSequencer`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

**Description:** This define is a trusted node (`true`) or a permission less (`false`). If you don't known
set to `false`

## <a name="Log"></a>2. Property `root > Log`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** Configure Log level for all the services, allow also to store the logs in a file

| Property                           | Pattern | Type             | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                                                                                               |
| ---------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| - [Environment](#Log_Environment ) | No      | enum (of string) | No         | -          | Environment defining the log format ("production" or "development").<br />In development mode enables development mode (which makes DPanicLevel logs panic), uses a console encoder, writes to standard error, and disables sampling. Stacktraces are automatically included on logs of WarnLevel and above.<br />Check [here](https://pkg.go.dev/go.uber.org/zap@v1.24.0#NewDevelopmentConfig) |
| - [Level](#Log_Level )             | No      | enum (of string) | No         | -          | Level of log. As lower value more logs are going to be generated                                                                                                                                                                                                                                                                                                                                |
| - [Outputs](#Log_Outputs )         | No      | array of string  | No         | -          | Outputs                                                                                                                                                                                                                                                                                                                                                                                         |

### <a name="Log_Environment"></a>2.1. Property `root > Log > Environment`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `enum (of string)` |
| **Required** | No                 |
| **Default**  | `"development"`    |

**Description:** Environment defining the log format ("production" or "development").
In development mode enables development mode (which makes DPanicLevel logs panic), uses a console encoder, writes to standard error, and disables sampling. Stacktraces are automatically included on logs of WarnLevel and above.
Check [here](https://pkg.go.dev/go.uber.org/zap@v1.24.0#NewDevelopmentConfig)

Must be one of:
* "production"
* "development"

### <a name="Log_Level"></a>2.2. Property `root > Log > Level`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `enum (of string)` |
| **Required** | No                 |
| **Default**  | `"info"`           |

**Description:** Level of log. As lower value more logs are going to be generated

Must be one of:
* "debug"
* "info"
* "warn"
* "error"
* "dpanic"
* "panic"
* "fatal"

### <a name="Log_Outputs"></a>2.3. Property `root > Log > Outputs`

|              |                   |
| ------------ | ----------------- |
| **Type**     | `array of string` |
| **Required** | No                |

**Description:** Outputs

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be     | Description |
| ----------------------------------- | ----------- |
| [Outputs items](#Log_Outputs_items) | -           |

#### <a name="autogenerated_heading_2"></a>2.3.1. root > Log > Outputs > Outputs items

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

## <a name="Etherman"></a>3. Property `root > Etherman`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** Configure service `Etherman` responsible to interact with L1.

| Property                                              | Pattern | Type    | Deprecated | Definition | Title/Description |
| ----------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------- |
| - [URL](#Etherman_URL )                               | No      | string  | No         | -          | -                 |
| - [PrivateKeyPath](#Etherman_PrivateKeyPath )         | No      | string  | No         | -          | -                 |
| - [PrivateKeyPassword](#Etherman_PrivateKeyPassword ) | No      | string  | No         | -          | -                 |
| - [MultiGasProvider](#Etherman_MultiGasProvider )     | No      | boolean | No         | -          | -                 |
| - [Etherscan](#Etherman_Etherscan )                   | No      | object  | No         | -          | -                 |

### <a name="Etherman_URL"></a>3.1. Property `root > Etherman > URL`

|              |                           |
| ------------ | ------------------------- |
| **Type**     | `string`                  |
| **Required** | No                        |
| **Default**  | `"http://localhost:8545"` |

### <a name="Etherman_PrivateKeyPath"></a>3.2. Property `root > Etherman > PrivateKeyPath`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

### <a name="Etherman_PrivateKeyPassword"></a>3.3. Property `root > Etherman > PrivateKeyPassword`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

### <a name="Etherman_MultiGasProvider"></a>3.4. Property `root > Etherman > MultiGasProvider`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

### <a name="Etherman_Etherscan"></a>3.5. Property `root > Etherman > Etherscan`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                | Pattern | Type   | Deprecated | Definition | Title/Description |
| --------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| - [ApiKey](#Etherman_Etherscan_ApiKey ) | No      | string | No         | -          | -                 |
| - [Url](#Etherman_Etherscan_Url )       | No      | string | No         | -          | -                 |

#### <a name="Etherman_Etherscan_ApiKey"></a>3.5.1. Property `root > Etherman > Etherscan > ApiKey`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

#### <a name="Etherman_Etherscan_Url"></a>3.5.2. Property `root > Etherman > Etherscan > Url`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

## <a name="EthTxManager"></a>4. Property `root > EthTxManager`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                        | Pattern | Type            | Deprecated | Definition | Title/Description                                                                                                                  |
| --------------------------------------------------------------- | ------- | --------------- | ---------- | ---------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| - [FrequencyToMonitorTxs](#EthTxManager_FrequencyToMonitorTxs ) | No      | string          | No         | -          | Duration                                                                                                                           |
| - [WaitTxToBeMined](#EthTxManager_WaitTxToBeMined )             | No      | string          | No         | -          | Duration                                                                                                                           |
| - [PrivateKeys](#EthTxManager_PrivateKeys )                     | No      | array of object | No         | -          | PrivateKeys defines all the key store files that are going<br />to be read in order to provide the private keys to sign the L1 txs |
| - [ForcedGas](#EthTxManager_ForcedGas )                         | No      | integer         | No         | -          | ForcedGas is the amount of gas to be forced in case of gas estimation error                                                        |

### <a name="EthTxManager_FrequencyToMonitorTxs"></a>4.1. Property `root > EthTxManager > FrequencyToMonitorTxs`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 1000000000}` |

**Description:** FrequencyToMonitorTxs frequency of the resending failed txs

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="EthTxManager_WaitTxToBeMined"></a>4.2. Property `root > EthTxManager > WaitTxToBeMined`

**Title:** Duration

|              |                              |
| ------------ | ---------------------------- |
| **Type**     | `string`                     |
| **Required** | No                           |
| **Default**  | `{"Duration": 120000000000}` |

**Description:** WaitTxToBeMined time to wait after transaction was sent to the ethereum

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="EthTxManager_PrivateKeys"></a>4.3. Property `root > EthTxManager > PrivateKeys`

|              |                   |
| ------------ | ----------------- |
| **Type**     | `array of object` |
| **Required** | No                |

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

#### <a name="autogenerated_heading_3"></a>4.3.1. root > EthTxManager > PrivateKeys > PrivateKeys items

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** KeystoreFileConfig has all the information needed to load a private key from a key store file

| Property                                                | Pattern | Type   | Deprecated | Definition | Title/Description                                      |
| ------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------ |
| - [Path](#EthTxManager_PrivateKeys_items_Path )         | No      | string | No         | -          | Path is the file path for the key store file           |
| - [Password](#EthTxManager_PrivateKeys_items_Password ) | No      | string | No         | -          | Password is the password to decrypt the key store file |

##### <a name="EthTxManager_PrivateKeys_items_Path"></a>4.3.1.1. Property `root > EthTxManager > PrivateKeys > PrivateKeys items > Path`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** Path is the file path for the key store file

##### <a name="EthTxManager_PrivateKeys_items_Password"></a>4.3.1.2. Property `root > EthTxManager > PrivateKeys > PrivateKeys items > Password`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** Password is the password to decrypt the key store file

### <a name="EthTxManager_ForcedGas"></a>4.4. Property `root > EthTxManager > ForcedGas`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

**Description:** ForcedGas is the amount of gas to be forced in case of gas estimation error

## <a name="Pool"></a>5. Property `root > Pool`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                        | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                    |
| ------------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------------------------------------------------- |
| - [IntervalToRefreshBlockedAddresses](#Pool_IntervalToRefreshBlockedAddresses ) | No      | string  | No         | -          | Duration                                                                                             |
| - [MaxTxBytesSize](#Pool_MaxTxBytesSize )                                       | No      | integer | No         | -          | MaxTxBytesSize is the max size of a transaction in bytes                                             |
| - [MaxTxDataBytesSize](#Pool_MaxTxDataBytesSize )                               | No      | integer | No         | -          | MaxTxDataBytesSize is the max size of the data field of a transaction in bytes                       |
| - [DB](#Pool_DB )                                                               | No      | object  | No         | -          | DB is the database configuration                                                                     |
| - [DefaultMinGasPriceAllowed](#Pool_DefaultMinGasPriceAllowed )                 | No      | integer | No         | -          | DefaultMinGasPriceAllowed is the default min gas price to suggest                                    |
| - [MinAllowedGasPriceInterval](#Pool_MinAllowedGasPriceInterval )               | No      | string  | No         | -          | Duration                                                                                             |
| - [PollMinAllowedGasPriceInterval](#Pool_PollMinAllowedGasPriceInterval )       | No      | string  | No         | -          | Duration                                                                                             |
| - [AccountQueue](#Pool_AccountQueue )                                           | No      | integer | No         | -          | AccountQueue represents the maximum number of non-executable transaction slots permitted per account |
| - [GlobalQueue](#Pool_GlobalQueue )                                             | No      | integer | No         | -          | GlobalQueue represents the maximum number of non-executable transaction slots for all accounts       |

### <a name="Pool_IntervalToRefreshBlockedAddresses"></a>5.1. Property `root > Pool > IntervalToRefreshBlockedAddresses`

**Title:** Duration

|              |                              |
| ------------ | ---------------------------- |
| **Type**     | `string`                     |
| **Required** | No                           |
| **Default**  | `{"Duration": 300000000000}` |

**Description:** IntervalToRefreshBlockedAddresses is the time it takes to sync the
blocked address list from db to memory

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Pool_MaxTxBytesSize"></a>5.2. Property `root > Pool > MaxTxBytesSize`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `100132`  |

**Description:** MaxTxBytesSize is the max size of a transaction in bytes

### <a name="Pool_MaxTxDataBytesSize"></a>5.3. Property `root > Pool > MaxTxDataBytesSize`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `100000`  |

**Description:** MaxTxDataBytesSize is the max size of the data field of a transaction in bytes

### <a name="Pool_DB"></a>5.4. Property `root > Pool > DB`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** DB is the database configuration

| Property                           | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| ---------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#Pool_DB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#Pool_DB_User )           | No      | string  | No         | -          | User name                                                  |
| - [Password](#Pool_DB_Password )   | No      | string  | No         | -          | Password of the user                                       |
| - [Host](#Pool_DB_Host )           | No      | string  | No         | -          | Host address                                               |
| - [Port](#Pool_DB_Port )           | No      | string  | No         | -          | Port Number                                                |
| - [EnableLog](#Pool_DB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#Pool_DB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

#### <a name="Pool_DB_Name"></a>5.4.1. Property `root > Pool > DB > Name`

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Default**  | `"pool_db"` |

**Description:** Database name

#### <a name="Pool_DB_User"></a>5.4.2. Property `root > Pool > DB > User`

|              |               |
| ------------ | ------------- |
| **Type**     | `string`      |
| **Required** | No            |
| **Default**  | `"pool_user"` |

**Description:** User name

#### <a name="Pool_DB_Password"></a>5.4.3. Property `root > Pool > DB > Password`

|              |                   |
| ------------ | ----------------- |
| **Type**     | `string`          |
| **Required** | No                |
| **Default**  | `"pool_password"` |

**Description:** Password of the user

#### <a name="Pool_DB_Host"></a>5.4.4. Property `root > Pool > DB > Host`

|              |                   |
| ------------ | ----------------- |
| **Type**     | `string`          |
| **Required** | No                |
| **Default**  | `"zkevm-pool-db"` |

**Description:** Host address

#### <a name="Pool_DB_Port"></a>5.4.5. Property `root > Pool > DB > Port`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `"5432"` |

**Description:** Port Number

#### <a name="Pool_DB_EnableLog"></a>5.4.6. Property `root > Pool > DB > EnableLog`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

**Description:** EnableLog

#### <a name="Pool_DB_MaxConns"></a>5.4.7. Property `root > Pool > DB > MaxConns`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `200`     |

**Description:** MaxConns is the maximum number of connections in the pool.

### <a name="Pool_DefaultMinGasPriceAllowed"></a>5.5. Property `root > Pool > DefaultMinGasPriceAllowed`

|              |              |
| ------------ | ------------ |
| **Type**     | `integer`    |
| **Required** | No           |
| **Default**  | `1000000000` |

**Description:** DefaultMinGasPriceAllowed is the default min gas price to suggest

### <a name="Pool_MinAllowedGasPriceInterval"></a>5.6. Property `root > Pool > MinAllowedGasPriceInterval`

**Title:** Duration

|              |                              |
| ------------ | ---------------------------- |
| **Type**     | `string`                     |
| **Required** | No                           |
| **Default**  | `{"Duration": 300000000000}` |

**Description:** MinAllowedGasPriceInterval is the interval to look back of the suggested min gas price for a tx

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Pool_PollMinAllowedGasPriceInterval"></a>5.7. Property `root > Pool > PollMinAllowedGasPriceInterval`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 15000000000}` |

**Description:** PollMinAllowedGasPriceInterval is the interval to poll the suggested min gas price for a tx

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Pool_AccountQueue"></a>5.8. Property `root > Pool > AccountQueue`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `64`      |

**Description:** AccountQueue represents the maximum number of non-executable transaction slots permitted per account

### <a name="Pool_GlobalQueue"></a>5.9. Property `root > Pool > GlobalQueue`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1024`    |

**Description:** GlobalQueue represents the maximum number of non-executable transaction slots for all accounts

## <a name="RPC"></a>6. Property `root > RPC`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                     | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                 |
| ---------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------- |
| - [Host](#RPC_Host )                                                         | No      | string  | No         | -          | Host defines the network adapter that will be used to serve the HTTP requests                                     |
| - [Port](#RPC_Port )                                                         | No      | integer | No         | -          | Port defines the port to serve the endpoints via HTTP                                                             |
| - [ReadTimeout](#RPC_ReadTimeout )                                           | No      | string  | No         | -          | Duration                                                                                                          |
| - [WriteTimeout](#RPC_WriteTimeout )                                         | No      | string  | No         | -          | Duration                                                                                                          |
| - [MaxRequestsPerIPAndSecond](#RPC_MaxRequestsPerIPAndSecond )               | No      | number  | No         | -          | MaxRequestsPerIPAndSecond defines how much requests a single IP can<br />send within a single second              |
| - [SequencerNodeURI](#RPC_SequencerNodeURI )                                 | No      | string  | No         | -          | SequencerNodeURI is used allow Non-Sequencer nodes<br />to relay transactions to the Sequencer node               |
| - [MaxCumulativeGasUsed](#RPC_MaxCumulativeGasUsed )                         | No      | integer | No         | -          | MaxCumulativeGasUsed is the max gas allowed per batch                                                             |
| - [WebSockets](#RPC_WebSockets )                                             | No      | object  | No         | -          | WebSockets configuration                                                                                          |
| - [EnableL2SuggestedGasPricePolling](#RPC_EnableL2SuggestedGasPricePolling ) | No      | boolean | No         | -          | EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price. |

### <a name="RPC_Host"></a>6.1. Property `root > RPC > Host`

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Default**  | `"0.0.0.0"` |

**Description:** Host defines the network adapter that will be used to serve the HTTP requests

### <a name="RPC_Port"></a>6.2. Property `root > RPC > Port`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `8545`    |

**Description:** Port defines the port to serve the endpoints via HTTP

### <a name="RPC_ReadTimeout"></a>6.3. Property `root > RPC > ReadTimeout`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 60000000000}` |

**Description:** ReadTimeout is the HTTP server read timeout
check net/http.server.ReadTimeout and net/http.server.ReadHeaderTimeout

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="RPC_WriteTimeout"></a>6.4. Property `root > RPC > WriteTimeout`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 60000000000}` |

**Description:** WriteTimeout is the HTTP server write timeout
check net/http.server.WriteTimeout

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="RPC_MaxRequestsPerIPAndSecond"></a>6.5. Property `root > RPC > MaxRequestsPerIPAndSecond`

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |
| **Default**  | `500`    |

**Description:** MaxRequestsPerIPAndSecond defines how much requests a single IP can
send within a single second

### <a name="RPC_SequencerNodeURI"></a>6.6. Property `root > RPC > SequencerNodeURI`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** SequencerNodeURI is used allow Non-Sequencer nodes
to relay transactions to the Sequencer node

### <a name="RPC_MaxCumulativeGasUsed"></a>6.7. Property `root > RPC > MaxCumulativeGasUsed`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

**Description:** MaxCumulativeGasUsed is the max gas allowed per batch

### <a name="RPC_WebSockets"></a>6.8. Property `root > RPC > WebSockets`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** WebSockets configuration

| Property                              | Pattern | Type    | Deprecated | Definition | Title/Description                                                           |
| ------------------------------------- | ------- | ------- | ---------- | ---------- | --------------------------------------------------------------------------- |
| - [Enabled](#RPC_WebSockets_Enabled ) | No      | boolean | No         | -          | Enabled defines if the WebSocket requests are enabled or disabled           |
| - [Host](#RPC_WebSockets_Host )       | No      | string  | No         | -          | Host defines the network adapter that will be used to serve the WS requests |
| - [Port](#RPC_WebSockets_Port )       | No      | integer | No         | -          | Port defines the port to serve the endpoints via WS                         |

#### <a name="RPC_WebSockets_Enabled"></a>6.8.1. Property `root > RPC > WebSockets > Enabled`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `true`    |

**Description:** Enabled defines if the WebSocket requests are enabled or disabled

#### <a name="RPC_WebSockets_Host"></a>6.8.2. Property `root > RPC > WebSockets > Host`

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Default**  | `"0.0.0.0"` |

**Description:** Host defines the network adapter that will be used to serve the WS requests

#### <a name="RPC_WebSockets_Port"></a>6.8.3. Property `root > RPC > WebSockets > Port`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `8546`    |

**Description:** Port defines the port to serve the endpoints via WS

### <a name="RPC_EnableL2SuggestedGasPricePolling"></a>6.9. Property `root > RPC > EnableL2SuggestedGasPricePolling`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `true`    |

**Description:** EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.

## <a name="Synchronizer"></a>7. Property `root > Synchronizer`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** Configuration of service `Syncrhonizer`. For this service is also important the value of `IsTrustedSequencer`

| Property                                                    | Pattern | Type    | Deprecated | Definition | Title/Description                                                        |
| ----------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ------------------------------------------------------------------------ |
| - [SyncInterval](#Synchronizer_SyncInterval )               | No      | string  | No         | -          | Duration                                                                 |
| - [SyncChunkSize](#Synchronizer_SyncChunkSize )             | No      | integer | No         | -          | SyncChunkSize is the number of blocks to sync on each chunk              |
| - [TrustedSequencerURL](#Synchronizer_TrustedSequencerURL ) | No      | string  | No         | -          | TrustedSequencerURL is the rpc url to connect and sync the trusted state |

### <a name="Synchronizer_SyncInterval"></a>7.1. Property `root > Synchronizer > SyncInterval`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 1000000000}` |

**Description:** SyncInterval is the delay interval between reading new rollup information

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Synchronizer_SyncChunkSize"></a>7.2. Property `root > Synchronizer > SyncChunkSize`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `100`     |

**Description:** SyncChunkSize is the number of blocks to sync on each chunk

### <a name="Synchronizer_TrustedSequencerURL"></a>7.3. Property `root > Synchronizer > TrustedSequencerURL`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** TrustedSequencerURL is the rpc url to connect and sync the trusted state

## <a name="Sequencer"></a>8. Property `root > Sequencer`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

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

### <a name="Sequencer_WaitPeriodPoolIsEmpty"></a>8.1. Property `root > Sequencer > WaitPeriodPoolIsEmpty`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 1000000000}` |

**Description:** WaitPeriodPoolIsEmpty is the time the sequencer waits until
trying to add new txs to the state

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Sequencer_BlocksAmountForTxsToBeDeleted"></a>8.2. Property `root > Sequencer > BlocksAmountForTxsToBeDeleted`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `100`     |

**Description:** BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool

### <a name="Sequencer_FrequencyToCheckTxsForDelete"></a>8.3. Property `root > Sequencer > FrequencyToCheckTxsForDelete`

**Title:** Duration

|              |                                |
| ------------ | ------------------------------ |
| **Type**     | `string`                       |
| **Required** | No                             |
| **Default**  | `{"Duration": 43200000000000}` |

**Description:** FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Sequencer_MaxTxsPerBatch"></a>8.4. Property `root > Sequencer > MaxTxsPerBatch`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `300`     |

**Description:** MaxTxsPerBatch is the maximum amount of transactions in the batch

### <a name="Sequencer_MaxBatchBytesSize"></a>8.5. Property `root > Sequencer > MaxBatchBytesSize`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `120000`  |

**Description:** MaxBatchBytesSize is the maximum batch size in bytes
(subtracted bits of all types.Sequence fields excluding BatchL2Data from MaxTxSizeForL1)

### <a name="Sequencer_MaxCumulativeGasUsed"></a>8.6. Property `root > Sequencer > MaxCumulativeGasUsed`

|              |            |
| ------------ | ---------- |
| **Type**     | `integer`  |
| **Required** | No         |
| **Default**  | `30000000` |

**Description:** MaxCumulativeGasUsed is max gas amount used by batch

### <a name="Sequencer_MaxKeccakHashes"></a>8.7. Property `root > Sequencer > MaxKeccakHashes`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `2145`    |

**Description:** MaxKeccakHashes is max keccak hashes used by batch

### <a name="Sequencer_MaxPoseidonHashes"></a>8.8. Property `root > Sequencer > MaxPoseidonHashes`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `252357`  |

**Description:** MaxPoseidonHashes is max poseidon hashes batch can handle

### <a name="Sequencer_MaxPoseidonPaddings"></a>8.9. Property `root > Sequencer > MaxPoseidonPaddings`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `135191`  |

**Description:** MaxPoseidonPaddings is max poseidon paddings batch can handle

### <a name="Sequencer_MaxMemAligns"></a>8.10. Property `root > Sequencer > MaxMemAligns`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `236585`  |

**Description:** MaxMemAligns is max mem aligns batch can handle

### <a name="Sequencer_MaxArithmetics"></a>8.11. Property `root > Sequencer > MaxArithmetics`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `236585`  |

**Description:** MaxArithmetics is max arithmetics batch can handle

### <a name="Sequencer_MaxBinaries"></a>8.12. Property `root > Sequencer > MaxBinaries`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `473170`  |

**Description:** MaxBinaries is max binaries batch can handle

### <a name="Sequencer_MaxSteps"></a>8.13. Property `root > Sequencer > MaxSteps`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `7570538` |

**Description:** MaxSteps is max steps batch can handle

### <a name="Sequencer_WeightBatchBytesSize"></a>8.14. Property `root > Sequencer > WeightBatchBytesSize`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightBatchBytesSize is the cost weight for the BatchBytesSize batch resource

### <a name="Sequencer_WeightCumulativeGasUsed"></a>8.15. Property `root > Sequencer > WeightCumulativeGasUsed`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightCumulativeGasUsed is the cost weight for the CumulativeGasUsed batch resource

### <a name="Sequencer_WeightKeccakHashes"></a>8.16. Property `root > Sequencer > WeightKeccakHashes`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightKeccakHashes is the cost weight for the KeccakHashes batch resource

### <a name="Sequencer_WeightPoseidonHashes"></a>8.17. Property `root > Sequencer > WeightPoseidonHashes`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightPoseidonHashes is the cost weight for the PoseidonHashes batch resource

### <a name="Sequencer_WeightPoseidonPaddings"></a>8.18. Property `root > Sequencer > WeightPoseidonPaddings`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightPoseidonPaddings is the cost weight for the PoseidonPaddings batch resource

### <a name="Sequencer_WeightMemAligns"></a>8.19. Property `root > Sequencer > WeightMemAligns`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightMemAligns is the cost weight for the MemAligns batch resource

### <a name="Sequencer_WeightArithmetics"></a>8.20. Property `root > Sequencer > WeightArithmetics`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightArithmetics is the cost weight for the Arithmetics batch resource

### <a name="Sequencer_WeightBinaries"></a>8.21. Property `root > Sequencer > WeightBinaries`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightBinaries is the cost weight for the Binaries batch resource

### <a name="Sequencer_WeightSteps"></a>8.22. Property `root > Sequencer > WeightSteps`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `1`       |

**Description:** WeightSteps is the cost weight for the Steps batch resource

### <a name="Sequencer_TxLifetimeCheckTimeout"></a>8.23. Property `root > Sequencer > TxLifetimeCheckTimeout`

**Title:** Duration

|              |                              |
| ------------ | ---------------------------- |
| **Type**     | `string`                     |
| **Required** | No                           |
| **Default**  | `{"Duration": 600000000000}` |

**Description:** TxLifetimeCheckTimeout is the time the sequencer waits to check txs lifetime

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Sequencer_MaxTxLifetime"></a>8.24. Property `root > Sequencer > MaxTxLifetime`

**Title:** Duration

|              |                                |
| ------------ | ------------------------------ |
| **Type**     | `string`                       |
| **Required** | No                             |
| **Default**  | `{"Duration": 10800000000000}` |

**Description:** MaxTxLifetime is the time a tx can be in the sequencer memory

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Sequencer_Finalizer"></a>8.25. Property `root > Sequencer > Finalizer`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** Finalizer's specific config properties

| Property                                                                                                                       | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                           |
| ------------------------------------------------------------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------- |
| - [GERDeadlineTimeout](#Sequencer_Finalizer_GERDeadlineTimeout )                                                               | No      | string  | No         | -          | Duration                                                                                                    |
| - [ForcedBatchDeadlineTimeout](#Sequencer_Finalizer_ForcedBatchDeadlineTimeout )                                               | No      | string  | No         | -          | Duration                                                                                                    |
| - [SleepDuration](#Sequencer_Finalizer_SleepDuration )                                                                         | No      | string  | No         | -          | Duration                                                                                                    |
| - [ResourcePercentageToCloseBatch](#Sequencer_Finalizer_ResourcePercentageToCloseBatch )                                       | No      | integer | No         | -          | ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed |
| - [GERFinalityNumberOfBlocks](#Sequencer_Finalizer_GERFinalityNumberOfBlocks )                                                 | No      | integer | No         | -          | GERFinalityNumberOfBlocks is number of blocks to consider GER final                                         |
| - [ClosingSignalsManagerWaitForCheckingL1Timeout](#Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingL1Timeout )         | No      | string  | No         | -          | Duration                                                                                                    |
| - [ClosingSignalsManagerWaitForCheckingGER](#Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingGER )                     | No      | string  | No         | -          | Duration                                                                                                    |
| - [ClosingSignalsManagerWaitForCheckingForcedBatches](#Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingForcedBatches ) | No      | string  | No         | -          | Duration                                                                                                    |
| - [ForcedBatchesFinalityNumberOfBlocks](#Sequencer_Finalizer_ForcedBatchesFinalityNumberOfBlocks )                             | No      | integer | No         | -          | ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final                               |
| - [TimestampResolution](#Sequencer_Finalizer_TimestampResolution )                                                             | No      | string  | No         | -          | Duration                                                                                                    |

#### <a name="Sequencer_Finalizer_GERDeadlineTimeout"></a>8.25.1. Property `root > Sequencer > Finalizer > GERDeadlineTimeout`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 5000000000}` |

**Description:** GERDeadlineTimeout is the time the finalizer waits after receiving closing signal to update Global Exit Root

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_Finalizer_ForcedBatchDeadlineTimeout"></a>8.25.2. Property `root > Sequencer > Finalizer > ForcedBatchDeadlineTimeout`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 60000000000}` |

**Description:** ForcedBatchDeadlineTimeout is the time the finalizer waits after receiving closing signal to process Forced Batches

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_Finalizer_SleepDuration"></a>8.25.3. Property `root > Sequencer > Finalizer > SleepDuration`

**Title:** Duration

|              |                           |
| ------------ | ------------------------- |
| **Type**     | `string`                  |
| **Required** | No                        |
| **Default**  | `{"Duration": 100000000}` |

**Description:** SleepDuration is the time the finalizer sleeps between each iteration, if there are no transactions to be processed

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_Finalizer_ResourcePercentageToCloseBatch"></a>8.25.4. Property `root > Sequencer > Finalizer > ResourcePercentageToCloseBatch`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `10`      |

**Description:** ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed

#### <a name="Sequencer_Finalizer_GERFinalityNumberOfBlocks"></a>8.25.5. Property `root > Sequencer > Finalizer > GERFinalityNumberOfBlocks`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `64`      |

**Description:** GERFinalityNumberOfBlocks is number of blocks to consider GER final

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingL1Timeout"></a>8.25.6. Property `root > Sequencer > Finalizer > ClosingSignalsManagerWaitForCheckingL1Timeout`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 10000000000}` |

**Description:** ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingGER"></a>8.25.7. Property `root > Sequencer > Finalizer > ClosingSignalsManagerWaitForCheckingGER`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 10000000000}` |

**Description:** ClosingSignalsManagerWaitForCheckingGER is used by the closing signals manager to wait for its operation

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_Finalizer_ClosingSignalsManagerWaitForCheckingForcedBatches"></a>8.25.8. Property `root > Sequencer > Finalizer > ClosingSignalsManagerWaitForCheckingForcedBatches`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 10000000000}` |

**Description:** ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_Finalizer_ForcedBatchesFinalityNumberOfBlocks"></a>8.25.9. Property `root > Sequencer > Finalizer > ForcedBatchesFinalityNumberOfBlocks`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `64`      |

**Description:** ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final

#### <a name="Sequencer_Finalizer_TimestampResolution"></a>8.25.10. Property `root > Sequencer > Finalizer > TimestampResolution`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 10000000000}` |

**Description:** TimestampResolution is the resolution of the timestamp used to close a batch

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Sequencer_DBManager"></a>8.26. Property `root > Sequencer > DBManager`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** DBManager's specific config properties

| Property                                                                     | Pattern | Type   | Deprecated | Definition | Title/Description |
| ---------------------------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| - [PoolRetrievalInterval](#Sequencer_DBManager_PoolRetrievalInterval )       | No      | string | No         | -          | Duration          |
| - [L2ReorgRetrievalInterval](#Sequencer_DBManager_L2ReorgRetrievalInterval ) | No      | string | No         | -          | Duration          |

#### <a name="Sequencer_DBManager_PoolRetrievalInterval"></a>8.26.1. Property `root > Sequencer > DBManager > PoolRetrievalInterval`

**Title:** Duration

|              |                           |
| ------------ | ------------------------- |
| **Type**     | `string`                  |
| **Required** | No                        |
| **Default**  | `{"Duration": 500000000}` |

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

#### <a name="Sequencer_DBManager_L2ReorgRetrievalInterval"></a>8.26.2. Property `root > Sequencer > DBManager > L2ReorgRetrievalInterval`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 5000000000}` |

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Sequencer_Worker"></a>8.27. Property `root > Sequencer > Worker`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** Worker's specific config properties

| Property                                                              | Pattern | Type   | Deprecated | Definition | Title/Description                                              |
| --------------------------------------------------------------------- | ------- | ------ | ---------- | ---------- | -------------------------------------------------------------- |
| - [ResourceCostMultiplier](#Sequencer_Worker_ResourceCostMultiplier ) | No      | number | No         | -          | ResourceCostMultiplier is the multiplier for the resource cost |

#### <a name="Sequencer_Worker_ResourceCostMultiplier"></a>8.27.1. Property `root > Sequencer > Worker > ResourceCostMultiplier`

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |
| **Default**  | `1000`   |

**Description:** ResourceCostMultiplier is the multiplier for the resource cost

## <a name="SequenceSender"></a>9. Property `root > SequenceSender`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                                                | Pattern | Type            | Deprecated | Definition | Title/Description                                                                                                                                                                                                                                                                                                  |
| ------------------------------------------------------------------------------------------------------- | ------- | --------------- | ---------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| - [WaitPeriodSendSequence](#SequenceSender_WaitPeriodSendSequence )                                     | No      | string          | No         | -          | Duration                                                                                                                                                                                                                                                                                                           |
| - [LastBatchVirtualizationTimeMaxWaitPeriod](#SequenceSender_LastBatchVirtualizationTimeMaxWaitPeriod ) | No      | string          | No         | -          | Duration                                                                                                                                                                                                                                                                                                           |
| - [MaxTxSizeForL1](#SequenceSender_MaxTxSizeForL1 )                                                     | No      | integer         | No         | -          | MaxTxSizeForL1 is the maximum size a single transaction can have. This field has<br />non-trivial consequences: larger transactions than 128KB are significantly harder and<br />more expensive to propagate; larger transactions also take more resources<br />to validate whether they fit into the pool or not. |
| - [SenderAddress](#SequenceSender_SenderAddress )                                                       | No      | string          | No         | -          | SenderAddress defines which private key the eth tx manager needs to use<br />to sign the L1 txs                                                                                                                                                                                                                    |
| - [PrivateKeys](#SequenceSender_PrivateKeys )                                                           | No      | array of object | No         | -          | PrivateKeys defines all the key store files that are going<br />to be read in order to provide the private keys to sign the L1 txs                                                                                                                                                                                 |

### <a name="SequenceSender_WaitPeriodSendSequence"></a>9.1. Property `root > SequenceSender > WaitPeriodSendSequence`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 5000000000}` |

**Description:** WaitPeriodSendSequence is the time the sequencer waits until
trying to send a sequence to L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="SequenceSender_LastBatchVirtualizationTimeMaxWaitPeriod"></a>9.2. Property `root > SequenceSender > LastBatchVirtualizationTimeMaxWaitPeriod`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 5000000000}` |

**Description:** LastBatchVirtualizationTimeMaxWaitPeriod is time since sequences should be sent

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="SequenceSender_MaxTxSizeForL1"></a>9.3. Property `root > SequenceSender > MaxTxSizeForL1`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `131072`  |

**Description:** MaxTxSizeForL1 is the maximum size a single transaction can have. This field has
non-trivial consequences: larger transactions than 128KB are significantly harder and
more expensive to propagate; larger transactions also take more resources
to validate whether they fit into the pool or not.

### <a name="SequenceSender_SenderAddress"></a>9.4. Property `root > SequenceSender > SenderAddress`

|              |                                                |
| ------------ | ---------------------------------------------- |
| **Type**     | `string`                                       |
| **Required** | No                                             |
| **Default**  | `"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"` |

**Description:** SenderAddress defines which private key the eth tx manager needs to use
to sign the L1 txs

### <a name="SequenceSender_PrivateKeys"></a>9.5. Property `root > SequenceSender > PrivateKeys`

|              |                   |
| ------------ | ----------------- |
| **Type**     | `array of object` |
| **Required** | No                |

**Description:** PrivateKeys defines all the key store files that are going
to be read in order to provide the private keys to sign the L1 txs

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

#### <a name="autogenerated_heading_4"></a>9.5.1. root > SequenceSender > PrivateKeys > PrivateKeys items

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** KeystoreFileConfig has all the information needed to load a private key from a key store file

| Property                                                  | Pattern | Type   | Deprecated | Definition | Title/Description                                      |
| --------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ------------------------------------------------------ |
| - [Path](#SequenceSender_PrivateKeys_items_Path )         | No      | string | No         | -          | Path is the file path for the key store file           |
| - [Password](#SequenceSender_PrivateKeys_items_Password ) | No      | string | No         | -          | Password is the password to decrypt the key store file |

##### <a name="SequenceSender_PrivateKeys_items_Path"></a>9.5.1.1. Property `root > SequenceSender > PrivateKeys > PrivateKeys items > Path`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** Path is the file path for the key store file

##### <a name="SequenceSender_PrivateKeys_items_Password"></a>9.5.1.2. Property `root > SequenceSender > PrivateKeys > PrivateKeys items > Password`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

**Description:** Password is the password to decrypt the key store file

## <a name="Aggregator"></a>10. Property `root > Aggregator`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

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

### <a name="Aggregator_Host"></a>10.1. Property `root > Aggregator > Host`

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Default**  | `"0.0.0.0"` |

**Description:** Host for the grpc server

### <a name="Aggregator_Port"></a>10.2. Property `root > Aggregator > Port`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `50081`   |

**Description:** Port for the grpc server

### <a name="Aggregator_RetryTime"></a>10.3. Property `root > Aggregator > RetryTime`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 5000000000}` |

**Description:** RetryTime is the time the aggregator main loop sleeps if there are no proofs to aggregate
or batches to generate proofs. It is also used in the isSynced loop

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Aggregator_VerifyProofInterval"></a>10.4. Property `root > Aggregator > VerifyProofInterval`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 90000000000}` |

**Description:** VerifyProofInterval is the interval of time to verify/send an proof in L1

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Aggregator_ProofStatePollingInterval"></a>10.5. Property `root > Aggregator > ProofStatePollingInterval`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 5000000000}` |

**Description:** ProofStatePollingInterval is the interval time to polling the prover about the generation state of a proof

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Aggregator_TxProfitabilityCheckerType"></a>10.6. Property `root > Aggregator > TxProfitabilityCheckerType`

|              |               |
| ------------ | ------------- |
| **Type**     | `string`      |
| **Required** | No            |
| **Default**  | `"acceptall"` |

**Description:** TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
possible values: base/acceptall

### <a name="Aggregator_TxProfitabilityMinReward"></a>10.7. Property `root > Aggregator > TxProfitabilityMinReward`

|                           |                                                                           |
| ------------------------- | ------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                  |
| **Required**              | No                                                                        |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |

**Description:** TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
this parameter is used for the base tx profitability checker

### <a name="Aggregator_IntervalAfterWhichBatchConsolidateAnyway"></a>10.8. Property `root > Aggregator > IntervalAfterWhichBatchConsolidateAnyway`

**Title:** Duration

|              |                   |
| ------------ | ----------------- |
| **Type**     | `string`          |
| **Required** | No                |
| **Default**  | `{"Duration": 0}` |

**Description:** IntervalAfterWhichBatchConsolidateAnyway this is interval for the main sequencer, that will check if there is no transactions

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Aggregator_ChainID"></a>10.9. Property `root > Aggregator > ChainID`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

**Description:** ChainID is the L2 ChainID provided by the Network Config

### <a name="Aggregator_ForkId"></a>10.10. Property `root > Aggregator > ForkId`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `2`       |

**Description:** ForkID is the L2 ForkID provided by the Network Config

### <a name="Aggregator_SenderAddress"></a>10.11. Property `root > Aggregator > SenderAddress`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** SenderAddress defines which private key the eth tx manager needs to use
to sign the L1 txs

### <a name="Aggregator_CleanupLockedProofsInterval"></a>10.12. Property `root > Aggregator > CleanupLockedProofsInterval`

**Title:** Duration

|              |                              |
| ------------ | ---------------------------- |
| **Type**     | `string`                     |
| **Required** | No                           |
| **Default**  | `{"Duration": 120000000000}` |

**Description:** CleanupLockedProofsInterval is the interval of time to clean up locked proofs.

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Aggregator_GeneratingProofCleanupThreshold"></a>10.13. Property `root > Aggregator > GeneratingProofCleanupThreshold`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `"10m"`  |

**Description:** GeneratingProofCleanupThreshold represents the time interval after
which a proof in generating state is considered to be stuck and
allowed to be cleared.

## <a name="NetworkConfig"></a>11. Property `root > NetworkConfig`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                     | Pattern | Type             | Deprecated | Definition | Title/Description |
| ---------------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------- |
| - [l1Config](#NetworkConfig_l1Config )                                       | No      | object           | No         | -          | -                 |
| - [L2GlobalExitRootManagerAddr](#NetworkConfig_L2GlobalExitRootManagerAddr ) | No      | array of integer | No         | -          | -                 |
| - [L2BridgeAddr](#NetworkConfig_L2BridgeAddr )                               | No      | array of integer | No         | -          | -                 |
| - [Genesis](#NetworkConfig_Genesis )                                         | No      | object           | No         | -          | -                 |
| - [MaxCumulativeGasUsed](#NetworkConfig_MaxCumulativeGasUsed )               | No      | integer          | No         | -          | -                 |

### <a name="NetworkConfig_l1Config"></a>11.1. Property `root > NetworkConfig > l1Config`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                                          | Pattern | Type             | Deprecated | Definition | Title/Description |
| ------------------------------------------------------------------------------------------------- | ------- | ---------------- | ---------- | ---------- | ----------------- |
| - [chainId](#NetworkConfig_l1Config_chainId )                                                     | No      | integer          | No         | -          | -                 |
| - [polygonZkEVMAddress](#NetworkConfig_l1Config_polygonZkEVMAddress )                             | No      | array of integer | No         | -          | -                 |
| - [maticTokenAddress](#NetworkConfig_l1Config_maticTokenAddress )                                 | No      | array of integer | No         | -          | -                 |
| - [polygonZkEVMGlobalExitRootAddress](#NetworkConfig_l1Config_polygonZkEVMGlobalExitRootAddress ) | No      | array of integer | No         | -          | -                 |

#### <a name="NetworkConfig_l1Config_chainId"></a>11.1.1. Property `root > NetworkConfig > l1Config > chainId`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

#### <a name="NetworkConfig_l1Config_polygonZkEVMAddress"></a>11.1.2. Property `root > NetworkConfig > l1Config > polygonZkEVMAddress`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `array of integer` |
| **Required** | No                 |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | 20                 |
| **Max items**        | 20                 |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                                                | Description |
| ------------------------------------------------------------------------------ | ----------- |
| [polygonZkEVMAddress items](#NetworkConfig_l1Config_polygonZkEVMAddress_items) | -           |

##### <a name="autogenerated_heading_5"></a>11.1.2.1. root > NetworkConfig > l1Config > polygonZkEVMAddress > polygonZkEVMAddress items

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

#### <a name="NetworkConfig_l1Config_maticTokenAddress"></a>11.1.3. Property `root > NetworkConfig > l1Config > maticTokenAddress`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `array of integer` |
| **Required** | No                 |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | 20                 |
| **Max items**        | 20                 |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                                            | Description |
| -------------------------------------------------------------------------- | ----------- |
| [maticTokenAddress items](#NetworkConfig_l1Config_maticTokenAddress_items) | -           |

##### <a name="autogenerated_heading_6"></a>11.1.3.1. root > NetworkConfig > l1Config > maticTokenAddress > maticTokenAddress items

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

#### <a name="NetworkConfig_l1Config_polygonZkEVMGlobalExitRootAddress"></a>11.1.4. Property `root > NetworkConfig > l1Config > polygonZkEVMGlobalExitRootAddress`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `array of integer` |
| **Required** | No                 |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | 20                 |
| **Max items**        | 20                 |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                                                                            | Description |
| ---------------------------------------------------------------------------------------------------------- | ----------- |
| [polygonZkEVMGlobalExitRootAddress items](#NetworkConfig_l1Config_polygonZkEVMGlobalExitRootAddress_items) | -           |

##### <a name="autogenerated_heading_7"></a>11.1.4.1. root > NetworkConfig > l1Config > polygonZkEVMGlobalExitRootAddress > polygonZkEVMGlobalExitRootAddress items

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

### <a name="NetworkConfig_L2GlobalExitRootManagerAddr"></a>11.2. Property `root > NetworkConfig > L2GlobalExitRootManagerAddr`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `array of integer` |
| **Required** | No                 |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | 20                 |
| **Max items**        | 20                 |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                                                       | Description |
| ------------------------------------------------------------------------------------- | ----------- |
| [L2GlobalExitRootManagerAddr items](#NetworkConfig_L2GlobalExitRootManagerAddr_items) | -           |

#### <a name="autogenerated_heading_8"></a>11.2.1. root > NetworkConfig > L2GlobalExitRootManagerAddr > L2GlobalExitRootManagerAddr items

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

### <a name="NetworkConfig_L2BridgeAddr"></a>11.3. Property `root > NetworkConfig > L2BridgeAddr`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `array of integer` |
| **Required** | No                 |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | 20                 |
| **Max items**        | 20                 |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                         | Description |
| ------------------------------------------------------- | ----------- |
| [L2BridgeAddr items](#NetworkConfig_L2BridgeAddr_items) | -           |

#### <a name="autogenerated_heading_9"></a>11.3.1. root > NetworkConfig > L2BridgeAddr > L2BridgeAddr items

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

### <a name="NetworkConfig_Genesis"></a>11.4. Property `root > NetworkConfig > Genesis`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                     | Pattern | Type             | Deprecated | Definition | Title/Description                                                           |
| ------------------------------------------------------------ | ------- | ---------------- | ---------- | ---------- | --------------------------------------------------------------------------- |
| - [GenesisBlockNum](#NetworkConfig_Genesis_GenesisBlockNum ) | No      | integer          | No         | -          | GenesisBlockNum is the block number where the polygonZKEVM smc was deployed |
| - [Root](#NetworkConfig_Genesis_Root )                       | No      | array of integer | No         | -          | -                                                                           |
| - [GenesisActions](#NetworkConfig_Genesis_GenesisActions )   | No      | array of object  | No         | -          | -                                                                           |

#### <a name="NetworkConfig_Genesis_GenesisBlockNum"></a>11.4.1. Property `root > NetworkConfig > Genesis > GenesisBlockNum`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

**Description:** GenesisBlockNum is the block number where the polygonZKEVM smc was deployed

#### <a name="NetworkConfig_Genesis_Root"></a>11.4.2. Property `root > NetworkConfig > Genesis > Root`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `array of integer` |
| **Required** | No                 |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | 32                 |
| **Max items**        | 32                 |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |

| Each item of this array must be                 | Description |
| ----------------------------------------------- | ----------- |
| [Root items](#NetworkConfig_Genesis_Root_items) | -           |

##### <a name="autogenerated_heading_10"></a>11.4.2.1. root > NetworkConfig > Genesis > Root > Root items

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

#### <a name="NetworkConfig_Genesis_GenesisActions"></a>11.4.3. Property `root > NetworkConfig > Genesis > GenesisActions`

|              |                   |
| ------------ | ----------------- |
| **Type**     | `array of object` |
| **Required** | No                |

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

##### <a name="autogenerated_heading_11"></a>11.4.3.1. root > NetworkConfig > Genesis > GenesisActions > GenesisActions items

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

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

##### <a name="NetworkConfig_Genesis_GenesisActions_items_address"></a>11.4.3.1.1. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > address`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_type"></a>11.4.3.1.2. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > type`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_storagePosition"></a>11.4.3.1.3. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > storagePosition`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_bytecode"></a>11.4.3.1.4. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > bytecode`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_key"></a>11.4.3.1.5. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > key`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_value"></a>11.4.3.1.6. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > value`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

##### <a name="NetworkConfig_Genesis_GenesisActions_items_root"></a>11.4.3.1.7. Property `root > NetworkConfig > Genesis > GenesisActions > GenesisActions items > root`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |

### <a name="NetworkConfig_MaxCumulativeGasUsed"></a>11.5. Property `root > NetworkConfig > MaxCumulativeGasUsed`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

## <a name="L2GasPriceSuggester"></a>12. Property `root > L2GasPriceSuggester`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                       | Pattern | Type    | Deprecated | Definition | Title/Description |
| ------------------------------------------------------------------------------ | ------- | ------- | ---------- | ---------- | ----------------- |
| - [Type](#L2GasPriceSuggester_Type )                                           | No      | string  | No         | -          | -                 |
| - [DefaultGasPriceWei](#L2GasPriceSuggester_DefaultGasPriceWei )               | No      | integer | No         | -          | -                 |
| - [MaxPrice](#L2GasPriceSuggester_MaxPrice )                                   | No      | object  | No         | -          | -                 |
| - [IgnorePrice](#L2GasPriceSuggester_IgnorePrice )                             | No      | object  | No         | -          | -                 |
| - [CheckBlocks](#L2GasPriceSuggester_CheckBlocks )                             | No      | integer | No         | -          | -                 |
| - [Percentile](#L2GasPriceSuggester_Percentile )                               | No      | integer | No         | -          | -                 |
| - [UpdatePeriod](#L2GasPriceSuggester_UpdatePeriod )                           | No      | string  | No         | -          | Duration          |
| - [CleanHistoryPeriod](#L2GasPriceSuggester_CleanHistoryPeriod )               | No      | string  | No         | -          | Duration          |
| - [CleanHistoryTimeRetention](#L2GasPriceSuggester_CleanHistoryTimeRetention ) | No      | string  | No         | -          | Duration          |
| - [Factor](#L2GasPriceSuggester_Factor )                                       | No      | number  | No         | -          | -                 |

### <a name="L2GasPriceSuggester_Type"></a>12.1. Property `root > L2GasPriceSuggester > Type`

|              |              |
| ------------ | ------------ |
| **Type**     | `string`     |
| **Required** | No           |
| **Default**  | `"follower"` |

### <a name="L2GasPriceSuggester_DefaultGasPriceWei"></a>12.2. Property `root > L2GasPriceSuggester > DefaultGasPriceWei`

|              |              |
| ------------ | ------------ |
| **Type**     | `integer`    |
| **Required** | No           |
| **Default**  | `2000000000` |

### <a name="L2GasPriceSuggester_MaxPrice"></a>12.3. Property `root > L2GasPriceSuggester > MaxPrice`

|                           |                                                                           |
| ------------------------- | ------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                  |
| **Required**              | No                                                                        |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |

### <a name="L2GasPriceSuggester_IgnorePrice"></a>12.4. Property `root > L2GasPriceSuggester > IgnorePrice`

|                           |                                                                           |
| ------------------------- | ------------------------------------------------------------------------- |
| **Type**                  | `object`                                                                  |
| **Required**              | No                                                                        |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |

### <a name="L2GasPriceSuggester_CheckBlocks"></a>12.5. Property `root > L2GasPriceSuggester > CheckBlocks`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

### <a name="L2GasPriceSuggester_Percentile"></a>12.6. Property `root > L2GasPriceSuggester > Percentile`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

### <a name="L2GasPriceSuggester_UpdatePeriod"></a>12.7. Property `root > L2GasPriceSuggester > UpdatePeriod`

**Title:** Duration

|              |                             |
| ------------ | --------------------------- |
| **Type**     | `string`                    |
| **Required** | No                          |
| **Default**  | `{"Duration": 10000000000}` |

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="L2GasPriceSuggester_CleanHistoryPeriod"></a>12.8. Property `root > L2GasPriceSuggester > CleanHistoryPeriod`

**Title:** Duration

|              |                               |
| ------------ | ----------------------------- |
| **Type**     | `string`                      |
| **Required** | No                            |
| **Default**  | `{"Duration": 3600000000000}` |

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="L2GasPriceSuggester_CleanHistoryTimeRetention"></a>12.9. Property `root > L2GasPriceSuggester > CleanHistoryTimeRetention`

**Title:** Duration

|              |                              |
| ------------ | ---------------------------- |
| **Type**     | `string`                     |
| **Required** | No                           |
| **Default**  | `{"Duration": 300000000000}` |

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="L2GasPriceSuggester_Factor"></a>12.10. Property `root > L2GasPriceSuggester > Factor`

|              |          |
| ------------ | -------- |
| **Type**     | `number` |
| **Required** | No       |
| **Default**  | `0.15`   |

## <a name="Executor"></a>13. Property `root > Executor`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                                                  | Pattern | Type    | Deprecated | Definition | Title/Description                                                                                                       |
| ------------------------------------------------------------------------- | ------- | ------- | ---------- | ---------- | ----------------------------------------------------------------------------------------------------------------------- |
| - [URI](#Executor_URI )                                                   | No      | string  | No         | -          | -                                                                                                                       |
| - [MaxResourceExhaustedAttempts](#Executor_MaxResourceExhaustedAttempts ) | No      | integer | No         | -          | MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion |
| - [WaitOnResourceExhaustion](#Executor_WaitOnResourceExhaustion )         | No      | string  | No         | -          | Duration                                                                                                                |
| - [MaxGRPCMessageSize](#Executor_MaxGRPCMessageSize )                     | No      | integer | No         | -          | -                                                                                                                       |

### <a name="Executor_URI"></a>13.1. Property `root > Executor > URI`

|              |                        |
| ------------ | ---------------------- |
| **Type**     | `string`               |
| **Required** | No                     |
| **Default**  | `"zkevm-prover:50071"` |

### <a name="Executor_MaxResourceExhaustedAttempts"></a>13.2. Property `root > Executor > MaxResourceExhaustedAttempts`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `3`       |

**Description:** MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion

### <a name="Executor_WaitOnResourceExhaustion"></a>13.3. Property `root > Executor > WaitOnResourceExhaustion`

**Title:** Duration

|              |                            |
| ------------ | -------------------------- |
| **Type**     | `string`                   |
| **Required** | No                         |
| **Default**  | `{"Duration": 1000000000}` |

**Description:** WaitOnResourceExhaustion is the time to wait before retrying a transaction because of resource exhaustion

**Examples:** 

```json
"1m"
```

```json
"300ms"
```

### <a name="Executor_MaxGRPCMessageSize"></a>13.4. Property `root > Executor > MaxGRPCMessageSize`

|              |             |
| ------------ | ----------- |
| **Type**     | `integer`   |
| **Required** | No          |
| **Default**  | `100000000` |

## <a name="MTClient"></a>14. Property `root > MTClient`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                | Pattern | Type   | Deprecated | Definition | Title/Description      |
| ----------------------- | ------- | ------ | ---------- | ---------- | ---------------------- |
| - [URI](#MTClient_URI ) | No      | string | No         | -          | URI is the server URI. |

### <a name="MTClient_URI"></a>14.1. Property `root > MTClient > URI`

|              |                        |
| ------------ | ---------------------- |
| **Type**     | `string`               |
| **Required** | No                     |
| **Default**  | `"zkevm-prover:50061"` |

**Description:** URI is the server URI.

## <a name="StateDB"></a>15. Property `root > StateDB`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                           | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| ---------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#StateDB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#StateDB_User )           | No      | string  | No         | -          | User name                                                  |
| - [Password](#StateDB_Password )   | No      | string  | No         | -          | Password of the user                                       |
| - [Host](#StateDB_Host )           | No      | string  | No         | -          | Host address                                               |
| - [Port](#StateDB_Port )           | No      | string  | No         | -          | Port Number                                                |
| - [EnableLog](#StateDB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#StateDB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

### <a name="StateDB_Name"></a>15.1. Property `root > StateDB > Name`

|              |              |
| ------------ | ------------ |
| **Type**     | `string`     |
| **Required** | No           |
| **Default**  | `"state_db"` |

**Description:** Database name

### <a name="StateDB_User"></a>15.2. Property `root > StateDB > User`

|              |                |
| ------------ | -------------- |
| **Type**     | `string`       |
| **Required** | No             |
| **Default**  | `"state_user"` |

**Description:** User name

### <a name="StateDB_Password"></a>15.3. Property `root > StateDB > Password`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `string`           |
| **Required** | No                 |
| **Default**  | `"state_password"` |

**Description:** Password of the user

### <a name="StateDB_Host"></a>15.4. Property `root > StateDB > Host`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `string`           |
| **Required** | No                 |
| **Default**  | `"zkevm-state-db"` |

**Description:** Host address

### <a name="StateDB_Port"></a>15.5. Property `root > StateDB > Port`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `"5432"` |

**Description:** Port Number

### <a name="StateDB_EnableLog"></a>15.6. Property `root > StateDB > EnableLog`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

**Description:** EnableLog

### <a name="StateDB_MaxConns"></a>15.7. Property `root > StateDB > MaxConns`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `200`     |

**Description:** MaxConns is the maximum number of connections in the pool.

## <a name="Metrics"></a>16. Property `root > Metrics`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                                         | Pattern | Type    | Deprecated | Definition | Title/Description |
| ------------------------------------------------ | ------- | ------- | ---------- | ---------- | ----------------- |
| - [Host](#Metrics_Host )                         | No      | string  | No         | -          | -                 |
| - [Port](#Metrics_Port )                         | No      | integer | No         | -          | -                 |
| - [Enabled](#Metrics_Enabled )                   | No      | boolean | No         | -          | -                 |
| - [ProfilingHost](#Metrics_ProfilingHost )       | No      | string  | No         | -          | -                 |
| - [ProfilingPort](#Metrics_ProfilingPort )       | No      | integer | No         | -          | -                 |
| - [ProfilingEnabled](#Metrics_ProfilingEnabled ) | No      | boolean | No         | -          | -                 |

### <a name="Metrics_Host"></a>16.1. Property `root > Metrics > Host`

|              |             |
| ------------ | ----------- |
| **Type**     | `string`    |
| **Required** | No          |
| **Default**  | `"0.0.0.0"` |

### <a name="Metrics_Port"></a>16.2. Property `root > Metrics > Port`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `9091`    |

### <a name="Metrics_Enabled"></a>16.3. Property `root > Metrics > Enabled`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

### <a name="Metrics_ProfilingHost"></a>16.4. Property `root > Metrics > ProfilingHost`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

### <a name="Metrics_ProfilingPort"></a>16.5. Property `root > Metrics > ProfilingPort`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

### <a name="Metrics_ProfilingEnabled"></a>16.6. Property `root > Metrics > ProfilingEnabled`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

## <a name="EventLog"></a>17. Property `root > EventLog`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property              | Pattern | Type   | Deprecated | Definition | Title/Description                |
| --------------------- | ------- | ------ | ---------- | ---------- | -------------------------------- |
| - [DB](#EventLog_DB ) | No      | object | No         | -          | DB is the database configuration |

### <a name="EventLog_DB"></a>17.1. Property `root > EventLog > DB`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

**Description:** DB is the database configuration

| Property                               | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| -------------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#EventLog_DB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#EventLog_DB_User )           | No      | string  | No         | -          | User name                                                  |
| - [Password](#EventLog_DB_Password )   | No      | string  | No         | -          | Password of the user                                       |
| - [Host](#EventLog_DB_Host )           | No      | string  | No         | -          | Host address                                               |
| - [Port](#EventLog_DB_Port )           | No      | string  | No         | -          | Port Number                                                |
| - [EnableLog](#EventLog_DB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#EventLog_DB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

#### <a name="EventLog_DB_Name"></a>17.1.1. Property `root > EventLog > DB > Name`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** Database name

#### <a name="EventLog_DB_User"></a>17.1.2. Property `root > EventLog > DB > User`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** User name

#### <a name="EventLog_DB_Password"></a>17.1.3. Property `root > EventLog > DB > Password`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** Password of the user

#### <a name="EventLog_DB_Host"></a>17.1.4. Property `root > EventLog > DB > Host`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** Host address

#### <a name="EventLog_DB_Port"></a>17.1.5. Property `root > EventLog > DB > Port`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `""`     |

**Description:** Port Number

#### <a name="EventLog_DB_EnableLog"></a>17.1.6. Property `root > EventLog > DB > EnableLog`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

**Description:** EnableLog

#### <a name="EventLog_DB_MaxConns"></a>17.1.7. Property `root > EventLog > DB > MaxConns`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `0`       |

**Description:** MaxConns is the maximum number of connections in the pool.

## <a name="HashDB"></a>18. Property `root > HashDB`

|                           |                                                         |
| ------------------------- | ------------------------------------------------------- |
| **Type**                  | `object`                                                |
| **Required**              | No                                                      |
| **Additional properties** | [[Not allowed]](# "Additional Properties not allowed.") |

| Property                          | Pattern | Type    | Deprecated | Definition | Title/Description                                          |
| --------------------------------- | ------- | ------- | ---------- | ---------- | ---------------------------------------------------------- |
| - [Name](#HashDB_Name )           | No      | string  | No         | -          | Database name                                              |
| - [User](#HashDB_User )           | No      | string  | No         | -          | User name                                                  |
| - [Password](#HashDB_Password )   | No      | string  | No         | -          | Password of the user                                       |
| - [Host](#HashDB_Host )           | No      | string  | No         | -          | Host address                                               |
| - [Port](#HashDB_Port )           | No      | string  | No         | -          | Port Number                                                |
| - [EnableLog](#HashDB_EnableLog ) | No      | boolean | No         | -          | EnableLog                                                  |
| - [MaxConns](#HashDB_MaxConns )   | No      | integer | No         | -          | MaxConns is the maximum number of connections in the pool. |

### <a name="HashDB_Name"></a>18.1. Property `root > HashDB > Name`

|              |               |
| ------------ | ------------- |
| **Type**     | `string`      |
| **Required** | No            |
| **Default**  | `"prover_db"` |

**Description:** Database name

### <a name="HashDB_User"></a>18.2. Property `root > HashDB > User`

|              |                 |
| ------------ | --------------- |
| **Type**     | `string`        |
| **Required** | No              |
| **Default**  | `"prover_user"` |

**Description:** User name

### <a name="HashDB_Password"></a>18.3. Property `root > HashDB > Password`

|              |                 |
| ------------ | --------------- |
| **Type**     | `string`        |
| **Required** | No              |
| **Default**  | `"prover_pass"` |

**Description:** Password of the user

### <a name="HashDB_Host"></a>18.4. Property `root > HashDB > Host`

|              |                    |
| ------------ | ------------------ |
| **Type**     | `string`           |
| **Required** | No                 |
| **Default**  | `"zkevm-state-db"` |

**Description:** Host address

### <a name="HashDB_Port"></a>18.5. Property `root > HashDB > Port`

|              |          |
| ------------ | -------- |
| **Type**     | `string` |
| **Required** | No       |
| **Default**  | `"5432"` |

**Description:** Port Number

### <a name="HashDB_EnableLog"></a>18.6. Property `root > HashDB > EnableLog`

|              |           |
| ------------ | --------- |
| **Type**     | `boolean` |
| **Required** | No        |
| **Default**  | `false`   |

**Description:** EnableLog

### <a name="HashDB_MaxConns"></a>18.7. Property `root > HashDB > MaxConns`

|              |           |
| ------------ | --------- |
| **Type**     | `integer` |
| **Required** | No        |
| **Default**  | `200`     |

**Description:** MaxConns is the maximum number of connections in the pool.

----------------------------------------------------------------------------------------------------------------------------
Generated using [json-schema-for-humans](https://github.com/coveooss/json-schema-for-humans)
