package sections

var (
	queryCreateSection      = "INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetOneSection      = "SELECT * FROM sections WHERE id = ?"
	queryGetAllSections     = "SELECT * FROM sections"
	queryDeleteSection      = "DELETE FROM sections WHERE id = ?"
	queryValidSectionNumber = "SELECT id, section_number FROM sections where section_number = ?"
	queryUpdateSection      = func(requestData map[string]interface{}, id int) (finalQuery string, valuesToUse []interface{}) {
		prefixQuery := "UPDATE sections SET"
		fieldsToUpdate := []string{}
		whereCase := "WHERE id = ?"

		for key := range requestData {
			switch key {
			case "section_number":
				fieldsToUpdate = append(fieldsToUpdate, " section_number = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "current_temperature":
				fieldsToUpdate = append(fieldsToUpdate, " current_temperature = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "minimum_temperature":
				fieldsToUpdate = append(fieldsToUpdate, " minimum_temperature = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "current_capacity":
				fieldsToUpdate = append(fieldsToUpdate, " current_capacity = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "minimum_capacity":
				fieldsToUpdate = append(fieldsToUpdate, " minimum_capacity = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "maximum_capacity":
				fieldsToUpdate = append(fieldsToUpdate, " maximum_capacity = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "warehouse_id":
				fieldsToUpdate = append(fieldsToUpdate, " warehouse_id = ?")
				valuesToUse = append(valuesToUse, int(requestData[key].(float64)))
			case "product_type_id":
				fieldsToUpdate = append(fieldsToUpdate, " product_type_id = ?")
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
