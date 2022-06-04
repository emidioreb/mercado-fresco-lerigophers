package sellers

import "fmt"

var sellers = []Seller{}
var globalID = 1

type Repository interface {
	Create(cid int, companyName, address, telephone string) (Seller, error)
	GetOne(id int) (Seller, error)
	GetAll() ([]Seller, error)
	Delete(id int) error
	Update(id, cid int, companyName, address, telephone string) (Seller, error)
	UpdateAddress(id int, address string) (Seller, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(cid int, companyName, address, telephone string) (Seller, error) {
	newSeller := Seller{
		Id:          globalID,
		Cid:         cid,
		CompanyName: companyName,
		Address:     address,
		Telephone:   telephone,
	}

	sellers = append(sellers, newSeller)
	globalID++

	return newSeller, nil
}

func (repository) GetOne(id int) (Seller, error) {
	for _, seller := range sellers {
		if seller.Id == id {
			return seller, nil
		}
	}

	return Seller{}, fmt.Errorf("can't find seller with id %d", id)
}
func (repository) GetAll() ([]Seller, error) {
	return sellers, nil
}
func (repository) Delete(id int) error {
	for i, seller := range sellers {
		if seller.Id == id {
			sellers = append(sellers[:i], sellers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("can't find seller with id %d", id)
}
func (repository) Update(id, cid int, companyName, address, telephone string) (Seller, error) {
	updatedSeller := Seller{id, cid, companyName, address, telephone}
	for i, seller := range sellers {
		if seller.Id == id {
			sellers[i] = updatedSeller
			return sellers[i], nil
		}
	}
	return Seller{}, fmt.Errorf("can't find seller with id %d", id)
}
func (repository) UpdateAddress(id int, address string) (Seller, error) {
	for i, seller := range sellers {
		if seller.Id == id {
			sellers[i].Address = address
			return sellers[i], nil
		}
	}

	return Seller{}, fmt.Errorf("can't find seller with id %d", id)
}
