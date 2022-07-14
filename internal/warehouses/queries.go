package warehouses

var (
	queryCreateWarehouse = `INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature) VALUES (?, ?, ?, ?, ?)`

	queryGetOneWarehouse  = "SELECT * FROM warehouses WHERE id = ?"
	queryGetAllWarehouses = "SELECT * FROM warehouses"
	queryDeleteWarehouse  = "DELETE FROM warehouses WHERE id = ?"
)
