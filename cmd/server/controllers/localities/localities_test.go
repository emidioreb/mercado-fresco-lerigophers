package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/localities"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponseLocality struct {
	Data localities.Locality
}

type ObjectResponseReports struct {
	Data []localities.ReportSellers
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func routerSellers() *gin.Engine {
	router := gin.Default()
	return router
}

func newLocalitiesController() (*mocks.Service, *controllers.LocalityController) {
	mockedService := new(mocks.Service)
	localitiesController := controllers.NewLocality(mockedService)
	return mockedService, localitiesController
}

var fakeLocalities = []localities.Locality{
	{
		Id:           "65760000",
		LocalityName: "Presidente Dutra",
		ProvinceName: "MA",
		CountryName:  "BR",
	},
	{
		Id:           "12345678",
		LocalityName: "Florianópolis",
		ProvinceName: "SC",
		CountryName:  "BR",
	},
}

var fakeReports = []localities.ReportSellers{
	{
		LocalityId:   "65760000",
		LocalityName: "Presidente Dutra",
		SellersCount: 1,
	},
	{
		LocalityId:   "12345678",
		LocalityName: "Florianópolis",
		SellersCount: 1,
	},
}

const (
	defaultURL       = "/api/v1/localities/"
	reportOne        = "/api/v1/localities/reportSellers?id=1"
	reportAll        = "/api/v1/localities/reportSellers"
	defaultReportURL = "/api/v1/localities/reportSellers"
)

func TestCreateLocality(t *testing.T) {
	t.Run("Successfully on create locality", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On("CreateLocality",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(fakeLocalities[0], web.ResponseCode{
			Code: http.StatusCreated,
		})

		parsedFakeLocality, err := json.Marshal(fakeLocalities[0])
		assert.NoError(t, err)

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeLocality),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Fail on create locality", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On("CreateLocality",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(
			localities.Locality{},
			web.ResponseCode{
				Code: http.StatusConflict,
				Err:  errors.New("locality_id already exists"),
			})

		parsedFakeLocality, err := json.Marshal(fakeLocalities[0])
		assert.NoError(t, err)

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeLocality),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Unprocessable entity 1 - CountryName", func(t *testing.T) {
		fakeLoc := localities.Locality{}
		fakeLoc.Id = "id"
		fakeLoc.LocalityName = "name"
		fakeLoc.ProvinceName = "province"
		for i := 0; i < 256; i++ {
			fakeLoc.CountryName = fakeLoc.CountryName + "a"
		}

		_, localityController := newLocalitiesController()

		parsedFakeLocality, err := json.Marshal(fakeLoc)
		assert.NoError(t, err)

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeLocality),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 2 - Id", func(t *testing.T) {
		fakeLoc := localities.Locality{}
		fakeLoc.LocalityName = "name"
		fakeLoc.ProvinceName = "province"
		fakeLoc.CountryName = "country"
		for i := 0; i < 256; i++ {
			fakeLoc.Id = fakeLoc.Id + "a"
		}

		_, localityController := newLocalitiesController()

		parsedFakeLocality, err := json.Marshal(fakeLoc)
		assert.NoError(t, err)

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeLocality),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 3 - ProvinceName", func(t *testing.T) {
		fakeLoc := localities.Locality{}
		fakeLoc.Id = "id"
		fakeLoc.LocalityName = "name"
		fakeLoc.CountryName = "country"
		for i := 0; i < 256; i++ {
			fakeLoc.ProvinceName = fakeLoc.ProvinceName + "a"
		}

		_, localityController := newLocalitiesController()

		parsedFakeLocality, err := json.Marshal(fakeLoc)
		assert.NoError(t, err)

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeLocality),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 4 - LocalityName", func(t *testing.T) {
		fakeLoc := localities.Locality{}
		fakeLoc.Id = "id"
		fakeLoc.ProvinceName = "province"
		fakeLoc.CountryName = "country"
		for i := 0; i < 256; i++ {
			fakeLoc.LocalityName = fakeLoc.LocalityName + "a"
		}

		_, localityController := newLocalitiesController()

		parsedFakeLocality, err := json.Marshal(fakeLoc)
		assert.NoError(t, err)

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeLocality),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		_, localityController := newLocalitiesController()

		r := routerSellers()
		r.POST(defaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer([]byte("")),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}

func TestGetReportSellers(t *testing.T) {
	t.Run("Test get report by one", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On(
			"GetReportOneSeller",
			mock.AnythingOfType("string"),
		).
			Return(
				[]localities.ReportSellers{fakeReports[0]},
				web.ResponseCode{Code: http.StatusOK},
			)

		r := routerSellers()
		r.GET(
			defaultReportURL,
			localityController.GetReportSellers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOne,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test get report all", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On("GetAllReportSellers").
			Return(
				fakeReports,
				web.ResponseCode{Code: http.StatusOK},
			)

		r := routerSellers()
		r.GET(
			defaultReportURL,
			localityController.GetReportSellers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportAll,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Fail when get report", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On("GetAllReportSellers").
			Return(
				[]localities.ReportSellers{},
				web.ResponseCode{
					Code: http.StatusInternalServerError,
					Err:  errors.New("error to get localities"),
				},
			)

		r := routerSellers()
		r.GET(
			defaultReportURL,
			localityController.GetReportSellers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportAll,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
