package inboundorders

var (
	queryCreate = `INSERT INTO inbound_orders(order_number, order_date, employee_id, product_batch_id, warehouse_id)
	VALUES(?, ?, ?, ?, ?)`

	queryReportGetAll = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, count(*) as inbound_orders_count FROM mercado_fresco.inbound_orders i
	JOIN employees e ON i.employee_id = e.id
	GROUP BY e.id, e.card_number_id`

	queryReportGetOne = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, count(*) as inbound_orders_count FROM mercado_fresco.inbound_orders i
	JOIN employees e ON i.employee_id = e.id WHERE e.id = ?
	GROUP BY e.id, e.card_number_id`
)