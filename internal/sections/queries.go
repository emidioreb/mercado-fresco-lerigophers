package sections

import "fmt"

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

		var fields = []string{"section_number", "current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "product_type_id"}
		for _, currField := range fields {
			if _, ok := requestData[currField]; ok {
				fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf(" %s = ?", currField))
				valuesToUse = append(valuesToUse, int(requestData[currField].(float64)))
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
