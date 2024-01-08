// Code generated by mockery. DO NOT EDIT.

package mock_syncinterfaces

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Halter is an autogenerated mock type for the Halter type
type Halter struct {
	mock.Mock
}

type Halter_Expecter struct {
	mock *mock.Mock
}

func (_m *Halter) EXPECT() *Halter_Expecter {
	return &Halter_Expecter{mock: &_m.Mock}
}

// Halt provides a mock function with given fields: ctx, err
func (_m *Halter) Halt(ctx context.Context, err error) {
	_m.Called(ctx, err)
}

// Halter_Halt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Halt'
type Halter_Halt_Call struct {
	*mock.Call
}

// Halt is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *Halter_Expecter) Halt(ctx interface{}, err interface{}) *Halter_Halt_Call {
	return &Halter_Halt_Call{Call: _e.mock.On("Halt", ctx, err)}
}

func (_c *Halter_Halt_Call) Run(run func(ctx context.Context, err error)) *Halter_Halt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *Halter_Halt_Call) Return() *Halter_Halt_Call {
	_c.Call.Return()
	return _c
}

func (_c *Halter_Halt_Call) RunAndReturn(run func(context.Context, error)) *Halter_Halt_Call {
	_c.Call.Return(run)
	return _c
}

// NewHalter creates a new instance of Halter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHalter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Halter {
	mock := &Halter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}