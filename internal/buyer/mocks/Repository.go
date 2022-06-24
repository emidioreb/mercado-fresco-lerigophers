// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	buyers "github.com/emidioreb/mercado-fresco-lerigophers/internal/buyer"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: cardNumberId, firstName, lastName
func (_m *Repository) Create(cardNumberId string, firstName string, lastName string) (buyers.Buyer, error) {
	ret := _m.Called(cardNumberId, firstName, lastName)

	var r0 buyers.Buyer
	if rf, ok := ret.Get(0).(func(string, string, string) buyers.Buyer); ok {
		r0 = rf(cardNumberId, firstName, lastName)
	} else {
		r0 = ret.Get(0).(buyers.Buyer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(cardNumberId, firstName, lastName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *Repository) Delete(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields:
func (_m *Repository) GetAll() ([]buyers.Buyer, error) {
	ret := _m.Called()

	var r0 []buyers.Buyer
	if rf, ok := ret.Get(0).(func() []buyers.Buyer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]buyers.Buyer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOne provides a mock function with given fields: id
func (_m *Repository) GetOne(id int) (buyers.Buyer, error) {
	ret := _m.Called(id)

	var r0 buyers.Buyer
	if rf, ok := ret.Get(0).(func(int) buyers.Buyer); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(buyers.Buyer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, requestData
func (_m *Repository) Update(id int, requestData map[string]interface{}) (buyers.Buyer, error) {
	ret := _m.Called(id, requestData)

	var r0 buyers.Buyer
	if rf, ok := ret.Get(0).(func(int, map[string]interface{}) buyers.Buyer); ok {
		r0 = rf(id, requestData)
	} else {
		r0 = ret.Get(0).(buyers.Buyer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, map[string]interface{}) error); ok {
		r1 = rf(id, requestData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
