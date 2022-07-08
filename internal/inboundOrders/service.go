package inboundorders

import (
	"fmt"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	CreateInboundOrders(orderNumber, orderDate string, employeeId, productBatchId, warehouseId int) (InboundOrder, web.ResponseCode)
	GetReportInboundOrders(employeeId string) ([]ReportInboundOrder, web.ResponseCode)
}

type service struct {
	repository          Repository
	warehouseRepository warehouses.Repository
	employeeRepository  employees.Repository
}

func NewService(r Repository, w warehouses.Repository, e employees.Repository) Service {
	return &service{
		repository:          r,
		warehouseRepository: w,
		employeeRepository:  e,
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
