package inboundorders_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	inboundorders "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/inboundOrders"
	inboundInternal "github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundOrders/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func routerInbounds() *gin.Engine {
	router := gin.Default()
	return router
}

func newInboundController() (*mocks.Service, *inboundorders.InboundOrdersController) {
	mockedService := new(mocks.Service)
	inboundController := inboundorders.NewInboud(mockedService)
	return mockedService, inboundController
}

var fakeReports = []inboundInternal.ReportInboundOrder{
	{
		Id:                 1,
		CardNumberId:       "456",
		FirstName:          "Iuri",
		LastName:           "Oi",
		WarehouseId:        1,
		InboundOrdersCount: 1,
	},
	{
		Id:                 2,
		CardNumberId:       "4565",
		FirstName:          "Iurizin",
		LastName:           "Tchau",
		WarehouseId:        2,
		InboundOrdersCount: 1,
	},
}

var fakeInbounds = []inboundInternal.InboundOrder{
	{
		Id:             1,
		OrderNumber:    "43",
		OrderDate:      "2006-01-02",
		EmployeeId:     1,
		ProductBatchId: 1,
		WarehouseId:    1,
	},
	{
		Id:             1,
		OrderNumber:    "434",
		OrderDate:      "2006-01-02",
		EmployeeId:     2,
		ProductBatchId: 1,
		WarehouseId:    2,
	},
}

const (
	inboundDefaultURL       = "/api/v1/inboundOrders"
	reportOneInbound        = "/api/v1/employees/reportInboundOrders?id=1"
	reportAllInbouds        = "/api/v1/employees/reportInboundOrders"
	defaultInboundReportURL = "/api/v1/employees/reportInboundOrders"
)

func TestGetReportInbound(t *testing.T) {
	t.Run("Test get report", func(t *testing.T) {
		mockedService, inboundController := newInboundController()
		mockedService.On(
			"GetReportInboundOrders",
			mock.AnythingOfType("string"),
		).
			Return(
				fakeReports,
				web.ResponseCode{Code: http.StatusOK},
			)

		r := routerInbounds()
		r.GET(
			defaultInboundReportURL,
			inboundController.GetReportInboundOrders(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOneInbound,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test get report", func(t *testing.T) {
		mockedService, inboundController := newInboundController()
		mockedService.On(
			"GetReportInboundOrders",
			mock.AnythingOfType("string"),
		).
			Return(
				[]inboundInternal.ReportInboundOrder{},
				web.ResponseCode{Code: http.StatusInternalServerError, Err: errors.New("error to get inbounds")},
			)

		r := routerInbounds()
		r.GET(
			defaultInboundReportURL,
			inboundController.GetReportInboundOrders(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOneInbound,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateInbound(t *testing.T) {
	t.Run("Successfully on create locality", func(t *testing.T) {
		mockedService, inboundController := newInboundController()
		mockedService.On("CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(fakeInbounds[0], web.ResponseCode{
			Code: http.StatusCreated,
		})

		parsedFakeInbound, err := json.Marshal(fakeInbounds[0])
		assert.NoError(t, err)

		r := routerInbounds()
		r.POST(inboundDefaultURL, inboundController.CreateInboundOrders())

		req, err := http.NewRequest(
			http.MethodPost,
			inboundDefaultURL,
			bytes.NewBuffer(parsedFakeInbound),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		_, inboundController := newInboundController()

		r := routerInbounds()
		r.POST(inboundDefaultURL, inboundController.CreateInboundOrders())

		req, err := http.NewRequest(
			http.MethodPost,
			inboundDefaultURL,
			bytes.NewBuffer([]byte("")),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Invalid order_number", func(t *testing.T) {
		_, inboundController := newInboundController()

		fakeInboundsFail := inboundInternal.InboundOrder{
			Id:             1,
			OrderNumber:    "434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444434444444443444444444344444444",
			OrderDate:      "2006-01-02",
			EmployeeId:     2,
			ProductBatchId: 1,
			WarehouseId:    2,
		}

		parsedFakeInbound, err := json.Marshal(fakeInboundsFail)
		assert.NoError(t, err)

		r := routerInbounds()
		r.POST(inboundDefaultURL, inboundController.CreateInboundOrders())

		req, err := http.NewRequest(
			http.MethodPost,
			inboundDefaultURL,
			bytes.NewBuffer(parsedFakeInbound),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Invalid order_date", func(t *testing.T) {
		_, inboundController := newInboundController()

		fakeInboundsFail := inboundInternal.InboundOrder{
			Id:             1,
			OrderNumber:    "434",
			OrderDate:      "207786-01-0266",
			EmployeeId:     2,
			ProductBatchId: 1,
			WarehouseId:    2,
		}

		parsedFakeInbound, err := json.Marshal(fakeInboundsFail)
		assert.NoError(t, err)

		r := routerInbounds()
		r.POST(inboundDefaultURL, inboundController.CreateInboundOrders())

		req, err := http.NewRequest(
			http.MethodPost,
			inboundDefaultURL,
			bytes.NewBuffer(parsedFakeInbound),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Fails on create locality", func(t *testing.T) {
		mockedService, inboundController := newInboundController()
		mockedService.On("CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundInternal.InboundOrder{}, web.ResponseCode{
			Code: http.StatusInternalServerError, Err: errors.New("any error"),
		})

		parsedFakeInbound, err := json.Marshal(fakeInbounds[0])
		assert.NoError(t, err)

		r := routerInbounds()
		r.POST(inboundDefaultURL, inboundController.CreateInboundOrders())

		req, err := http.NewRequest(
			http.MethodPost,
			inboundDefaultURL,
			bytes.NewBuffer(parsedFakeInbound),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
