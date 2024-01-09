package l1_parallel_sync

import (
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestSORMulticaseWithReset(t *testing.T) {
	tcs := []struct {
		description                     string
		lastBlock                       uint64
		packages                        []L1SyncMessage
		expected                        []L1SyncMessage
		expectedlastBlockOnSynchronizer uint64
		resetOnPackageNumber            int
		resetToBlock                    uint64
	}{
		{
			description: "inverse_br",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(131, 141),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []L1SyncMessage{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
			},
			expectedlastBlockOnSynchronizer: 130,
			resetOnPackageNumber:            1,
			resetToBlock:                    100,
		},
		{
			description: "crtl_linked_to_br",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(131, 141),
				*newActionPackage(eventNone),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []L1SyncMessage{
				*newActionPackage(eventNone),
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
			},
			expectedlastBlockOnSynchronizer: 130,
			resetOnPackageNumber:            1,
			resetToBlock:                    100,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			sut := newFilterToSendOrdererResultsToConsumer(tc.lastBlock)
			sendData := []L1SyncMessage{}
			for i, p := range tc.packages {
				if i == tc.resetOnPackageNumber {
					sut.Reset(tc.resetToBlock)
				}
				dataToSend := sut.Filter(p)
				sendData = append(sendData, dataToSend...)
			}

			require.Equal(t, tc.expected, sendData)
			require.Equal(t, tc.expectedlastBlockOnSynchronizer, sut.lastBlockOnSynchronizer)
		})
	}
}

func TestSORMulticase(t *testing.T) {
	tcs := []struct {
		description                      string
		lastBlock                        uint64
		packages                         []L1SyncMessage
		expected                         []L1SyncMessage
		excpectedLastBlockOnSynchronizer uint64
	}{
		{
			description:                      "empty_case",
			lastBlock:                        100,
			packages:                         []L1SyncMessage{},
			expected:                         []L1SyncMessage{},
			excpectedLastBlockOnSynchronizer: 100,
		},
		{
			description:                      "just_ctrl",
			lastBlock:                        100,
			packages:                         []L1SyncMessage{*newActionPackage(eventNone)},
			expected:                         []L1SyncMessage{*newActionPackage(eventNone)},
			excpectedLastBlockOnSynchronizer: 100,
		},
		{
			description:                      "just_br",
			lastBlock:                        100,
			packages:                         []L1SyncMessage{*newDataPackage(101, 119)},
			expected:                         []L1SyncMessage{*newDataPackage(101, 119)},
			excpectedLastBlockOnSynchronizer: 119,
		},
		{
			description:                      "just_br_missing_intermediate_block",
			lastBlock:                        100,
			packages:                         []L1SyncMessage{*newDataPackage(102, 119)},
			expected:                         []L1SyncMessage{},
			excpectedLastBlockOnSynchronizer: 100,
		},
		{
			description: "inverse_br",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(131, 141),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []L1SyncMessage{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
				*newDataPackage(131, 141),
			},
			excpectedLastBlockOnSynchronizer: 141,
		},
		{
			description: "crtl_linked_to_br",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(131, 141),
				*newActionPackage(eventNone),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []L1SyncMessage{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
				*newDataPackage(131, 141),
				*newActionPackage(eventNone),
			},
			excpectedLastBlockOnSynchronizer: 141,
		},
		{
			description: "crtl_linked_to_last_br",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, 140),
				*newActionPackage(eventNone),
				*newDataPackage(101, 110)},
			expected: []L1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, 140),
				*newActionPackage(eventNone),
			},
			excpectedLastBlockOnSynchronizer: 140,
		},
		{
			description: "latest with no data doesnt change last block",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(101, 110)},
			expected: []L1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
			},
			excpectedLastBlockOnSynchronizer: 130,
		},
		{
			description: "two latest one empty and one with data change to highest block in rollupinfo",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(101, 110),
				*newDataPackageWithData(131, latestBlockNumber, 140),
			},
			expected: []L1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackageWithData(131, latestBlockNumber, 140),
			},
			excpectedLastBlockOnSynchronizer: 140,
		},
		{
			description: "one latest one normal",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newDataPackage(131, 140),
				*newActionPackage(eventNone),
				*newDataPackage(101, 110),
			},
			expected: []L1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newDataPackage(131, 140),
				*newActionPackage(eventNone),
			},
			excpectedLastBlockOnSynchronizer: 140,
		},
		{
			description: "a rollupinfo with data",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackageWithData(121, 130, 125),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(101, 110),
				*newDataPackage(131, 140),
			},
			expected: []L1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackageWithData(121, 130, 125),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(131, 140),
			},
			excpectedLastBlockOnSynchronizer: 140,
		},
		{
			description: "two latest empty with control in between",
			lastBlock:   100,
			packages: []L1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(101, 110),
				*newDataPackage(131, 140),
			},
			expected: []L1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(131, latestBlockNumber),
				*newActionPackage(eventNone),
				*newDataPackage(131, 140),
			},
			excpectedLastBlockOnSynchronizer: 140,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			sut := newFilterToSendOrdererResultsToConsumer(tc.lastBlock)
			sendData := []L1SyncMessage{}
			for _, p := range tc.packages {
				dataToSend := sut.Filter(p)
				sendData = append(sendData, dataToSend...)
			}
			require.Equal(t, len(tc.expected), len(sendData))
			require.Equal(t, tc.expected, sendData)
			require.Equal(t, tc.excpectedLastBlockOnSynchronizer, sut.lastBlockOnSynchronizer)
		})
	}
}

func newDataPackage(fromBlock, toBlock uint64) *L1SyncMessage {
	res := L1SyncMessage{
		data: rollupInfoByBlockRangeResult{
			blockRange: blockRange{
				fromBlock: fromBlock,
				toBlock:   toBlock,
			},
			lastBlockOfRange: types.NewBlock(&types.Header{Number: big.NewInt(int64(toBlock))}, nil, nil, nil, nil),
		},
		dataIsValid: true,
		ctrlIsValid: false,
	}
	if toBlock == latestBlockNumber {
		res.data.lastBlockOfRange = nil
	}
	return &res
}

func newDataPackageWithData(fromBlock, toBlock uint64, blockWithData uint64) *L1SyncMessage {
	res := L1SyncMessage{
		data: rollupInfoByBlockRangeResult{
			blockRange: blockRange{
				fromBlock: fromBlock,
				toBlock:   toBlock,
			},
			blocks: []etherman.Block{{BlockNumber: blockWithData}},
		},
		dataIsValid: true,
		ctrlIsValid: false,
	}

	return &res
}

func newActionPackage(action eventEnum) *L1SyncMessage {
	return &L1SyncMessage{
		dataIsValid: false,
		data: rollupInfoByBlockRangeResult{
			blockRange: blockRange{
				fromBlock: 0,
				toBlock:   0,
			},
		},

		ctrlIsValid: true,
		ctrl: l1ConsumerControl{
			event: action,
		},
	}
}
