package product_batches

var (
	queryGetReportAll = `SELECT s.id as section_id, s.section_number, count(*) as sections_count
	FROM sections s
	JOIN product_batches pb ON s.id = pb.section_id GROUP BY s.id, s.section_number;`

	queryGetReportOne = `SELECT s.id as section_id, s.section_number, count(*) as sections_count
	FROM sections s
	JOIN product_batches pb ON s.id = pb.section_id WHERE s.id = ? GROUP BY s.id, s.section_number;`

	queryCreateProductBatch = `INSERT INTO product_batches (batch_number, current_quatity, current_temperature, initial_quantity, manufacturing_hour, minimum_temperature, product_id, section_id, due_date, manufacturing_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	queryGetOneProductBatch = `SELECT * FROM product_batches WHERE id = ?;`
)
