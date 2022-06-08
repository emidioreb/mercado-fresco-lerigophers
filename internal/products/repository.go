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
	Update(id int, requestData map[string]interface{}) (Product, error)
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
func (repository) Update(id int, requestData map[string]interface{}) (Product, error) {
	var p *Product

	for i, product := range Products {
		if product.Id == id {
			p = &Products[i]

			for key, value := range requestData {

				switch key {
				case "product_code":
					p.ProductCode = value.(string)
				case "description":
					p.Description = value.(string)
				case "width":
					p.Width = value.(float64)
				case "heigth":
					p.Height = value.(float64)
				case "length":
					p.Length = value.(float64)
				case "net_weight":
					p.NetWeight = value.(float64)
				case "expiration_rate":
					p.ExpirationRate = value.(float64)
				case "recommended_freezing_temperature":
					p.RecommendedFreezingTemperature = value.(float64)
				case "freezing_rate":
					p.FreezingRate = value.(float64)
				case "product_type_id":
					p.ProductTypeId = int(value.(float64))
				}
			}
			return *p, nil
		}
	}

	return Product{}, fmt.Errorf("can't find product with id %d", id)

}
