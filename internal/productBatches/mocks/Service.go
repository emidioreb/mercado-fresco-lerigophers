// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	time "time"

	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	mock "github.com/stretchr/testify/mock"

	web "github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateProductBatch provides a mock function with given fields: BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId, DueDate, ManufacturingDate
func (_m *Service) CreateProductBatch(BatchNumber int, CurrentQuantity int, CurrentTemperature int, InitialQuantity int, ManufacturingHour int, MinimumTemperature int, ProductId int, SectionId int, DueDate time.Time, ManufacturingDate time.Time) (product_batches.ProductBatches, web.ResponseCode) {
	ret := _m.Called(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId, DueDate, ManufacturingDate)

	var r0 product_batches.ProductBatches
	if rf, ok := ret.Get(0).(func(int, int, int, int, int, int, int, int, time.Time, time.Time) product_batches.ProductBatches); ok {
		r0 = rf(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId, DueDate, ManufacturingDate)
	} else {
		r0 = ret.Get(0).(product_batches.ProductBatches)
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(int, int, int, int, int, int, int, int, time.Time, time.Time) web.ResponseCode); ok {
		r1 = rf(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId, DueDate, ManufacturingDate)
	} else {
		r1 = ret.Get(1).(web.ResponseCode)
	}

	return r0, r1
}

// GetReportSection provides a mock function with given fields: SectionId
func (_m *Service) GetReportSection(SectionId int) ([]product_batches.ProductsQuantity, web.ResponseCode) {
	ret := _m.Called(SectionId)

	var r0 []product_batches.ProductsQuantity
	if rf, ok := ret.Get(0).(func(int) []product_batches.ProductsQuantity); ok {
		r0 = rf(SectionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]product_batches.ProductsQuantity)
		}
	}

	var r1 web.ResponseCode
	if rf, ok := ret.Get(1).(func(int) web.ResponseCode); ok {
		r1 = rf(SectionId)
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
