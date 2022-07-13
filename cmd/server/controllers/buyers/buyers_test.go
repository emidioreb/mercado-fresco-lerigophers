package buyers_controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/buyers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponseArr struct {
	Data []buyers.Buyer
}

type ObjectResponse struct {
	Data buyers.Buyer
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func routerBuyers() *gin.Engine {
	router := gin.Default()
	return router
}

func newBuyerController() (*mocks.Service, *controller.BuyerController) {
	mockedService := new(mocks.Service)
	buyerController := controller.NewBuyer(mockedService)
	return mockedService, buyerController
}

var fakeBuyers = []buyers.Buyer{{
	Id:           1,
	CardNumberId: "12345",
	FirstName:    "Fulano",
	LastName:     "Beltrano",
}, {
	Id:           2,
	CardNumberId: "12345",
	FirstName:    "Fulano",
	LastName:     "Beltrano"},
}

const (
	defaultURL = "/api/v1/buyers/"
	idString   = "/api/v1/buyers/string"
	idNumber1  = "/api/v1/buyers/1"
	idRequest  = "/api/v1/buyers/:id"
)

var (
	errServer             = errors.New("internal server error")
	errBuyerNotFound      = errors.New("buyer with id 1 not found")
	errIdNotNumber        = errors.New("id must be a number")
	errInvalidRequest     = errors.New("invalid request data")
	errNeedBody           = errors.New("invalid request data - body needed")
	errCardNumberIdEmpty  = errors.New("empty card_number_id not allowed")
	errTypeData           = errors.New("invalid type of data")
	errInvalidInput       = errors.New("invalid request input")
	errCardNumberIdExists = errors.New("CardNumberId already exists")
)

func TestGetBuyer(t *testing.T) {
	t.Run("Get all buyers", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("GetAll").Return(fakeBuyers, web.ResponseCode{})

		r := routerBuyers()
		r.GET(defaultURL, buyerController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponseArr
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeBuyers[0], currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errServer,
		})

		r := routerBuyers()
		r.GET(defaultURL, buyerController.GetAll())

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
		mockedService, buyerController := newBuyerController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeBuyers[0], web.ResponseCode{})

		r := routerBuyers()
		r.GET(idRequest, buyerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fakeBuyers[0], currentResponse.Data)
	})

	t.Run("Not exist case", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errBuyerNotFound,
		})

		r := routerBuyers()
		r.GET(idRequest, buyerController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errBuyerNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, web.ResponseCode{})

		r := routerBuyers()
		r.GET(idRequest, buyerController.GetOne())

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

func TestDeleteBuyer(t *testing.T) {
	t.Run("Success case if exists", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		r := routerBuyers()
		r.DELETE(idRequest, buyerController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.True(t, "" == rec.Body.String())
	})

	t.Run("Error case if not exists", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  errBuyerNotFound,
		})

		r := routerBuyers()
		r.DELETE(idRequest, buyerController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, errBuyerNotFound.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(buyers.Buyer{}, web.ResponseCode{})

		r := routerBuyers()
		r.DELETE(idRequest, buyerController.Delete())

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

func TestUpdateBuyer(t *testing.T) {
	t.Run("Sucessfully case", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeBuyers[0], web.ResponseCode{})

		parsedFakeBuyer, err := json.Marshal(fakeBuyers[0])
		assert.Nil(t, err)

		r := routerBuyers()
		r.PATCH(idRequest, buyerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeBuyer))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeBuyers[0], bodyResponse.Data)
	})

	t.Run("Not found case", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(buyers.Buyer{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  errBuyerNotFound,
			})

		r := routerBuyers()
		r.PATCH(idRequest, buyerController.Update())

		parsedFakeBuyer, err := json.Marshal(fakeBuyers[0])
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeBuyer))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errBuyerNotFound.Error(), bodyResponse.Error)
	})

	t.Run("Id must be a number", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(buyers.Buyer{}, web.ResponseCode{})

		parsedFakeBuyer, err := json.Marshal(fakeBuyers[0])
		assert.Nil(t, err)

		r := routerBuyers()
		r.PATCH(idRequest, buyerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idString, bytes.NewBuffer(parsedFakeBuyer))
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
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(buyers.Buyer{}, web.ResponseCode{})

		r := routerBuyers()
		r.PATCH(idRequest, buyerController.Update())

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
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(buyers.Buyer{}, web.ResponseCode{})

		router := routerBuyers()
		router.PATCH(idRequest, buyerController.Update())

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

	t.Run("CardNumberId can't be empty", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(buyers.Buyer{}, web.ResponseCode{})

		router := routerBuyers()
		router.PATCH(idRequest, buyerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"card_number_id": "" }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCardNumberIdEmpty.Error(), bodyResponse.Error)
	})

	t.Run("Syntax error on body", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(buyers.Buyer{}, web.ResponseCode{})

		router := routerBuyers()
		router.PATCH(idRequest, buyerController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"first_name": 0}`)))
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

func TestCreateBuyer(t *testing.T) {
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(fakeBuyers[0], web.ResponseCode{
				Code: http.StatusCreated,
			})

		parsedFakeBuyer, err := json.Marshal(fakeBuyers[0])
		assert.Nil(t, err)

		r := gin.Default()
		r.POST(defaultURL, buyerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer(parsedFakeBuyer))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeBuyers[0], bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(buyers.Buyer{}, web.ResponseCode{})

		r := routerBuyers()
		r.POST(defaultURL, buyerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"first_name": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidInput.Error(), bodyResponse.Error)
	})

	t.Run("CardNumberId can't be empty", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(buyers.Buyer{}, web.ResponseCode{})

		router := routerBuyers()
		router.POST(defaultURL, buyerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"card_number_id": "" }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCardNumberIdEmpty.Error(), bodyResponse.Error)
	})

	t.Run("Conflict Card Number Id", func(t *testing.T) {
		mockedService, buyerController := newBuyerController()
		mockedService.On("GetAll").Return(fakeBuyers, nil)
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(buyers.Buyer{}, web.ResponseCode{
			Code: http.StatusConflict,
			Err:  errCardNumberIdExists,
		})

		r := routerBuyers()
		r.POST(defaultURL, buyerController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"card_number_id": "54321"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCardNumberIdExists.Error(), bodyResponse.Error)
	})
}

// TODO - implement test GetReportPurchaseOrders
