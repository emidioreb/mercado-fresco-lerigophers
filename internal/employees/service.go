package employees

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
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
	repository    Repository
	warehouseRepo warehouses.Repository
}

func NewService(r Repository, w warehouses.Repository) Service {
	return &service{
		repository:    r,
		warehouseRepo: w,
	}
}

func (s service) Create(cardNumber string, firstName, lastName string, warehouseId int) (Employee, web.ResponseCode) {

	errGetByCardNumber := s.repository.GetOneByCardNumber(0, cardNumber)

	if errGetByCardNumber != nil {
		if errGetByCardNumber.Error() == "card_number_id already exists" {
			return Employee{}, web.NewCodeResponse(http.StatusConflict, errors.New("card_number_id already exists"))
		}
		return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, errGetByCardNumber)
	}

	_, warehouseErr := s.warehouseRepo.GetOne(warehouseId)
	if warehouseErr != nil {
		if warehouseErr.Error() == fmt.Sprintf("warehouse with id %d not found", warehouseId) {
			return Employee{}, web.NewCodeResponse(http.StatusConflict, warehouseErr)
		}
		return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, warehouseErr)
	}

	employee, errCreate := s.repository.Create(cardNumber, firstName, lastName, warehouseId)

	if errCreate != nil {
		return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, errCreate)
	}

	return employee, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Employee, web.ResponseCode) {
	employee, err := s.repository.GetOne(id)

	if err != nil {
		if err.Error() == fmt.Sprintf("employee with id %d not found", id) {
			return Employee{}, web.NewCodeResponse(http.StatusNotFound, err)
		}
		return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}
	return employee, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) GetAll() ([]Employee, web.ResponseCode) {
	employees, err := s.repository.GetAll()

	if err != nil {
		return []Employee{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return employees, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		if err.Error() == fmt.Sprintf("employee with id %d not found", id) {
			return web.NewCodeResponse(http.StatusNotFound, err)
		}
		return web.NewCodeResponse(http.StatusInternalServerError, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Employee, web.ResponseCode) {
	_, responseCode := s.GetOne(id)

	cardNumberReqData := requestData["card_number_id"]

	if responseCode.Err != nil {
		return Employee{}, web.NewCodeResponse(http.StatusNotFound, responseCode.Err)
	}

	if cardNumberReqData != nil {
		errGetByCardNumber := s.repository.GetOneByCardNumber(id, cardNumberReqData.(string))
		if errGetByCardNumber != nil {
			if errGetByCardNumber.Error() == "card_number_id already exists" {
				return Employee{}, web.NewCodeResponse(http.StatusConflict, errGetByCardNumber)
			}
			return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, errGetByCardNumber)
		}
	}

	if requestData["warehouse_id"] != nil {
		if int(requestData["warehouse_id"].(float64)) != 0 {
			_, warehouseErr := s.warehouseRepo.GetOne(int(requestData["warehouse_id"].(float64)))
			if warehouseErr != nil {
				if warehouseErr.Error() == fmt.Sprintf("warehouse with id %d not found", int(requestData["warehouse_id"].(float64))) {
					return Employee{}, web.NewCodeResponse(http.StatusNotFound, warehouseErr)
				}
				return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, warehouseErr)
			}
		}
	}

	employee, err := s.repository.Update(id, requestData)
	if err != nil {
		return Employee{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return employee, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
