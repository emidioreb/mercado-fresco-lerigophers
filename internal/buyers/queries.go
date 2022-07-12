package buyers

var (
	QueryGetReportAll = `SELECT
							b.id,
							b.card_number_id,
							b.first_name,
							b.last_name,
							count(po.id) as purchase_orders_count
						FROM purchase_orders po
						RIGHT JOIN buyers b ON po.buyer_id = b.id
						GROUP BY 	
							b.id,
							b.card_number_id,
							b.first_name,
							b.last_name;`
	QueryGetReportOne = `SELECT
							b.id,
							b.card_number_id,
							b.first_name,
							b.last_name,
							count(*) as purchase_orders_count
						FROM purchase_orders po
						JOIN buyers b ON po.buyer_id = b.id
						WHERE buyer_id = ?
						GROUP BY 	
							b.id,
							b.card_number_id,
							b.first_name,
							b.last_name;`
)
