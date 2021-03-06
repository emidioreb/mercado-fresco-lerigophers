package localities_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fakeLocalities = []localities.Locality{
	{
		Id:           "65760000",
		LocalityName: "Presidente Dutra",
		ProvinceName: "MA",
		CountryName:  "BR",
	},
	{
		Id:           "12345678",
		LocalityName: "Florianópolis",
		ProvinceName: "SC",
		CountryName:  "BR",
	},
}

var fakeReports = []localities.ReportSellers{
	{
		LocalityId:   "65760000",
		LocalityName: "Presidente Dutra",
		SellersCount: 1,
	},
	{
		LocalityId:   "12345678",
		LocalityName: "Florianópolis",
		SellersCount: 1,
	},
}

var fakeCarriersReports = []localities.ReportCarriers{
	{
		LocalityId:    "123",
		LocalityName:  "Locality",
		CarriersCount: 1,
	},
	{
		LocalityId:    "456",
		LocalityName:  "Locality",
		CarriersCount: 1,
	},
}

func TestCreateLocality(t *testing.T) {
	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On(
			"GetOne",
			mock.AnythingOfType("string"),
		).Return(localities.Locality{}, errors.New(""))

		mockedRepository.On(
			"CreateLocality",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(fakeLocalities[0], nil)

		service := localities.NewService(mockedRepository)
		result, err := service.CreateLocality(
			fakeLocalities[0].Id,
			fakeLocalities[0].LocalityName,
			fakeLocalities[0].ProvinceName,
			fakeLocalities[0].CountryName,
		)

		assert.Nil(t, err.Err)

		assert.Equal(t, fakeLocalities[0], result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Internal server error case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetOne", mock.AnythingOfType("string")).
			Return(localities.Locality{}, errors.New(""))

		mockedRepository.On(
			"CreateLocality",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(localities.Locality{}, errors.New("any error"))

		service := localities.NewService(mockedRepository)
		_, err := service.CreateLocality(
			fakeLocalities[0].Id,
			fakeLocalities[0].LocalityName,
			fakeLocalities[0].ProvinceName,
			fakeLocalities[0].CountryName,
		)

		assert.Error(t, err.Err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Conflict locality", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetOne", mock.AnythingOfType("string")).
			Return(localities.Locality{}, nil)

		service := localities.NewService(mockedRepository)
		_, err := service.CreateLocality(
			fakeLocalities[0].Id,
			fakeLocalities[0].LocalityName,
			fakeLocalities[0].ProvinceName,
			fakeLocalities[0].CountryName,
		)

		assert.Error(t, err.Err)
		assert.Equal(t, http.StatusConflict, err.Code)
		mockedRepository.AssertExpectations(t)
	})
}

func TestGetOneReportSeller(t *testing.T) {
	t.Run("Test if get successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetOne", mock.AnythingOfType("string")).
			Return(localities.Locality{}, nil)

		mockedRepository.On("GetReportSellers", mock.AnythingOfType("string")).
			Return(fakeReports, nil)

		service := localities.NewService(mockedRepository)
		result, err := service.GetReportOneSeller("")

		assert.Nil(t, err.Err)
		assert.Len(t, result, 2)
		assert.Equal(t, http.StatusOK, err.Code)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test fail case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetOne", mock.AnythingOfType("string")).
			Return(localities.Locality{}, nil)

		mockedRepository.On("GetReportSellers", mock.AnythingOfType("string")).
			Return([]localities.ReportSellers{}, errors.New("any error"))

		service := localities.NewService(mockedRepository)
		result, err := service.GetReportOneSeller("")

		assert.Error(t, err.Err)
		assert.Len(t, result, 0)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test fail case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetOne", mock.AnythingOfType("string")).
			Return(localities.Locality{}, errors.New("locality with id 1 not found"))

		service := localities.NewService(mockedRepository)
		result, err := service.GetReportOneSeller("1")

		assert.Error(t, err.Err)
		assert.Len(t, result, 0)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockedRepository.AssertExpectations(t)
	})
}

func TestGetAllReportSellers(t *testing.T) {
	t.Run("Test if get successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetReportSellers", mock.AnythingOfType("string")).
			Return(fakeReports, nil)

		service := localities.NewService(mockedRepository)
		result, err := service.GetAllReportSellers()

		assert.Nil(t, err.Err)
		assert.Len(t, result, 2)
		assert.Equal(t, http.StatusOK, err.Code)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test fail case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetReportSellers", mock.AnythingOfType("string")).
			Return([]localities.ReportSellers{}, errors.New("any error"))

		service := localities.NewService(mockedRepository)
		result, err := service.GetAllReportSellers()

		assert.Error(t, err.Err)
		assert.Len(t, result, 0)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockedRepository.AssertExpectations(t)
	})
}

func TestGetAllReportCarriers(t *testing.T) {
	t.Run("Test if get successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetReportCarriers", mock.AnythingOfType("string")).
			Return(fakeCarriersReports, nil)

		service := localities.NewService(mockedRepository)
		result, err := service.GetReportCarriers("")

		assert.Nil(t, err.Err)
		assert.Len(t, result, 2)
		assert.Equal(t, http.StatusOK, err.Code)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test fail case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetReportCarriers", mock.AnythingOfType("string")).
			Return([]localities.ReportCarriers{}, errors.New("any error"))

		service := localities.NewService(mockedRepository)
		result, err := service.GetReportCarriers("")

		assert.Error(t, err.Err)
		assert.Len(t, result, 0)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockedRepository.AssertExpectations(t)
	})
}
