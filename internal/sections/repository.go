package sections

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedSection             = errors.New("ocurred an error while updating the section")
	errCreateSection              = errors.New("ocurred an error to create section")
	errGetSections                = errors.New("couldn't get sections")
	errGetOneSection              = errors.New("unexpected error to get section")
	errDeleteSection              = errors.New("unexpected error to delete section")
	errVerifySectionNumber        = errors.New("failed to verify if section_number already exists")
	errSectionNumberAlreadyExists = errors.New("section number already exists")
)

func GetErrSectionNotFound(id int) error {
	return fmt.Errorf("section with id %d not found", id)
}

type Repository interface {
	Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error)
	GetOne(id int) (Section, error)
	GetAll() ([]Section, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Section, error)
	GetBySectionNumber(sectionNumber int) (int, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error) {
	newSection := Section{
		SectionNumber:      sectionNumber,
		CurrentTemperature: currentTemperature,
		MinimumTemperature: minimumTemperature,
		CurrentCapacity:    currentCapacity,
		MininumCapacity:    mininumCapacity,
		MaximumCapacity:    maximumCapacity,
		WarehouseId:        warehouseId,
		ProductTypeId:      productTypeId,
	}

	result, err := mariaDb.db.Exec(
		queryCreateSection,
		sectionNumber,
		currentTemperature,
		minimumTemperature,
		currentCapacity,
		mininumCapacity,
		maximumCapacity,
		warehouseId,
		productTypeId,
	)

	if err != nil {
		return Section{}, errCreateSection
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Section{}, errCreateSection
	}

	newSection.Id = int(lastId)

	return newSection, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Section, error) {
	currentSection := Section{}

	row := mariaDb.db.QueryRow(queryGetOneSection, id)
	err := row.Scan(
		&currentSection.Id,
		&currentSection.SectionNumber,
		&currentSection.CurrentTemperature,
		&currentSection.MinimumTemperature,
		&currentSection.CurrentCapacity,
		&currentSection.MininumCapacity,
		&currentSection.MaximumCapacity,
		&currentSection.WarehouseId,
		&currentSection.ProductTypeId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Section{}, GetErrSectionNotFound(id)
	}

	if err != nil {
		return Section{}, errGetOneSection
	}

	return currentSection, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Section, error) {
	sections := []Section{}

	rows, err := mariaDb.db.Query(queryGetAllSections)
	if err != nil {
		return []Section{}, errGetSections
	}

	for rows.Next() {
		var currentSection Section
		if err := rows.Scan(
			&currentSection.Id,
			&currentSection.SectionNumber,
			&currentSection.CurrentTemperature,
			&currentSection.MinimumTemperature,
			&currentSection.CurrentCapacity,
			&currentSection.MininumCapacity,
			&currentSection.MaximumCapacity,
			&currentSection.WarehouseId,
			&currentSection.ProductTypeId,
		); err != nil {
			return []Section{}, errGetSections
		}
		sections = append(sections, currentSection)
	}

	return sections, nil

}
func (mariaDb mariaDbRepository) Delete(id int) error {
	result, err := mariaDb.db.Exec(queryDeleteSection, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 {
		return GetErrSectionNotFound(id)
	}

	if err != nil {
		return errDeleteSection
	}

	return nil
}

func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Section, error) {
	finalQuery, valuesToUse := queryUpdateSection(requestData, id)

	result, err := mariaDb.db.Exec(finalQuery, valuesToUse...)
	if err != nil {
		return Section{}, errUpdatedSection
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 && err != nil {
		return Section{}, errUpdatedSection
	}

	currentSection, err := mariaDb.GetOne(id)
	if err != nil {
		return Section{}, errUpdatedSection
	}

	return currentSection, nil
}

func (mariaDb mariaDbRepository) GetBySectionNumber(sectionNumber int) (int, error) {
	var selectedSectionNumber, selectedId int

	row := mariaDb.db.QueryRow(queryValidSectionNumber, sectionNumber)
	err := row.Scan(&selectedId, &selectedSectionNumber)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, errVerifySectionNumber
	}

	return selectedId, errSectionNumberAlreadyExists
}
