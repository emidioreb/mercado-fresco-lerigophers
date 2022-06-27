package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type objectResponseArrEmp struct {
	Data []employees.Employee
}

type objectResponseEmp struct {
	Data employees.Employee
}

type objectErrorResponseEmp struct {
	Error string `json:"error"`
}

func newEmployeeController() (*mocks.Service, *controllers.EmployeeController) {
	mockedService := new(mocks.Service)
	employeeController := controllers.NewEmployee(mockedService)
	return mockedService, employeeController
}

var fakeEmployee = []employees.Employee{
	{
		Id:           1,
		CardNumberId: "100",
		FirstName:    "",
		LastName:     "",
		WarehouseId:  1,
	},
	{
		Id:           1,
		CardNumberId: "320",
		FirstName:    "",
		LastName:     "",
		WarehouseId:  3,
	},
}

const (
	defaultURL = "/api/v1/employees/"
	idString   = "/api/v1/employees/string"
	idNumber1  = "/api/v1/employees/1"
	idRequest  = "/api/v1/employees/:id"
)

var (
	errServer           = errors.New("internal server error")
	errEmployeeNotFound = errors.New("employee with id 1 not found")
	errIdNotNumber      = errors.New("id must be a number")
	errInvalidRequest   = errors.New("invalid request data")
	errNeedBody         = errors.New("invalid request data - body needed")
	errCardIdNeeded     = errors.New("empty card_number_id not allowed")
	errTypeData         = errors.New("invalid type of data")
	errInvalidInput     = errors.New("invalid request input")
	errCardIdExists     = errors.New("card_number_id already exists")
)

func TestCreateSeller(t *testing.T) {
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).
			Return(fakeEmployee[0], web.ResponseCode{
				Code: http.StatusCreated,
			})

		parsedFakeEmployee, err := json.Marshal(fakeEmployee[0])
		assert.Nil(t, err)

		r := gin.Default()
		r.POST(defaultURL, employeeController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer(parsedFakeEmployee))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse objectResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeEmployee[0], bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).
			Return(employees.Employee{}, web.ResponseCode{})

		r := gin.Default()
		r.POST(defaultURL, employeeController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"card_number_id": 132}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidInput.Error(), bodyResponse.Error)
	})

	t.Run("card_number_id must be informed", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(employees.Employee{}, web.ResponseCode{})

		router := gin.Default()
		router.POST(defaultURL, employeeController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"card_number_id": ""}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCardIdNeeded.Error(), bodyResponse.Error)
	})

	t.Run("Conflict card_number_id", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("GetAll").Return(fakeEmployee, nil)
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, web.ResponseCode{
			Code: http.StatusConflict,
			Err:  errCardIdExists,
		})

		r := gin.Default()
		r.POST(defaultURL, employeeController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"card_number_id": "100"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCardIdExists.Error(), bodyResponse.Error)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Get all employees", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("GetAll").Return(fakeEmployee, web.ResponseCode{})

		r := gin.Default()
		r.GET(defaultURL, employeeController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse objectResponseArrEmp
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeEmployee[0], currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errServer,
		})

		r := gin.Default()
		r.GET(defaultURL, employeeController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse web.ResponseCode
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetOne(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeEmployee[0], web.ResponseCode{})

		r := gin.Default()
		r.GET(idRequest, employeeController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse objectResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fakeEmployee[0], currentResponse.Data)
	})

	t.Run("Not exist case", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errEmployeeNotFound,
		})

		r := gin.Default()
		r.GET(idRequest, employeeController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse objectErrorResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errEmployeeNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(employees.Employee{}, web.ResponseCode{})

		r := gin.Default()
		r.GET(idRequest, employeeController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idString, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse objectErrorResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, errIdNotNumber.Error(), currentResponse.Error)
	})
}

func TestDeleteSeller(t *testing.T) {
	t.Run("Success case if exists", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		r := gin.Default()
		r.DELETE(idRequest, employeeController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.True(t, rec.Body.String() == "")
	})

	t.Run("Error case if not exists", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errEmployeeNotFound,
		})

		r := gin.Default()
		r.DELETE(idRequest, employeeController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse objectErrorResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errEmployeeNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(employees.Employee{}, web.ResponseCode{})

		r := gin.Default()
		r.DELETE(idRequest, employeeController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idString, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse objectErrorResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, errIdNotNumber.Error(), currentResponse.Error)
	})
}

func TestUpdateSeller(t *testing.T) {
	t.Run("Sucessfully case", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeEmployee[0], web.ResponseCode{})

		parsedFakeEmployee, err := json.Marshal(fakeEmployee[0])
		assert.Nil(t, err)

		r := gin.Default()
		r.PATCH(idRequest, employeeController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeEmployee))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var bodyResponse objectResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeEmployee[0], bodyResponse.Data)
	})

	t.Run("Not found case", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(employees.Employee{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  errEmployeeNotFound,
			})

		r := gin.Default()
		r.PATCH(idRequest, employeeController.Update())

		parsedFakeEmployee, err := json.Marshal(fakeEmployee[0])
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeEmployee))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errEmployeeNotFound.Error(), bodyResponse.Error)
	})

	t.Run("Id must be a number", func(t *testing.T) {
		mockedService, sellerController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(employees.Employee{}, web.ResponseCode{})

		parsedFakeEmployee, err := json.Marshal(fakeEmployee[0])
		assert.Nil(t, err)

		r := gin.Default()
		r.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idString, bytes.NewBuffer(parsedFakeEmployee))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errIdNotNumber.Error(), bodyResponse.Error)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(employees.Employee{}, web.ResponseCode{})

		r := gin.Default()
		r.PATCH(idRequest, employeeController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte{}))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidRequest.Error(), bodyResponse.Error)
	})

	t.Run("Body needed", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(employees.Employee{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH(idRequest, employeeController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte("{}")))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errNeedBody.Error(), bodyResponse.Error)
	})

	t.Run("card_number_id not be empty", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(employees.Employee{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH(idRequest, employeeController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"card_number_id": "" }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCardIdNeeded.Error(), bodyResponse.Error)
	})

	t.Run("Syntax error on body", func(t *testing.T) {
		mockedService, employeeController := newEmployeeController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(employees.Employee{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH(idRequest, employeeController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"first_name": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse objectErrorResponseEmp
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errTypeData.Error(), bodyResponse.Error)
	})
}
