package locality

import "database/sql"

type Repository interface {
	// POST /api/v1/localities
	CreateLocality(
		id,
		localityName,
		provinceName,
		countryName string,
	) (Locality, error)

	// GET /api/v1/localities/reportSellers?id=6700 <- ID é a ID da localidade
	GetReportSellers(localityId []string) []ReportSellers
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) CreateLocality(id, localityName, provinceName, countryName string) (Locality, error) {
	insert := `INSERT INTO localities (id, locality_name, province_name, conuntry_name) VALUES (?, ?, ?, ?)`
	_, err := mariaDb.db.Exec(insert, id, localityName, provinceName, countryName)
	if err != nil {
		return Locality{}, err
	}

	return Locality{}, nil
}

func (mariaDb mariaDbRepository) GetReportSellers(localityId []string) []ReportSellers {
	return []ReportSellers{}
}
