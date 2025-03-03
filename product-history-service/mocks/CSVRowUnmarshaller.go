// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// CSVRowUnmarshaller is an autogenerated mock type for the CSVRowUnmarshaller type
type CSVRowUnmarshaller[T any] struct {
	mock.Mock
}

type CSVRowUnmarshaller_Expecter[T any] struct {
	mock *mock.Mock
}

func (_m *CSVRowUnmarshaller[T]) EXPECT() *CSVRowUnmarshaller_Expecter[T] {
	return &CSVRowUnmarshaller_Expecter[T]{mock: &_m.Mock}
}

// ReadUnmarshalCSVRow provides a mock function with given fields:
func (_m *CSVRowUnmarshaller[T]) ReadUnmarshalCSVRow() (any, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadUnmarshalCSVRow")
	}

	var r0 any
	var r1 error
	if rf, ok := ret.Get(0).(func() (any, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() any); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(any)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadUnmarshalCSVRow'
type CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T any] struct {
	*mock.Call
}

// ReadUnmarshalCSVRow is a helper method to define mock.On call
func (_e *CSVRowUnmarshaller_Expecter[T]) ReadUnmarshalCSVRow() *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T] {
	return &CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T]{Call: _e.mock.On("ReadUnmarshalCSVRow")}
}

func (_c *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T]) Run(run func()) *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T]) Return(_a0 any, _a1 error) *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T]) RunAndReturn(run func() (any, error)) *CSVRowUnmarshaller_ReadUnmarshalCSVRow_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewCSVRowUnmarshaller creates a new instance of CSVRowUnmarshaller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCSVRowUnmarshaller[T any](t interface {
	mock.TestingT
	Cleanup(func())
}) *CSVRowUnmarshaller[T] {
	mock := &CSVRowUnmarshaller[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
