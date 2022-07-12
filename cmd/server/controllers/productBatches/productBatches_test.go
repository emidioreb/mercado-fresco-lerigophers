package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/productBatches"
	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponseProductBatch struct {
	Data product_batches.ProductBatches
}

type ObjectResponseReports struct {
	Data []product_batches.ProductsQuantity
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func router() *gin.Engine {
	router := gin.Default()
	return router
}

func newProductBatcheController() (*mocks.Service, *controllers.ProductBatchController) {
	mockedService := new(mocks.Service)
	productBatchesController := controllers.NewProductBatch(mockedService)
	return mockedService, productBatchesController
}

const layout = "2006-01-02"

var date, _ = time.Parse(layout, layout)

var fakeInput = controllers.ReqProductBatch{
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 8,
	ProductId:          23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}

var successfullyResponse = product_batches.ProductBatches{
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 8,
	ProductId:          23,
	SectionId:          56,
	DueDate:            date,
	ManufacturingDate:  date,
}

var fakeProductBatches = []controllers.ReqProductBatch{{
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 8,
	ProductId:          23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}, {
	BatchNumber:        -1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}, {
	BatchNumber:        1,
	CurrentQuantity:    -10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    -500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          -23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          -56,
	DueDate:            layout,
	ManufacturingDate:  layout,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            "2006-2",
	ManufacturingDate:  layout,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            layout,
	ManufacturingDate:  "200",
}}

var fakeReports = []product_batches.ProductsQuantity{
	{
		SectionId:     1,
		SectionNumber: 1,
		ProductsCount: 1,
	},
	{
		SectionId:     2,
		SectionNumber: 2,
		ProductsCount: 2,
	},
}

const (
	defaultURL       = "/api/v1/productBatches/"
	reportOne        = "/api/v1/productBatches/reportProducts?id=1"
	reportAll        = "/api/v1/productBatches/reportProducts/"
	defaultReportURL = "/api/v1/productBatches/reportProducts"
)

func TestCreateProductBatch(t *testing.T) {
	t.Run("Successfully on create product_batch", func(t *testing.T) {
		mockedService, ProductBatchController := newProductBatcheController()

		mockedService.On("CreateProductBatch",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"),
		).Return(successfullyResponse, web.ResponseCode{
			Code: http.StatusCreated, Err: nil,
		})

		parsedFakeProductBatch, err := json.Marshal(fakeInput)
		assert.NoError(t, err)

		router := gin.Default()
		router.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Fail on create product_batch", func(t *testing.T) {
		mockedService, ProductBatchController := newProductBatcheController()
		mockedService.On("CreateProductBatch",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"),
		).Return(
			product_batches.ProductBatches{},
			web.ResponseCode{
				Code: http.StatusConflict,
				Err:  errors.New("batch_number already exists"),
			})

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[0])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Unprocessable entity 1 - batch_number", func(t *testing.T) {

		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[1])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 2 - current_quantity", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[2])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 3 - initial_quantity", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[3])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 4 - product_id", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[4])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 5 - section_id", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[5])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 6 - due_date", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[6])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unprocessable entity 7 - manufacturing_date", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		parsedFakeProductBatch, err := json.Marshal(fakeProductBatches[7])
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductBatch),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		_, ProductBatchController := newProductBatcheController()

		r := router()
		r.POST(defaultURL, ProductBatchController.CreateProductBatch())

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

func TestGetReportSections(t *testing.T) {
	t.Run("Test get report by one", func(t *testing.T) {
		mockedService, ProductBatchController := newProductBatcheController()
		mockedService.On(
			"GetReportSection",
			mock.AnythingOfType("int"),
		).Return([]product_batches.ProductsQuantity{}, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errors.New("any error"),
		},
		)

		r := router()
		r.GET(
			defaultReportURL,
			ProductBatchController.GetReportSection(),
		)

		req, err := http.NewRequest(
			http.MethodGet,
			reportOne,
			nil,
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Test get report", func(t *testing.T) {
		mockedService, ProductBatchController := newProductBatcheController()
		mockedService.On(
			"GetReportSection",
			mock.AnythingOfType("int"),
		).Return(
			[]product_batches.ProductsQuantity{fakeReports[0]},
			web.ResponseCode{Code: http.StatusOK},
		)

		r := router()
		r.GET(
			defaultReportURL,
			ProductBatchController.GetReportSection(),
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
}
