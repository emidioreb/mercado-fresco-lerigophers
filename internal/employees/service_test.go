package employees_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	mockWarehouseRepository "github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {

	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouse := new(mockWarehouseRepository.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
		mockedWarehouse.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, nil).Once()
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := employees.NewService(mockedRepository, mockedWarehouse)

		result, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)
		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test error case if employee's cardnumberid already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouse := new(mockWarehouseRepository.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		errEmployeeCarId := errors.New("card_number_id already exists")

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(errEmployeeCarId)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := employees.NewService(mockedRepository, mockedWarehouse)

		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), errEmployeeCarId.Error())
	})

	t.Run("Test error case if employee's verification cardnumberid already exists fail", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouse := new(mockWarehouseRepository.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		errEmployeeCarId := errors.New("any error")

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(errEmployeeCarId)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := employees.NewService(mockedRepository, mockedWarehouse)

		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), errEmployeeCarId.Error())
	})

	t.Run("Test error case if warehouse's warehouse_id exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouse := new(mockWarehouseRepository.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		errWarehouseId := errors.New("any error")

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
		mockedWarehouse.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, errWarehouseId).Once()
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := employees.NewService(mockedRepository, mockedWarehouse)

		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), errWarehouseId.Error())
	})

	t.Run("Test error case if warehouse's verification warehouse_id exists fail", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouse := new(mockWarehouseRepository.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		errWarehouseId := errors.New("warehouse with id 1 not found")

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
		mockedWarehouse.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, errWarehouseId).Once()
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := employees.NewService(mockedRepository, mockedWarehouse)

		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), errWarehouseId.Error())
	})

	t.Run("Test if create fails", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouse := new(mockWarehouseRepository.Repository)

		input := employees.Employee{
			Id:           1,
			CardNumberId: "123",
			FirstName:    "",
			LastName:     "",
			WarehouseId:  1,
		}

		errCreate := errors.New("ocurred an error to create employee")

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
		mockedWarehouse.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, nil).Once()
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, errCreate)

		service := employees.NewService(mockedRepository, mockedWarehouse)

		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName, input.WarehouseId)
		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), errCreate.Error())
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

	t.Run("Tests if getAll returns error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		expectedErr := errors.New("any error")

		mockedRepository.On("GetAll").Return([]employees.Employee{}, expectedErr)

		service := employees.NewService(mockedRepository, nil)

		_, err := service.GetAll()

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr.Error(), err.Err.Error())
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

	t.Run("Tests if getOne returns not found", func(t *testing.T) {
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

	t.Run("Tests if getOne returns error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		expectedError := errors.New("unexpected error to get employee")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, nil)

		_, err := service.GetOne(1)

		assert.NotNil(t, err.Err)
		assert.Equal(t, expectedError.Error(), err.Err.Error())

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

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(expectedEmployee, nil)

		service := employees.NewService(mockedRepository, nil)
		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)
		assert.Equal(t, expectedEmployee, result)
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

		expectedError := errors.New("card_number_id already exists")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(expectedError)

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

	t.Run("Return error when CardNumberId's validation fail", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"card_number_id": "1",
		}

		expectedError := errors.New("unexpected error to get employee")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(expectedError)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, nil)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("Return the updated employee when fail", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"first_name": "Now",
			"last_name":  "Are",
		}

		errUpdate := errors.New("ocurred an error while updating the employee")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedRepository.On("GetOneByCardNumber", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(employees.Employee{}, errUpdate)

		service := employees.NewService(mockedRepository, nil)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, errUpdate, err.Err)
	})

	t.Run("Return error when warehouse_id's validation fail", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouseRepo := new(mockWarehouseRepository.Repository)

		requestData := map[string]interface{}{
			"warehouse_id": 1.0,
		}

		expectedError := errors.New("unexpected error to get warehouse")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedWarehouseRepo.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, expectedError)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, mockedWarehouseRepo)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("Return error when warehouse_id's not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedWarehouseRepo := new(mockWarehouseRepository.Repository)

		requestData := map[string]interface{}{
			"warehouse_id": 1.0,
		}

		expectedError := errors.New("warehouse with id 1 not found")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, nil)

		mockedWarehouseRepo.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, expectedError)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(employees.Employee{}, expectedError)

		service := employees.NewService(mockedRepository, mockedWarehouseRepo)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
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

	t.Run("Verify the error case delete return error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("unexpected error to delete employee")

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(expectedError)

		service := employees.NewService(mockedRepository, nil)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusInternalServerError)
		assert.Equal(t, result.Err.Error(), expectedError.Error())
		mockedRepository.AssertExpectations(t)
	})
}
