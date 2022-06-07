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

	return Section{}, fmt.Errorf("can't find section with id %d", id)
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
	return fmt.Errorf("can't find section with id %d", id)
}
func (repository) Update(id int, requestData map[string]interface{}) (Section, error) {
	var s *Section

	for i, section := range sections {
		if section.Id == id {
			s = &sections[i]

			for key, value := range requestData {
				valueParsed := int(value.(float64))

				switch key {
				case "section_number":
					s.SectionNumber = valueParsed
				case "current_temperature":
					s.CurrentTemperature = valueParsed
				case "minimum_temperature":
					s.MinimumTemperature = valueParsed
				case "current_capacity":
					s.CurrentCapacity = valueParsed
				case "minimum_capacity":
					s.MininumCapacity = valueParsed
				case "maximum_capacity":
					s.MaximumCapacity = valueParsed
				case "warehouse_id":
					s.WarehouseId = valueParsed
				case "product_type_id":
					s.ProductTypeId = valueParsed
				}
			}
			return *s, nil
		}

	}

	return Section{}, fmt.Errorf("can't find section with id %d", id)
}
