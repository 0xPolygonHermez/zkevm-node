package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := NewQueue[int]()

	q.Push(10)
	q.Push(20)
	q.Push(30)

	top, err := q.Top()
	require.NoError(t, err)
	assert.Equal(t, 10, top)
	assert.Equal(t, 3, q.Len())
	assert.Equal(t, false, q.IsEmpty())

	pop, err := q.Pop()
	require.NoError(t, err)
	assert.Equal(t, 10, pop)
	assert.Equal(t, 2, q.Len())
	assert.Equal(t, false, q.IsEmpty())

	top, err = q.Top()
	require.NoError(t, err)
	assert.Equal(t, 20, top)
	assert.Equal(t, 2, q.Len())
	assert.Equal(t, false, q.IsEmpty())

	pop, err = q.Pop()
	require.NoError(t, err)
	assert.Equal(t, 20, pop)
	assert.Equal(t, 1, q.Len())
	assert.Equal(t, false, q.IsEmpty())

	pop, err = q.Pop()
	require.NoError(t, err)
	assert.Equal(t, 30, pop)
	assert.Equal(t, 0, q.Len())
	assert.Equal(t, true, q.IsEmpty())

	_, err = q.Top()
	require.Error(t, ErrQueueEmpty, err)

	_, err = q.Pop()
	require.Error(t, ErrQueueEmpty, err)
}
