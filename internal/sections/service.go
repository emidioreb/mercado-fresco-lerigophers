package sections

import (
	"errors"
	"net/http"

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
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId int) (Section, web.ResponseCode) {
	allSections, _ := s.GetAll()

	for _, section := range allSections {
		if section.SectionNumber == sectionNumber {
			return Section{}, web.NewCodeResponse(http.StatusConflict, errors.New("section number already exists"))
		}
	}

	section, _ := s.repository.Create(sectionNumber, currentTemperature, minimumTemperature, currentCapacity, mininumCapacity, maximumCapacity, warehouseId, productTypeId)

	return section, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Section, web.ResponseCode) {
	section, err := s.repository.GetOne(id)

	if err != nil {
		return Section{}, web.NewCodeResponse(http.StatusNotFound, err)
	}
	return section, web.NewCodeResponse(http.StatusNotFound, nil)
}

func (s service) GetAll() ([]Section, web.ResponseCode) {
	sections, err := s.repository.GetAll()
	return sections, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Section, web.ResponseCode) {
	_, responseCode := s.GetOne(id)

	if responseCode.Err != nil {
		return Section{}, web.NewCodeResponse(http.StatusNotFound, errors.New("section not found"))
	}

	allSections, _ := s.GetAll()
	var sectionNumberReqData = requestData["section_number"]

	if sectionNumberReqData != nil {
		for _, section := range allSections {
			if float64(section.SectionNumber) == sectionNumberReqData && section.Id != id {
				return Section{}, web.NewCodeResponse(http.StatusConflict, errors.New("section number already exists"))
			}
		}

	}
	section, _ := s.repository.Update(id, requestData)

	return section, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
