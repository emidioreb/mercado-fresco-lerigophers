package warehouses_test

import (
	"errors"
	"net/http"
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
			MinimumTemperature: 30,
		}

		mockedRepository.On("GetAll").Return([]warehouses.Warehouse{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := warehouses.NewService(mockedRepository)

		result, _ := service.Create(input.WarehouseCode, input.Address, input.Telephone, input.MinimumCapacity, input.MinimumTemperature)

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
			MinimumTemperature: 30,
		}

		expectedError := errors.New("warehouse_code already exists")

		mockedRepository.On("GetAll").Return([]warehouses.Warehouse{input}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, expectedError)

		service := warehouses.NewService(mockedRepository)

		_, err := service.Create(input.WarehouseCode, input.Address, input.Telephone, input.MinimumCapacity, input.MinimumTemperature)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, http.StatusConflict, err.Code)
	})

	t.Run("must return an error from the repository", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      "212",
			Address:            "rua do bobo",
			Telephone:          "0",
			MinimumCapacity:    10,
			MinimumTemperature: 30,
		}

		expectedError := errors.New("ocurred an error to create warehouse")

		mockedRepository.On("GetAll").Return([]warehouses.Warehouse{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, expectedError)

		service := warehouses.NewService(mockedRepository)

		_, err := service.Create(input.WarehouseCode, input.Address, input.Telephone, input.MinimumCapacity, input.MinimumTemperature)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, http.StatusConflict, err.Code)
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
				MinimumTemperature: 30,
			},
			{
				Id:                 2,
				WarehouseCode:      "212a",
				Address:            "rua do bobo",
				Telephone:          "0",
				MinimumCapacity:    10,
				MinimumTemperature: 30,
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

	t.Run("must return an error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		expectedError := errors.New("couldn't get warehouses")

		mockedRepository.On("GetAll").Return([]warehouses.Warehouse{}, expectedError).Once()

		service := warehouses.NewService(mockedRepository)

		_, err := service.GetAll()

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, http.StatusInternalServerError, err.Code)
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
			MinimumTemperature: 30,
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

func TestServiceDelete(t *testing.T) {
	t.Run("Verify the successfully case if the warehouse is deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(nil)

		service := warehouses.NewService(mockedRepository)
		result := service.Delete(1)
		assert.Nil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNoContent)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Verify the error case if warehouse do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("warehouse with id 1 not found")

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(expectedError)

		service := warehouses.NewService(mockedRepository)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNotFound)
		assert.Equal(t, result.Err, expectedError)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("return the updated warehouse when successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"warehouse_code":      "212",
			"address":             "rua do bobo",
			"telephone":           "0",
			"minimum_capacity":    10,
			"minimum_temperature": 30,
		}

		expectedWarehouse := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      "212",
			Address:            "rua do bobo",
			Telephone:          "0",
			MinimumCapacity:    10,
			MinimumTemperature: 30,
		}

		input := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      "1",
			Address:            "melicidade",
			Telephone:          "1111111111",
			MinimumCapacity:    150,
			MinimumTemperature: 1000,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input, nil)

		mockedRepository.On("GetAll").
			Return([]warehouses.Warehouse{}, nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(expectedWarehouse, nil)

		service := warehouses.NewService(mockedRepository)
		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)
		assert.Equal(t, expectedWarehouse, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("return error when warehouse_code already exists and id doesn't match", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"warehouse_code": "1",
		}

		expectedError := errors.New("warehouse_code already exists")

		input := []warehouses.Warehouse{
			{
				Id:                 1,
				WarehouseCode:      "1",
				Address:            "melicidade",
				Telephone:          "1111111111",
				MinimumCapacity:    150,
				MinimumTemperature: 1000,
			},
			{
				Id:                 2,
				WarehouseCode:      "2",
				Address:            "melicidade",
				Telephone:          "1111111111",
				MinimumCapacity:    150,
				MinimumTemperature: 1000,
			},
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input[1], nil).Once()

		mockedRepository.On("GetAll").
			Return(input, nil).Once()

		mockedRepository.On("Update", mock.AnythingOfType("int"), mock.Anything).Return(warehouses.Warehouse{}, nil).Once()

		service := warehouses.NewService(mockedRepository)
		_, err := service.Update(2, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, expectedError.Error(), err.Err.Error())
		assert.Equal(t, http.StatusConflict, err.Code)
	})

	t.Run("return null when warehouse id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("warehouse not found")
		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(warehouses.Warehouse{}, expectedError).Once()

		mockedRepository.On("GetAll").
			Return([]warehouses.Warehouse{}, nil).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(warehouses.Warehouse{}, nil).Once()

		service := warehouses.NewService(mockedRepository)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("return error from repository update", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("ocurred an error while updating the warehouse")
		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(warehouses.Warehouse{}, nil).Once()

		mockedRepository.On("GetAll").
			Return([]warehouses.Warehouse{}, nil).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(warehouses.Warehouse{}, expectedError).Once()

		service := warehouses.NewService(mockedRepository)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})
}
