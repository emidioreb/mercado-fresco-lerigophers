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
	GetAllReportSellers() ([]ReportSellers, web.ResponseCode)
	GetReportOneSeller(localityId string) ([]ReportSellers, web.ResponseCode)
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

func (s service) GetAllReportSellers() ([]ReportSellers, web.ResponseCode) {
	report, err := s.repository.GetReportSellers("")
	if err != nil {
		return []ReportSellers{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return report, web.NewCodeResponse(http.StatusOK, nil)
}

func (s service) GetReportOneSeller(localityId string) ([]ReportSellers, web.ResponseCode) {
	if _, err := s.repository.GetOne(localityId); err != nil {
		return []ReportSellers{}, web.NewCodeResponse(http.StatusNotFound, err)
	}

	report, err := s.repository.GetReportSellers(localityId)
	if err != nil {
		return []ReportSellers{}, web.NewCodeResponse(http.StatusInternalServerError, err)
	}

	return report, web.NewCodeResponse(http.StatusOK, nil)
}
