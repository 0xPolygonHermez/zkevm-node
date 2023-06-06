package state

import (
	"errors"
	"sync"
)

// ErrStackEmpty returned when Pop is called and the stack is empty
var ErrStackEmpty = errors.New("Empty Stack")

// Stack is a thread safe stack data structure implementation implementing generics
type Stack[T any] struct {
	lock  sync.Mutex
	items []T
}

// NewStack creates a new stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{sync.Mutex{}, make([]T, 0)}
}

// Push adds an item to the stack
func (s *Stack[T]) Push(v T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, v)
}

// Pop removes and returns the last item added to the stack
func (s *Stack[T]) Pop() (T, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	size := len(s.items)
	if size == 0 {
		var r T
		return r, ErrStackEmpty
	}

	res := s.items[size-1]
	s.items = s.items[:size-1]
	return res, nil
}
