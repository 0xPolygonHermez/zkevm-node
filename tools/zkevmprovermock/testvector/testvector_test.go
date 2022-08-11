package testvector_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/testvector"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewContainer(t *testing.T) {
	const defaultSourceDir = "/a/b/c"

	tcs := []struct {
		description       string
		sourceFiles       map[string]string
		testVectorPath    string
		expectedContainer *testvector.Container
		expectedError     bool
		expectedErrorMsg  string
	}{
		{
			description: "happy path, single file",
			sourceFiles: map[string]string{
				filepath.Join(defaultSourceDir, "a.json"): `[
{
  "batchL2Data": "0xabc123456",
  "globalExitRoot": "0x1234abcd",
  "traces": {
    "batchHash": "batchHash",
    "old_state_root": "old_state_root",
    "globalHash": "globalHash",
    "numBatch": 1,
    "timestamp": 1944498031,
    "sequencerAddr": "sequencerAddr",
    "responses": [
      {
        "tx_hash": "0xabc",
        "type": 0,
        "gas_left": "28099",
        "gas_used": "71901",
        "gas_refunded": "0",
        "state_root": "0x2031e0233b733481aa0e8c1056b874d68731fc0c673248538f7acfb81d2d7764",
        "logs": [
          {
            "data": [
              "000000000000000000000002540be400"
            ],
            "topics": [
              "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
              "00000000000000000000000000000000",
              "4d5cf5032b2a844602278b01199ed191a86c93ff"
            ],
            "address": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
            "batch_number": 1,
            "tx_hash": "0xeeb51664fd2b6dcf865de752589f59f29b8398bdc38b5f556715ae88615c4641",
            "tx_index": 0,
            "batch_hash": "0x7624c022e923e798a6682171fd27b54912c426cfd92d49fb6be3bf300ba27679",
            "index": 0
          }
        ],
        "unprocessed_transaction": false,
        "call_trace": {
          "context": {
            "from": "0x617b3a3528f9cdd6630fd3301b9c8911f7bf063d",
            "to": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
            "type": "CALL",
            "data": "0x40c10f190000000000000000000000004d5cf5032b2a844602278b01199ed191a86c93ff00000000000000000000000000000000000000000000000000000002540be400",
            "gas": "100000",
            "value": "0",
            "batch": "0x7624c022e923e798a6682171fd27b54912c426cfd92d49fb6be3bf300ba27679",
            "output": "",
            "gas_used": "71901",
            "execution_time": "",
            "old_state_root": "0x2031e0233b733481aa0e8c1056b874d68731fc0c673248538f7acfb81d2d7764",
            "nonce": 0,
            "gasPrice": "1000000000",
            "chainId": 1000,
            "return_value": []
          },
          "steps": [
            {
              "depth": 1,
              "pc": 0,
              "remaining_gas": "78392",
              "opcode": "PUSH1",
              "gas_refund": "0",
              "op": "0x60",
              "error": "",
              "state_root": "0xdf403125ab76e36f1fb313619e704e4630b984bd93c3121b3775b71d1b72f9c6",
              "contract": {
                "address": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
                "caller": "0x617b3a3528f9cdd6630fd3301b9c8911f7bf063d",
                "value": "0",
                "data": "0x40c10f190000000000000000000000004d5cf5032b2a844602278b01199ed191a86c93ff00000000000000000000000000000000000000000000000000000002540be40000000000000000000000000000000000000000000000000000000000",
                "gas": "100000"
              },
              "return_data": [],
              "gas_cost": "3",
              "stack": [
                "0x80"
              ],
              "memory": []
            }
          ]
        }
      }
    ],
    "cumulative_gas_used": "171380",
    "new_state_root": "0xb6f5c8f596b7130c7c7eb3257e80badfd7a89aa2977cc6d521a8b97f764230f5",
    "new_local_exit_root": "0x00",
    "cnt_keccak_hashes": 1,
    "cnt_poseidon_hashes": 1,
    "cnt_poseidon_paddings": 1,
    "cnt_mem_aligns": 1,
    "cnt_arithmetics": 1,
    "cnt_binaries": 1,
    "cnt_steps": 1
  },
  "genesisRaw": [
    {
      "address": "addressRaw0",
      "type": 0,
      "key": "keyRaw0",
      "value": "valueRaw0"
    },
    {
      "address": "addressRaw1",
      "type": 1,
      "key": "keyRaw1",
      "value": "valueRaw1"
    },
    {
      "address": "addressRaw2",
      "type": 2,
      "key": "keyRaw2",
      "value": "valueRaw2",
      "bytecode": "bytecodeRaw2"
    },
    {
      "address": "addressRaw3",
      "type": 3,
      "key": "keyRaw3",
      "value": "valueRaw3",
      "storagePosition": "storagePositionRaw3"
    },
    {
      "address": "addressRaw4",
      "type": 4,
      "key": "keyRaw4",
      "value": "valueRaw4"
    }
  ]
}
]`,
			},
			testVectorPath: defaultSourceDir,
			expectedContainer: &testvector.Container{
				E2E: &testvector.E2E{
					Items: []*testvector.E2EItem{
						{
							BatchL2Data:    "0xabc123456",
							GlobalExitRoot: "0x1234abcd",
							Traces: &testvector.Traces{
								BatchHash:     "batchHash",
								OldStateRoot:  "old_state_root",
								GlobalHash:    "globalHash",
								NumBatch:      1,
								Timestamp:     1944498031,
								SequencerAddr: "sequencerAddr",
								ProcessBatchResponse: &testvector.ProcessBatchResponse{
									Responses: []*testvector.ProcessTransactionResponse{
										{
											TxHash:      "0xabc",
											Type:        0,
											GasLeft:     "28099",
											GasUsed:     "71901",
											GasRefunded: "0",
											StateRoot:   "0x2031e0233b733481aa0e8c1056b874d68731fc0c673248538f7acfb81d2d7764",
											Logs: []*testvector.Log{
												{
													Data: []string{"000000000000000000000002540be400"},
													Topics: []string{
														"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
														"00000000000000000000000000000000",
														"4d5cf5032b2a844602278b01199ed191a86c93ff",
													},
													Address:     "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
													BatchNumber: 1,
													TxHash:      "0xeeb51664fd2b6dcf865de752589f59f29b8398bdc38b5f556715ae88615c4641",
													TxIndex:     0,
													BatchHash:   "0x7624c022e923e798a6682171fd27b54912c426cfd92d49fb6be3bf300ba27679",
													Index:       0,
												},
											},
											UnprocessedTransaction: false,
											CallTrace: &testvector.CallTrace{
												Context: &testvector.TransactionContext{
													From:          "0x617b3a3528f9cdd6630fd3301b9c8911f7bf063d",
													To:            "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
													Type:          "CALL",
													Data:          "0x40c10f190000000000000000000000004d5cf5032b2a844602278b01199ed191a86c93ff00000000000000000000000000000000000000000000000000000002540be400",
													Gas:           "100000",
													Value:         "0",
													Batch:         "0x7624c022e923e798a6682171fd27b54912c426cfd92d49fb6be3bf300ba27679",
													Output:        "",
													GasUsed:       "71901",
													ExecutionTime: "",
													OldStateRoot:  "0x2031e0233b733481aa0e8c1056b874d68731fc0c673248538f7acfb81d2d7764",
													GasPrice:      "1000000000",
												},

												Steps: []*testvector.TransactionStep{
													{
														Depth:        1,
														Pc:           0,
														RemainingGas: "78392",
														OpCode:       "PUSH1",
														GasRefund:    "0",
														Op:           "0x60",
														StateRoot:    "0xdf403125ab76e36f1fb313619e704e4630b984bd93c3121b3775b71d1b72f9c6",
														Contract: &testvector.Contract{
															Address: "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
															Caller:  "0x617b3a3528f9cdd6630fd3301b9c8911f7bf063d",
															Value:   "0",
															Data:    "0x40c10f190000000000000000000000004d5cf5032b2a844602278b01199ed191a86c93ff00000000000000000000000000000000000000000000000000000002540be40000000000000000000000000000000000000000000000000000000000",
															Gas:     "100000",
														},
														ReturnData: []string{},
														GasCost:    "3",
														Stack:      []string{"0x80"},
														Memory:     []string{},
													},
												},
											},
										},
									},
									CumulativeGasUsed:   "171380",
									NewStateRoot:        "0xb6f5c8f596b7130c7c7eb3257e80badfd7a89aa2977cc6d521a8b97f764230f5",
									NewLocalExitRoot:    "0x00",
									CntKeccakHashes:     1,
									CntPoseidonHashes:   1,
									CntPoseidonPaddings: 1,
									CntMemAligns:        1,
									CntArithmetics:      1,
									CntBinaries:         1,
									CntSteps:            1,
								},
							},
							GenesisRaw: []*state.GenesisAction{
								{
									Address: "addressRaw0",
									Type:    0,
									Key:     "keyRaw0",
									Value:   "valueRaw0",
								},
								{
									Address: "addressRaw1",
									Type:    1,
									Key:     "keyRaw1",
									Value:   "valueRaw1",
								},
								{
									Address:  "addressRaw2",
									Type:     2,
									Key:      "keyRaw2",
									Value:    "valueRaw2",
									Bytecode: "bytecodeRaw2",
								},
								{
									Address:         "addressRaw3",
									Type:            3,
									Key:             "keyRaw3",
									Value:           "valueRaw3",
									StoragePosition: "storagePositionRaw3",
								},
								{
									Address: "addressRaw4",
									Type:    4,
									Key:     "keyRaw4",
									Value:   "valueRaw4",
								},
							},
						},
					},
				},
			},
		},
		{
			description: "happy path, multiple files",
			sourceFiles: map[string]string{
				filepath.Join(defaultSourceDir, "a.json"): `[
{
  "genesisRaw": [
    {
      "address": "addressRaw0",
      "type": 0,
      "key": "keyRaw0",
      "value": "valueRaw0"
    }
  ]
}
]`,
				filepath.Join(defaultSourceDir, "b.json"): `[
{
  "genesisRaw": [
    {
      "address": "addressRaw1",
      "type": 1,
      "key": "keyRaw1",
      "value": "valueRaw1"
    }
  ]
}
]`,
			},
			testVectorPath: defaultSourceDir,
			expectedContainer: &testvector.Container{
				E2E: &testvector.E2E{
					Items: []*testvector.E2EItem{
						{
							GenesisRaw: []*state.GenesisAction{
								{
									Address: "addressRaw0",
									Type:    0,
									Key:     "keyRaw0",
									Value:   "valueRaw0",
								},
							},
						},
						{
							GenesisRaw: []*state.GenesisAction{
								{
									Address: "addressRaw1",
									Type:    1,
									Key:     "keyRaw1",
									Value:   "valueRaw1",
								},
							},
						},
					},
				},
			},
		},
		{
			description: "invalid test vector causes error",
			sourceFiles: map[string]string{
				filepath.Join(defaultSourceDir, "a.json"): "not a real json",
			},
			testVectorPath:   defaultSourceDir,
			expectedError:    true,
			expectedErrorMsg: "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			description:      "empty test vector path returns empty object",
			sourceFiles:      map[string]string{},
			testVectorPath:   defaultSourceDir,
			expectedError:    true,
			expectedErrorMsg: fmt.Sprintf("open %s: file does not exist", defaultSourceDir),
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			var appFs = afero.NewMemMapFs()

			require.NoError(t, testutils.CreateTestFiles(appFs, tc.sourceFiles))

			actualContainer, err := testvector.NewContainer(tc.testVectorPath, appFs)

			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			if err == nil {
				require.Equal(t, tc.expectedContainer.E2E.Items, actualContainer.E2E.Items)
			}
		})
	}
}

