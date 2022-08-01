package testvector_test

import (
	"fmt"
	"path/filepath"
	"testing"

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
  "stateTransition": {
    "genesis": [
      {
        "address": "address1",
        "balance": "balance1",
        "nonce": "nonce1",
        "bytecode": "bytecode1",
        "storage": {
          "storageKey1_1": "storageValue1_1",
          "storageKey1_2": "storageValue1_2"
        }
      },
      {
        "address": "address2",
        "balance": "balance2",
        "nonce": "nonce2",
        "bytecode": "bytecode2",
        "storage": {
          "storageKey2_1": "storageValue2_1",
          "storageKey2_2": "storageValue2_2"
        }
      }
    ]
  },
  "contractsBytecode": {
    "contract1Key": "contract1Bytecode",
    "contract2Key": "contract2Bytecode"
  },
  "genesisRaw": {
    "keys": [
      "a","b","c"
    ],
    "values": [
      "1","2","3"
    ],
    "expectedRoots": [
      "root1", "root2", "root3"
    ]
  }
}
]`,
			},
			testVectorPath: defaultSourceDir,
			expectedContainer: &testvector.Container{
				E2E: &testvector.E2E{
					Items: []*testvector.E2EItem{
						{
							StateTransition: &testvector.StateTransition{
								Genesis: []*testvector.GenesisItem{
									{
										Address:  "address1",
										Balance:  "balance1",
										Nonce:    "nonce1",
										Bytecode: "bytecode1",
										Storage: map[string]string{
											"storageKey1_1": "storageValue1_1",
											"storageKey1_2": "storageValue1_2",
										},
									},
									{
										Address:  "address2",
										Balance:  "balance2",
										Nonce:    "nonce2",
										Bytecode: "bytecode2",
										Storage: map[string]string{
											"storageKey2_1": "storageValue2_1",
											"storageKey2_2": "storageValue2_2",
										},
									},
								},
							},
							ContractsBytecode: map[string]string{
								"contract1Key": "contract1Bytecode",
								"contract2Key": "contract2Bytecode",
							},
							GenesisRaw: &testvector.GenesisRaw{
								Keys:          []string{"a", "b", "c"},
								Values:        []string{"1", "2", "3"},
								ExpectedRoots: []string{"root1", "root2", "root3"},
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
  "stateTransition": {
    "genesis": [
      {
        "address": "address1",
        "balance": "balance1",
        "nonce": "nonce1",
        "bytecode": "bytecode1",
        "storage": {
          "storageKey1_1": "storageValue1_1",
          "storageKey1_2": "storageValue1_2"
        }
      }
    ]
  },
  "contractsBytecode": {
    "contract1Key": "contract1Bytecode"
  },
  "genesisRaw": {
    "keys": [
      "a"
    ],
    "values": [
      "1"
    ],
    "expectedRoots": [
      "root1"
    ]
  }
}
]`,
				filepath.Join(defaultSourceDir, "b.json"): `[
{
  "stateTransition": {
    "genesis": [
      {
        "address": "address2",
        "balance": "balance2",
        "nonce": "nonce2",
        "bytecode": "bytecode2",
        "storage": {
          "storageKey2_1": "storageValue2_1",
          "storageKey2_2": "storageValue2_2"
        }
      }
    ]
  },
  "contractsBytecode": {
    "contract2Key": "contract2Bytecode"
  },
  "genesisRaw": {
    "keys": [
      "b"
    ],
    "values": [
      "2"
    ],
    "expectedRoots": [
      "root2"
    ]
  }
}
]`,
			},
			testVectorPath: defaultSourceDir,
			expectedContainer: &testvector.Container{
				E2E: &testvector.E2E{
					Items: []*testvector.E2EItem{
						{
							StateTransition: &testvector.StateTransition{
								Genesis: []*testvector.GenesisItem{
									{
										Address:  "address1",
										Balance:  "balance1",
										Nonce:    "nonce1",
										Bytecode: "bytecode1",
										Storage: map[string]string{
											"storageKey1_1": "storageValue1_1",
											"storageKey1_2": "storageValue1_2",
										},
									},
								},
							},
							ContractsBytecode: map[string]string{
								"contract1Key": "contract1Bytecode",
							},
							GenesisRaw: &testvector.GenesisRaw{
								Keys:          []string{"a"},
								Values:        []string{"1"},
								ExpectedRoots: []string{"root1"},
							},
						},
						{
							StateTransition: &testvector.StateTransition{
								Genesis: []*testvector.GenesisItem{
									{
										Address:  "address2",
										Balance:  "balance2",
										Nonce:    "nonce2",
										Bytecode: "bytecode2",
										Storage: map[string]string{
											"storageKey2_1": "storageValue2_1",
											"storageKey2_2": "storageValue2_2",
										},
									},
								},
							},
							ContractsBytecode: map[string]string{
								"contract2Key": "contract2Bytecode",
							},
							GenesisRaw: &testvector.GenesisRaw{
								Keys:          []string{"b"},
								Values:        []string{"2"},
								ExpectedRoots: []string{"root2"},
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

			require.Equal(t, tc.expectedContainer, actualContainer)
		})
	}
}

func TestFindE2EGenesisRaw(t *testing.T) {
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
						GenesisRaw: &testvector.GenesisRaw{
							Keys:          []string{"key1", "key2"},
							Values:        []string{"value1", "value2"},
							ExpectedRoots: []string{"root1", "root2"},
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
						GenesisRaw: &testvector.GenesisRaw{
							Keys:          []string{"key3", "key4"},
							Values:        []string{"value3", "value4"},
							ExpectedRoots: []string{"root3", "root4"},
						},
					},
					{
						GenesisRaw: &testvector.GenesisRaw{
							Keys:          []string{"key1", "key2"},
							Values:        []string{"value1", "value2"},
							ExpectedRoots: []string{"root1", "root2"},
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
			description: "happy path, first item requires empty root",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: &testvector.GenesisRaw{
							Keys:          []string{"key1", "key2"},
							Values:        []string{"value1", "value2"},
							ExpectedRoots: []string{"root1", "root2"},
						},
					},
				},
			},
			key:             "key1",
			oldRoot:         "",
			expectedValue:   "value1",
			expectedNewRoot: "root1",
		},
		{
			description: "happy path, bytecode",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						ContractsBytecode: map[string]string{
							"key1": "value1",
						},
					},
				},
			},
			key:             "key1",
			oldRoot:         "",
			expectedValue:   "value1",
			expectedNewRoot: "",
		},
		{
			description: "unexisting key gives error",
			e2e: &testvector.E2E{
				Items: []*testvector.E2EItem{
					{
						GenesisRaw: &testvector.GenesisRaw{
							Keys:          []string{"key1", "key2"},
							Values:        []string{"value1", "value2"},
							ExpectedRoots: []string{"root1", "root2"},
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
						GenesisRaw: &testvector.GenesisRaw{
							Keys:          []string{"key1", "key2"},
							Values:        []string{"value1", "value2"},
							ExpectedRoots: []string{"root1", "root2"},
						},
					},
				},
			},
			key:              "key1",
			oldRoot:          "root2",
			expectedError:    true,
			expectedErrorMsg: `key "key1" not found for oldRoot "root2"`,
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

			actualValue, actualRoot, err := subject.FindE2EGenesisRaw(tc.key, tc.oldRoot)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			if err == nil {
				require.Equal(t, tc.expectedValue, actualValue)
				require.Equal(t, tc.expectedNewRoot, actualRoot)
			}
		})
	}
}
