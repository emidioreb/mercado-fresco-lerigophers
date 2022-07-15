package employees

import "fmt"

var (
	queryCreate             = `INSERT INTO employees(card_number_id, first_name, last_name, warehouse_id) VALUES(?, ?, ?, ?)`
	queryGetOne             = `SELECT * FROM employees WHERE id = ?`
	queryGetAll             = `SELECT * FROM employees`
	queryDelete             = "DELETE FROM employees WHERE id = ?"
	queryByCardNumberUpdate = `SELECT id FROM employees WHERE card_number_id = ? AND id <> ?`
	queryByCardNumberCreate = `SELECT id FROM employees WHERE card_number_id = ?`
	queryUpdate             = func(
		requestData map[string]interface{},
		id int) (
		finalQuery string,
		valuesToUse []interface{}) {
		prefixQuery := "UPDATE employees SET"
		fieldsToUpdate := []string{}
		whereCase := "WHERE id = ?"

		var fields = []string{
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
		}
		for _, currField := range fields {
			if _, ok := requestData[currField]; ok {
				fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf(" %s = ?", currField))
				if currField == "warehouse_id" {
					valuesToUse = append(valuesToUse, int(requestData[currField].(float64)))
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
