package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SOR_Multicase(t *testing.T) {
	tcs := []struct {
		description             string
		lastBlock               uint64
		packages                []l1PackageData
		expected                []l1PackageData
		lastBlockOnSynchronizer uint64
	}{
		{
			description:             "empty_case",
			lastBlock:               100,
			packages:                []l1PackageData{},
			expected:                []l1PackageData{},
			lastBlockOnSynchronizer: 100,
		},
		{
			description:             "just_ctrl",
			lastBlock:               100,
			packages:                []l1PackageData{*newActionPackage(actionNone)},
			expected:                []l1PackageData{*newActionPackage(actionNone)},
			lastBlockOnSynchronizer: 100,
		},
		{
			description:             "just_br",
			lastBlock:               100,
			packages:                []l1PackageData{*newDataPackage(101, 119)},
			expected:                []l1PackageData{*newDataPackage(101, 119)},
			lastBlockOnSynchronizer: 119,
		},
		{
			description:             "just_br_missing_intermediate_block",
			lastBlock:               100,
			packages:                []l1PackageData{*newDataPackage(102, 119)},
			expected:                []l1PackageData{},
			lastBlockOnSynchronizer: 100,
		},
		{
			description: "inverse_br",
			lastBlock:   100,
			packages: []l1PackageData{
				*newDataPackage(131, 141),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []l1PackageData{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
				*newDataPackage(131, 141),
			},
			lastBlockOnSynchronizer: 141,
		},
		{
			description: "crtl_linked_to_br",
			lastBlock:   100,
			packages: []l1PackageData{
				*newDataPackage(131, 141),
				*newActionPackage(actionNone),
				*newDataPackage(120, 130),
				*newDataPackage(101, 119)},
			expected: []l1PackageData{
				*newDataPackage(101, 119),
				*newDataPackage(120, 130),
				*newDataPackage(131, 141),
				*newActionPackage(actionNone),
			},
			lastBlockOnSynchronizer: 141,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			ch := createChannel(100)
			sut := newFilterToSendOrdererResultsToConsumer(ch, tc.lastBlock)
			for _, p := range tc.packages {
				sut.addResultAndSendToConsumer(&p)
			}
			sendData := getAllDataFromChannel(ch)
			require.Equal(t, tc.expected, sendData)
			require.Equal(t, tc.lastBlockOnSynchronizer, sut.lastBlockOnSynchronizer)
		})
	}

}

func Test_SOR_ControlDataPackage(t *testing.T) {
	ch := createChannel(100)
	lastBlock := uint64(100)
	sut := newFilterToSendOrdererResultsToConsumer(ch, lastBlock)
	result := newL1PackageDataControl(actionNone)
	sut.addResultAndSendToConsumer(result)
	require.Equal(t, 0, len(sut.pendingResults), "an control package must be send is order, so it must be send immediately")
	require.Equal(t, uint64(100), sut.lastBlockOnSynchronizer, "a control package doesn't change lastBlockOnSynchronizer")
}

func Test_SOR_TrivialCaseThatArrivesNextBlock(t *testing.T) {
	ch := createChannel(100)
	lastBlock := uint64(100)
	sut := newFilterToSendOrdererResultsToConsumer(ch, lastBlock)
	result := newDataPackage(101, 110)
	sut.addResultAndSendToConsumer(result)
	require.Equal(t, 0, len(sut.pendingResults))
	require.Equal(t, uint64(110), sut.lastBlockOnSynchronizer)
}

func Test_SOR_ReceivedABlockThatIsNotNextOne(t *testing.T) {
	ch := createChannel(100)
	lastBlock := uint64(100)
	sut := newFilterToSendOrdererResultsToConsumer(ch, lastBlock)
	result := newDataPackage(111, 120)
	sut.addResultAndSendToConsumer(result)
	require.Equal(t, 1, len(sut.pendingResults))
	require.Equal(t, lastBlock, sut.lastBlockOnSynchronizer)
}

func Test_SOR_ControlDataPackageBlockedAtTheEndBecauseThereAreAMissingBR(t *testing.T) {
	ch := createChannel(100)
	lastBlock := uint64(100)
	sut := newFilterToSendOrdererResultsToConsumer(ch, lastBlock)
	sut.addResultAndSendToConsumer(newDataPackage(111, 120))
	sut.addResultAndSendToConsumer(newDataPackage(121, 130))
	sut.addResultAndSendToConsumer(newDataPackage(131, 140))
	sut.addResultAndSendToConsumer(newActionPackage(actionNone))
	require.Equal(t, 4, len(sut.pendingResults))
	require.Equal(t, lastBlock, sut.lastBlockOnSynchronizer)

	sut.addResultAndSendToConsumer(newDataPackage(101, 110))
	require.Equal(t, 0, len(sut.pendingResults))
	require.Equal(t, uint64(140), sut.lastBlockOnSynchronizer)

	sendData := getAllDataFromChannel(ch)
	last := sendData[len(sendData)-1]
	require.Equal(t, newL1PackageDataControl(actionNone).toStringBrief(), last.toStringBrief())
}

func getAllDataFromChannel[T any](ch chan T) []T {
	res := []T{}
	for len(ch) > 0 {
		res = append(res, <-ch)
	}
	return res
}

func Test_SOR_ThereAreSomePendingBlocksAndArriveTheMissingOne(t *testing.T) {
	ch := createChannel(100)
	lastBlock := uint64(100)
	sut := newFilterToSendOrdererResultsToConsumer(ch, lastBlock)
	sut.addResultAndSendToConsumer(newDataPackage(111, 120))
	sut.addResultAndSendToConsumer(newDataPackage(121, 130))
	sut.addResultAndSendToConsumer(newDataPackage(131, 140))
	require.Equal(t, 3, len(sut.pendingResults))
	require.Equal(t, lastBlock, sut.lastBlockOnSynchronizer)

	sut.addResultAndSendToConsumer(newDataPackage(101, 110))
	require.Equal(t, 0, len(sut.pendingResults))
	require.Equal(t, uint64(140), sut.lastBlockOnSynchronizer)
}

func createChannel(size int) chan l1PackageData {
	return make(chan l1PackageData, size)
}

func newDataPackage(fromBlock, toBlock uint64) *l1PackageData {
	return &l1PackageData{
		data: getRollupInfoByBlockRangeResult{
			blockRange: blockRange{
				fromBlock: fromBlock,
				toBlock:   toBlock,
			},
		},
		dataIsValid: true,
		ctrlIsValid: false,
	}
}

func newActionPackage(action actionsEnum) *l1PackageData {
	return &l1PackageData{
		dataIsValid: false,
		data: getRollupInfoByBlockRangeResult{
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
