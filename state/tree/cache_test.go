package tree

import (
	"sync"
	"testing"

	"github.com/hermeznetwork/hermez-core/test/testutils"
	"github.com/stretchr/testify/require"
)

func TestMTNodeCacheGet(t *testing.T) {
	tcs := []struct {
		description    string
		data           map[string][]uint64
		key            [][]uint64
		expected       [][]uint64
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			description: "single matching item",
			data: map[string][]uint64{
				"0x0000000000000001000000000000000100000000000000010000000000000001": {15},
			},
			key:      [][]uint64{{1, 1, 1, 1}},
			expected: [][]uint64{{15}},
		},
		{
			description: "single non-matching item",
			data: map[string][]uint64{
				"0x0000000000000001000000000000000100000000000000010000000000000001": {15},
			},
			key:            [][]uint64{{1, 1, 1, 0}},
			expectedErr:    true,
			expectedErrMsg: errMTNodeCacheItemNotFound.Error(),
		},
		{
			description: "multiple matching items",
			data: map[string][]uint64{
				"0x0000000000000001000000000000000100000000000000010000000000000001": {15},
				"0x0000000000000001000000000000000100000000000000010000000000000002": {16},
				"0x0000000000000001000000000000000100000000000000010000000000000003": {17},
			},
			key:      [][]uint64{{1, 1, 1, 1}, {2, 1, 1, 1}, {3, 1, 1, 1}},
			expected: [][]uint64{{15}, {16}, {17}},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			subject := newNodeCache()

			subject.data = tc.data

			for i := 0; i < len(tc.key); i++ {
				actual, err := subject.get(tc.key[i])
				require.NoError(t, testutils.CheckError(err, tc.expectedErr, tc.expectedErrMsg))

				if !tc.expectedErr {
					require.Equal(t, tc.expected[i], actual)
				}
			}
		})
	}
}

func TestMTNodeCacheSet(t *testing.T) {
	tcs := []struct {
		description    string
		key            [][]uint64
		value          [][]uint64
		expectedData   map[string][]uint64
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			description: "single item set",
			key:         [][]uint64{{1, 1, 1, 1}},
			value:       [][]uint64{{15}},
			expectedData: map[string][]uint64{
				"0x0000000000000001000000000000000100000000000000010000000000000001": {15},
			},
		},
		{
			description: "mutiple items set",
			key:         [][]uint64{{1, 1, 1, 1}, {1, 1, 1, 2}, {1, 2, 1, 3}},
			value:       [][]uint64{{15}, {16}, {17}},
			expectedData: map[string][]uint64{
				"0x0000000000000001000000000000000100000000000000010000000000000001": {15},
				"0x0000000000000002000000000000000100000000000000010000000000000001": {16},
				"0x0000000000000003000000000000000100000000000000020000000000000001": {17},
			},
		},
		{
			description: "keys can be updated",
			key:         [][]uint64{{1, 1, 1, 1}, {1, 1, 1, 2}, {1, 1, 1, 1}},
			value:       [][]uint64{{15}, {16}, {1500}},
			expectedData: map[string][]uint64{
				"0x0000000000000001000000000000000100000000000000010000000000000001": {1500},
				"0x0000000000000002000000000000000100000000000000010000000000000001": {16},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			subject := newNodeCache()
			for i := 0; i < len(tc.key); i++ {
				err := subject.set(tc.key[i], tc.value[i])
				require.NoError(t, testutils.CheckError(err, tc.expectedErr, tc.expectedErrMsg))
			}
			if !tc.expectedErr {
				require.Equal(t, tc.expectedData, subject.data)
			}
		})
	}
}

func TestMTNodeCacheMaxItems(t *testing.T) {
	subject := newNodeCache()
	for i := 0; i < maxMTNodeCacheEntries; i++ {
		err := subject.set([]uint64{1, 1, 1, uint64(i)}, []uint64{1})
		require.NoError(t, err)
	}

	err := subject.set([]uint64{1, 1, 1, uint64(maxMTNodeCacheEntries + 1)}, []uint64{1})
	require.Error(t, err)

	require.Equal(t, "MT node cache is full", err.Error())
}

func TestMTNodeCacheClear(t *testing.T) {
	subject := newNodeCache()
	for i := 0; i < 5; i++ {
		err := subject.set([]uint64{1, 1, 1, uint64(i)}, []uint64{1})
		require.NoError(t, err)
	}

	subject.clear()

	require.Zero(t, len(subject.data))
}

func TestConcurrentAccess(t *testing.T) {
	subject := newNodeCache()
	var wg sync.WaitGroup

	const totalItems = 10
	for i := 0; i < totalItems; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			err := subject.set([]uint64{1, 1, 1, uint64(i)}, []uint64{uint64(i)})
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	for i := 0; i < totalItems; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			value, err := subject.get([]uint64{1, 1, 1, uint64(i)})
			require.NoError(t, err)

			require.Equal(t, value, []uint64{uint64(i)})
		}(i)
	}
	wg.Wait()
}
