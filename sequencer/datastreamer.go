package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

//TODO: create interface to access datastreamer functions

func (f *finalizer) DSSendL2Block(batchNumber uint64, blockResponse *state.ProcessBlockResponse) error {
	forkID := f.state.GetForkIDByBatchNumber(batchNumber)

	// Send data to streamer
	if f.streamServer != nil {
		l2Block := state.DSL2Block{
			BatchNumber:    batchNumber,
			L2BlockNumber:  blockResponse.BlockNumber,
			Timestamp:      int64(blockResponse.Timestamp),
			GlobalExitRoot: blockResponse.BlockInfoRoot, //TODO: is it ok?
			Coinbase:       f.sequencerAddress,
			ForkID:         uint16(forkID),
			BlockHash:      blockResponse.BlockHash,
			StateRoot:      blockResponse.BlockHash, //TODO: in etrog the blockhash is the block root
		}

		l2Transactions := []state.DSL2Transaction{}

		for _, txResponse := range blockResponse.TransactionResponses {
			binaryTxData, err := txResponse.Tx.MarshalBinary()
			if err != nil {
				return err
			}

			l2Transaction := state.DSL2Transaction{
				L2BlockNumber:               blockResponse.BlockNumber,
				EffectiveGasPricePercentage: uint8(txResponse.EffectivePercentage),
				IsValid:                     1,
				EncodedLength:               uint32(len(binaryTxData)),
				Encoded:                     binaryTxData,
			}

			l2Transactions = append(l2Transactions, l2Transaction)
		}

		f.dataToStream <- state.DSL2FullBlock{
			DSL2Block: l2Block,
			Txs:       l2Transactions,
		}
	}

	return nil
}

func (f *finalizer) DSSendUpdateGER(batchNumber uint64, timestamp int64, GER common.Hash, stateRoot common.Hash) {
	//TODO: review this datastream event
	updateGer := state.DSUpdateGER{
		BatchNumber:    batchNumber,
		Timestamp:      timestamp,
		GlobalExitRoot: GER,
		Coinbase:       f.sequencerAddress,
		ForkID:         uint16(f.state.GetForkIDByBatchNumber(batchNumber)),
		StateRoot:      stateRoot,
	}

	err := f.streamServer.StartAtomicOp()
	if err != nil {
		log.Errorf("failed to start atomic op for batch %v: %v", batchNumber, err)
		return
	}

	_, err = f.streamServer.AddStreamEntry(state.EntryTypeUpdateGER, updateGer.Encode())
	if err != nil {
		log.Errorf("failed to add stream entry for batch %v: %v", batchNumber, err)
		return
	}

	err = f.streamServer.CommitAtomicOp()
	if err != nil {
		log.Errorf("failed to commit atomic op for batch %v: %v", batchNumber, err)
		return
	}
}
