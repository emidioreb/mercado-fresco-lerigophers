package employees

var (
	queryCreate       = `INSERT INTO employees(card_number_id, first_name, last_name, warehouse_id) VALUES(?, ?, ?, ?)`
	queryGetOne       = `SELECT * FROM employees WHERE id = ?`
	queryGetAll       = `SELECT * FROM employees`
	queryDelete       = "DELETE FROM employees WHERE id = ?"
	queryByCardNumber = `SELECT * FROM mercado_fresco.employees WHERE card_number_id = ? AND id <> ?`
	queryUpdate       = func(
		requestData map[string]interface{},
		id int) (
		finalQuery string,
		valuesToUse []interface{}) {
		prefixQuery := "UPDATE employees SET"
		fieldsToUpdate := []string{}
		whereCase := "WHERE id = ?"

		for key := range requestData {
			switch key {
			case "card_number_id":
				fieldsToUpdate = append(fieldsToUpdate, " card_number_id = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "first_name":
				fieldsToUpdate = append(fieldsToUpdate, " first_name = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "last_name":
				fieldsToUpdate = append(fieldsToUpdate, " last_name = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "warehouse_id":
				fieldsToUpdate = append(fieldsToUpdate, " warehouse_id = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
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
