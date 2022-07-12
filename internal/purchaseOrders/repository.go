package purchase_orders

import (
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	CreatePurchaseOrders(OrderNumber string, OrderDate time.Time, TrackingCode string, BuyerId, ProductRecordId, OrderStatusId int) (PurchaseOrders, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

var (
	errCreatePurchaseOrders = errors.New("couldn't create purchase order")
)

func (mariaDb mariaDbRepository) CreatePurchaseOrders(OrderNumber string, OrderDate time.Time, TrackingCode string, BuyerId, ProductRecordId, OrderStatusId int) (PurchaseOrders, error) {
	result, err := mariaDb.db.Exec(QueryCreatePurchaseOrder, OrderNumber, OrderDate, TrackingCode, BuyerId, ProductRecordId, OrderStatusId)
	if err != nil {
		return PurchaseOrders{}, errCreatePurchaseOrders
	}

	newPurchaseOrder := PurchaseOrders{
		OrderNumber:     OrderNumber,
		OrderDate:       OrderDate,
		TrackingCode:    TrackingCode,
		BuyerId:         BuyerId,
		ProductRecordId: ProductRecordId,
		OrderStatusId:   OrderStatusId,
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return PurchaseOrders{}, errors.New("ocurred an error to create purchase order")
	}

	newPurchaseOrder.Id = int(lastId)

	return newPurchaseOrder, nil
}
