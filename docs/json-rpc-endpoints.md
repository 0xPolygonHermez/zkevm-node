# JSON RPC Endpoints

## ETH

### Supported

- eth_blockNumber
- eth_chainId
- eth_gasPrice
- eth_getBlockByHash
- eth_getBlockByNumber
- eth_getFilterChanges
- eth_getFilterLogs
- eth_getLogs
- eth_getTransactionByBlockHashAndIndex
- eth_getTransactionByHash
- eth_getTransactionCount
- eth_getBlockTransactionCountByHash
- eth_getBlockTransactionCountByNumber
- eth_getTransactionReceipt
- eth_newBlockFilter
- eth_newFilter
- eth_uninstallFilter
- eth_syncing
- eth_subscribe
- eth_unsubscribe

### Differences

- eth_call _* doesn't support state override at the moment and pending block_ 
- eth_estimateGas _* if the block number is set to pending we assume it is the latest_
- eth_getBalance _* if the block number is set to pending we assume it is the latest_
- eth_getCode _* if the block number is set to pending we assume it is the latest_
- eth_getCompilers _* response is always empty_
- eth_getStorageAt _* if the block number is set to pending we assume it is the latest_
- eth_getTransactionByBlockNumberAndIndex _* if the block number is set to pending we assume it is the latest_
- eth_newPendingTransactionFilter _* not supported yet_
- eth_sendRawTransaction _* can relay TXs to another node_
- eth_getUncleByBlockHashAndIndex _* response is always empty_
- eth_getUncleByBlockNumberAndIndex _* response is always empty_
- eth_getUncleCountByBlockHash _* response is always zero_
- eth_getUncleCountByBlockNumber _* response is always zero_
- eth_protocolVersion _* response is always zero_

### Not supported

- eth_accounts _* not supported_
- eth_coinbase _* not supported_
- eth_compileSolidity _* not supported_
- eth_compileLLL _* not supported_
- eth_compileSerpent _* not supported_
- eth_getWork _* not supported_
- eth_hashrate _* not supported_
- eth_mining _* not supported_
- eth_sendTransaction _* not supported_
- eth_sign _* not supported_
- eth_signTransaction _* not supported_
- eth_submitWork _* not supported_
- eth_submitHashrate _* not supported_

## DEBUG

### Supported

- debug_traceTransaction
- debug_traceBlockByNumber
- debug_traceBlockByHash

### Differences

TBD

### Not supported

TBD

## NET

### Supported

- net_version

### Differences

TBD

### Not supported

- net_listening _* not supported_
- net_peerCount _* not supported_

## TXPOOL

### Supported

TBD

### Differences

- txpool_content _* response is always empty_

### Not supported

TBD

## WEB3

### Supported

- web3_clientVersion
- web3_sha3

### Differences

TBD

### Not supported

TBD

## ZKEVM

- zkevm_consolidatedBlockNumber
- zkevm_isBlockConsolidated
- zkevm_isBlockVirtualized
- zkevm_batchNumberByBlockNumber
- zkevm_batchNumber
- zkevm_virtualBatchNumber
- zkevm_verifiedBatchNumber
- zkevm_getBatchByNumber
- zkevm_getBroadcastURI
