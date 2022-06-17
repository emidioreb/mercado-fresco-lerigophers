package service_test

import (
	"errors"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {
	t.Run("must return the created warehouse", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      "212",
			Address:            "rua do bobo",
			Telephone:          "0",
			MinimumCapacity:    10,
			MaximumTemperature: 30,
		}

		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := warehouses.NewService(mockedRepository)

		result, _ := service.Create(input.WarehouseCode, input.Address, input.Telephone, input.MinimumCapacity, input.MaximumTemperature)

		assert.Equal(t, result, input)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("must return an error when the warehousecode already existes", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      "212",
			Address:            "rua do bobo",
			Telephone:          "0",
			MinimumCapacity:    10,
			MaximumTemperature: 30,
		}

		expectedError := errors.New("warehouse_code already exists")

		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, expectedError)

		service := warehouses.NewService(mockedRepository)

		_, err := service.Create(input.WarehouseCode, input.Address, input.Telephone, input.MinimumCapacity, input.MaximumTemperature)

		assert.Equal(t, err.Err.Error(), expectedError.Error())
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceGetAll(t *testing.T) {
	t.Run("must return a list of warehouses", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := []warehouses.Warehouse{
			{
				Id:                 1,
				WarehouseCode:      "212",
				Address:            "rua do bobo",
				Telephone:          "0",
				MinimumCapacity:    10,
				MaximumTemperature: 30,
			},
			{
				Id:                 2,
				WarehouseCode:      "212a",
				Address:            "rua do bobo",
				Telephone:          "0",
				MinimumCapacity:    10,
				MaximumTemperature: 30,
			},
		}

		mockedRepository.On("GetAll").Return(input, nil).Once()

		service := warehouses.NewService(mockedRepository)

		result, _ := service.GetAll()

		for i, warehouse := range input {
			assert.Equal(t, result[i], warehouse)
		}
		mockedRepository.AssertExpectations(t)

	})
}

func TestServiceGetOne(t *testing.T) {
	mockedRepository := new(mocks.Repository)

	t.Run("must return an warehouse with the given id", func(t *testing.T) {
		input := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      "212",
			Address:            "rua do bobo",
			Telephone:          "0",
			MinimumCapacity:    10,
			MaximumTemperature: 30,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(input, nil).Once()

		service := warehouses.NewService(mockedRepository)

		result, _ := service.GetOne(123)

		assert.Equal(t, result, input)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("must return an error when searching for an unregistered id", func(t *testing.T) {
		expectedError := errors.New("warehouse with id 1 not found")

		service := warehouses.NewService(mockedRepository)

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, expectedError).Once()

		_, err := service.GetOne(1)

		assert.NotNil(t, err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		mockedRepository.AssertExpectations(t)
	})
}
