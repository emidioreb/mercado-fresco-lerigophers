package employees_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {

	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		mockedRepository.On("GetAll").Return([]employees.Employee{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := employees.NewService(mockedRepository, nil)

		result, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)
		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test error case if employee's cardnumberid already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}
		listaEmployees := []employees.Employee{}
		listaEmployees = append(listaEmployees, input)

		expectedError := errors.New("card_number_id already exists")

		mockedRepository.On("GetAll").Return(listaEmployees, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, nil)

		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
	})
}

func TestServiceGetAll(t *testing.T) {
	t.Run("Tests if getAll returns employees", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}
		listaEmployees := []employees.Employee{}
		listaEmployees = append(listaEmployees, input)

		mockedRepository.On("GetAll").Return(listaEmployees, nil)

		service := employees.NewService(mockedRepository, nil)

		result, _ := service.GetAll()

		assert.Equal(t, listaEmployees, result)
	})
}

func TestServiceGetOne(t *testing.T) {
	t.Run("Tests if getOne returns employee", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(input, nil)

		service := employees.NewService(mockedRepository, nil)

		result, err := service.GetOne(1)
		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Tests if getOne returns error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		expectedError := errors.New("employee with id 1 not found")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, nil)

		_, err := service.GetOne(1)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, expectedError, err.Err)

		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("Return the updated employee when successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"first_name": "Now",
			"last_name":  "Are",
		}

		expectedEmployee := employees.Employee{
			Id:           1,
			CardNumberId: "1",
			FirstName:    "Now",
			LastName:     "Are",
			WarehouseId:  1,
		}

		input := employees.Employee{
			Id:           1,
			CardNumberId: "1",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(input, nil)

		mockedRepository.On("GetAll").Return([]employees.Employee{}, nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(expectedEmployee, nil)

		service := employees.NewService(mockedRepository, nil)
		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)
		assert.Equal(t, expectedEmployee, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Return null when employee id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("employee with id 1 not found")
		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(employees.Employee{}, expectedError).Once()

		mockedRepository.On("GetAll").
			Return([]employees.Employee{}, nil).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(employees.Employee{}, nil).Once()

		service := employees.NewService(mockedRepository, nil)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("Return conflict when CardNumberId alread exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"card_number_id": "1",
		}

		input := []employees.Employee{{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}, {
			Id:           2,
			CardNumberId: "1",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		},
		}

		expectedError := errors.New("card_number_id already exists")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedRepository.On("GetAll").Return(input, nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, nil)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusConflict, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("Verify the successfully case if the employee is deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(nil)

		service := employees.NewService(mockedRepository, nil)
		result := service.Delete(1)
		assert.Nil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNoContent)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Verify the error case if seller do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("employee with id 1 not found")

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(expectedError)

		service := employees.NewService(mockedRepository, nil)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNotFound)
		assert.Equal(t, result.Err, expectedError)
		mockedRepository.AssertExpectations(t)
	})
}
