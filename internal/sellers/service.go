package sellers

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(cid int, companyName, address, telephone, localityId string) (Seller, web.ResponseCode)
	GetOne(id int) (Seller, web.ResponseCode)
	GetAll() ([]Seller, web.ResponseCode)
	Delete(id int) web.ResponseCode
	Update(id int, requestData map[string]interface{}) (Seller, web.ResponseCode)
}

type service struct {
	repository         Repository
	localityRepository localities.Repository
}

func NewService(r Repository, lr localities.Repository) Service {
	return &service{
		repository:         r,
		localityRepository: lr,
	}
}

func (s service) Create(cid int, companyName, address, telephone, localityId string) (Seller, web.ResponseCode) {
	allSellers, _ := s.repository.GetAll()
	for _, seller := range allSellers {
		if seller.Cid == cid {
			return Seller{}, web.NewCodeResponse(http.StatusConflict, errors.New("cid already exists"))
		}
	}

	_, localityErr := s.localityRepository.GetOne(localityId)
	if localityErr != nil {
		return Seller{}, web.NewCodeResponse(http.StatusConflict, errors.New("informed locality_id don't exists"))
	}

	seller, err := s.repository.Create(cid, companyName, address, telephone, localityId)
	if err != nil {
		return Seller{}, web.NewCodeResponse(
			http.StatusInternalServerError,
			err,
		)
	}

	return seller, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetOne(id int) (Seller, web.ResponseCode) {
	seller, err := s.repository.GetOne(id)

	if err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusNotFound, err)
	}

	return seller, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) GetAll() ([]Seller, web.ResponseCode) {
	sellers, err := s.repository.GetAll()

	if err != nil {
		return []Seller{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return sellers, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) Delete(id int) web.ResponseCode {
	err := s.repository.Delete(id)

	if err != nil {
		return web.NewCodeResponse(http.StatusNotFound, err)
	}
	return web.NewCodeResponse(http.StatusNoContent, nil)
}

func (s service) Update(id int, requestData map[string]interface{}) (Seller, web.ResponseCode) {
	if _, err := s.repository.GetOne(id); err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusNotFound, err)
	}

	allSellers, _ := s.GetAll()
	if currCid := requestData["cid"]; currCid != nil {
		for _, seller := range allSellers {
			if float64(seller.Cid) == currCid && seller.Id != id {
				return Seller{}, web.NewCodeResponse(http.StatusConflict, errors.New("cid already exists"))
			}
		}
	}

	if currLocalityId := requestData["locality_id"]; currLocalityId != nil {
		parsedLocalityId, _ := currLocalityId.(string)
		if _, err := s.localityRepository.GetOne(parsedLocalityId); err != nil {
			return Seller{}, web.NewCodeResponse(
				http.StatusConflict,
				errors.New("informed locality_id don't exists"))
		}
	}

	seller, err := s.repository.Update(id, requestData)
	if err != nil {
		return Seller{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return seller, web.ResponseCode{Code: http.StatusOK, Err: nil}
}
