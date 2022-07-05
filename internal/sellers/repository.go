package sellers

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedSeller = errors.New("ocurred an error while updating the seller")
	errCreateSeller  = errors.New("ocurred an error to create seller")
	errGetSellers    = errors.New("couldn't get sellers")
	errGetOneSeller  = errors.New("unexpected error to get seller")
	errDeleteSeller  = errors.New("unexpected error to delete seller")
)

type Repository interface {
	Create(cid int, companyName, address, telephone, localityId string) (Seller, error)
	GetOne(id int) (Seller, error)
	GetAll() ([]Seller, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Seller, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(cid int, companyName, address, telephone, localityId string) (Seller, error) {
	newSeller := Seller{
		Cid:         cid,
		CompanyName: companyName,
		Address:     address,
		Telephone:   telephone,
		LocalityId:  localityId,
	}

	result, err := mariaDb.db.Exec(
		queryCreateSeller,
		cid,
		companyName,
		address,
		telephone,
		localityId,
	)

	if err != nil {
		return Seller{}, errCreateSeller
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Seller{}, errCreateSeller
	}

	newSeller.Id = int(lastId)

	return newSeller, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Seller, error) {
	currentSeller := Seller{}

	row := mariaDb.db.QueryRow(queryGetOneSeller, id)
	err := row.Scan(
		&currentSeller.Id,
		&currentSeller.Cid,
		&currentSeller.CompanyName,
		&currentSeller.Address,
		&currentSeller.Telephone,
		&currentSeller.LocalityId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Seller{}, fmt.Errorf("seller with id %d not found", id)
	}

	if err != nil {
		return Seller{}, errGetOneSeller
	}

	return currentSeller, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Seller, error) {
	sellers := []Seller{}

	rows, err := mariaDb.db.Query(queryGetAllSellers)
	if err != nil {
		return []Seller{}, errGetSellers
	}

	for rows.Next() {
		var currentSeller Seller
		if err := rows.Scan(
			&currentSeller.Id,
			&currentSeller.Cid,
			&currentSeller.CompanyName,
			&currentSeller.Address,
			&currentSeller.Telephone,
			&currentSeller.LocalityId,
		); err != nil {
			return []Seller{}, errGetSellers
		}
		sellers = append(sellers, currentSeller)
	}
	return sellers, nil
}

func (mariaDb mariaDbRepository) Delete(id int) error {
	result, err := mariaDb.db.Exec(queryDeleteSeller, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 {
		return fmt.Errorf("seller with id %d not found", id)
	}

	if err != nil {
		return errDeleteSeller
	}

	return nil
}
func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Seller, error) {
	finalQuery, valuesToUse := queryUpdateSeller(requestData, id)

	result, err := mariaDb.db.Exec(finalQuery, valuesToUse...)
	if err != nil {
		return Seller{}, errUpdatedSeller
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 && err != nil {
		return Seller{}, errUpdatedSeller
	}

	currentSeller, err := mariaDb.GetOne(id)
	if err != nil {
		return Seller{}, errUpdatedSeller
	}

	return currentSeller, nil
}
