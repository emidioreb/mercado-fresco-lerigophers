package sellers

var (
	queryCreateSeller  = "INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
	queryGetOneSeller  = "SELECT * FROM sellers WHERE id = ?"
	queryGetAllSellers = "SELECT * FROM sellers"
	queryDeleteSeller  = "DELETE FROM sellers WHERE id = ?"
	queryUpdateSeller  = func(
		requestData map[string]interface{},
		id int) (
		finalQuery string,
		valuesToUse []interface{}) {
		prefixQuery := "UPDATE sellers SET"
		fieldsToUpdate := []string{}
		whereCase := "WHERE id = ?"

		for key, _ := range requestData {
			switch key {
			case "company_name":
				fieldsToUpdate = append(fieldsToUpdate, " company_name = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "address":
				fieldsToUpdate = append(fieldsToUpdate, " address = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "telephone":
				fieldsToUpdate = append(fieldsToUpdate, " telephone = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "locality_id":
				fieldsToUpdate = append(fieldsToUpdate, " locality_id = ?")
				valuesToUse = append(valuesToUse, requestData[key])
			case "cid":
				fieldsToUpdate = append(fieldsToUpdate, " cid = ?")
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

	queryFindByCID = "SELECT id, cid FROM sellers WHERE cid = ?"
)
