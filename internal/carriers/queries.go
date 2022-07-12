package carriers

var (
	queryCreateCarry  = "INSERT INTO carriers (cid, company_name, address, telephone,locality_id) VALUES (?, ?, ?, ?,?)"
	queryGetOneCarry  = "SELECT * FROM carriers WHERE cid=?"
)