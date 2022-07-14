package buyers

import "fmt"

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
	QueryCreateBuyer = `INSERT INTO buyers(card_number_id, first_name, last_name) VALUES(?, ?, ?);`
	QueryGetOneBuyer = `SELECT * FROM buyers WHERE id = ?`
	QueryGetAllBuyer = `SELECT * FROM buyers`
	QueryDeleteBuyer = "DELETE FROM buyers WHERE id = ?"
	QueryUpdateBuyer = func(
		requestData map[string]interface{},
		id int) (
		finalQuery string,
		valuesToUse []interface{}) {
		prefixQuery := "UPDATE buyers SET"
		fieldsToUpdate := []string{}
		whereCase := "WHERE id = ?"

		var fields = []string{
			"card_number_id",
			"first_name",
			"last_name",
		}
		for _, currField := range fields {
			if _, ok := requestData[currField]; ok {
				fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf(" %s = ?", currField))
				if currField == "card_number_id" {
					valuesToUse = append(valuesToUse, requestData[currField])
				} else {
					valuesToUse = append(valuesToUse, requestData[currField])
				}
			}
		}

		valuesToUse = append(valuesToUse, id)
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
