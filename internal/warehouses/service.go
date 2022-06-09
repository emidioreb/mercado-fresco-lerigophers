package warehouses

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, web.ResponseCode)
	GetOne(id int) (Warehouse, web.ResponseCode)
	GetAll() ([]Warehouse, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]interface{}) (Warehouse, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(warehouseCode, adress, telephone string, minimumCapacity, maxmumCapacity int) (Warehouse, web.ResponseCode) {
	allWarehouses, _ := s.GetAll()

	for _, warehouse := range allWarehouses {
		if warehouse.WarehouseCode == warehouseCode {
			return Warehouse{}, web.NewCodeResponse(http.StatusConflict, errors.New("warehouse_code already exists"))
		}
	}

	warehouse, _ := s.repository.Create(warehouseCode, adress, telephone, minimumCapacity, maxmumCapacity)

	return warehouse, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Warehouse, web.ResponseCode) {
	warehouse, err := s.repository.GetOne(id)

	if err != nil {
		return Warehouse{}, web.NewCodeResponse(http.StatusNotFound, err)
	}
	return warehouse, web.NewCodeResponse(http.StatusNotFound, nil)
}

func (s service) GetAll() ([]Warehouse, web.ResponseCode) {
	warehouse, err := s.repository.GetAll()
	return warehouse, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Warehouse, web.ResponseCode) {
	_, responseCode := s.GetOne(id)
	allWarehouses, _ := s.GetAll()
	warehouseCodeReqData := requestData["warehouse_code"]

	if responseCode.Err != nil {
		return Warehouse{}, web.NewCodeResponse(http.StatusNotFound, errors.New("warehouse not found"))
	}

	for _, warehouse := range allWarehouses {
		if warehouse.WarehouseCode == warehouseCodeReqData && warehouse.Id != id {
			return Warehouse{}, web.NewCodeResponse(http.StatusConflict, errors.New("warehouse_code already exists"))
		}
	}

	warehouse, _ := s.repository.Update(id, requestData)

	return warehouse, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
