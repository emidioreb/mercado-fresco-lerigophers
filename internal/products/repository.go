package products

import "fmt"

var Products = []Product{}
var globalID = 1

type Repository interface {
	Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate float64, productTypeId int) (Product, error)
	GetOne(id int) (Product, error)
	GetAll() ([]Product, error)
	Delete(id int) error
	Update(id int, productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate float64, productTypeId int) (Product, error)
	UpdateExpirationRate(id int, expirationRate float64) (Product, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (repository) Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
	freezingRate float64, productTypeId int) (Product, error) {
	newProduct := Product{
		Id:                             globalID,
		ProductCode:                    productCode,
		Description:                    description,
		Width:                          width,
		Height:                         height,
		Length:                         length,
		NetWeight:                      netWeight,
		ExpirationRate:                 expirationRate,
		RecommendedFreezingTemperature: recommendedFreezingTemperaturechan,
		FreezingRate:                   freezingRate,
		ProductTypeId:                  productTypeId,
	}

	Products = append(Products, newProduct)
	globalID++

	return newProduct, nil
}

func (repository) GetOne(id int) (Product, error) {
	for _, Product := range Products {
		if Product.Id == id {
			return Product, nil
		}
	}

	return Product{}, fmt.Errorf("can't find Product with id %d", id)
}
func (repository) GetAll() ([]Product, error) {
	return Products, nil
}
func (repository) Delete(id int) error {
	for i, Product := range Products {
		if Product.Id == id {
			Products = append(Products[:i], Products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("can't find Product with id %d", id)
}
func (repository) Update(id int, productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
	freezingRate float64, productTypeId int) (Product, error) {
	updatedProduct := Product{id, productCode, description, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate, productTypeId}
	for i, Product := range Products {
		if Product.Id == id {
			Products[i] = updatedProduct
			return Products[i], nil
		}
	}
	return Product{}, fmt.Errorf("can't find Product with id %d", id)
}
func (repository) UpdateExpirationRate(id int, expirationRate float64) (Product, error) {
	for i, Product := range Products {
		if Product.Id == id {
			Products[i].ExpirationRate = expirationRate
			return Products[i], nil
		}
	}

	return Product{}, fmt.Errorf("can't find Product with id %d", id)
}
