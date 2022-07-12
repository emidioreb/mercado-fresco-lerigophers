package sellers

import "fmt"

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

		var fields = []string{"company_name", "address", "telephone", "locality_id", "cid"}
		for _, currField := range fields {
			if _, ok := requestData[currField]; ok {
				fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf(" %s = ?", currField))
				if currField == "cid" {
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

	queryFindByCID = "SELECT id, cid FROM sellers WHERE cid = ?"
)
