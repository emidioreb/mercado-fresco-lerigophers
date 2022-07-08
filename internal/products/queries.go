package products

var (
	queryGetReportOne = `SELECT p.id as product_id, p.description, count(*) as products_count
						FROM products p
						JOIN product_records pr ON p.id = pr.product_id 
						WHERE pr.product_id = ? GROUP BY p.id`

	queryGetReportAll = `SELECT p.id as product_id, p.description, count(*) as products_count
						FROM products p
						JOIN product_records pr ON p.id = pr.product_id GROUP BY p.id`
)
