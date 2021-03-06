// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	warehouses "github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	web "github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Create provides a mock function with given fields: warehouseCode, adress, telephone, minimumCapacity, maxmumCapacity
func (_m *Service) Create(warehouseCode string, adress string, telephone string, minimumCapacity int, maxmumCapacity int) (warehouses.Warehouse, web.ResponseCode) {
	ret := _m.Called(warehouseCode, adress, telephone, minimumCapacity, maxmumCapacity)

	var r0 warehouses.Warehouse
	if rf, ok := ret.Get(0).(func(string, string, string, int, int) warehouses.Warehouse); ok {
		r0 = rf(warehouseCode, adress, telephone, minimumCapacity, maxmumCapacity)
	} else {
		r0 = ret.Get(0).(warehouses.Warehouse)
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(string, string, string, int, int) web.ResponseCode); ok {
		r1 = rf(warehouseCode, adress, telephone, minimumCapacity, maxmumCapacity)
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *Service) Delete(id int) web.ResponseCode {
	ret := _m.Called(id)

	var r0 web.ResponseCode
	if rf, ok := ret.Get(0).(func(int) web.ResponseCode); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(web.ResponseCode)
	}

	return r0
}

// GetAll provides a mock function with given fields:
func (_m *Service) GetAll() ([]warehouses.Warehouse, web.ResponseCode) {
	ret := _m.Called()

	var r0 []warehouses.Warehouse
	if rf, ok := ret.Get(0).(func() []warehouses.Warehouse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]warehouses.Warehouse)
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

// GetOne provides a mock function with given fields: id
func (_m *Service) GetOne(id int) (warehouses.Warehouse, web.ResponseCode) {
	ret := _m.Called(id)

	var r0 warehouses.Warehouse
	if rf, ok := ret.Get(0).(func(int) warehouses.Warehouse); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(warehouses.Warehouse)
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(int) web.ResponseCode); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, requestData
func (_m *Service) Update(id int, requestData map[string]interface{}) (warehouses.Warehouse, web.ResponseCode) {
	ret := _m.Called(id, requestData)

	var r0 warehouses.Warehouse
	if rf, ok := ret.Get(0).(func(int, map[string]interface{}) warehouses.Warehouse); ok {
		r0 = rf(id, requestData)
	} else {
		r0 = ret.Get(0).(warehouses.Warehouse)
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(int, map[string]interface{}) web.ResponseCode); ok {
		r1 = rf(id, requestData)
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
