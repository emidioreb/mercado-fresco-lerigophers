package sellers

import "fmt"

var sellers = []Seller{}
var globalID = 1

type Repository interface {
	Create(cid int, companyName, address, telephone string) (Seller, error)
	GetOne(id int) (Seller, error)
	GetAll() ([]Seller, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Seller, error)
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

	return Seller{}, fmt.Errorf("seller with id %d not found", id)
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
	return fmt.Errorf("seller with id %d not found", id)
}
func (repository) Update(id int, requestData map[string]interface{}) (Seller, error) {
	var s *Seller

	for i, seller := range sellers {
		if seller.Id == id {
			s = &sellers[i]

			for key, _ := range requestData {
				valueString, _ := requestData[key].(string)
				switch key {
				case "company_name":
					s.CompanyName = valueString
				case "address":
					s.Address = valueString

				case "telephone":
					s.Telephone = valueString
				case "cid":
					s.Cid = int(requestData[key].(float64))
				}
			}
			return *s, nil
		}
	}
	return Seller{}, fmt.Errorf("seller with id %d not found", id)
}
