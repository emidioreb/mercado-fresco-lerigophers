package controllers_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponse struct {
	Data []products.Product
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func Test_Get_Product_OK(t *testing.T) {
	t.Run("OK Case - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)
		mockeProductList := make([]products.Product, 0)

		productController := controllers.NewProduct(mockedService)

		fakeProduct := products.Product{
			Id:                             1,
			ProductCode:                    "BS0003",
			Description:                    "Batata",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}

		mockeProductList = append(mockeProductList, fakeProduct)

		mockedService.On("GetAll").Return(mockeProductList, web.ResponseCode{})

		router := gin.Default()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/products/", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.GET("/api/v1/products/", productController.GetAll())
		router.ServeHTTP(rec, req)

		responseData, _ := ioutil.ReadAll(rec.Body)

		var currentResponse ObjectResponse

		err = json.Unmarshal(responseData, &currentResponse)

		assert.Nil(t, err)
		assert.Equal(t, fakeProduct, currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case - 500", func(t *testing.T) {
		mockedService := new(mocks.Service)

		productController := controllers.NewProduct(mockedService)

		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errors.New("internal server error"),
		})

		router := gin.Default()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/products/", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.Handle(http.MethodGet, "/api/v1/products/",
			productController.GetAll())
		router.ServeHTTP(rec, req)

		responseData, err := ioutil.ReadAll(rec.Body)
		assert.Nil(t, err)

		var currentResponse web.ResponseCode

		err = json.Unmarshal(responseData, &currentResponse)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test_Get_One_Product(t *testing.T) {
	t.Run("OK Case if exists - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)

		productController := controllers.NewProduct(mockedService)

		fakeProduct := products.Product{
			Id:                             1,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}

		mockedService.On("GetOne",
			mock.AnythingOfType("int")).Return(fakeProduct, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/products/:id", productController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		type objResponse struct {
			Data products.Product
		}

		var currentResponse objResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fakeProduct, currentResponse.Data)
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService := new(mocks.Service)

		productController := controllers.NewProduct(mockedService)

		expectedError := errors.New("product with id 1 not found")
		mockedService.On("GetOne",
			mock.AnythingOfType("int")).Return(products.Product{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		router := gin.Default()
		router.GET("/api/v1/products/:id", productController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})
}
