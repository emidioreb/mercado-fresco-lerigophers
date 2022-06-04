package warehouses

import "fmt"

var warehouses = []Warehouse{}
var globalID = 1

type Repository interface {
	Create(warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, error)
	GetOne(id int) (Warehouse, error)
	GetAll() ([]Warehouse, error)
	Delete(id int) error
	Update(id int, warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, error)
	UpdateTelephone(id int, telephone string) (Warehouse, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, error) {
	newWarehouse := Warehouse{
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

	return Warehouse{}, fmt.Errorf("can't find warehouse with id %d", id)
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
	return fmt.Errorf("can't find warehouse with id %d", id)
}
func (repository) Update(id int, warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, error) {
	updatedWarehouse := Warehouse{id, warehouseCode, adress, telephone, minimumCapacity, maxmumCapacity}
	for i, warehouse := range warehouses {
		if warehouse.Id == id {
			warehouses[i] = updatedWarehouse
			return warehouses[i], nil
		}
	}
	return Warehouse{}, fmt.Errorf("can't find warehouse with id %d", id)
}
func (repository) UpdateTelephone(id int, telephone string) (Warehouse, error) {
	for i, warehouse := range warehouses {
		if warehouse.Id == id {
			warehouses[i].Telephone = telephone
			return warehouses[i], nil
		}
	}

	return Warehouse{}, fmt.Errorf("can't find warehouse with id %d", id)
}
