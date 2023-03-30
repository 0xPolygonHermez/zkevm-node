# JSON RPC Endpoints

Here you will find the list of all supported JSON RPC endpoints and any differences between them in comparison to the default behavior of an ethereum node.

If the endpoint is not in the list below, it means this specific endpoint is not supported yet, feel free to open an issue requesting it to be added and please explain the reason why you need it. 

<!-- DEBUG -->
- `debug_traceBlockByHash`
- `debug_traceBlockByNumber`
- `debug_traceTransaction`

<!-- ETH -->
- `eth_blockNumber`
- `eth_call` _* doesn't support state override at the moment and pending block_ 
- `eth_chainId`
- `eth_estimateGas` _* if the block number is set to pending we assume it is the latest_
- `eth_gasPrice`
- `eth_getBalance` _* if the block number is set to pending we assume it is the latest_
- `eth_getBlockByHash`
- `eth_getBlockByNumber`
- `eth_getBlockTransactionCountByHash`
- `eth_getBlockTransactionCountByNumber`
- `eth_getCode` _* if the block number is set to pending we assume it is the latest_
- `eth_getCompilers` _* response is always empty_
- `eth_getFilterChanges`
- `eth_getFilterLogs`
- `eth_getLogs`
- `eth_getStorageAt` _* if the block number is set to pending we assume it is the latest_
- `eth_getTransactionByBlockHashAndIndex`
- `eth_getTransactionByHash`
- `eth_getTransactionCount`
- `eth_getTransactionByBlockNumberAndIndex` _* if the block number is set to pending we assume it is the latest_
- `eth_getTransactionReceipt`
- `eth_getUncleByBlockHashAndIndex` _* response is always empty_
- `eth_getUncleByBlockNumberAndIndex` _* response is always empty_
- `eth_getUncleCountByBlockHash` _* response is always zero_
- `eth_getUncleCountByBlockNumber` _* response is always zero_
- `eth_newBlockFilter`
- `eth_newFilter`
- `eth_protocolVersion` _* response is always zero_
- `eth_sendRawTransaction` _* can relay TXs to another node_
- `eth_syncing`
- `eth_subscribe`
- `eth_uninstallFilter`
- `eth_unsubscribe`

<!-- NET -->
- `net_version`

<!-- TXPOOL -->
- `txpool_content` _* response is always empty_

<!-- WEB3 -->
- `web3_clientVersion`
- `web3_sha3`

<!-- ZKEVM -->
- `zkevm_batchNumber`
- `zkevm_batchNumberByBlockNumber`
- `zkevm_getBatchByNumber`
- `zkevm_getBroadcastURI` _* deprecated, will be removed in the future_
- `zkevm_consolidatedBlockNumber`
- `zkevm_isBlockConsolidated`
- `zkevm_isBlockVirtualized`
- `zkevm_verifiedBatchNumber`
- `zkevm_virtualBatchNumber`
