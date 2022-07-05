package product_batches

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	CreateProductBatch(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId int, DueDate, ManufacturingDate string) (ProductBatches, web.ResponseCode)
	GetReportSection(SectionId int) ([]ProductsQuantity, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) CreateProductBatch(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId int, DueDate, ManufacturingDate string) (ProductBatches, web.ResponseCode) {
	_, err := s.repository.GetOne(BatchNumber)
	if err == nil {
		return ProductBatches{}, web.NewCodeResponse(http.StatusConflict, errors.New("product_batch already exists"))
	}

	result, err := s.repository.CreateProductBatch(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId, DueDate, ManufacturingDate)
	if err != nil {
		return ProductBatches{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return result, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetReportSection(SectionId int) ([]ProductsQuantity, web.ResponseCode) {
	report, err := s.repository.GetReportSection(SectionId)

	if err != nil {
		return []ProductsQuantity{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return report, web.NewCodeResponse(http.StatusOK, nil)
}
