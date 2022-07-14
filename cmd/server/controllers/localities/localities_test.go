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

type ObjectResponseCarriers struct {
	Data []localities.ReportCarriers
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

var fakeSellerReports = []localities.ReportSellers{
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

var fakeCarriersReports = []localities.ReportCarriers{
	{
		LocalityId:    "123",
		LocalityName:  "Locality",
		CarriersCount: 1,
	},
	{
		LocalityId:    "456",
		LocalityName:  "Locality",
		CarriersCount: 1,
	},
}

const (
	localitiesDefaultURL     = "/api/v1/localities/"
	reportOneSeller          = "/api/v1/localities/reportSellers?id=1"
	reportAllSellers         = "/api/v1/localities/reportSellers"
	defaultSellersReportURL  = "/api/v1/localities/reportSellers"
	defaultCarriersReportURL = "/api/v1/localities/reportCarriers"
	reportOneCarry           = "/api/v1/localities/reportCarriers?id=1"
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
		r.POST(localitiesDefaultURL, localityController.CreateLocality())

		req, err := http.NewRequest(
			http.MethodPost,
			localitiesDefaultURL,
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
				[]localities.ReportSellers{fakeSellerReports[0]},
				web.ResponseCode{Code: http.StatusOK},
			)

		r := routerSellers()
		r.GET(
			defaultSellersReportURL,
			localityController.GetReportSellers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOneSeller,
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
				fakeSellerReports,
				web.ResponseCode{Code: http.StatusOK},
			)

		r := routerSellers()
		r.GET(
			defaultSellersReportURL,
			localityController.GetReportSellers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportAllSellers,
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
			defaultSellersReportURL,
			localityController.GetReportSellers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportAllSellers,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetReportCarriers(t *testing.T) {
	t.Run("Test get report by one", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On(
			"GetReportCarriers",
			mock.AnythingOfType("string"),
		).
			Return(
				[]localities.ReportCarriers{fakeCarriersReports[0]},
				web.ResponseCode{Code: http.StatusOK},
			)

		r := routerSellers()
		r.GET(
			defaultCarriersReportURL,
			localityController.GetReportCarriers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOneCarry,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Locality id not found case", func(t *testing.T) {
		mockedService, localityController := newLocalitiesController()
		mockedService.On(
			"GetReportCarriers",
			mock.AnythingOfType("string"),
		).
			Return(
				[]localities.ReportCarriers{},
				web.ResponseCode{
					Code: http.StatusNotFound,
					Err:  errors.New("locality with id 1 not found"),
				},
			)

		r := routerSellers()
		r.GET(
			defaultCarriersReportURL,
			localityController.GetReportCarriers(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOneCarry,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var objectRespo ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &objectRespo)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "locality with id 1 not found", objectRespo.Error)
	})
}
