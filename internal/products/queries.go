package products

import "fmt"

var (
	queryGetReportOne = `SELECT p.id as product_id, p.description, count(*) as products_count
						FROM products p
						JOIN product_records pr ON p.id = pr.product_id 
						WHERE pr.product_id = ? GROUP BY p.id`

	queryGetReportAll = `SELECT p.id as product_id, p.description, count(*) as products_count
						FROM products p
						JOIN product_records pr ON p.id = pr.product_id GROUP BY p.id`

	queryCreateProduct = `INSERT INTO products (product_code,
		description,
		width,
		height,
		length,
		net_weight,
		expiration_rate,
		recommended_freezing_temperature,
		freezing_rate,
		product_type_id,
		seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	queryGetOneProduct = `SELECT * FROM products WHERE id = ?`

	queryGetAllProducts = `SELECT * FROM products`

	queryDeleteProduct = "DELETE FROM products WHERE id = ?"

	queryUpdateProduct = func(
		requestData map[string]interface{},
		id int) (
		finalQuery string,
		valuesToUse []interface{}) {
		prefixQuery := "UPDATE products SET"
		fieldsToUpdate := []string{}
		whereCase := "WHERE id = ?"

		var fields = []string{
			"product_code",
			"description",
			"width",
			"height",
			"length",
			"net_weight",
			"expiration_rate",
			"recommended_freezing_temperature",
			"freezing_rate",
			"product_type_id",
			"seller_id",
		}
		for _, currField := range fields {
			if _, ok := requestData[currField]; ok {
				fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf(" %s = ?", currField))
				if currField == "seller_id" || currField == "product_type_id" {
					valuesToUse = append(valuesToUse, requestData[currField].(int))
				} else if currField == "product_code" || currField == "description" {
					valuesToUse = append(valuesToUse, requestData[currField])
				} else {
					valuesToUse = append(valuesToUse, requestData[currField].(float64))
				}
			}
		}

		finalQuery += prefixQuery
		for index, field := range fieldsToUpdate {
			if index+1 == len(fieldsToUpdate) {
				finalQuery += field + " "
			} else {
				finalQuery += field + ", "
			}
		}
		finalQuery += whereCase

		return finalQuery, valuesToUse
	}
)
