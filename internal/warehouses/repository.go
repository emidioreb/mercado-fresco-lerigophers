package warehouses

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedWarehouse = errors.New("ocurred an error while updating the warehouse")
	errCreateWarehouse  = errors.New("ocurred an error to create warehouse")
	errGetWarehouses    = errors.New("couldn't get warehouses")
	errGetOneWarehouse  = errors.New("unexpected error to get warehouse")
	errDeleteWarehouse  = errors.New("unexpected error to delete warehouse")
)

type Repository interface {
	Create(warehouseCode, adress, telephone string, minimumCapacity, minimumTemperature int) (Warehouse, error)
	GetOne(id int) (Warehouse, error)
	GetAll() ([]Warehouse, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Warehouse, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(warehouseCode, adress, telephone string, minimumCapacity, minimumTemperature int) (Warehouse, error) {
	insert := `INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature) VALUES (?, ?, ?, ?, ?)`

	newWarehouse := Warehouse{
		WarehouseCode:      warehouseCode,
		Address:            adress,
		Telephone:          telephone,
		MinimumCapacity:    minimumCapacity,
		MinimumTemperature: minimumTemperature,
	}

	result, err := mariaDb.db.Exec(
		insert,
		warehouseCode,
		adress,
		telephone,
		minimumCapacity,
		minimumTemperature,
	)

	if err != nil {
		return Warehouse{}, errCreateWarehouse
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Warehouse{}, errCreateWarehouse
	}

	newWarehouse.Id = int(lastId)

	return newWarehouse, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Warehouse, error) {
	getOne := "SELECT * FROM warehouses WHERE id = ?"
	currentWarehouse := Warehouse{}

	row := mariaDb.db.QueryRow(getOne, id)

	err := row.Scan(
		&currentWarehouse.Id,
		&currentWarehouse.WarehouseCode,
		&currentWarehouse.Address,
		&currentWarehouse.Telephone,
		&currentWarehouse.MinimumCapacity,
		&currentWarehouse.MinimumTemperature,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Warehouse{}, fmt.Errorf("warehouse with id %d not found", id)
	}

	if err != nil {
		return Warehouse{}, errGetOneWarehouse
	}

	return currentWarehouse, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Warehouse, error) {
	query := `SELECT * FROM warehouses`
	warehouses := []Warehouse{}

	rows, err := mariaDb.db.Query(query)
	if err != nil {
		return []Warehouse{}, errGetWarehouses
	}
	for rows.Next() {
		var currentWarehouse Warehouse
		if err := rows.Scan(
			&currentWarehouse.Id,
			&currentWarehouse.WarehouseCode,
			&currentWarehouse.Address,
			&currentWarehouse.Telephone,
			&currentWarehouse.MinimumCapacity,
			&currentWarehouse.MinimumTemperature,
		); err != nil {
			return []Warehouse{}, errGetWarehouses
		}
		warehouses = append(warehouses, currentWarehouse)

	}
	return warehouses, nil

}

func (mariaDb mariaDbRepository) Delete(id int) error {
	delete := "DELETE FROM warehouses WHERE id = ?"
	result, err := mariaDb.db.Exec(delete, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 {
		return fmt.Errorf("warehouse with id %d not found", id)
	}

	if err != nil {
		return errDeleteWarehouse
	}

	return nil
}

func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Warehouse, error) {
	prefixQuery := "UPDATE warehouses SET"
	fieldsToUpdate := []string{}
	valuesToUse := []interface{}{}
	whereCase := "WHERE id = ?"
	var finalQuery string

	for key := range requestData {
		switch key {
		case "warehouse_code":
			fieldsToUpdate = append(fieldsToUpdate, " warehouse_code = ?")
			valuesToUse = append(valuesToUse, requestData[key])
		case "address":
			fieldsToUpdate = append(fieldsToUpdate, " address = ?")
			valuesToUse = append(valuesToUse, requestData[key])
		case "telephone":
			fieldsToUpdate = append(fieldsToUpdate, " telephone = ?")
			valuesToUse = append(valuesToUse, requestData[key])
		case "minimum_capacity":
			fieldsToUpdate = append(fieldsToUpdate, " minimum_capacity = ?")
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
		return Warehouse{}, errUpdatedWarehouse
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 && err != nil {
		return Warehouse{}, errUpdatedWarehouse
	}

	currentWarehouse, err := mariaDb.GetOne(id)
	if err != nil {
		return Warehouse{}, errUpdatedWarehouse
	}

	return currentWarehouse, nil
}
