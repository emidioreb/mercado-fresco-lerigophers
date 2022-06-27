package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponseArr struct {
	Data []products.Product
}

type ObjectResponse struct {
	Data products.Product
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func routerProducts() *gin.Engine {
	router := gin.Default()
	return router
}

func newProductController() (*mocks.Service, *controllers.ProductController) {
	mockedService := new(mocks.Service)
	productController := controllers.NewProduct(mockedService)
	return mockedService, productController
}

var fakeProducts = []products.Product{
	{
		Id:                             1,
		ProductCode:                    "BS0001",
		Description:                    "Batata",
		Width:                          23,
		Height:                         62,
		Length:                         101,
		NetWeight:                      27,
		ExpirationRate:                 88,
		RecommendedFreezingTemperature: 17,
		FreezingRate:                   23,
		ProductTypeId:                  7,
	},
	{
		Id:                             2,
		ProductCode:                    "BS0002",
		Description:                    "Batata",
		Width:                          23,
		Height:                         62,
		Length:                         101,
		NetWeight:                      27,
		ExpirationRate:                 88,
		RecommendedFreezingTemperature: 17,
		FreezingRate:                   23,
		ProductTypeId:                  7,
	},
}

const (
	defaultURL = "/api/v1/products/"
	idString   = "/api/v1/products/string"
	idNumber1  = "/api/v1/products/1"
	idRequest  = "api/v1/products/:id"
)

var (
	errServer                 = errors.New("internal server error")
	errProductNotFound        = errors.New("product with id 1 not found")
	errIdNotNumber            = errors.New("id must be a number")
	errInvalidRequest         = errors.New("invalid request data")
	errNeedBody               = errors.New("invalid request data - body needed")
	errProductWithBlankSpaces = errors.New("empty product_code not allowed")
	errProductCodeNeeded      = errors.New("product code must be informed and greather than 0")
	errTypeData               = errors.New("invalid type of data")
	errInvalidInput           = errors.New("invalid request input")
	errProductCodeExists      = errors.New("product code already exists")
)

func TestGetProduct(t *testing.T) {
	t.Run("Get all products", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("GetAll").Return(fakeProducts, web.ResponseCode{})

		r := routerProducts()
		r.GET(defaultURL, productController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponseArr
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeProducts[0], currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errServer,
		})

		r := routerProducts()
		r.GET(defaultURL, productController.GetAll())

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
		mockedService, productController := newProductController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeProducts[0], web.ResponseCode{})

		r := routerProducts()
		r.GET(idRequest, productController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fakeProducts[0], currentResponse.Data)
	})

	t.Run("Not exist case", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(products.Product{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errProductNotFound,
		})

		r := routerProducts()
		r.GET(idRequest, productController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errProductNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.GET(idRequest, productController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idString, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, errIdNotNumber.Error(), currentResponse.Error)
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("Success case if exists", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		r := routerProducts()
		r.DELETE(idRequest, productController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.True(t, "" == string(rec.Body.String()))
	})

	t.Run("Error case if not exists", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errProductNotFound,
		})

		r := routerProducts()
		r.DELETE(idRequest, productController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errProductNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.DELETE(idRequest, productController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idString, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, errIdNotNumber.Error(), currentResponse.Error)
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Sucessfully case", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeProducts[0], web.ResponseCode{})

		parsedFakeProduct, err := json.Marshal(fakeProducts[0])
		assert.Nil(t, err)

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeProduct))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeProducts[0], bodyResponse.Data)

	})

	t.Run("Not found case", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(products.Product{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  errProductNotFound,
			})

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		parsedFakeProduct, err := json.Marshal(fakeProducts[0])
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeProduct))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errProductNotFound.Error(), bodyResponse.Error)
	})

	t.Run("Id must be a number", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(products.Product{}, web.ResponseCode{})

		parsedFakeProduct, err := json.Marshal(fakeProducts[0])
		assert.Nil(t, err)

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		req, err := http.NewRequest(http.MethodPatch, idString, bytes.NewBuffer(parsedFakeProduct))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errIdNotNumber.Error(), bodyResponse.Error)

	})

	t.Run("Invalid request data", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte{}))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidRequest.Error(), bodyResponse.Error)
	})

	t.Run("Body needed", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte("{}")))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errNeedBody.Error(), bodyResponse.Error)
	})

	t.Run("Product Code just with blank spaces", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"product_code": "   " }`)))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errProductWithBlankSpaces.Error(), bodyResponse.Error)
	})

	t.Run("Sysntax error on body", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.PATCH(idRequest, productController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"product_code": 0}`)))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errTypeData.Error(), bodyResponse.Error)
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("Success on Create", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).
			Return(fakeProducts[0], web.ResponseCode{
				Code: http.StatusCreated,
			})

		parsedFakeProduct, err := json.Marshal(fakeProducts[0])
		assert.Nil(t, err)

		r := gin.Default()
		r.POST(defaultURL, productController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer(parsedFakeProduct))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeProducts[0], bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).
			Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.POST(defaultURL, productController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"product_code":0}`)))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidInput.Error(), bodyResponse.Error)
	})

	t.Run("Product Code with only blank spaces", func(t *testing.T) {
		mockedService, productController := newProductController()
		mockedService.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(products.Product{}, web.ResponseCode{})

		r := routerProducts()
		r.POST(defaultURL, productController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"product_code":"   "}`)))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errProductWithBlankSpaces.Error(), bodyResponse.Error)
	})

	t.Run("Conflict Product Code", func(t *testing.T) {
		mockedService, sellerController := newProductController()
		mockedService.On("GetAll").Return(fakeProducts, nil)
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(products.Product{}, web.ResponseCode{
			Code: http.StatusConflict,
			Err:  errProductCodeExists,
		})

		r := routerProducts()
		r.POST(defaultURL, sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"product_code": "BS0001"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errProductCodeExists.Error(), bodyResponse.Error)
	})
}
