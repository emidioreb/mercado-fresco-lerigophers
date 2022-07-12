package localities

var (
	queryGetReportAll = `SELECT l.id as locality_id, l.locality_name, count(s.id) as sellers_count
	FROM sellers s
	RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id, l.locality_name;`

	queryGetReportOne = `SELECT l.id as locality_id, l.locality_name, count(s.id) as sellers_count
	FROM sellers s
	RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id, l.locality_name;`

	queryCreateLocality = `INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)`
	queryGetOneLocality = `SELECT * FROM localities WHERE id = ?`

	queryGetReportCarriersAll = `SELECT l.id as locality_id, l.locality_name, count(*) as carriers_count 
	FROM localities as l
	INNER JOIN carriers as c
	ON l.id = c.locality_id 
	GROUP BY l.id, l.locality_name;`

	queryGetReportCarriersOne = `SELECT l.id as locality_id, l.locality_name, count(*) as carriers_count
	FROM carriers c
	JOIN localities l ON c.locality_id = l.id
	WHERE l.id = ?
	GROUP BY l.id, l.locality_name;`
)
