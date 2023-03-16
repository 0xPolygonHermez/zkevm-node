// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	context "context"

	pb "github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	mock "github.com/stretchr/testify/mock"
)

// ProverMock is an autogenerated mock type for the proverInterface type
type ProverMock struct {
	mock.Mock
}

// Addr provides a mock function with given fields:
func (_m *ProverMock) Addr() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AggregatedProof provides a mock function with given fields: inputProof1, inputProof2
func (_m *ProverMock) AggregatedProof(inputProof1 string, inputProof2 string) (*string, error) {
	ret := _m.Called(inputProof1, inputProof2)

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*string, error)); ok {
		return rf(inputProof1, inputProof2)
	}
	if rf, ok := ret.Get(0).(func(string, string) *string); ok {
		r0 = rf(inputProof1, inputProof2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(inputProof1, inputProof2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchProof provides a mock function with given fields: input
func (_m *ProverMock) BatchProof(input *pb.InputProver) (*string, error) {
	ret := _m.Called(input)

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(*pb.InputProver) (*string, error)); ok {
		return rf(input)
	}
	if rf, ok := ret.Get(0).(func(*pb.InputProver) *string); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(*pb.InputProver) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FinalProof provides a mock function with given fields: inputProof, aggregatorAddr
func (_m *ProverMock) FinalProof(inputProof string, aggregatorAddr string) (*string, error) {
	ret := _m.Called(inputProof, aggregatorAddr)

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*string, error)); ok {
		return rf(inputProof, aggregatorAddr)
	}
	if rf, ok := ret.Get(0).(func(string, string) *string); ok {
		r0 = rf(inputProof, aggregatorAddr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(inputProof, aggregatorAddr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ID provides a mock function with given fields:
func (_m *ProverMock) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsIdle provides a mock function with given fields:
func (_m *ProverMock) IsIdle() (bool, error) {
	ret := _m.Called()

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields:
func (_m *ProverMock) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// WaitFinalProof provides a mock function with given fields: ctx, proofID
func (_m *ProverMock) WaitFinalProof(ctx context.Context, proofID string) (*pb.FinalProof, error) {
	ret := _m.Called(ctx, proofID)

	var r0 *pb.FinalProof
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*pb.FinalProof, error)); ok {
		return rf(ctx, proofID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *pb.FinalProof); ok {
		r0 = rf(ctx, proofID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.FinalProof)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, proofID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WaitRecursiveProof provides a mock function with given fields: ctx, proofID
func (_m *ProverMock) WaitRecursiveProof(ctx context.Context, proofID string) (string, error) {
	ret := _m.Called(ctx, proofID)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, proofID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, proofID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, proofID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProverMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewProverMock creates a new instance of ProverMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProverMock(t mockConstructorTestingTNewProverMock) *ProverMock {
	mock := &ProverMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
