package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-node/state"
)

func (f *finalizer) DSSendL2Block(batchNumber uint64, blockResponse *state.ProcessBlockResponse) error {
	forkID := f.stateIntf.GetForkIDByBatchNumber(batchNumber)

	// Send data to streamer
	if f.streamServer != nil {
		l2Block := state.DSL2Block{
			BatchNumber:   batchNumber,
			L2BlockNumber: blockResponse.BlockNumber,
			Timestamp:     int64(blockResponse.Timestamp),
			GERorInfoRoot: blockResponse.BlockInfoRoot,
			Coinbase:      f.sequencerAddress,
			ForkID:        uint16(forkID),
			BlockHash:     blockResponse.BlockHash,
			StateRoot:     blockResponse.BlockHash, //From etrog the blockhash is the block root
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
