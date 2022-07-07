package carriers

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	Create(cid, companyName, address, telephone, localityId string) (Carry, web.ResponseCode)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(cid, companyName, address, telephone, localityId string) (Carry, web.ResponseCode) {

	_, err := s.repository.GetOne(cid)

	if err == nil {
		return Carry{}, web.NewCodeResponse(http.StatusConflict, errors.New("CID already exists"))
	}

	CarryResult, err := s.repository.Create(cid, companyName, address, telephone, localityId)
	if err != nil {
		return Carry{}, web.NewCodeResponse(
			http.StatusInternalServerError,
			err,
		)
	}

	return CarryResult, web.NewCodeResponse(http.StatusCreated, nil)
}
