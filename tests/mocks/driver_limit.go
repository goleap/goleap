package mocks

import (
	specs "github.com/kitstack/dbkit/specs"
	mock "github.com/stretchr/testify/mock"
)

// FakeDriverLimit is an mock type for the FakeDriverLimit type
type FakeDriverLimit struct {
	mock.Mock
}

// Formatted provides a mock function with given fields:
func (_m *FakeDriverLimit) Formatted() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Limit provides a mock function with given fields:
func (_m *FakeDriverLimit) Limit() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Offset provides a mock function with given fields:
func (_m *FakeDriverLimit) Offset() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// SetLimit provides a mock function with given fields: index
func (_m *FakeDriverLimit) SetLimit(index int) specs.DriverLimit {
	ret := _m.Called(index)

	var r0 specs.DriverLimit
	if rf, ok := ret.Get(0).(func(int) specs.DriverLimit); ok {
		r0 = rf(index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(specs.DriverLimit)
		}
	}

	return r0
}

// SetOffset provides a mock function with given fields: index
func (_m *FakeDriverLimit) SetOffset(index int) specs.DriverLimit {
	ret := _m.Called(index)

	var r0 specs.DriverLimit
	if rf, ok := ret.Get(0).(func(int) specs.DriverLimit); ok {
		r0 = rf(index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(specs.DriverLimit)
		}
	}

	return r0
}

type mockConstructorTestingTNewDriverLimit interface {
	mock.TestingT
	Cleanup(func())
}

// NewFakeDriverLimit creates a new instance of FakeDriverLimit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFakeDriverLimit(t mockConstructorTestingTNewDriverLimit) *FakeDriverLimit {
	fakeDriverLimit := &FakeDriverLimit{}
	fakeDriverLimit.Mock.Test(t)

	t.Cleanup(func() { fakeDriverLimit.AssertExpectations(t) })

	return fakeDriverLimit
}