func TestFindSMTValue(t *testing.T) {
	tcs := []struct {
		description      string
		e2e              *testvector.E2E
		key              string
		oldRoot          string
		expectedValue    string
		expectedNewRoot  string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description: "happy path, single item",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key1",
								Value: "value1",
								Root:  "root1",
							},
							{
								Key:   "key2",
								Value: "value2",
								Root:  "root2",
							},
						},
					},
				},
			},
			key:             "key2",
			oldRoot:         "root1",
			expectedValue:   "value2",
			expectedNewRoot: "root2",
		},
		{
			description: "happy path, multiple items",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key3",
								Value: "value3",
								Root:  "root3",
							},
							{
								Key:   "key4",
								Value: "value4",
								Root:  "root4",
							},
						},
					},
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key1",
								Value: "value1",
								Root:  "root1",
							},
							{
								Key:   "key2",
								Value: "value2",
								Root:  "root2",
							},
						},
					},
				},
			},
			key:             "key2",
			oldRoot:         "root1",
			expectedValue:   "value2",
			expectedNewRoot: "root2",
		},
		{
			description: "happy path, first item requires zero root",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key1",
								Value: "value1",
								Root:  "root1",
							},
							{
								Key:   "key2",
								Value: "value2",
								Root:  "root2",
							},
						},
					},
				},
			},
			key:             "key1",
			oldRoot:         "0x0000000000000000000000000000000000000000000000000000000000000000",
			expectedValue:   "value1",
			expectedNewRoot: "root1",
		},
		{
			description: "happy path, querying existing key and last root returns last value of the key",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key1",
								Value: "value1",
								Root:  "root1",
							},
							{
								Key:   "key2",
								Value: "value2",
								Root:  "root2",
							},
							{
								Key:   "key1",
								Value: "value3",
								Root:  "root2",
							},
							{
								Key:   "key4",
								Value: "value4",
								Root:  "root4",
							},
						},
					},
				},
			},
			key:             "key1",
			oldRoot:         "root4",
			expectedValue:   "value3",
			expectedNewRoot: "root4",
		},
		{
			description: "unexisting key gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key1",
								Value: "value1",
								Root:  "root1",
							},
							{
								Key:   "key2",
								Value: "value2",
								Root:  "root2",
							},
						},
					},
				},
			},
			key:              "key10",
			oldRoot:          "root1",
			expectedError:    true,
			expectedErrorMsg: `key "key10" not found for oldRoot "root1"`,
		},
		{
			description: "unmatched root gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:   "key1",
								Value: "value1",
								Root:  "root1",
							},
							{
								Key:   "key2",
								Value: "value2",
								Root:  "root2",
							},
						},
					},
				},
			},
			key:              "key1",
			oldRoot:          "root3",
			expectedError:    true,
			expectedErrorMsg: `key "key1" not found for oldRoot "root3"`,
		},
		{
			description: "empty GenesisRaw gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{},
				},
			},
			key:              "key1",
			oldRoot:          "root2",
			expectedError:    true,
			expectedErrorMsg: `key "key1" not found for oldRoot "root2"`,
		},
	}

	subject := &testvector.Container{}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			subject.E2E = tc.e2e

			actualValue, actualRoot, err := subject.FindSMTValue(tc.key, tc.oldRoot)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			if err == nil {
				require.Equal(t, tc.expectedValue, actualValue)
				require.Equal(t, tc.expectedNewRoot, actualRoot)
			}
		})
	}
}

