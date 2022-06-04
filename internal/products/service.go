package products

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate float64, productTypeId int) (Product, web.ResponseCode)
	GetOne(id int) (Product, web.ResponseCode)
	GetAll() ([]Product, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate float64, productTypeId int) (Product, web.ResponseCode)
	UpdateExpirationRate(id int, expiration_rate float64) (Product, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
	freezingRate float64, productTypeId int) (Product, web.ResponseCode) {
	allProducts, _ := s.GetAll()

	for _, product := range allProducts {
		if product.ProductCode == productCode {
			return Product{}, web.NewCodeResponse(http.StatusConflict, errors.New("Product_code already exists"))
		}
	}

	Product, _ := s.repository.Create(productCode, description, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate, productTypeId)

	return Product, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Product, web.ResponseCode) {
	product, err := s.repository.GetOne(id)

	if err != nil {
		return Product{}, web.NewCodeResponse(http.StatusNotFound, err)
	}
	return product, web.NewCodeResponse(http.StatusNotFound, nil)
}

func (s service) GetAll() ([]Product, web.ResponseCode) {
	products, err := s.repository.GetAll()
	return products, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan, freezingRate float64, productTypeId int) (Product, web.ResponseCode) {
	allProducts, _ := s.GetAll()

	for _, product := range allProducts {
		if product.ProductCode == productCode && product.Id != id {
			return Product{}, web.NewCodeResponse(http.StatusConflict, errors.New("Product_code already exists"))
		}
	}

	product, err := s.repository.Update(id, productCode, description, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan, freezingRate, productTypeId)

	if err != nil {
		return Product{}, web.NewCodeResponse(http.StatusNotFound, errors.New("Product not found"))
	}

	return product, web.ResponseCode{Code: http.StatusOK, Err: nil}
}

func (s service) UpdateExpirationRate(id int, expirationRate float64) (Product, web.ResponseCode) {
	product, err := s.repository.UpdateExpirationRate(id, expirationRate)

	if err != nil {
		return Product{}, web.NewCodeResponse(http.StatusNotFound, errors.New("Product not found"))
	}

	return product, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
