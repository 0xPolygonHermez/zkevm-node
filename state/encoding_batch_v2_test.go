package state

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/stretchr/testify/require"
)

const (
	// changeL2Block + deltaTimeStamp + indexL1InfoTree
	codedL2BlockHeader = "0b73e6af6f00000000"
	// 2 x [ tx coded in RLP + r,s,v,efficiencyPercentage]
	codedRLP2Txs1 = "ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bffee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff"
	// 2 x [ tx coded in RLP + r,s,v,efficiencyPercentage]
	codedRLP2Txs2 = "ee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880801cee7e01dc62f69a12c3510c6d64de04ee6346d84b6a017f3e786c7d87f963e75d8cc91fa983cd6d9cf55fff80d73bd26cd333b0f098acc1e58edb1fd484ad731bffee01843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880803ee20a0764440b016c4a2ee4e7e4eb3a5a97f1e6a6c9f40bf5ecf50f95ff636d63878ddb3e997e519826c7bb26fb7c5950a208e1ec722a9f1c568c4e479b40341cff"
	codedL2Block1 = codedL2BlockHeader + codedRLP2Txs1
	codedL2Block2 = codedL2BlockHeader + codedRLP2Txs2
)

func TestDecodeEmptyBatchV2(t *testing.T) {
	//batchL2Data, err := hex.DecodeString("0b73e6af6f00000000ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bffee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff")
	batchL2Data, err := hex.DecodeString("")
	require.NoError(t, err)

	batch, err := DecodeBatchV2(batchL2Data)
	require.NoError(t, err)
	require.Equal(t, 0, len(batch.Blocks))
}

func TestDecodeShorterBatchV2(t *testing.T) {
	//batchL2Data, err := hex.DecodeString("0b73e6af6f00000000ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bffee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff")
	batchL2Data, err := hex.DecodeString("0b010203")
	require.NoError(t, err)

	_, err = DecodeBatchV2(batchL2Data)
	log.Debug(err)
	require.ErrorIs(t, err, ErrInvalidBatchV2)
}

func TestDecodeWrongBatch(t *testing.T) {
	type testCase struct {
		name          string
		batchL2Data   string
		expectedError bool
	}
	testCases := []testCase{
		{
			name:          "batch dont start with 0x0b (changeL2Block)",
			batchL2Data:   "0c",
			expectedError: true,
		},
		{
			name:          "batch no enough  data to decode block.deltaTimestamp",
			batchL2Data:   "0b010203",
			expectedError: true,
		},
		{
			name:          "batch no enough  data to decode block.index",
			batchL2Data:   "0b01020304010203",
			expectedError: true,
		},
		{
			name:          "batch no enough  data to decode block.index",
			batchL2Data:   "0b01020304010203",
			expectedError: true,
		},
		{
			name:          "valid batch no trx, just L2Block",
			batchL2Data:   "0b0102030401020304",
			expectedError: false,
		},
		{
			name:          "invalid batch bad RLP codification",
			batchL2Data:   "0b" + "01020304" + "01020304" + "7f",
			expectedError: true,
		},
		{
			name: "1 block + 2 txs",
			batchL2Data: "0b" + "73e6af6f" + "00000000" +
				"ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bff" +
				"ee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff",
			expectedError: false,
		},
		{
			name: "1 block + 1 txs",
			batchL2Data: "0b" + "73e6af6f" + "00000000" +
				"ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bff",
			expectedError: false,
		},
		{
			name: "1 block + 1 txs, missiging efficiencyPercentage",
			batchL2Data: "0b" + "73e6af6f" + "00000000" +
				"ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51b",
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")
			data, err := hex.DecodeString(tc.batchL2Data)
			require.NoError(t, err)
			_, err = DecodeBatchV2(data)
			if err != nil {
				log.Debugf("[%s] %v", tc.name, err)
			}
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDecodeBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedL2Block1)
	require.NoError(t, err)
	batchL2Data2, err := hex.DecodeString(codedL2Block2)
	require.NoError(t, err)
	batch := append(batchL2Data, batchL2Data2...)
	decodedBatch, err := DecodeBatchV2(batch)
	require.NoError(t, err)
	require.Equal(t, 2, len(decodedBatch.Blocks))
	require.Equal(t, uint32(0x73e6af6f), decodedBatch.Blocks[0].DeltaTimestamp)
	require.Equal(t, uint32(0x00000000), decodedBatch.Blocks[0].IndexL1InfoTree)
}

func TestDecodeRLPLength(t *testing.T) {
	type testCase struct {
		name           string
		data           string
		expectedError  error
		expectedResult uint64
	}
	testCases := []testCase{
		{
			name:          "must start >= 0xc0",
			data:          "bf",
			expectedError: ErrInvalidRLP,
		},
		{
			name:           "shortRLP: c0 -> len=0",
			data:           "c0",
			expectedResult: 1,
		},
		{
			name:           "shortRLP: c1 -> len=1",
			data:           "c1",
			expectedResult: 2, // 1 byte header + 1 byte of data
		},
		{
			name:           "shortRLP: byte>0xf7",
			data:           "f7",
			expectedResult: 56, // 1 byte header + 55 bytes of data
		},
		{
			name:          "longRLP: f8: 1 extra byte, missing data",
			data:          "f8",
			expectedError: ErrInvalidRLP,
		},
		{
			name:           "longRLP: f8:size is stored in next byte ->0x01 (code add the length of bytes of the size??)",
			data:           "f8" + "01",
			expectedResult: 3, // 2 bytes of header + 1 byte of data
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")
			data, err := hex.DecodeString(tc.data)
			require.NoError(t, err)
			length, err := decodeRLPListLengthFromOffset(data, 0)
			if err != nil {
				log.Debugf("[%s] %v", tc.name, err)
			}
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, length)
			}
		})
	}
}

