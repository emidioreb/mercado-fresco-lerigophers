package inboundorders

var (
	QueryCreate = `INSERT INTO inbound_orders(order_number, order_date, employee_id, product_batch_id, warehouse_id)
	VALUES(?, ?, ?, ?, ?)`

	QueryReportGetAll = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, count(*) as inbound_orders_count FROM inbound_orders i
	JOIN employees e ON i.employee_id = e.id
	GROUP BY e.id, e.card_number_id`

	QueryReportGetOne = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, count(*) as inbound_orders_count FROM inbound_orders i
	JOIN employees e ON i.employee_id = e.id WHERE e.id = ?
	GROUP BY e.id, e.card_number_id`
)
