package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/sellers"
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

func routerSellers() *gin.Engine {
	router := gin.Default()
	return router
}

func newSellerController() (*mocks.Service, *controllers.SellerController) {
	mockedService := new(mocks.Service)
	sellerController := controllers.NewSeller(mockedService)
	return mockedService, sellerController
}

var fakeSellers = []sellers.Seller{{
	Id:          1,
	Cid:         1,
	CompanyName: "Fake Business",
	Address:     "Fake Address",
	Telephone:   "Fake Number",
	LocalityId:  "12345",
}, {
	Id:          2,
	Cid:         2,
	CompanyName: "Fake Business",
	Address:     "Fake Address",
	Telephone:   "Fake Number",
	LocalityId:  "67890",
},
}

const (
	defaultURL = "/api/v1/sellers/"
	idString   = "/api/v1/sellers/string"
	idNumber1  = "/api/v1/sellers/1"
	idRequest  = "/api/v1/sellers/:id"
)

var (
	errServer         = errors.New("internal server error")
	errSellerNotFound = errors.New("seller with id 1 not found")
	errIdNotNumber    = errors.New("id must be a number")
	errInvalidRequest = errors.New("invalid request data")
	errNeedBody       = errors.New("invalid request data - body needed")
	errCidZero        = errors.New("cid must be greather than 0")
	errCidNeeded      = errors.New("cid must be informed and greather than 0")
	errTypeData       = errors.New("invalid type of data")
	errInvalidInput   = errors.New("invalid request input")
	errCidExists      = errors.New("cid already exists")
)

func TestGetSeller(t *testing.T) {
	t.Run("Get all sellers", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("GetAll").Return(fakeSellers, web.ResponseCode{})

		r := routerSellers()
		r.GET(defaultURL, sellerController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponseArr
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSellers[0], currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errServer,
		})

		r := routerSellers()
		r.GET(defaultURL, sellerController.GetAll())

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
		mockedService, sellerController := newSellerController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeSellers[0], web.ResponseCode{})

		r := routerSellers()
		r.GET(idRequest, sellerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fakeSellers[0], currentResponse.Data)
	})

	t.Run("Not exist case", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errSellerNotFound,
		})

		r := routerSellers()
		r.GET(idRequest, sellerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errSellerNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, web.ResponseCode{})

		r := routerSellers()
		r.GET(idRequest, sellerController.GetOne())

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

func TestDeleteSeller(t *testing.T) {
	t.Run("Success case if exists", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		r := routerSellers()
		r.DELETE(idRequest, sellerController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.True(t, "" == rec.Body.String())
	})

	t.Run("Error case if not exists", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errSellerNotFound,
		})

		r := routerSellers()
		r.DELETE(idRequest, sellerController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errSellerNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(sellers.Seller{}, web.ResponseCode{})

		r := routerSellers()
		r.DELETE(idRequest, sellerController.Delete())

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

func TestUpdateSeller(t *testing.T) {
	t.Run("Sucessfully case", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeSellers[0], web.ResponseCode{})

		parsedFakeSeller, err := json.Marshal(fakeSellers[0])
		assert.Nil(t, err)

		r := routerSellers()
		r.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSellers[0], bodyResponse.Data)
	})

	t.Run("Not found case", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  errSellerNotFound,
			})

		r := routerSellers()
		r.PATCH(idRequest, sellerController.Update())

		parsedFakeSeller, err := json.Marshal(fakeSellers[0])
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errSellerNotFound.Error(), bodyResponse.Error)
	})

	t.Run("Id must be a number", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		parsedFakeSeller, err := json.Marshal(fakeSellers[0])
		assert.Nil(t, err)

		r := routerSellers()
		r.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idString, bytes.NewBuffer(parsedFakeSeller))
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
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		r := routerSellers()
		r.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte{}))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidRequest.Error(), bodyResponse.Error)
	})

	t.Run("Body needed", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := routerSellers()
		router.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte("{}")))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errNeedBody.Error(), bodyResponse.Error)
	})

	t.Run("CID greather than 0", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := routerSellers()
		router.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"cid": 0 }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCidZero.Error(), bodyResponse.Error)
	})

	t.Run("Syntax error on body", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sellers.Seller{}, web.ResponseCode{})

		router := routerSellers()
		router.PATCH(idRequest, sellerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"address": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errTypeData.Error(), bodyResponse.Error)
	})
}

func TestCreateSeller(t *testing.T) {
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(fakeSellers[0], web.ResponseCode{
				Code: http.StatusCreated,
			})

		parsedFakeSeller, err := json.Marshal(fakeSellers[0])
		assert.Nil(t, err)

		r := gin.Default()
		r.POST(defaultURL, sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer(parsedFakeSeller))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSellers[0], bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sellers.Seller{}, web.ResponseCode{})

		r := routerSellers()
		r.POST(defaultURL, sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"cid": "vinicius"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidInput.Error(), bodyResponse.Error)
	})

	t.Run("CID must be greather than 0", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(sellers.Seller{}, web.ResponseCode{})

		router := routerSellers()
		router.POST(defaultURL, sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"cid": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCidNeeded.Error(), bodyResponse.Error)
	})

	t.Run("Conflict CID", func(t *testing.T) {
		mockedService, sellerController := newSellerController()
		mockedService.On("GetAll").Return(fakeSellers, nil)
		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(sellers.Seller{}, web.ResponseCode{
			Code: http.StatusConflict,
			Err:  errCidExists,
		})

		r := routerSellers()
		r.POST(defaultURL, sellerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"cid": 2}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCidExists.Error(), bodyResponse.Error)
	})
}
