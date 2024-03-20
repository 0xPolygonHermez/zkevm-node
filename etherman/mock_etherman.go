package etherman

import (
	"errors"
	"fmt"
	"math/big"

	polygonzkevm "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonvalidium_xlayer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BuildMockSequenceBatchesTxData builds a []bytes to be sent to the PoE SC method SequenceBatches.
func (etherMan *Client) BuildMockSequenceBatchesTxData(sender common.Address,
	validiumBatchData []polygonzkevm.PolygonValidiumEtrogValidiumBatchData,
	maxSequenceTimestamp uint64,
	initSequencedBatch uint64,
	l2Coinbase common.Address,
	dataAvailabilityMessage []byte) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if errors.Is(err, ErrNotFound) {
		return nil, nil, fmt.Errorf("failed to build sequence batches, err: %w", ErrPrivateKeyNotFound)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	tx, err := etherMan.sequenceMockBatches(opts, validiumBatchData, maxSequenceTimestamp, initSequencedBatch, l2Coinbase, dataAvailabilityMessage)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func (etherMan *Client) sequenceMockBatches(opts bind.TransactOpts,
	validiumBatchData []polygonzkevm.PolygonValidiumEtrogValidiumBatchData,
	maxSequenceTimestamp uint64,
	initSequencedBatch uint64,
	l2Coinbase common.Address,
	dataAvailabilityMessage []byte) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error
	tx, err = etherMan.ZkEVM.SequenceBatchesValidium(&opts, validiumBatchData, maxSequenceTimestamp, initSequencedBatch, l2Coinbase, dataAvailabilityMessage)

	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		err = fmt.Errorf(
			"error sequencing batches: %w, dataAvailabilityMessage: %s",
			err, common.Bytes2Hex(dataAvailabilityMessage),
		)
	}

	return tx, err
}
