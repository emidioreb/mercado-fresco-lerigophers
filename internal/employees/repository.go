package employees

import "fmt"

var Employees = []Employee{}
var globalID = 1

type Repository interface {
	Create(cardNumber, firstName, lastName string, warehouseId int) (Employee, error)
	GetOne(id int) (Employee, error)
	GetAll() ([]Employee, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Employee, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(cardNumber, firstName, lastName string, warehouseId int) (Employee, error) {
	newEmployee := Employee{
		Id:           globalID,
		CardNumberId: cardNumber,
		FirstName:    firstName,
		LastName:     lastName,
		WarehouseId:  warehouseId,
	}

	Employees = append(Employees, newEmployee)
	globalID++

	return newEmployee, nil
}

func (repository) GetOne(id int) (Employee, error) {
	for _, Employee := range Employees {
		if Employee.Id == id {
			return Employee, nil
		}
	}

	return Employee{}, fmt.Errorf("can't find Employee with id %d", id)
}
func (repository) GetAll() ([]Employee, error) {
	return Employees, nil
}
func (repository) Delete(id int) error {
	for i, Employee := range Employees {
		if Employee.Id == id {
			Employees = append(Employees[:i], Employees[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("can't find Employee with id %d", id)
}
func (repository) Update(id int, requestData map[string]interface{}) (Employee, error) {
	var e *Employee

	for i, employee := range Employees {
		if employee.Id == id {
			e = &Employees[i]

			for key, value := range requestData {

				switch key {
				case "card_number_id":
					e.CardNumberId = value.(string)
				case "first_name":
					e.FirstName = value.(string)
				case "last_name":
					e.LastName = value.(string)
				case "warehouse_id":
					e.WarehouseId = int(value.(float64))
				}
			}
			return *e, nil
		}

	}
	return Employee{}, fmt.Errorf("can't find Employee with id %d", id)
}
