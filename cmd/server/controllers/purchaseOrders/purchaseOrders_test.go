package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/purchaseOrders"
	purchase_orders "github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponsePurchaseOrder struct {
	Data purchase_orders.PurchaseOrders
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func router() *gin.Engine {
	router := gin.Default()
	return router
}

func newPurchaseOrdersController() (*mocks.Service, *controllers.PurchaseOrdersController) {
	mockedService := new(mocks.Service)
	purchaseOrdersController := controllers.NewPurchaseOrder(mockedService)
	return mockedService, purchaseOrdersController
}

const layout = "2006-01-02"

var date, _ = time.Parse(layout, layout)

var fakeInput = controllers.ReqPurchaseOrders{
	OrderNumber:     "#order1",
	OrderDate:       layout,
	TrackingCode:    "QB123400",
	BuyerId:         1,
	ProductRecordId: 1,
	OrderStatusId:   1,
}

var successfullyResponse = purchase_orders.PurchaseOrders{
	OrderNumber:     "#order1",
	OrderDate:       date,
	TrackingCode:    "QB123400",
	BuyerId:         1,
	ProductRecordId: 1,
	OrderStatusId:   1,
}

var fakePurchaseOrder = controllers.ReqPurchaseOrders{
	OrderNumber:     "#order1",
	OrderDate:       "2006-2",
	TrackingCode:    "QB123400",
	BuyerId:         1,
	ProductRecordId: 1,
	OrderStatusId:   1,
}

const (
	defaultURL = "/api/v1/purchaseOrders/"
)

func TestCreatePurchaseOrder(t *testing.T) {
	t.Run("Successfully on create purchase order", func(t *testing.T) {
		mockedService, PurchaseOrderController := newPurchaseOrdersController()

		mockedService.On("CreatePurchaseOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(successfullyResponse, web.ResponseCode{
			Code: http.StatusCreated, Err: nil,
		})

		parsedFakePurchaseOrder, err := json.Marshal(fakeInput)
		assert.NoError(t, err)

		router := gin.Default()
		router.POST(defaultURL, PurchaseOrderController.CreatePurchaseOrder())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakePurchaseOrder),
		)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Invalid request data", func(t *testing.T) {
		_, PurchaseOrderController := newPurchaseOrdersController()

		r := router()
		r.POST(defaultURL, PurchaseOrderController.CreatePurchaseOrder())

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

	t.Run("Unprocessable entity - order_date", func(t *testing.T) {
		_, PurchaseOrderController := newPurchaseOrdersController()

		parsedFakePurchaseOrder, err := json.Marshal(fakePurchaseOrder)
		assert.NoError(t, err)

		r := router()
		r.POST(defaultURL, PurchaseOrderController.CreatePurchaseOrder())

		req, err := http.NewRequest(
			http.MethodPost,
			defaultURL,
			bytes.NewBuffer(parsedFakePurchaseOrder),
		)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

}
