// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	carriers "github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers"
	mock "github.com/stretchr/testify/mock"

	web "github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Create provides a mock function with given fields: cid, companyName, address, telephone, localityId
func (_m *Service) Create(cid string, companyName string, address string, telephone string, localityId string) (carriers.Carry, web.ResponseCode) {
	ret := _m.Called(cid, companyName, address, telephone, localityId)

	var r0 carriers.Carry
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) carriers.Carry); ok {
		r0 = rf(cid, companyName, address, telephone, localityId)
	} else {
		r0 = ret.Get(0).(carriers.Carry)
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(string, string, string, string, string) web.ResponseCode); ok {
		r1 = rf(cid, companyName, address, telephone, localityId)
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
