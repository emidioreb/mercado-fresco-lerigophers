package sections

import (
	"net/http"

	product_types "github.com/emidioreb/mercado-fresco-lerigophers/internal/productTypes"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, web.ResponseCode)
	GetOne(id int) (Section, web.ResponseCode)
	GetAll() ([]Section, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]interface{}) (Section, web.ResponseCode)
}

type service struct {
	repository            Repository
	warehouseRepository   warehouses.Repository
	productTypeRepository product_types.Repository
}

func NewService(r Repository, wr warehouses.Repository, pr product_types.Repository) Service {
	return &service{
		repository:            r,
		warehouseRepository:   wr,
		productTypeRepository: pr,
	}
}

func (s service) Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, web.ResponseCode) {
	if _, err := s.repository.GetBySectionNumber(sectionNumber); err != nil {
		return Section{}, web.NewCodeResponse(http.StatusConflict, err)
	}

	_, err := s.warehouseRepository.GetOne(warehouseId)
	if err != nil {
		return Section{}, web.NewCodeResponse(http.StatusConflict, err)
	}

	err = s.productTypeRepository.GetOne(productTypeId)
	if err != nil {
		return Section{}, web.NewCodeResponse(http.StatusConflict, err)
	}

	section, err := s.repository.Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId)
	if err != nil {
		return Section{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return section, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Section, web.ResponseCode) {
	section, err := s.repository.GetOne(id)

	if err != nil {
		return Section{}, web.NewCodeResponse(http.StatusNotFound, err)
	}

	return section, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) GetAll() ([]Section, web.ResponseCode) {
	sections, err := s.repository.GetAll()

	if err != nil {
		return []Section{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return sections, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) Delete(id int) web.ResponseCode {
	_, err := s.repository.GetOne(id)
	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}

	err = s.repository.Delete(id)
	if err != nil {
		return web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Section, web.ResponseCode) {
	_, responseCode := s.GetOne(id)

	if responseCode.Err != nil {
		return Section{}, responseCode
	}

	if sectionNumberReqData := requestData["section_number"]; sectionNumberReqData != nil {
		if selectedId, err := s.repository.GetBySectionNumber(int(sectionNumberReqData.(float64))); err != nil && id != selectedId {
			return Section{}, web.NewCodeResponse(http.StatusConflict, err)
		}
	}

	if warehouseId := requestData["warehouse_id"]; warehouseId != nil {
		_, err := s.warehouseRepository.GetOne(int(warehouseId.(float64)))
		if err != nil {
			return Section{}, web.NewCodeResponse(http.StatusConflict, err)
		}
	}

	if productTypeId := requestData["product_type_id"]; productTypeId != nil {
		err := s.productTypeRepository.GetOne(int(productTypeId.(float64)))
		if err != nil {
			return Section{}, web.NewCodeResponse(http.StatusConflict, err)
		}
	}

	section, err := s.repository.Update(id, requestData)
	if err != nil {
		return Section{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return section, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
