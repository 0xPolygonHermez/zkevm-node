
This is a refactor of L1 synchronization to improve speed.
- It ask data in parallel  to L1 meanwhile another goroutine is executing the rollup info.
- It makes that the executor be occupied 100% of the time.

## Pending to do  
- Some test on ` synchronizer/synchronizer_test.go` are based on this feature, so are running against legacy code

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
