package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponse struct {
	Data warehouses.Warehouse
}

type ObjectResponseArr struct {
	Data []warehouses.Warehouse
}

var fakeWarehouse = []warehouses.Warehouse{{
	Id:                 1,
	WarehouseCode:      "1",
	Address:            "Rua do bobo",
	Telephone:          "11111111111",
	MinimumCapacity:    0,
	MaximumTemperature: 20,
}, {
	Id:                 2,
	WarehouseCode:      "2",
	Address:            "Terra do nunca",
	Telephone:          "Fake Address",
	MinimumCapacity:    0,
	MaximumTemperature: 20,
}}

func TestControllerWarehouseCreate(t *testing.T) {
	t.Run("return error when warehouse_code already exists", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		input := map[string]interface{}{
			"warehouse_code":     "a1ß",
			"address":            "rua do bobo",
			"telephone":          "0",
			"minimumCapacity":    0,
			"maximumTemperature": 30,
		}

		parsedInput, err := json.Marshal(input)
		assert.Nil(t, err)

		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, web.NewCodeResponse(http.StatusConflict, errors.New("warehouse_code already exists")))

		router := gin.Default()
		router.POST("/api/v1/warehouses", warehouseController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusConflict, rec.Code)

		var currentResponse ObjectResponse

		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, warehouses.Warehouse{}, currentResponse.Data)
	})

	t.Run("return error when warehouse_code is empty", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		input := map[string]interface{}{
			"warehouse_code":     "",
			"address":            "rua do bobo",
			"telephone":          "0",
			"minimumCapacity":    0,
			"maximumTemperature": 30,
		}

		parsedInput, err := json.Marshal(input)
		assert.Nil(t, err)

		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, web.NewCodeResponse(0, nil))

		router := gin.Default()
		router.POST("/api/v1/warehouses", warehouseController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("return error when gives invalid input", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		input := map[string]interface{}{
			"warehouse_code":     1,
			"address":            1.02,
			"telephone":          "0",
			"minimumCapacity":    0,
			"maximumTemperature": 30,
		}

		parsedInput, err := json.Marshal(input)
		assert.Nil(t, err)

		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, web.NewCodeResponse(0, nil))

		router := gin.Default()
		router.POST("/api/v1/warehouses", warehouseController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("create, success case", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		input := controllers.ReqWarehouses{
			WarehouseCode:      "1",
			Address:            "rua do bobo",
			Telephone:          "0",
			MinimumCapacity:    0,
			MaximumTemperature: 30,
		}

		parsedInput, err := json.Marshal(input)
		assert.Nil(t, err)

		expectedReturnData := warehouses.Warehouse{
			Id:                 1,
			WarehouseCode:      input.WarehouseCode,
			Address:            input.Address,
			Telephone:          input.Telephone,
			MinimumCapacity:    input.MinimumCapacity,
			MaximumTemperature: input.MaximumTemperature,
		}

		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(expectedReturnData, web.NewCodeResponse(http.StatusCreated, nil))

		router := gin.Default()
		router.POST("/api/v1/warehouses", warehouseController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/warehouses", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var currentResponse ObjectResponse

		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, expectedReturnData, currentResponse.Data)
	})
}

func TestControllerWarehouseGetAll(t *testing.T) {
	t.Run("success getAll", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("GetAll").Return([]warehouses.Warehouse{}, web.NewCodeResponse(http.StatusOK, nil))

		router := gin.Default()
		router.GET("/api/v1/warehouses", warehouseController.GetAll())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/warehouses", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var currentResponse ObjectResponseArr

		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, []warehouses.Warehouse{}, currentResponse.Data)
	})

	t.Run("error getAll", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("GetAll").Return([]warehouses.Warehouse{}, web.NewCodeResponse(http.StatusInternalServerError, errors.New("")))

		router := gin.Default()
		router.GET("/api/v1/warehouses", warehouseController.GetAll())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/warehouses", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var currentResponse web.ResponseCode

		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestControllerWarehouseGetOne(t *testing.T) {
	t.Run("success getOne", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/warehouses/:id", warehouseController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/warehouses/1", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("return error when id not be a number", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/warehouses/:id", warehouseController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/warehouses/a", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when wharehouse not found", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(warehouses.Warehouse{}, web.ResponseCode{Code: http.StatusNotFound, Err: errors.New("warehouse with id 2 not found")})

		router := gin.Default()
		router.GET("/api/v1/warehouses/:id", warehouseController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/warehouses/2", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestControllerWarehouseDelete(t *testing.T) {
	t.Run("success delete", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.NewCodeResponse(http.StatusNoContent, nil))

		router := gin.Default()
		router.DELETE("/api/v1/warehouses/:id", warehouseController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/warehouses/2", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("error delete", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.NewCodeResponse(http.StatusNoContent, nil))

		router := gin.Default()
		router.DELETE("/api/v1/warehouses/:id", warehouseController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/warehouses/a", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when warehouse not found", func(t *testing.T) {
		mockedService := new(mocks.Service)
		warehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.NewCodeResponse(http.StatusNotFound, errors.New("warehouse not found")))

		router := gin.Default()
		router.DELETE("/api/v1/warehouses/:id", warehouseController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/warehouses/2", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestControllerWarehouseUpdate(t *testing.T) {
	t.Run("success update", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).Return(fakeWarehouse[0], web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		parsedInput, _ := json.Marshal(fakeWarehouse[0])

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/1", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeWarehouse[0], currentResponse.Data)
	})

	t.Run("return error when gives an invalid id", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		parsedInput, _ := json.Marshal(fakeWarehouse[0])

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/a", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when send a bad request", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/1", bytes.NewBuffer([]byte("")))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when send a bad request", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		input := map[string]interface{}{}
		parsedInput, _ := json.Marshal(input)

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/1", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when send a bad request", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		input := map[string]interface{}{"warehouse_code": true}
		parsedInput, _ := json.Marshal(input)

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/1", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when send a bad request", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		input := map[string]interface{}{"warehouse_code": " "}
		parsedInput, _ := json.Marshal(input)

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/1", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("return error when send a bad request", func(t *testing.T) {
		mockedService := new(mocks.Service)
		WarehouseController := controllers.NewWarehouse(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).Return(warehouses.Warehouse{}, web.ResponseCode{Code: http.StatusNotFound, Err: errors.New("warehouse not found")})

		router := gin.Default()
		router.PATCH("/api/v1/warehouses/:id", WarehouseController.Update())

		input := map[string]interface{}{"warehouse_code": "1"}
		parsedInput, _ := json.Marshal(input)

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/warehouses/1", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
