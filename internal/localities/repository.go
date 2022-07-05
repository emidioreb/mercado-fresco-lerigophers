package localities

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	CreateLocality(id, localityName, provinceName, countryName string) (Locality, error)
	GetReportSellers(localityId string) ([]ReportSellers, error)
	GetOne(id string) (Locality, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) GetOne(id string) (Locality, error) {
	currentLocality := Locality{}

	row := mariaDb.db.QueryRow(queryGetOneLocality, id)
	err := row.Scan(
		&currentLocality.Id,
		&currentLocality.LocalityName,
		&currentLocality.ProvinceName,
		&currentLocality.CountryName,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Locality{}, fmt.Errorf("locality with id %s not found", id)
	}

	if err != nil {
		return Locality{}, errors.New("error to find locality")
	}

	return currentLocality, nil
}

func (mariaDb mariaDbRepository) CreateLocality(
	id,
	localityName,
	provinceName,
	countryName string,
) (Locality, error) {

	_, err := mariaDb.db.Exec(
		queryCreateLocality,
		id,
		localityName,
		provinceName,
		countryName,
	)

	if err != nil {
		return Locality{}, errors.New("couldn't create a locality")
	}

	newLocality := Locality{
		Id:           id,
		LocalityName: localityName,
		ProvinceName: provinceName,
		CountryName:  countryName,
	}

	return newLocality, nil
}

func (mariaDb mariaDbRepository) GetReportSellers(localityId string) ([]ReportSellers, error) {
	reports := []ReportSellers{}

	var (
		rows *sql.Rows
		err  error
	)

	if localityId != "" {
		rows, err = mariaDb.db.Query(queryGetReportOne, localityId)
	} else {
		rows, err = mariaDb.db.Query(queryGetReportAll)
	}

	if err != nil {
		return []ReportSellers{}, errors.New("error to report sellers by locality")
	}

	for rows.Next() {
		var currentReport ReportSellers
		if err := rows.Scan(
			&currentReport.LocalityId,
			&currentReport.LocalityName,
			&currentReport.SellersCount,
		); err != nil {
			return []ReportSellers{}, errors.New("error to report sellers by locality")
		}
		reports = append(reports, currentReport)
	}

	return reports, nil
}
