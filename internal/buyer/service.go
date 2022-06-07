package buyers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(cardNumberId string, firstName, lastName string) (Buyer, web.ResponseCode)
	GetOne(id int) (Buyer, web.ResponseCode)
	GetAll() ([]Buyer, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]string) (Buyer, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(cardNumberId string, firstName string, lastName string) (Buyer, web.ResponseCode) {
	allBuyers, _ := s.GetAll()

	for _, buyer := range allBuyers {
		if buyer.CardNumberId == cardNumberId  {
			return Buyer{}, web.NewCodeResponse(http.StatusConflict, errors.New("CardNumberId already exists"))
		}
	}
	fmt.Printf("%T, %s, cardNumberId",cardNumberId, cardNumberId)

	Buyer, _ := s.repository.Create(cardNumberId , firstName, lastName )

	return Buyer, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Buyer, web.ResponseCode) {
	buyer, err := s.repository.GetOne(id)

	if err != nil {
		return Buyer{}, web.NewCodeResponse(http.StatusNotFound, err)
	}
	return buyer, web.NewCodeResponse(http.StatusNotFound, nil)
}

func (s service) GetAll() ([]Buyer, web.ResponseCode) {
	buyers, err := s.repository.GetAll()
	return buyers, web.NewCodeResponse(http.StatusOK, err)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]string) (Buyer, web.ResponseCode) {
	_, responseCode := s.GetOne(id)
	allBuyers, _ := s.GetAll()
	buyerNumberReqData := requestData["card_number_id"]

	if responseCode.Err != nil {
		return Buyer{}, web.NewCodeResponse(http.StatusNotFound, errors.New("buyer not found"))
	}

	for _, buyer := range allBuyers {
		if buyer.CardNumberId == buyerNumberReqData && buyer.Id != id {
			return Buyer{}, web.NewCodeResponse(http.StatusConflict, errors.New("buyer number already exists"))
		}
	}

	buyer, _ := s.repository.Update(id, requestData)

	return buyer, web.ResponseCode{Code: http.StatusOK, Err: nil}
}