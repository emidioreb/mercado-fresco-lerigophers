package purchase_orders

import (
	"net/http"
	"time"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	order_status "github.com/emidioreb/mercado-fresco-lerigophers/internal/orderStatus"
	product_records "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	CreatePurchaseOrders(OrderNumber string, OrderDate time.Time, TrackingCode string, BuyerId, ProductRecordId, OrderStatusId int) (PurchaseOrders, web.ResponseCode)
}

type service struct {
	repository               Repository
	buyerRepository          buyers.Repository
	productRecordsRepository product_records.Repository
	orderStatusRepository    order_status.Repository
}

func NewService(r Repository, br buyers.Repository, prr product_records.Repository, osr order_status.Repository) Service {
	return &service{
		repository:               r,
		buyerRepository:          br,
		productRecordsRepository: prr,
		orderStatusRepository:    osr,
	}
}

func (s service) CreatePurchaseOrders(OrderNumber string, OrderDate time.Time, TrackingCode string, BuyerId, ProductRecordId, OrderStatusId int) (PurchaseOrders, web.ResponseCode) {

	_, err := s.buyerRepository.GetOne(BuyerId)
	if err != nil {
		return PurchaseOrders{}, web.NewCodeResponse(http.StatusConflict, err)
	}

	err = s.productRecordsRepository.GetOne(ProductRecordId)
	if err != nil {
		return PurchaseOrders{}, web.NewCodeResponse(http.StatusConflict, err)
	}

	err = s.orderStatusRepository.GetOne(OrderStatusId)
	if err != nil {
		return PurchaseOrders{}, web.NewCodeResponse(http.StatusConflict, err)
	}

	result, err := s.repository.CreatePurchaseOrders(OrderNumber, OrderDate, TrackingCode, BuyerId, ProductRecordId, OrderStatusId)
	if err != nil {
		return PurchaseOrders{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return result, web.NewCodeResponse(http.StatusCreated, nil)
}
