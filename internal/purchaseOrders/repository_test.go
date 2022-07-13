package purchase_orders_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	purchase_orders "github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders"
	"github.com/stretchr/testify/assert"
)

const layout = "2006-01-02"

var date, _ = time.Parse(layout, layout)

var mockPurchaseOrder = purchase_orders.PurchaseOrders{
	OrderNumber:     "#order-1",
	OrderDate:       date,
	TrackingCode:    "A1234",
	BuyerId:         1,
	ProductRecordId: 1,
	OrderStatusId:   1,
}

func TestCreate(t *testing.T) {
	query := `INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id, order_status_id) VALUES (?, ?, ?, ?, ?, ?);`

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				mockPurchaseOrder.OrderNumber,
				mockPurchaseOrder.OrderDate,
				mockPurchaseOrder.TrackingCode,
				mockPurchaseOrder.BuyerId,
				mockPurchaseOrder.ProductRecordId,
				mockPurchaseOrder.OrderStatusId,
			).WillReturnResult(sqlmock.NewResult(1, 1)) // last id, // rows affected

		purchaseOrderRepo := purchase_orders.NewMariaDbRepository(db)

		po, err := purchaseOrderRepo.CreatePurchaseOrders(
			mockPurchaseOrder.OrderNumber,
			mockPurchaseOrder.OrderDate,
			mockPurchaseOrder.TrackingCode,
			mockPurchaseOrder.BuyerId,
			mockPurchaseOrder.ProductRecordId,
			mockPurchaseOrder.OrderStatusId)
		assert.NoError(t, err)

		expectedTrackingCode := "A1234"

		assert.Equal(t, expectedTrackingCode, po.TrackingCode)
	})

	t.Run("failed to create", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(0, 0, 0, 0, 0, 0).
			WillReturnResult(sqlmock.NewResult(1, 1)) // last id, // rows affected

		purchaseOrderRepo := purchase_orders.NewMariaDbRepository(db)
		_, err = purchaseOrderRepo.CreatePurchaseOrders(
			mockPurchaseOrder.OrderNumber,
			mockPurchaseOrder.OrderDate,
			mockPurchaseOrder.TrackingCode,
			mockPurchaseOrder.BuyerId,
			mockPurchaseOrder.ProductRecordId,
			mockPurchaseOrder.OrderStatusId)

		assert.Error(t, err)
	})
}
