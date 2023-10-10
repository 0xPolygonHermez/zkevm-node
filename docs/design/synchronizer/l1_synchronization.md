# L1 parallel synchronization
This is a refactor of L1 synchronization to improve speed.
- It ask data in parallel  to L1 meanwhile another goroutine is executing the rollup info.
- It makes that the executor be occupied 100% of the time.

## Pending to do  
- Some test on ` synchronizer/synchronizer_test.go` are based on this feature, so are running against legacy code

## Configuration
You could choose between new L1 parallel sync or sequential one (legacy): 
```
[Synchronizer]
UseParallelModeForL1Synchronization = false
```
If you activate this feature you can configure:
- `NumberOfParallelOfEthereumClients`: how many parallel request can be done. You must consider that 1 is just for requesting the last block on L1, and the rest for rollup info
- `CapacityOfBufferingRollupInfoFromL1`:  buffer of data pending to be processed. This is the queue data to be executed by consumer.

For a full description of fields please check config-file documentation.

Example: 
```
UseParallelModeForL1Synchronization = true
	[Synchronizer.L1ParallelSynchronization]
		NumberOfParallelOfEthereumClients = 2
		CapacityOfBufferingRollupInfoFromL1 = 10
		TimeForCheckLastBlockOnL1Time = "5s"
		TimeoutForRequestLastBlockOnL1 = "5s"
		MaxNumberOfRetriesForRequestLastBlockOnL1 = 3
		TimeForShowUpStatisticsLog = "5m"
		TimeOutMainLoop = "5m"
		MinTimeBetweenRetriesForRollupInfo = "5s"
		[Synchronizer.L1ParallelSynchronization.PerformanceCheck]
			AcceptableTimeWaitingForNewRollupInfo = "5s"
			NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo = 10

```
## Remakable logs
### How to known the occupation of executor
To check that executor are fully ocuppied you can check next log:
```
INFO	synchronizer/l1_rollup_info_consumer.go:128	consumer: processing rollupInfo #1553: range:[8720385, 8720485] num_blocks [37] statistics:wasted_time_waiting_for_data [0s] last_process_time [6m2.635208117s] block_per_second [2.766837]
```
The `wasted_time_waiting_for_data` show the waiting time between this call and the previous to executor. It could show a warning configuring `Synchronizer.L1ParallelSynchronization.PerformanceCheck`

### Estimated time to be fully synchronizer with L1
This log show the estimated time (**ETA**) to reach the block goal. You can configure the frequency with var `TimeForShowUpStatisticsLog`
```
INFO	synchronizer/l1_rollup_info_producer.go:357	producer: Statistics:ETA: 54h7m47.594422312s percent:12.26  blocks_per_seconds:5.48 pending_block:149278/1217939 num_errors:8
```

## Flow of data
![l1_sync_channels_flow_v2 drawio](l1_sync_channels_flow_v2.drawio.png)


### The main objects are:
- `l1SyncOrchestration`: is the entry point and the reponsable to launch the producer and consumer
- `l1RollupInfoProducer`: this object send rollup data through the channel to the consumer
- `l1RollupInfoConsumer`: that receive the data and execute it


## Future changes
- Configure multiples servers for L1 information: instead of calling the same server,it make sense to configure individually each URL to allow to have multiples sources
