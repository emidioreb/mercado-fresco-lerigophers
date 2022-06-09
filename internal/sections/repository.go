package sections

import "fmt"

var sections = []Section{}
var globalID = 1

type Repository interface {
	Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error)
	GetOne(id int) (Section, error)
	GetAll() ([]Section, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Section, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error) {
	newSection := Section{
		Id:                 globalID,
		SectionNumber:      sectionNumber,
		CurrentTemperature: currentTemperature,
		MinimumTemperature: minimumTemperature,
		CurrentCapacity:    currentCapacity,
		MininumCapacity:    mininumCapacity,
		MaximumCapacity:    maximumCapacity,
		WarehouseId:        warehouseId,
		ProductTypeId:      productTypeId,
	}

	sections = append(sections, newSection)
	globalID++

	return newSection, nil
}

func (repository) GetOne(id int) (Section, error) {
	for _, section := range sections {
		if section.Id == id {
			return section, nil
		}
	}

	return Section{}, fmt.Errorf("section with id %d not found", id)
}
func (repository) GetAll() ([]Section, error) {
	return sections, nil
}
func (repository) Delete(id int) error {
	for i, section := range sections {
		if section.Id == id {
			sections = append(sections[:i], sections[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("section with id %d not found", id)
}
func (repository) Update(id int, requestData map[string]interface{}) (Section, error) {
	var s *Section

	for i, section := range sections {
		if section.Id == id {
			s = &sections[i]

			for key, value := range requestData {
				switch key {
				case "section_number":
					s.SectionNumber = int(value.(float64))
				case "current_temperature":
					s.CurrentTemperature = int(value.(float64))
				case "minimum_temperature":
					s.MinimumTemperature = int(value.(float64))
				case "current_capacity":
					s.CurrentCapacity = int(value.(float64))
				case "minimum_capacity":
					s.MininumCapacity = int(value.(float64))
				case "maximum_capacity":
					s.MaximumCapacity = int(value.(float64))
				case "warehouse_id":
					s.WarehouseId = int(value.(float64))
				case "product_type_id":
					s.ProductTypeId = int(value.(float64))
				}
			}
			return *s, nil
		}

	}

	return Section{}, fmt.Errorf("section with id %d not found", id)
}
