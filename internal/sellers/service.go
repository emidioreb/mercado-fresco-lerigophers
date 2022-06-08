package sellers

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(cid int, companyName, address, telephone string) (Seller, web.ResponseCode)
	GetOne(id int) (Seller, web.ResponseCode)
	GetAll() ([]Seller, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]interface{}) (Seller, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(cid int, companyName, address, telephone string) (Seller, web.ResponseCode) {
	allSellers, _ := s.GetAll()

	for _, seller := range allSellers {
		if seller.Cid == cid {
			return Seller{}, web.NewCodeResponse(http.StatusConflict, errors.New("cid already exists"))
		}
	}

	seller, _ := s.repository.Create(cid, companyName, address, telephone)

	return seller, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Seller, web.ResponseCode) {
	seller, err := s.repository.GetOne(id)

	if err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusNotFound, err)
	}
	return seller, web.NewCodeResponse(http.StatusNotFound, nil)
}

func (s service) GetAll() ([]Seller, web.ResponseCode) {
	sellers, err := s.repository.GetAll()
	return sellers, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Seller, web.ResponseCode) {
	_, responseCode := s.GetOne(id)
	if responseCode.Err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusNotFound, errors.New("seller not found"))
	}

	allSellers, _ := s.GetAll()
	currentCid := int(requestData["cid"].(float64))

	for _, seller := range allSellers {
		if seller.Cid == currentCid && seller.Id != id {
			return Seller{}, web.NewCodeResponse(http.StatusConflict, errors.New("cid already exists"))
		}
	}

	seller, _ := s.repository.Update(id, requestData)

	return seller, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
