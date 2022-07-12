package purchase_orders

import "time"

type PurchaseOrders struct {
	Id              int       `json:"id"`
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerId         int       `json:"buyer_id"`
	ProductRecordId int       `json:"product_record_id"`
	OrderStatusId   int       `json:"order_status_id"`
}
