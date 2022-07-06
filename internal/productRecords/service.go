package product_records

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	CreateProductRecord(LastUpdateDate string, PurchasePrice float64, SalePrice float64, ProductId int) (ProductRecords, web.ResponseCode)
}

type service struct {
	repository        Repository
	productRepository products.Repository
}

func NewService(r Repository, pr products.Repository) Service {
	return &service{
		repository:        r,
		productRepository: pr,
	}
}

func (s service) CreateProductRecord(LastUpdateDate string, PurchasePrice float64, SalePrice float64, ProductId int) (ProductRecords, web.ResponseCode) {
	_, err := s.productRepository.GetOne(ProductId)
	if err != nil {
		fmt.Println(err, 30)
		return ProductRecords{}, web.NewCodeResponse(http.StatusConflict, errors.New("product_id don`t exists"))
	}

	result, err := s.repository.CreateProductRecord(LastUpdateDate, PurchasePrice, SalePrice, ProductId)
	if err != nil {
		return ProductRecords{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return result, web.NewCodeResponse(http.StatusCreated, nil)
}
