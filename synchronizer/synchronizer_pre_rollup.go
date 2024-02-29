package synchronizer

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/jackc/pgx/v4"
)

// SyncPreRollup is the struct for synchronizing pre genesis rollup events.
// Implements: syncinterfaces.SyncPreRollupSyncer
type SyncPreRollup struct {
	etherman            syncinterfaces.EthermanPreRollup
	state               syncinterfaces.StateLastBlockGetter
	blockRangeProcessor syncinterfaces.BlockRangeProcessor
	SyncChunkSize       uint64
	GenesisBlockNumber  uint64
}

func NewSyncPreRollup(
	etherman syncinterfaces.EthermanPreRollup,
	state syncinterfaces.StateLastBlockGetter,
	blockRangeProcessor syncinterfaces.BlockRangeProcessor,
	syncChunkSize uint64,
	genesisBlockNumber uint64,
) *SyncPreRollup {
	return &SyncPreRollup{
		etherman:            etherman,
		state:               state,
		blockRangeProcessor: blockRangeProcessor,
		SyncChunkSize:       syncChunkSize,
		GenesisBlockNumber:  genesisBlockNumber,
	}
}

// SynchronizePreGenesisRollupEvents sync pre-rollup events
func (s *SyncPreRollup) SynchronizePreGenesisRollupEvents(ctx context.Context) error {
	// Sync events from RollupManager that happen before rollup creation
	log.Info("synchronizing events from RollupManager that happen before rollup creation")
	//genesisBlockNumber := s.genesis.BlockNumber //uint64(19218658)
	genesisBlockNumber := s.GenesisBlockNumber
	lxLyUpgradeBlock, err := s.etherman.GetL1BlockUpgradeLxLy(ctx, genesisBlockNumber)
	if err != nil && errors.Is(err, etherman.ErrNotFound) {
		log.Infof("LxLy upgrade not detected before genesis block %d, it'll be sync as usual. Nothing to do yet", genesisBlockNumber)
		//s.ProcessL1InfoRootEvents(ctx, uint64(19331715), uint64(19333413), s.SyncChunkSize, dbTx)
		return nil
	}
	if err != nil {
		log.Errorf("error getting LxLy upgrade block. Error: %v", err)
		return err
	}
	fromBlock, err := s.getStartingL1Block(ctx, lxLyUpgradeBlock, nil)
	if err != nil {
		log.Errorf("error getting starting L1 block. Error: %v", err)
		return err
	}
	return s.ProcessL1InfoRootEvents(ctx, fromBlock, genesisBlockNumber-1, s.SyncChunkSize)
}

func (s *SyncPreRollup) getStartingL1Block(ctx context.Context, upgradeLxLyBlockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	lastBlock, err := s.state.GetLastBlock(ctx, dbTx)
	if err != nil && errors.Is(err, state.ErrStateNotSynchronized) {
		// No block on DB
		log.Infof("No block on DB, starting from LxLy upgrade block %d", upgradeLxLyBlockNumber)
		return upgradeLxLyBlockNumber, nil
	} else if err != nil {
		log.Errorf("Error getting last Block on DB err:%v", err)
		return 0, err
	}
	if lastBlock.BlockNumber >= s.GenesisBlockNumber {
		log.Warnf("Last block processed is %d, which is greater or equal than the genesis block %d", lastBlock, s.GenesisBlockNumber)
		return 0, errors.New("last block processed is greater or equal than the genesis block")
	}
	log.Infof("Pre genesis LxLy upgrade at block %d, last block processed on DB is %d", upgradeLxLyBlockNumber, lastBlock.BlockNumber)
	return lastBlock.BlockNumber, nil
}

// ProcessL1InfoRootEvents processes the L1InfoRoot events for a range for L1 blocks
func (s *SyncPreRollup) ProcessL1InfoRootEvents(ctx context.Context, fromBlock uint64, toBlock uint64, syncChunkSize uint64) error {
	log.Info("synchronizing L1InfoRoot events")
	log.Infof("Starting syncing pre genesis LxLy events from block %d to block %d (total %d blocks)",
		fromBlock, toBlock, toBlock-fromBlock+1)
	for i := fromBlock; true; i += syncChunkSize {
		toBlockReq := min(i+syncChunkSize-1, toBlock)
		percent := float32(toBlockReq-fromBlock+1) * 100.0 / float32(toBlock-fromBlock+1)
		log.Infof("sync L1InfoTree events from %d to %d percent:%3.1f pending_blocks:%d", i, toBlockReq, percent, toBlock-toBlockReq)
		blocks, order, err := s.etherman.GetRollupInfoByBlockRangePreviousRollupGenesis(ctx, i, &toBlockReq)
		if err != nil {
			log.Error("error getting rollupInfoByBlockRange before rollup genesis: ", err)
			return err
		}
		err = s.blockRangeProcessor.ProcessBlockRange(ctx, blocks, order)
		if err != nil {
			log.Error("error processing blocks before the genesis: ", err)
			return err
		}
		if toBlockReq == toBlock {
			break
		}
	}
	return nil
}
