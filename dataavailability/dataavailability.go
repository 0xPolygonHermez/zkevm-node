package dataavailability

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type DataAvailability struct {
	isTrustedSequencer bool

	state       stateInterface
	zkEVMClient syncinterfaces.ZKEVMClientInterface
	backend     DABackender

	ctx context.Context
}

func New(
	IsTrustedSequencer bool,
	backend DABackender,
	state stateInterface,
	zkEVMClient syncinterfaces.ZKEVMClientInterface,
) (*DataAvailability, error) {
	da := &DataAvailability{
		isTrustedSequencer: IsTrustedSequencer,
		backend:            backend,
		state:              state,
		zkEVMClient:        zkEVMClient,
		ctx:                context.Background(),
	}
	err := da.backend.Init()
	return da, err
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (d *DataAvailability) PostSequence(ctx context.Context, sequences []types.Sequence) ([]byte, error) {
	return d.backend.PostSequence(ctx, sequences)
}

func (d *DataAvailability) GetBatchL2Data(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	found := true
	transactionsData, err := d.state.GetBatchL2DataByNumber(d.ctx, batchNum, nil)
	if err != nil {
		if err == state.ErrNotFound {
			found = false
		} else {
			return nil, fmt.Errorf("failed to get batch data from state for batch num %d: %w", batchNum, err)
		}
	}
	actualTransactionsHash := crypto.Keccak256Hash(transactionsData)
	if !found || expectedTransactionsHash != actualTransactionsHash {
		if found {
			log.Warnf(unexpectedHashTemplate, batchNum, expectedTransactionsHash, actualTransactionsHash)
		}

		if !d.isTrustedSequencer {
			log.Info("trying to get data from trusted sequencer")
			data, err := d.getDataFromTrustedSequencer(batchNum, expectedTransactionsHash)
			if err != nil {
				log.Error(err)
			} else {
				return data, nil
			}
		}

		log.Info("trying to get data from the data availability backend")
		data, err := d.backend.GetData(batchNum, expectedTransactionsHash)
		if err != nil {
			log.Error(err)
			if d.isTrustedSequencer {
				return nil, fmt.Errorf("data not found on the local DB nor on any data committee member")
			} else {
				return nil, fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member")
			}
		}
		return data, nil
	}
	return transactionsData, nil
}

const unexpectedHashTemplate = "missmatch on transaction data for batch num %d. Expected hash %s, actual hash: %s"

func (d *DataAvailability) getDataFromTrustedSequencer(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	b, err := d.zkEVMClient.BatchByNumber(d.ctx, big.NewInt(int64(batchNum)))
	if err != nil {
		return nil, fmt.Errorf("failed to get batch num %d from trusted sequencer: %w", batchNum, err)
	}
	actualTransactionsHash := crypto.Keccak256Hash(b.BatchL2Data)
	if expectedTransactionsHash != actualTransactionsHash {
		return nil, fmt.Errorf(
			unexpectedHashTemplate, batchNum, expectedTransactionsHash, actualTransactionsHash,
		)
	}
	return b.BatchL2Data, nil
}
