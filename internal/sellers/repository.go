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
	Create(cid int, companyName, address, telephone string) (Seller, error)
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

func (mariaDb mariaDbRepository) Create(cid int, companyName, address, telephone string) (Seller, error) {
	insert := `INSERT INTO sellers (cid, company_name, address, telephone) VALUES (?, ?, ?, ?)`

	newSeller := Seller{
		Cid:         cid,
		CompanyName: companyName,
		Address:     address,
		Telephone:   telephone,
	}

	result, err := mariaDb.db.Exec(
		insert,
		cid,
		companyName,
		address,
		telephone,
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
	getOne := `SELECT * FROM sellers WHERE id = ?`
	currentSeller := Seller{}

	row := mariaDb.db.QueryRow(getOne, id)
	err := row.Scan(
		&currentSeller.Id,
		&currentSeller.Cid,
		&currentSeller.CompanyName,
		&currentSeller.Address,
		&currentSeller.Telephone,
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
	query := `SELECT * FROM sellers`
	sellers := []Seller{}

	rows, err := mariaDb.db.Query(query)
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
		); err != nil {
			return []Seller{}, errGetSellers
		}
		sellers = append(sellers, currentSeller)
	}
	return sellers, nil
}

func (mariaDb mariaDbRepository) Delete(id int) error {
	delete := "DELETE FROM sellers WHERE id = ?"
	result, err := mariaDb.db.Exec(delete, id)
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
	prefixQuery := "UPDATE sellers SET"
	fieldsToUpdate := []string{}
	valuesToUse := []interface{}{}
	whereCase := "WHERE id = ?"
	var finalQuery string

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
