package inboundorders

type InboundOrder struct {
	Id             int    `json:"id"`
	OrderNumber    string `json:"order_number"`
	OrderDate      string `json:"order_date"`
	EmployeeId     int    `json:"employee_id"`
	ProductBatchId int    `json:"product_batch_id"`
	WarehouseId    int    `json:"warehouse_id"`
}

type ReportInboundOrder struct {
	Id                 int    `json:"id"`
	CardNumberId       string `json:"card_number_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	WarehouseId        int    `json:"warehouse_id"`
	InboundOrdersCount int    `json:"inbound_orders_count"`
}
