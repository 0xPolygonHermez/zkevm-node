
This is a refactor of L1 synchronization to improve speed.
- It ask data in parallel  to L1 meanwhile another goroutine is execution the rollup info.
- It makes that executor be ocupied 100% of time.

## Pending to do

  - All the stuff related to updating last block on L1 could be moved to another class
  - Check context usage:
    It need a context to cancel itself and create another context to cancel workers?
  - Emit metrics
  - if nothing to update reduce  code to be executed (not sure, because functionality  to keep update beyond last block on L1)
  - Improve the unittest of all objects
  - Check all log.fatals to remove it or add a status before the panic
  - Missing **feature update beyond last block on L1**: Old syncBlocks method try to ask for blocks over last L1 block, I suppose that is to keep   synchronizing even a long the synchronization have new blocks. This is not implemented here
    This is the behaviour of ethman in that situation:
           - GetRollupInfoByBlockRange returns no errors, zero blocks...
           - EthBlockByNumber returns error:  "not found"
- Some test on ` synchronizer/synchronizer_test.go` are based on this feature, so are running against legacy code
- Move to configuration file some 'hardcoded' values

## Configuration
This feature is experimental for that reason you can configure to use old sequential one: 
```
[Synchronizer]
UseParallelModeForL1Synchronization = false
```
If you activate this feature you can configure:
- `NumberOfParallelOfEthereumClients`: how many parallel request can be done. Currently this create multiples instances of etherman over same server, in the future maybe make sense to use differents servers
- `CapacityOfBufferingRollupInfoFromL1`:  buffer of data pending to be processed
```
UseParallelModeForL1Synchronization = true
	[Synchronizer.L1ParallelSynchronization]
		NumberOfParallelOfEthereumClients = 2
		CapacityOfBufferingRollupInfoFromL1 = 10
```
## Remakable logs
### How to known the occupation of executor
To check that executor are fully ocuppied you can check next log:
```
INFO	synchronizer/l1_processor_consumer.go:110	consumer: processing rollupInfo #1291: range:[188064, 188164] num_blocks [0] wasted_time_waiting_for_data [74.17575ms] last_process_time [2.534115ms] block_per_second [0.000000]
```
The `wasted_time_waiting_for_data` show the waiting time between this call and the previous to executor. If this value (after 20 interations) are greater to 1 seconds a warning is show.

### Estimated time to be fully synchronizer with L1
This log show the estimated time (**ETA**) to reach the block goal
```
INFO	synchronizer/l1_data_retriever_producer.go:255	producer: Statistics:ETA: 3h40m1.311379085s percent:1.35  blocks_per_seconds:706.80 pending_block:127563/9458271 num_errors:0
```

## Flow of data
![l1_sync_channels_flow_v2 drawio](https://github.com/0xPolygonHermez/zkevm-node/assets/129153821/430abeb3-13b2-4c13-8d5e-4996a134a353)

## Class diagram
This is a class diagram of principal class an relationships.
The entry point is `synchronizer.go:276` function `syncBlocksParallel`.
- It create all objects needed and launch `l1SyncOrchestration` that wait until the job is done to return 

### The main objects are:
- `l1RollupInfoProducer`: is the object that send rollup data through the channel
- `l1RollupInfoConsumer`: that receive the data and execute it

![image](https://github.com/0xPolygonHermez/zkevm-node/assets/129153821/957a3e95-77c7-446b-a6ec-ef28cc44cb18)
