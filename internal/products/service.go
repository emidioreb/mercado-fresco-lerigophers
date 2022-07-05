package products

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate float64, productTypeId, sellerId int) (Product, web.ResponseCode)
	GetOne(id int) (Product, web.ResponseCode)
	GetAll() ([]Product, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]interface{}) (Product, web.ResponseCode)
}

type service struct {
	repository       Repository
	sellerRepository sellers.Repository
}

func NewService(r Repository, sr sellers.Repository) Service {
	return &service{
		repository:       r,
		sellerRepository: sr,
	}
}

func (s service) Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperature,
	freezingRate float64, productTypeId, sellerId int) (Product, web.ResponseCode) {
	allProducts, _ := s.repository.GetAll()

	for _, product := range allProducts {
		if product.ProductCode == productCode {
			return Product{}, web.NewCodeResponse(http.StatusConflict, errors.New("Product_code already exists"))
		}
	}

	if _, err := s.sellerRepository.GetOne(sellerId); err != nil {
		return Product{}, web.NewCodeResponse(http.StatusConflict, errors.New("informed seller_id don't exists"))
	}

	product, err := s.repository.Create(productCode, description, width, height, length, netWeight, expirationRate, recommendedFreezingTemperature,
		freezingRate, productTypeId, sellerId)

	fmt.Println(err)

	return product, web.NewCodeResponse(http.StatusCreated, nil)
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

func (s service) Update(id int, requestData map[string]interface{}) (Product, web.ResponseCode) {
	_, err := s.repository.GetOne(id)
	if err != nil {
		return Product{}, web.NewCodeResponse(http.StatusNotFound, err)
	}

	allProducts, _ := s.GetAll()
	productCodeReqData := requestData["product_code"]

	for _, product := range allProducts {
		if product.ProductCode == productCodeReqData && product.Id != id {
			return Product{}, web.NewCodeResponse(http.StatusConflict, errors.New("product_code already exists"))
		}
	}

	product, _ := s.repository.Update(id, requestData)

	return product, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
