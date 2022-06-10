package employees

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(cardNumber, firstName, lastName string, warehouseId int) (Employee, web.ResponseCode)
	GetOne(id int) (Employee, web.ResponseCode)
	GetAll() ([]Employee, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]interface{}) (Employee, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(cardNumber string, firstName, lastName string, warehouseId int) (Employee, web.ResponseCode) {
	allEmployees, _ := s.GetAll()

	for _, employee := range allEmployees {
		if employee.CardNumberId == cardNumber {
			return Employee{}, web.NewCodeResponse(http.StatusConflict, errors.New("card_number_id already exists"))
		}
	}

	employee, _ := s.repository.Create(cardNumber, firstName, lastName, warehouseId)

	return employee, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Employee, web.ResponseCode) {
	employee, err := s.repository.GetOne(id)

	if err != nil {
		return Employee{}, web.NewCodeResponse(http.StatusNotFound, err)
	}
	return employee, web.NewCodeResponse(http.StatusNotFound, nil)
}

func (s service) GetAll() ([]Employee, web.ResponseCode) {
	employees, err := s.repository.GetAll()
	return employees, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Employee, web.ResponseCode) {
	_, responseCode := s.GetOne(id)
	allEmployees, _ := s.GetAll()

	cardNumberReqData := requestData["card_number_id"]

	if responseCode.Err != nil {
		return Employee{}, web.NewCodeResponse(http.StatusNotFound, errors.New("employee not found"))
	}

	for _, employee := range allEmployees {
		if employee.CardNumberId == cardNumberReqData && employee.Id != id {
			return Employee{}, web.NewCodeResponse(http.StatusConflict, errors.New("card_number_id already exists"))
		}
	}

	employee, _ := s.repository.Update(id, requestData)

	return employee, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
