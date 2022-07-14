package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/productRecords"
	product_records "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponse struct {
	Data product_records.ProductRecords
}

type ObjectResponseArr struct {
	Data []product_records.ProductRecords
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

var fakeInput = controllers.ReqProductRecord{
	LastUpdateDate: "2022-02-07",
	PurchasePrice:  3.0,
	SalePrice:      2.0,
	ProductId:      1,
}

const (
	defaultURL = "/api/v1/productRecords/"
)

func router() *gin.Engine {
	router := gin.Default()
	return router
}

func newProductRecordController() (*mocks.Service, *controllers.ProductRecordController) {
	mockedService := new(mocks.Service)
	productRecordsController := controllers.NewProductRecord(mockedService)
	return mockedService, productRecordsController
}

func TestCreateProductRecord(t *testing.T) {
	t.Run("Successfully on create product record", func(t *testing.T) {
		mockedService := new(mocks.Service)
		productRecordsController := controllers.NewProductRecord(mockedService)

		input := product_records.ProductRecords{
			Id:             1,
			LastUpdateDate: "2022-07-02",
			PurchasePrice:  3.0,
			SalePrice:      4.0,
			ProductId:      1,
		}

		parsedInput, err := json.Marshal(input)
		assert.Nil(t, err)

		expectedReturnData := product_records.ProductRecords{
			Id:             1,
			LastUpdateDate: input.LastUpdateDate,
			PurchasePrice:  input.PurchasePrice,
			SalePrice:      input.SalePrice,
			ProductId:      input.ProductId,
		}

		mockedService.On("CreateProductRecord",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(expectedReturnData, web.NewCodeResponse(http.StatusCreated, nil))

		router := gin.Default()
		router.POST("/api/v1/productRecords", productRecordsController.CreateProductRecord())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/productRecords", bytes.NewBuffer(parsedInput))
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var currentResponse ObjectResponse
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)
		assert.Equal(t, expectedReturnData, currentResponse.Data)

	})

	t.Run("Invalid request data", func(t *testing.T) {
		mockedService := new(mocks.Service)
		productRecordsController := controllers.NewProductRecord(mockedService)

		mockedService.On("CreateProductRecord",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(
			product_records.ProductRecords{},
			web.ResponseCode{})

		router := gin.Default()
		router.POST("/api/v1/productRecords", productRecordsController.CreateProductRecord())

		req, _ := http.NewRequest(
			http.MethodPost,
			"/api/v1/productRecords",
			bytes.NewBuffer([]byte(`sale_price: "a"`)),
		)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	})

	t.Run("Fail on create product_record", func(t *testing.T) {
		mockedService, ProductRecordController := newProductRecordController()
		mockedService.On("CreateProductRecord",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(
			product_records.ProductRecords{},
			web.ResponseCode{
				Code: http.StatusConflict,
				Err:  errors.New("any error"),
			})

		parsedFakeProductRecord, err := json.Marshal(fakeInput)
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, ProductRecordController.CreateProductRecord())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakeProductRecord),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
