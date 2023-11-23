package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/processor_manager"
)

func defaultsL1EventProcessors(sync *ClientSynchronizer) *processor_manager.L1EventProcessors {
	p := processor_manager.NewL1EventProcessorsBuilder()
	p.AddEventProcessor(etherman.GlobalExitRootsOrder, incaberry.NewProcessorL1GlobalExitRoot(sync.state))
	p.AddEventProcessor(etherman.SequenceBatchesOrder, incaberry.NewProcessorL1SequenceBatches(sync.state, sync.etherMan, sync.pool, sync.eventLog, sync))
	p.AddEventProcessor(etherman.ForcedBatchesOrder, incaberry.NewProcessL1ForcedBatches(sync.state))
	p.AddEventProcessor(etherman.SequenceForceBatchesOrder, incaberry.NewProcessL1SequenceForcedBatches(sync.state, sync))
	p.AddEventProcessor(etherman.TrustedVerifyBatchOrder, incaberry.NewProcessorL1VerifyBatch(sync.state))
	p.AddEventProcessor(etherman.VerifyBatchOrder, incaberry.NewProcessorL1VerifyBatch(sync.state))
	p.AddEventProcessor(etherman.ForkIDsOrder, incaberry.NewProcessorForkId(sync.state, sync))
	return p.Build()
}
