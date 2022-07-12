package employees

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedEmployee = errors.New("ocurred an error while updating the employee")
	errCreateEmployee  = errors.New("ocurred an error to create employee")
	errGetEmployees    = errors.New("couldn't get employees")
	errGetOneEmployees = errors.New("unexpected error to get employee")
	errDeleteEmployee  = errors.New("unexpected error to delete employee")
)

type Repository interface {
	Create(cardNumber, firstName, lastName string, warehouseId int) (Employee, error)
	GetOne(id int) (Employee, error)
	GetAll() ([]Employee, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Employee, error)
	GetOneByCardNumber(id int, cardNumber string) error
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(cardNumber, firstName, lastName string, warehouseId int) (Employee, error) {

	newEmployee := Employee{
		CardNumberId: cardNumber,
		FirstName:    firstName,
		LastName:     lastName,
		WarehouseId:  warehouseId,
	}

	result, err := mariaDb.db.Exec(
		queryCreate,
		cardNumber,
		firstName,
		lastName,
		warehouseId,
	)

	if err != nil {
		return Employee{}, errCreateEmployee
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Employee{}, errCreateEmployee
	}

	newEmployee.Id = int(lastId)

	return newEmployee, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Employee, error) {

	currentEmployee := Employee{}

	row := mariaDb.db.QueryRow(queryGetOne, id)
	err := row.Scan(
		&currentEmployee.Id,
		&currentEmployee.CardNumberId,
		&currentEmployee.FirstName,
		&currentEmployee.LastName,
		&currentEmployee.WarehouseId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Employee{}, fmt.Errorf("employee with id %d not found", id)
	}

	if err != nil {
		return Employee{}, errGetOneEmployees
	}

	return currentEmployee, nil
}

func (mariaDb mariaDbRepository) GetOneByCardNumber(id int, cardNumber string) error {

	var err error
	currentEmployee := Employee{}

	if id != 0 {
		row := mariaDb.db.QueryRow(queryByCardNumberUpdate, cardNumber, id)
		err = row.Scan(
			&currentEmployee.Id,
		)
	} else {
		row := mariaDb.db.QueryRow(queryByCardNumberCreate, cardNumber)
		err = row.Scan(
			&currentEmployee.Id,
		)
	}

	if err != nil {
		return errors.New("unexpected error to get employee")
	}

	if currentEmployee.Id != 0 {
		return errors.New("card_number_id already exists")
	}

	return nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Employee, error) {

	employees := []Employee{}

	rows, err := mariaDb.db.Query(queryGetAll)
	if err != nil {
		return []Employee{}, errGetEmployees
	}

	for rows.Next() {
		var currentEmployee Employee
		if err := rows.Scan(
			&currentEmployee.Id,
			&currentEmployee.CardNumberId,
			&currentEmployee.FirstName,
			&currentEmployee.LastName,
			&currentEmployee.WarehouseId,
		); err != nil {
			return []Employee{}, errGetEmployees
		}
		employees = append(employees, currentEmployee)
	}
	return employees, nil
}

func (mariaDb mariaDbRepository) Delete(id int) error {

	result, err := mariaDb.db.Exec(queryDelete, id)
	if err != nil {
		return errDeleteEmployee
	}

	_, err = result.RowsAffected()
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("employee with id %d not found", id)
	}

	if err != nil {
		return errDeleteEmployee
	}

	return nil
}

func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Employee, error) {

	finalQuery, valuesToUse := queryUpdate(requestData, id)

	result, err := mariaDb.db.Exec(finalQuery, valuesToUse...)
	if err != nil {
		return Employee{}, errUpdatedEmployee
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 && err != nil {
		return Employee{}, errUpdatedEmployee
	}

	currentEmployee, err := mariaDb.GetOne(id)
	if err != nil {
		return Employee{}, errUpdatedEmployee
	}

	return currentEmployee, nil
}
