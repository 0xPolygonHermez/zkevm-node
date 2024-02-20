package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/elderberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/etrog"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
)

func defaultsL1EventProcessors(sync *ClientSynchronizer) *processor_manager.L1EventProcessors {
	p := processor_manager.NewL1EventProcessorsBuilder()
	p.Register(incaberry.NewProcessorL1GlobalExitRoot(sync.state))
	p.Register(incaberry.NewProcessorL1SequenceBatches(sync.state, sync.etherMan, sync.pool, sync.eventLog, sync))
	p.Register(incaberry.NewProcessL1ForcedBatches(sync.state))
	p.Register(incaberry.NewProcessL1SequenceForcedBatches(sync.state, sync))
	p.Register(incaberry.NewProcessorForkId(sync.state, sync))
	p.Register(etrog.NewProcessorL1InfoTreeUpdate(sync.state))
	sequenceBatchesProcessor := etrog.NewProcessorL1SequenceBatches(sync.state, sync, common.DefaultTimeProvider{}, sync.halter)
	p.Register(sequenceBatchesProcessor)
	p.Register(incaberry.NewProcessorL1VerifyBatch(sync.state))
	p.Register(etrog.NewProcessorL1UpdateEtrogSequence(sync.state, sync, common.DefaultTimeProvider{}))
	p.Register(elderberry.NewProcessorL1SequenceBatchesElderberry(sequenceBatchesProcessor, sync.state))
	return p.Build()
}
