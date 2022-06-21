package controllers_test

import (
	"bytes"
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

type ObjectResponseArr struct {
	Data []sellers.Seller
}

type ObjectResponse struct {
	Data sellers.Seller
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

		var currentResponse ObjectResponseArr

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
		assert.Nil(t, err)

		w := httptest.NewRecorder()

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

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)
		expectedError := errors.New("id must be a number")

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/sellers/:id", sellerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/string", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
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

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)
		expectedError := errors.New("id must be a number")

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.DELETE("/api/v1/sellers/:id", sellerController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/sellers/string", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})
}

func Test_Update_One_Seller(t *testing.T) {
	t.Run("OK Case if update sucessfully", func(t *testing.T) {
		mockedService := new(mocks.Service)

		fakeSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		parsedFakeSeller, err := json.Marshal(fakeSeller)
		assert.Nil(t, err)

		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeSeller, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSeller, bodyResponse.Data)
	})

	t.Run("Not found case", func(t *testing.T) {
		fakeSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		expectedError := errors.New("seller with id 1 not found")

		parsedFakeSeller, err := json.Marshal(fakeSeller)
		assert.Nil(t, err)
		mockedService := new(mocks.Service)

		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  expectedError,
			})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Id must be a number", func(t *testing.T) {
		fakeSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		expectedError := errors.New("id must be a number")

		parsedFakeSeller, err := json.Marshal(fakeSeller)
		assert.Nil(t, err)

		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/aaaa", bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Invalid request data", func(t *testing.T) {
		expectedError := errors.New("invalid request data")

		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer([]byte{}))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Body needed", func(t *testing.T) {
		expectedError := errors.New("invalid request data - body needed")

		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer([]byte("{}")))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("CID greather than 0", func(t *testing.T) {
		expectedError := errors.New("cid must be greather than 0")

		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer([]byte(`{"cid": 0 }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Syntax error on body", func(t *testing.T) {
		expectedError := errors.New("invalid type of data")

		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sellers/:id", sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer([]byte(`{"address": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})
}

func Test_Create_Seller(t *testing.T) {
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService := new(mocks.Service)

		fakeSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		parsedFakeSeller, err := json.Marshal(fakeSeller)
		assert.Nil(t, err)

		sellerController := controllers.NewSeller(mockedService)

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(fakeSeller, web.ResponseCode{
				Code: http.StatusCreated,
			})

		router := gin.Default()
		router.POST("/api/v1/sellers", sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSeller, bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		expectedError := errors.New("invalid request input")

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.POST("/api/v1/sellers", sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer([]byte(`{"cid": "vinicius"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("CID must be greather than 0", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		expectedError := errors.New("cid must be informed and greather than 0")

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := gin.Default()
		router.POST("/api/v1/sellers", sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer([]byte(`{"cid": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Conflict CID", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sellerController := controllers.NewSeller(mockedService)

		fakeSeller := []sellers.Seller{{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}, {
			Id:          2,
			Cid:         2,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number"},
		}

		expectedError := errors.New("cid already exists")

		mockedService.On("GetAll").Return(fakeSeller, nil)

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sellers.Seller{}, web.ResponseCode{
				Code: http.StatusConflict,
				Err:  expectedError,
			})

		router := gin.Default()
		router.POST("/api/v1/sellers", sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer([]byte(`{"cid": 2}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})
}
