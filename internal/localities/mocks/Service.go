// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	localities "github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
	mock "github.com/stretchr/testify/mock"

	web "github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateLocality provides a mock function with given fields: id, localityName, provinceName, countryName
func (_m *Service) CreateLocality(id string, localityName string, provinceName string, countryName string) (localities.Locality, web.ResponseCode) {
	ret := _m.Called(id, localityName, provinceName, countryName)

	var r0 localities.Locality
	if rf, ok := ret.Get(0).(func(string, string, string, string) localities.Locality); ok {
		r0 = rf(id, localityName, provinceName, countryName)
	} else {
		r0 = ret.Get(0).(localities.Locality)
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(string, string, string, string) web.ResponseCode); ok {
		r1 = rf(id, localityName, provinceName, countryName)
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

// GetAllReportSellers provides a mock function with given fields:
func (_m *Service) GetAllReportSellers() ([]localities.ReportSellers, web.ResponseCode) {
	ret := _m.Called()

	var r0 []localities.ReportSellers
	if rf, ok := ret.Get(0).(func() []localities.ReportSellers); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]localities.ReportSellers)
		}
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func() web.ResponseCode); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

// GetReportCarriers provides a mock function with given fields: localityId
func (_m *Service) GetReportCarriers(localityId string) ([]localities.ReportCarriers, web.ResponseCode) {
	ret := _m.Called(localityId)

	var r0 []localities.ReportCarriers
	if rf, ok := ret.Get(0).(func(string) []localities.ReportCarriers); ok {
		r0 = rf(localityId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]localities.ReportCarriers)
		}
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(string) web.ResponseCode); ok {
		r1 = rf(localityId)
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

// GetReportOneSeller provides a mock function with given fields: localityId
func (_m *Service) GetReportOneSeller(localityId string) ([]localities.ReportSellers, web.ResponseCode) {
	ret := _m.Called(localityId)

	var r0 []localities.ReportSellers
	if rf, ok := ret.Get(0).(func(string) []localities.ReportSellers); ok {
		r0 = rf(localityId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]localities.ReportSellers)
		}
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(string) web.ResponseCode); ok {
		r1 = rf(localityId)
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
