package buyers

import (
"fmt"
)
var Buyers = []Buyer{}
var globalID = 1

type Repository interface {
	Create(cardNumberId string, firstName, lastName string) (Buyer, error)
	GetOne(id int) (Buyer, error)
	GetAll() ([]Buyer, error)
	Delete(id int) error
	Update(id int, requestData map[string]string) (Buyer, error)
	// UpdateLastName(id int, lastName string) (Buyer, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(cardNumberId string, firstName, lastName string) (Buyer, error) {
	newBuyer := Buyer{
		Id:           globalID,
		CardNumberId: cardNumberId,
		FirstName:    firstName,
		LastName:     lastName,
	}

	Buyers = append(Buyers, newBuyer)
	globalID++

	return newBuyer, nil
}

func (repository) GetOne(id int) (Buyer, error) {
	for _, Buyer := range Buyers {
		if Buyer.Id == id {
			return Buyer, nil
		}
	}

	return Buyer{}, fmt.Errorf("can't find Buyer with id %d", id)
}
func (repository) GetAll() ([]Buyer, error) {
	return Buyers, nil
}
func (repository) Delete(id int) error {
	for i, Buyer := range Buyers {
		if Buyer.Id == id {
			Buyers = append(Buyers[:i], Buyers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("can't find Buyer with id %d", id)
}
func (repository) Update(id int, requestData map[string]string) (Buyer, error) {
	var s *Buyer

	for i, buyer := range Buyers {
		if buyer.Id == id {
			s = &Buyers[i]

			for key, value := range requestData {
				valueParsed := value

				switch key {
				case "card_number_id":
					s.CardNumberId = valueParsed
				case "first_name":
					s.FirstName = valueParsed
				case "last_name":
					s.LastName = valueParsed
				}
			}
			return *s, nil
		}

	}

	return Buyer{}, fmt.Errorf("can't find buyer with id %d", id)
}
func (repository) UpdateLastName(id int, lastName string) (Buyer, error) {
	for i, Buyer := range Buyers {
		if Buyer.Id == id {
			Buyers[i].LastName= lastName
			return Buyers[i], nil
		}
	}

	return Buyer{}, fmt.Errorf("can't find Buyer with id %d", id)
}