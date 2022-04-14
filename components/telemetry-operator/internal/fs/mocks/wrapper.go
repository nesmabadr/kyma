// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	"github.com/kyma-project/kyma/components/telemetry-operator/internal/fs"
	"github.com/stretchr/testify/mock"
)

// Wrapper is an autogenerated mock type for the Wrapper type
type Wrapper struct {
	mock.Mock
}

// CreateAndWrite provides a mock function with given fields: s
func (_m *Wrapper) CreateAndWrite(s fs.File) error {
	ret := _m.Called(s)

	var r0 error
	if rf, ok := ret.Get(0).(func(fs.File) error); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveDirectory provides a mock function with given fields: path
func (_m *Wrapper) RemoveDirectory(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