func TestEncodeBatchV2(t *testing.T) {
	block1 := L2BlockRaw{
		DeltaTimestamp:  123,
		IndexL1InfoTree: 456,
		Transactions:    []L2Tx{},
	}
	block2 := L2BlockRaw{
		DeltaTimestamp:  789,
		IndexL1InfoTree: 101112,
		Transactions:    []L2Tx{},
	}
	blocks := []L2BlockRaw{block1, block2}

	expectedBatchData := []byte{
		0xb, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x1, 0xc8, 0xb, 0x0, 0x0, 0x3, 0x15, 0x0, 0x1, 0x8a, 0xf8,
	}

	batchData, err := EncodeBatchV2(&BatchV2{Blocks: blocks})
	require.NoError(t, err)
	require.Equal(t, expectedBatchData, batchData)
}

func TestDecodeEncodeBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedL2Block1 + codedL2Block2)
	require.NoError(t, err)
	decodedBatch, err := DecodeBatchV2(batchL2Data)
	require.NoError(t, err)
	require.Equal(t, 2, len(decodedBatch.Blocks))
	encoded, err := EncodeBatchV2(decodedBatch)
	require.NoError(t, err)
	require.Equal(t, batchL2Data, encoded)
}

func TestEncodeEmptyBatchV2Fails(t *testing.T) {
	l2Batch := BatchV2{}
	_, err := EncodeBatchV2(&l2Batch)
	require.ErrorIs(t, err, ErrInvalidBatchV2)
	_, err = EncodeBatchV2(nil)
	require.ErrorIs(t, err, ErrInvalidBatchV2)
}

func TestDecodeForcedBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedRLP2Txs1)
	require.NoError(t, err)
	decodedBatch, err := DecodeForcedBatchV2(batchL2Data)
	require.NoError(t, err)
	require.Equal(t, 2, len(decodedBatch.Transactions))
}

func TestDecodeForcedBatchV2WithaNonForcedBatch(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedL2Block1)
	require.NoError(t, err)
	_, err = DecodeForcedBatchV2(batchL2Data)
	require.Error(t, err)
}