func TestFindBytecode(t *testing.T) {
	tcs := []struct {
		description      string
		e2e              *testvector.E2E
		key              string
		expectedBytecode string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description: "happy path, single item",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Value:    "key1",
								Bytecode: "bytecode1",
							},
							{
								Value:    "key2",
								Bytecode: "bytecode2",
							},
						},
					},
				},
			},
			key:              "key2",
			expectedBytecode: "bytecode2",
		},
		{
			description: "happy path, multiple items",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Value:    "key3",
								Bytecode: "bytecode3",
							},
							{
								Value:    "key4",
								Bytecode: "bytecode4",
							},
						},
					},
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Value:    "key1",
								Bytecode: "bytecode1",
							},
							{
								Value:    "key2",
								Bytecode: "bytecode2",
							},
						},
					},
				},
			},
			key:              "key2",
			expectedBytecode: "bytecode2",
		},
		{
			description: "unexisting key gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Value:    "key1",
								Bytecode: "bytecode1",
							},
							{
								Value:    "key2",
								Bytecode: "bytecode2",
							},
						},
					},
				},
			},
			key:              "key10",
			expectedError:    true,
			expectedErrorMsg: `bytecode for key "key10" not found`,
		},
		{
			description: "empty GenesisRaw gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{},
				},
			},
			key:              "key1",
			expectedError:    true,
			expectedErrorMsg: `bytecode for key "key1" not found`,
		},
		{
			description: "empty bytecode for matching key gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Value: "key1",
							},
						},
					},
				},
			},
			key:              "key1",
			expectedError:    true,
			expectedErrorMsg: `bytecode for key "key1" not found`,
		},
	}

	subject := &testvector.Container{}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			subject.E2E = tc.e2e

			actualBytecode, err := subject.FindBytecode(tc.key)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			if err == nil {
				require.Equal(t, tc.expectedBytecode, actualBytecode)
			}
		})
	}
}

