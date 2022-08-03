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

			require.Equal(t, tc.expectedContainer, actualContainer)
		})
	}
}

func TestFindValue(t *testing.T) {
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
			description: "happy path, first item requires empty root",
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
			oldRoot:         "",
			expectedValue:   "value1",
			expectedNewRoot: "root1",
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

			actualValue, actualRoot, err := subject.FindValue(tc.key, tc.oldRoot)
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
								Key:      "key1",
								Bytecode: "bytecode1",
							},
							{
								Key:      "key2",
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
								Key:      "key3",
								Bytecode: "bytecode3",
							},
							{
								Key:      "key4",
								Bytecode: "bytecode4",
							},
						},
					},
					{
						GenesisRaw: []*state.GenesisAction{
							{
								Key:      "key1",
								Bytecode: "bytecode1",
							},
							{
								Key:      "key2",
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
								Key:      "key1",
								Bytecode: "bytecode1",
							},
							{
								Key:      "key2",
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
								Key:   "key1",
								Value: "value1",
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
