package inboundorders

import (
	"fmt"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	CreateInboundOrders(orderNumber, orderDate string, employeeId, productBatchId, warehouseId int) (InboundOrder, web.ResponseCode)
	GetReportInboundOrders(employeeId string) ([]ReportInboundOrder, web.ResponseCode)
}

type service struct {
	repository               Repository
	warehouseRepository      warehouses.Repository
	employeeRepository       employees.Repository
	productBatchesRepository product_batches.Repository
}

func NewService(r Repository, w warehouses.Repository, e employees.Repository, pb product_batches.Repository) Service {
	return &service{
		repository:               r,
		warehouseRepository:      w,
		employeeRepository:       e,
		productBatchesRepository: pb,
	}
}

func (s service) CreateInboundOrders(orderNumber, orderDate string, employeeId, productBatchId, warehouseId int) (InboundOrder, web.ResponseCode) {
	_, errEmployee := s.employeeRepository.GetOne(employeeId)
	if errEmployee != nil {
		if errEmployee.Error() == fmt.Sprintf("employee with id %d not found", employeeId) {
			return InboundOrder{}, web.NewCodeResponse(http.StatusUnprocessableEntity, errEmployee)
		}
		return InboundOrder{}, web.NewCodeResponse(http.StatusInternalServerError, errEmployee)
	}

	_, errWarehouse := s.warehouseRepository.GetOne(warehouseId)
	if errWarehouse != nil {
		if errWarehouse.Error() == fmt.Sprintf("warehouse with id %d not found", warehouseId) {
			return InboundOrder{}, web.NewCodeResponse(http.StatusUnprocessableEntity, errWarehouse)
		}
		return InboundOrder{}, web.NewCodeResponse(http.StatusInternalServerError, errWarehouse)
	}

	_, errProductBat := s.productBatchesRepository.GetOne(productBatchId)
	if errProductBat != nil {
		if errProductBat.Error() == fmt.Sprintf("product_batch with batch_number %d not found", productBatchId) {
			return InboundOrder{}, web.NewCodeResponse(http.StatusUnprocessableEntity, errProductBat)
		}
		return InboundOrder{}, web.NewCodeResponse(http.StatusInternalServerError, errProductBat)
	}

	result, err := s.repository.CreateInboundOrders(orderNumber, orderDate, employeeId, productBatchId, warehouseId)
	if err != nil {
		return InboundOrder{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return result, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetReportInboundOrders(employeeId string) ([]ReportInboundOrder, web.ResponseCode) {
	report, err := s.repository.GetReportInboundOrders(employeeId)

	if err != nil {
		return []ReportInboundOrder{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return report, web.NewCodeResponse(http.StatusOK, nil)
}