func TestFindProcessBatchResponse(t *testing.T) {
	tcs := []struct {
		description      string
		e2e              *testvector.E2E
		batchL2Data      string
		expectedResponse *testvector.ProcessBatchResponse
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description: "happy path, single item",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						BatchL2Data: "0xabc",
						Traces: &testvector.Traces{
							ProcessBatchResponse: &testvector.ProcessBatchResponse{
								CumulativeGasUsed: "100",
								CntKeccakHashes:   200,
								CntMemAligns:      300,
							},
						},
					},
				},
			},
			batchL2Data: "0xabc",
			expectedResponse: &testvector.ProcessBatchResponse{
				CumulativeGasUsed: "100",
				CntKeccakHashes:   200,
				CntMemAligns:      300,
			},
		},
		{
			description: "happy path, no leading 0x in id",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						BatchL2Data: "0xabc",
						Traces: &testvector.Traces{
							ProcessBatchResponse: &testvector.ProcessBatchResponse{
								CumulativeGasUsed: "100",
								CntKeccakHashes:   200,
								CntMemAligns:      300,
							},
						},
					},
				},
			},
			batchL2Data: "abc",
			expectedResponse: &testvector.ProcessBatchResponse{
				CumulativeGasUsed: "100",
				CntKeccakHashes:   200,
				CntMemAligns:      300,
			},
		},
		{
			description: "happy path, multiple item",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						BatchL2Data: "0xabc1234",
						Traces: &testvector.Traces{
							ProcessBatchResponse: &testvector.ProcessBatchResponse{
								CumulativeGasUsed: "1100",
								CntKeccakHashes:   1200,
								CntMemAligns:      1300,
							},
						},
					},
					{
						BatchL2Data: "0xabc",
						Traces: &testvector.Traces{
							ProcessBatchResponse: &testvector.ProcessBatchResponse{
								CumulativeGasUsed: "100",
								CntKeccakHashes:   200,
								CntMemAligns:      300,
							},
						},
					},
				},
			},
			batchL2Data: "0xabc",
			expectedResponse: &testvector.ProcessBatchResponse{
				CumulativeGasUsed: "100",
				CntKeccakHashes:   200,
				CntMemAligns:      300,
			},
		},
		{
			description: "unhappy path, id not found",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						BatchL2Data: "0xabc",
						Traces: &testvector.Traces{
							ProcessBatchResponse: &testvector.ProcessBatchResponse{
								CumulativeGasUsed: "100",
								CntKeccakHashes:   200,
								CntMemAligns:      300,
							},
						},
					},
				},
			},
			batchL2Data:      "0x123",
			expectedError:    true,
			expectedErrorMsg: `ProcessBatchResponse for batchL2Data "0x123" not found`,
		},
	}

	subject := &testvector.Container{}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			subject.E2E = tc.e2e

			actualResponse, err := subject.FindProcessBatchResponse(tc.batchL2Data)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			if err == nil {
				require.Equal(t, tc.expectedResponse, actualResponse)
			}
		})
	}
}
