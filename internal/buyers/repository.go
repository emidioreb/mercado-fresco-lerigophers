package buyers

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedBuyer = errors.New("ocurred an error while updating the buyer")
	errCreateBuyer  = errors.New("ocurred an error to create buyer")
	errGetBuyers    = errors.New("couldn't get buyers")
	errGetOneBuyer  = errors.New("unexpected error to get buyer")
	errDeleteBuyer  = errors.New("unexpected error to delete buyer")
)

type Repository interface {
	Create(cardNumberId string, firstName, lastName string) (Buyer, error)
	GetOne(id int) (Buyer, error)
	GetAll() ([]Buyer, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Buyer, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(cardNumberId, firstName, lastName string) (Buyer, error) {
	insert := `INSERT INTO buyers(card_number_id, first_name, last_name) VALUES(?, ?, ?);`

	newBuyer := Buyer{
		CardNumberId: cardNumberId,
		FirstName:    firstName,
		LastName:     lastName,
	}

	result, err := mariaDb.db.Exec(
		insert,
		cardNumberId,
		firstName,
		lastName,
	)

	if err != nil {
		return Buyer{}, errCreateBuyer
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Buyer{}, errCreateBuyer
	}

	newBuyer.Id = int(lastId)

	return newBuyer, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Buyer, error) {
	getOne := `SELECT * FROM buyers WHERE id = ?`
	currentBuyer := Buyer{}

	row := mariaDb.db.QueryRow(getOne, id)
	err := row.Scan(
		&currentBuyer.Id,
		&currentBuyer.CardNumberId,
		&currentBuyer.FirstName,
		&currentBuyer.LastName,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Buyer{}, fmt.Errorf("buyer with id %d not found", id)
	}

	if err != nil {
		return Buyer{}, errGetOneBuyer
	}

	return currentBuyer, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Buyer, error) {
	query := `SELECT * FROM buyers`
	buyers := []Buyer{}

	rows, err := mariaDb.db.Query(query)
	if err != nil {
		return []Buyer{}, errGetBuyers
	}

	for rows.Next() {
		var currentBuyer Buyer
		if err := rows.Scan(
			&currentBuyer.Id,
			&currentBuyer.CardNumberId,
			&currentBuyer.FirstName,
			&currentBuyer.LastName,
		); err != nil {
			return []Buyer{}, errGetBuyers
		}
		buyers = append(buyers, currentBuyer)
	}
	return buyers, nil
}
func (mariaDb mariaDbRepository) Delete(id int) error {
	delete := "DELETE FROM buyers WHERE id = ?"
	result, err := mariaDb.db.Exec(delete, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 {
		return fmt.Errorf("buyer with id %d not found", id)
	}

	if err != nil {
		return errDeleteBuyer
	}

	return nil
}
func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Buyer, error) {
	prefixQuery := "UPDATE buyers SET"
	fieldsToUpdate := []string{}
	valuesToUse := []interface{}{}
	whereCase := "WHERE id = ?"
	var finalQuery string

	for key, _ := range requestData {
		switch key {
		case "card_number_id":
			fieldsToUpdate = append(fieldsToUpdate, " card_number_id = ?")
			valuesToUse = append(valuesToUse, requestData[key])
		case "first_name":
			fieldsToUpdate = append(fieldsToUpdate, " first_name = ?")
			valuesToUse = append(valuesToUse, requestData[key])
		case "last_name":
			fieldsToUpdate = append(fieldsToUpdate, " last_name = ?")
			valuesToUse = append(valuesToUse, requestData[key])
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
		return Buyer{}, errUpdatedBuyer
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 && err != nil {
		return Buyer{}, errUpdatedBuyer
	}

	currentBuyer, err := mariaDb.GetOne(id)
	if err != nil {
		return Buyer{}, errUpdatedBuyer
	}

	return currentBuyer, nil
}
