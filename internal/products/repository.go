package products

import (
	"database/sql"
	"errors"
	"fmt"
)

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

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
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

func (mariaDb mariaDbRepository) GetOne(id int) (Product, error) {
	query := `SELECT * FROM products WHERE ID = ?`
	row := mariaDb.db.QueryRow(query, id)

	product := Product{}
	err := row.Scan(&product.Id,
		&product.ProductCode,
		&product.Description,
		&product.Width,
		&product.Height,
		&product.Length,
		&product.NetWeight,
		&product.ExpirationRate,
		&product.RecommendedFreezingTemperature,
		&product.FreezingRate,
		&product.ProductTypeId,
		&product.SellerId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, fmt.Errorf("product with id %d not found", id)
	}

	if err != nil {
		fmt.Printf("erro aqui", err)
		return Product{}, errors.New("")
	}

	return product, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Product, error) {
	return Products, nil
}
func (mariaDb mariaDbRepository) Delete(id int) error {
	for i, Product := range Products {
		if Product.Id == id {
			Products = append(Products[:i], Products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product with id %d not found", id)
}
func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Product, error) {
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

	return Product{}, fmt.Errorf("product with id %d not found", id)

}
