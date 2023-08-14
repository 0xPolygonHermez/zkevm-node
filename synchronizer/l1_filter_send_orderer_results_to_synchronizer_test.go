package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SOR_Multicase(t *testing.T) {
	tcs := []struct {
		description             string
		lastBlock               uint64
		packages                []l1SyncMessage
		expected                []l1SyncMessage
		lastBlockOnSynchronizer uint64
	}{
		{
			description:             "empty_case",
			lastBlock:               100,
			packages:                []l1SyncMessage{},
			expected:                []l1SyncMessage{},
			lastBlockOnSynchronizer: 100,
		},
		{
			description:             "just_ctrl",
			lastBlock:               100,
			packages:                []l1SyncMessage{*newActionPackage(actionNone)},
			expected:                []l1SyncMessage{*newActionPackage(actionNone)},
			lastBlockOnSynchronizer: 100,
		},
		{
			description:             "just_br",
			lastBlock:               100,
			packages:                []l1SyncMessage{*newDataPackage(101, 119)},
			expected:                []l1SyncMessage{*newDataPackage(101, 119)},
			lastBlockOnSynchronizer: 119,
		},
		{
			description:             "just_br_missing_intermediate_block",
			lastBlock:               100,
			packages:                []l1SyncMessage{*newDataPackage(102, 119)},
			expected:                []l1SyncMessage{},
			lastBlockOnSynchronizer: 100,
		},
		{
			description: "inverse_br",
			lastBlock:   100,
			packages: []l1SyncMessage{
				*newDataPackage(131, 141),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []l1SyncMessage{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
				*newDataPackage(131, 141),
			},
			lastBlockOnSynchronizer: 141,
		},
		{
			description: "crtl_linked_to_br",
			lastBlock:   100,
			packages: []l1SyncMessage{
				*newDataPackage(131, 141),
				*newActionPackage(actionNone),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []l1SyncMessage{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
				*newDataPackage(131, 141),
				*newActionPackage(actionNone),
			},
			lastBlockOnSynchronizer: 141,
		},
		{
			description: "crtl_linked_to_last_br",
			lastBlock:   100,
			packages: []l1SyncMessage{
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, 140),
				*newActionPackage(actionNone),
				*newDataPackage(101, 110)},
			expected: []l1SyncMessage{
				*newDataPackage(101, 110),
				*newDataPackage(111, 120),
				*newDataPackage(121, 130),
				*newDataPackage(131, 140),
				*newActionPackage(actionNone),
			},
			lastBlockOnSynchronizer: 140,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			sut := newFilterToSendOrdererResultsToConsumer(tc.lastBlock)
			sendData := []l1SyncMessage{}
			for _, p := range tc.packages {
				dataToSend := sut.filter(p)
				sendData = append(sendData, dataToSend...)
			}

			require.Equal(t, tc.expected, sendData)
			require.Equal(t, tc.lastBlockOnSynchronizer, sut.lastBlockOnSynchronizer)
		})
	}
}

func newDataPackage(fromBlock, toBlock uint64) *l1SyncMessage {
	return &l1SyncMessage{
		data: responseRollupInfoByBlockRange{
			blockRange: blockRange{
				fromBlock: fromBlock,
				toBlock:   toBlock,
			},
		},
		dataIsValid: true,
		ctrlIsValid: false,
	}
}

func newActionPackage(action actionsEnum) *l1SyncMessage {
	return &l1SyncMessage{
		dataIsValid: false,
		data: responseRollupInfoByBlockRange{
			blockRange: blockRange{
				fromBlock: 0,
				toBlock:   0,
			},
		},

		ctrlIsValid: true,
		ctrl: l1ConsumerControl{
			action: action,
		},
	}
}
