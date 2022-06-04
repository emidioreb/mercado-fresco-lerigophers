package sections

import "fmt"

var sections = []Section{}
var globalID = 1

type Repository interface {
	Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error)
	GetOne(id int) (Section, error)
	GetAll() ([]Section, error)
	Delete(id int) error
	Update(id, sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error)
	UpdateCurrCapacity(id int, currentCapacity int) (Section, error)
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
func (repository) Update(id, sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, error) {
	updatedSection := Section{id, sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId}
	for i, section := range sections {
		if section.Id == id {
			sections[i] = updatedSection
			return sections[i], nil
		}
	}
	return Section{}, fmt.Errorf("can't find section with id %d", id)
}

func (repository) UpdateCurrCapacity(id int, currentCapacity int) (Section, error) {
	for i, section := range sections {
		if section.Id == id {
			sections[i].CurrentCapacity = currentCapacity
			return sections[i], nil
		}
	}

	return Section{}, fmt.Errorf("can't find section with id %d", id)
}
