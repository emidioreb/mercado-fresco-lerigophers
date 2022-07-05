package localities

import (
	"errors"
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
)

type Service interface {
	CreateLocality(
		id,
		localityName,
		provinceName,
		countryName string,
	) (Locality, web.ResponseCode)
	GetReportSellers(localityId string) ([]ReportSellers, web.ResponseCode)
	GetReportCarriers(localityId string) ([]ReportCarriers, web.ResponseCode) 
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) CreateLocality(id, localityName, provinceName, countryName string) (Locality, web.ResponseCode) {
	_, err := s.repository.GetOne(id)
	if err == nil {
		return Locality{}, web.NewCodeResponse(http.StatusConflict, errors.New("locality already exists"))
	}

	result, err := s.repository.CreateLocality(id, localityName, provinceName, countryName)
	if err != nil {
		return Locality{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return result, web.NewCodeResponse(http.StatusCreated, nil)
}

func (s service) GetReportSellers(localityId string) ([]ReportSellers, web.ResponseCode) {
	report, err := s.repository.GetReportSellers(localityId)

	if err != nil {
		return []ReportSellers{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return report, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) GetReportCarriers(localityId string) ([]ReportCarriers, web.ResponseCode) {
	report, err := s.repository.GetReportCarriers(localityId)

	if err != nil {
		return []ReportCarriers{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return report, web.NewCodeResponse(http.StatusOK, nil)
}


