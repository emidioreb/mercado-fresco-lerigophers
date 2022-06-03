package sellers

import (
	"errors"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"net/http"
)

type Service interface {
	Create(cid int, companyName, address, telephone string) (Seller, web.ResponseCode)
	GetOne(id int) (Seller, error)
	GetAll() ([]Seller, error)
	Delete(id int) error
	Update(id, cid int, companyName, address, telephone string) (Seller, web.ResponseCode)
	UpdateAddress(id int, address string) (Seller, web.ResponseCode)
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

func (s service) GetOne(id int) (Seller, error) {
	seller, err := s.repository.GetOne(id)

	if err != nil {
		return Seller{}, err
	}
	return seller, nil
}

func (s service) GetAll() ([]Seller, error) {
	sellers, err := s.repository.GetAll()

	if err != nil {
		return []Seller{}, err
	}
	return sellers, nil
}

func (s service) Delete(id int) error {
	err := s.repository.Delete(id)

	if err != nil {
		return err
	}
	return nil
}

func (s service) Update(id, cid int, companyName, address, telephone string) (Seller, web.ResponseCode) {
	allSellers, _ := s.GetAll()

	for _, seller := range allSellers {

		if seller.Cid == cid && seller.Id != id {
			return Seller{}, web.NewCodeResponse(http.StatusConflict, errors.New("cid already exists"))
		}
	}

	seller, err := s.repository.Update(id, cid, companyName, address, telephone)

	if err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusNotFound, errors.New("seller not found"))
	}

	return seller, web.ResponseCode{Code: http.StatusOK, Err: nil}
}

func (s service) UpdateAddress(id int, address string) (Seller, web.ResponseCode) {
	seller, err := s.repository.UpdateAddress(id, address)

	if err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusNotFound, errors.New("seller not found"))
	}

	return seller, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
