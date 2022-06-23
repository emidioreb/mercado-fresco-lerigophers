package warehouses

import (
	"fmt"
)

var warehouses = []Warehouse{}
var globalID = 1

type Repository interface {
	Create(warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, error)
	GetOne(id int) (Warehouse, error)
	GetAll() ([]Warehouse, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Warehouse, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, error) {

	newWarehouse := Warehouse{
		Id:                 globalID,
		WarehouseCode:      warehouseCode,
		Address:            adress,
		Telephone:          telephone,
		MinimumCapacity:    minimumCapacity,
		MaximumTemperature: maxmumCapacity,
	}

	warehouses = append(warehouses, newWarehouse)
	globalID++

	return newWarehouse, nil
}

func (repository) GetOne(id int) (Warehouse, error) {
	for _, warehouse := range warehouses {
		if warehouse.Id == id {
			return warehouse, nil
		}
	}

	return Warehouse{}, fmt.Errorf("warehouse with id %d not found", id)
}
func (repository) GetAll() ([]Warehouse, error) {
	return warehouses, nil
}
func (repository) Delete(id int) error {
	for i, warehouse := range warehouses {
		if warehouse.Id == id {
			warehouses = append(warehouses[:i], warehouses[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("warehouse with id %d not found", id)
}

func (repository) Update(id int, requestData map[string]interface{}) (Warehouse, error) {
	var w *Warehouse

	for i, warehouse := range warehouses {
		if warehouse.Id == id {
			w = &warehouses[i]

			for key, value := range requestData {
				switch key {
				case "warehouse_code":
					w.WarehouseCode = value.(string)
				case "adress":
					w.Address = value.(string)
				case "telephone":
					w.Telephone = value.(string)
				case "minimum_capacity":
					w.MinimumCapacity = int(value.(float64))
				case "maximum_temperature":
					w.MaximumTemperature = int(value.(float64))
				}
			}
			return *w, nil
		}

	}
	return Warehouse{}, fmt.Errorf("warehouse with id %d not found", id)
}
