package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/carriers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponseArr struct {
	Data []carriers.Carry
}

type ObjectResponse struct {
	Data carriers.Carry
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func routerCarry() *gin.Engine {
	router := gin.Default()
	return router
}

func newCarryController() (*mocks.Service, *controllers.CarryController) {
	mockedService := new(mocks.Service)
	carryController := controllers.NewCarry(mockedService)
	return mockedService, carryController
}

const (
	defaultURL = "/api/v1/carries/"
)

var (
	errInvalidInput = errors.New("invalid request input")
	errCidTooLong = errors.New("CID too long: max 255 characters")
	errCompanyName = errors.New("company_name too long: max 255 characters")
	errAddress = errors.New("address too long: max 255 characters")
	errTelephone = errors.New("telephone too long: max 20 characters")
)

func TestCreateBuyer(t *testing.T) {
	fakeCarriers := carriers.Carry{
		Cid:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		LocalityId:  "456",
	}
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService, carryController := newCarryController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(fakeCarriers, web.ResponseCode{
				Code: http.StatusCreated,
			})

		parsedFakeCarry, err := json.Marshal(fakeCarriers)
		assert.Nil(t, err)

		r := gin.Default()
		r.POST(defaultURL, carryController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer(parsedFakeCarry))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeCarriers, bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService, carryController := newCarryController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(carriers.Carry{}, web.ResponseCode{})

		r := gin.Default()
		r.POST(defaultURL, carryController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"cid": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errInvalidInput.Error(), bodyResponse.Error)
	})

	t.Run("cid too long", func(t *testing.T) {
		mockedService, carryController := newCarryController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(carriers.Carry{}, web.ResponseCode{})

		r := gin.Default()
		r.POST(defaultURL, carryController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"cid": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCidTooLong.Error(), bodyResponse.Error)
	})

	t.Run("company_name too long", func(t *testing.T) {
		mockedService, carryController := newCarryController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(carriers.Carry{}, web.ResponseCode{})

		r := gin.Default()
		r.POST(defaultURL, carryController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"company_name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errCompanyName.Error(), bodyResponse.Error)
	})

	t.Run("address too long", func(t *testing.T) {
		mockedService, carryController := newCarryController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(carriers.Carry{}, web.ResponseCode{})

		r := gin.Default()
		r.POST(defaultURL, carryController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"address": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errAddress.Error(), bodyResponse.Error)
	})

	t.Run("telephone too long", func(t *testing.T) {
		mockedService, carryController := newCarryController()
		mockedService.On(
			"Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(carriers.Carry{}, web.ResponseCode{})

		r := gin.Default()
		r.POST(defaultURL, carryController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"telephone": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, errTelephone.Error(), bodyResponse.Error)
	})
}
