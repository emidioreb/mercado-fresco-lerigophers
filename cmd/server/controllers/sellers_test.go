package controllers_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponse struct {
	Data []sellers.Seller
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func Test_Get_Seller_OK(t *testing.T) {
	t.Run("OK Case - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)
		mockSellerList := make([]sellers.Seller, 0)

		sellerController := controllers.NewSeller(mockedService)

		fakeSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		mockSellerList = append(mockSellerList, fakeSeller)

		mockedService.On("GetAll").Return(mockSellerList, web.ResponseCode{})

		router := gin.Default()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.GET("/api/v1/sellers/", sellerController.GetAll())
		router.ServeHTTP(rec, req)

		responseData, _ := ioutil.ReadAll(rec.Body)

		var currentResponse ObjectResponse

		err = json.Unmarshal(responseData, &currentResponse)

		assert.Nil(t, err)
		assert.Equal(t, fakeSeller, currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case - 500", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errors.New("internal server error"),
		})

		router := gin.Default()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.Handle(http.MethodGet, "/api/v1/sellers/", sellerController.GetAll())
		router.ServeHTTP(rec, req)

		responseData, err := ioutil.ReadAll(rec.Body)
		assert.Nil(t, err)

		var currentResponse web.ResponseCode

		err = json.Unmarshal(responseData, &currentResponse)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test_Get_One_Seller(t *testing.T) {
	t.Run("OK Case if exists - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sellerController := controllers.NewSeller(mockedService)

		fakeSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeSeller, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/sellers/:id", sellerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		type objResponse struct {
			Data sellers.Seller
		}

		var currentResponse objResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fakeSeller, currentResponse.Data)
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sellerController := controllers.NewSeller(mockedService)

		expectedError := errors.New("seller with id 1 not found")
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		router := gin.Default()
		router.GET("/api/v1/sellers/:id", sellerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/1", nil)
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

func Test_Delete_One_Seller(t *testing.T) {
	t.Run("OK Case if exists - 204", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		router := gin.Default()
		router.DELETE("/api/v1/sellers/:id", sellerController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/sellers/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Nil(t, err)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, "" == string(w.Body.String()))
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sellerController := controllers.NewSeller(mockedService)

		expectedError := errors.New("seller with id 1 not found")
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		router := gin.Default()
		router.GET("/api/v1/sellers/:id", sellerController.Delete())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/1", nil)
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
