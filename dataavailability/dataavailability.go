package dataavailability

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const unexpectedHashTemplate = "mismatch on transaction data for batch num %d. Expected hash %s, actual hash: %s"

// DataAvailability implements an abstract data availability integration
type DataAvailability struct {
	isTrustedSequencer bool

	state       stateInterface
	zkEVMClient ZKEVMClientTrustedBatchesGetter
	backend     DABackender

	ctx context.Context
}

// New creates a DataAvailability instance
func New(
	isTrustedSequencer bool,
	backend DABackender,
	state stateInterface,
	zkEVMClient ZKEVMClientTrustedBatchesGetter,
) (*DataAvailability, error) {
	da := &DataAvailability{
		isTrustedSequencer: isTrustedSequencer,
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
	batchesData := [][]byte{}
	for _, batch := range sequences {
		// Do not send to the DA backend data that will be stored to L1
		if batch.ForcedBatchTimestamp == 0 {
			batchesData = append(batchesData, batch.BatchL2Data)
		}
	}
	return d.backend.PostSequence(ctx, batchesData)
}

// GetBatchL2Data tries to return the data from a batch, in the following priorities
// 1. From local DB
// 2. From Sequencer
// 3. From DA backend
func (d *DataAvailability) GetBatchL2Data(batchNums []uint64, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	if len(batchNums) != len(batchHashes) {
		return nil, fmt.Errorf("invalid L2 batch data retrieval arguments, %d != %d", len(batchNums), len(batchHashes))
	}

	data, err := d.localData(batchNums, batchHashes)
	if err != nil {
		log.Warnf("error retrieving local data for batches %v: %s", batchNums, err.Error())
	} else {
		return data, nil
	}

	if !d.isTrustedSequencer {
		data, err = d.trustedSeqData(batchNums, batchHashes)
		if err != nil {
			log.Warnf("error retrieving trusted sequencer data for batches %v: %s", batchNums, err.Error())
		} else {
			return data, nil
		}
	}

	return d.backend.GetSequence(d.ctx, batchHashes, dataAvailabilityMessage)
}

func (d *DataAvailability) localData(numbers []uint64, hashes []common.Hash) ([][]byte, error) {
	data, err := d.state.GetBatchL2DataByNumbers(d.ctx, numbers, nil)
	if err != nil {
		return nil, err
	}
	var batches [][]byte
	for i := 0; i < len(numbers); i++ {
		batchNumber := numbers[i]
		expectedHash := hashes[i]
		batchData, ok := data[batchNumber]
		if !ok {
			return nil, fmt.Errorf("could not get data locally for batch numbers %v", numbers)
		}
		actualHash := crypto.Keccak256Hash(batchData)
		if actualHash != expectedHash {
			log.Warnf(unexpectedHashTemplate, batchNumber, expectedHash, actualHash)
		} else {
			batches = append(batches, batchData)
		}
	}
	return batches, nil
}

func (d *DataAvailability) trustedSeqData(numbers []uint64, hashes []common.Hash) ([][]byte, error) {
	data, err := d.getBatchesDataFromTrustedSequencer(numbers, hashes)
	if err != nil {
		return nil, err
	}
	var batches [][]byte
	for i := 0; i < len(numbers); i++ {
		batchNumber := numbers[i]
		// hash has already been checked
		batchData, ok := data[batchNumber]
		if !ok {
			continue
		}
		batches[i] = batchData
	}
	return batches, nil
}

func (d *DataAvailability) getDataFromTrustedSequencer(batchNum uint64, expectedTransactionsHash common.Hash) ([]byte, error) {
	b, err := d.zkEVMClient.BatchByNumber(d.ctx, new(big.Int).SetUint64(batchNum))
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

func (d *DataAvailability) getBatchesDataFromTrustedSequencer(batchNums []uint64, expectedHashes []common.Hash) (map[uint64][]byte, error) {
	if len(batchNums) != len(expectedHashes) {
		return nil, fmt.Errorf("invalid arguments, len of batch numbers does not equal length of expected hashes: %d != %d",
			len(batchNums), len(expectedHashes))
	}
	var nums []*big.Int
	for _, n := range batchNums {
		nums = append(nums, new(big.Int).SetUint64(n))
	}
	batchData, err := d.zkEVMClient.BatchesByNumbers(d.ctx, nums)
	if err != nil {
		return nil, fmt.Errorf("failed to get batches %v data from trusted sequencer: %w", batchNums, err)
	}
	result := make(map[uint64][]byte)
	for i := 0; i < len(batchNums); i++ {
		number := batchNums[i]
		batch := batchData[i]
		if batch.Empty {
			continue
		}
		expectedTransactionsHash := expectedHashes[i]
		actualTransactionsHash := crypto.Keccak256Hash(batch.BatchL2Data)
		if expectedTransactionsHash != actualTransactionsHash {
			log.Warnf(unexpectedHashTemplate, number, expectedTransactionsHash, actualTransactionsHash)
			continue
		}
		result[number] = batch.BatchL2Data
	}
	return result, nil
}
