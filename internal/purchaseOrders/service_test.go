package purchase_orders_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	buyers_mock "github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers/mocks"
	order_status_mock "github.com/emidioreb/mercado-fresco-lerigophers/internal/orderStatus/mocks"
	product_records_mock "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords/mocks"
	purchase_orders "github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fakePurchaseOrders = []purchase_orders.PurchaseOrders{{
	OrderNumber:     "#order-1",
	OrderDate:       date,
	TrackingCode:    "A1234",
	BuyerId:         1,
	ProductRecordId: 1,
	OrderStatusId:   1,
}, {
	OrderNumber:     "#order-2",
	OrderDate:       date,
	TrackingCode:    "A1235",
	BuyerId:         2,
	ProductRecordId: 2,
	OrderStatusId:   2,
}}

func TestServiceCreate(t *testing.T) {
	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedBuyersRepository := new(buyers_mock.Repository)
		mockedProductRecordsRepository := new(product_records_mock.Repository)
		mockedOrderStatusRepository := new(order_status_mock.Repository)

		mockedBuyersRepository.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, nil)
		mockedProductRecordsRepository.On("GetOne", mock.AnythingOfType("int")).Return(nil)
		mockedOrderStatusRepository.On("GetOne", mock.AnythingOfType("int")).Return(nil)

		mockedRepository.On("CreatePurchaseOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(fakePurchaseOrders[0], nil)

		service := purchase_orders.NewService(mockedRepository, mockedBuyersRepository, mockedProductRecordsRepository, mockedOrderStatusRepository)

		result, err := service.CreatePurchaseOrders(
			fakePurchaseOrders[0].OrderNumber,
			fakePurchaseOrders[0].OrderDate,
			fakePurchaseOrders[0].TrackingCode,
			fakePurchaseOrders[0].BuyerId,
			fakePurchaseOrders[0].ProductRecordId,
			fakePurchaseOrders[0].OrderStatusId)
		assert.Nil(t, err.Err)

		assert.Equal(t, fakePurchaseOrders[0], result)
	})

	t.Run("Test conflict if buyer_id do not exist", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedBuyersRepository := new(buyers_mock.Repository)
		mockedProductRecordsRepository := new(product_records_mock.Repository)
		mockedOrderStatusRepository := new(order_status_mock.Repository)

		expectedError := errors.New("some error")
		mockedBuyersRepository.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, expectedError)

		service := purchase_orders.NewService(mockedRepository, mockedBuyersRepository, mockedProductRecordsRepository, mockedOrderStatusRepository)
		_, resp := service.CreatePurchaseOrders(
			fakePurchaseOrders[0].OrderNumber,
			fakePurchaseOrders[0].OrderDate,
			fakePurchaseOrders[0].TrackingCode,
			fakePurchaseOrders[0].BuyerId,
			fakePurchaseOrders[0].ProductRecordId,
			fakePurchaseOrders[0].OrderStatusId,
		)

		assert.Error(t, resp.Err)
		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("Test conflict if product_records do not exist", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedBuyersRepository := new(buyers_mock.Repository)
		mockedProductRecordsRepository := new(product_records_mock.Repository)
		mockedOrderStatusRepository := new(order_status_mock.Repository)

		expectedError := errors.New("some error")
		mockedBuyersRepository.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, nil)
		mockedProductRecordsRepository.On("GetOne", mock.AnythingOfType("int")).Return(expectedError)

		service := purchase_orders.NewService(mockedRepository, mockedBuyersRepository, mockedProductRecordsRepository, mockedOrderStatusRepository)
		_, resp := service.CreatePurchaseOrders(
			fakePurchaseOrders[0].OrderNumber,
			fakePurchaseOrders[0].OrderDate,
			fakePurchaseOrders[0].TrackingCode,
			fakePurchaseOrders[0].BuyerId,
			fakePurchaseOrders[0].ProductRecordId,
			fakePurchaseOrders[0].OrderStatusId,
		)

		assert.Error(t, resp.Err)
		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})
}
