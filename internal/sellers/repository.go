package sellers

import (
	"database/sql"
	"errors"
	"fmt"
)

var sellers = []Seller{}
var globalID = 1

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
		return Seller{}, errors.New("ocurred an error to create seller")
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Seller{}, errors.New("ocurred an error to create seller")
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
		return Seller{}, errors.New("unexpected error")
	}

	return currentSeller, nil
}
func (mariaDb mariaDbRepository) GetAll() ([]Seller, error) {
	getOne := `SELECT * FROM sellers`
	sellers := []Seller{}

	rows, err := mariaDb.db.Query(getOne)
	if err != nil {
		return []Seller{}, errors.New("couldn't get sellers")
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
			return []Seller{}, errors.New("couldn't get sellers")
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
		return fmt.Errorf("error to find seller")
	}

	return nil
}
func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Seller, error) {
	var s *Seller

	for i, seller := range sellers {
		if seller.Id == id {
			s = &sellers[i]

			for key, _ := range requestData {
				valueString, _ := requestData[key].(string)
				switch key {
				case "company_name":
					s.CompanyName = valueString
				case "address":
					s.Address = valueString

				case "telephone":
					s.Telephone = valueString
				case "cid":
					s.Cid = int(requestData[key].(float64))
				}
			}
			return *s, nil
		}
	}
	return Seller{}, fmt.Errorf("seller with id %d not found", id)
}
