// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	graphql "github.com/kyma-project/kyma/tests/components/application-connector/test/compass-runtime-agent/testkit/third_party/machinebox/graphql"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Do provides a mock function with given fields: req, res
func (_m *Client) Do(req *graphql.Request, res interface{}) error {
	ret := _m.Called(req, res)

	var r0 error
	if rf, ok := ret.Get(0).(func(*graphql.Request, interface{}) error); ok {
		r0 = rf(req, res)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClient(t mockConstructorTestingTNewClient) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}