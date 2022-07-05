package sellers_test

import (
	"errors"
	"net/http"
	"testing"

	mockLocalityRepository "github.com/emidioreb/mercado-fresco-lerigophers/internal/localities/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fakeSellers = []sellers.Seller{{
	Id:          1,
	Cid:         1,
	CompanyName: "Gouveia empreendimentos",
	Address:     "Av. Nações Unidas",
	Telephone:   "3003",
	LocalityId:  "65760-000",
}, {
	Id:          2,
	Cid:         2,
	CompanyName: "Gouveia empreendimentos",
	Address:     "Av. Nações Unidas",
	Telephone:   "3003",
	LocalityId:  "65760-123",
}}

func TestServiceCreate(t *testing.T) {
	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)
		mockedRepository.On("GetAll").Return([]sellers.Seller{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(fakeSellers[0], nil)

		service := sellers.NewService(mockedRepository, mockedLocality)

		result, err := service.Create(
			fakeSellers[0].Cid,
			fakeSellers[0].CompanyName,
			fakeSellers[0].Address,
			fakeSellers[0].Telephone,
			fakeSellers[0].LocalityId)
		assert.Nil(t, err.Err)

		assert.Equal(t, fakeSellers[0], result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test error case if seller CID already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		expectedError := errors.New("cid already exists")

		mockedRepository.On("GetAll").Return(fakeSellers, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(sellers.Seller{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)

		_, err := service.Create(
			fakeSellers[0].Cid,
			fakeSellers[0].CompanyName,
			fakeSellers[0].Address,
			fakeSellers[0].Telephone,
			fakeSellers[0].LocalityId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, err.Code, http.StatusConflict)
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("Verify the successfully case if the seller is deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(nil)

		service := sellers.NewService(mockedRepository, mockedLocality)
		result := service.Delete(1)
		assert.Nil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNoContent)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Verify the error case if seller do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)
		expectedError := errors.New("seller with id 1 not found")

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNotFound)
		assert.Equal(t, result.Err, expectedError)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceGetAll(t *testing.T) {
	t.Run("Test if an array of sellers is returned when GetAll", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		mockedRepository.On("GetAll").Return(fakeSellers, nil)

		service := sellers.NewService(mockedRepository, mockedLocality)

		result, err := service.GetAll()
		assert.Nil(t, err.Err)

		assert.Len(t, result, 2)
		assert.Equal(t, fakeSellers[1].Cid, result[1].Cid)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceGetOne(t *testing.T) {
	t.Run("Test if seller is returned based on valid id", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(fakeSellers[0], nil)

		service := sellers.NewService(mockedRepository, mockedLocality)

		result, err := service.GetOne(1)
		assert.Nil(t, err.Err)

		assert.Equal(t, fakeSellers[0], result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test if requested seller do not exists based on invalid id", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)
		expectedError := errors.New("seller with id 1 not found")

		mockedRepository.On(
			"GetOne",
			mock.AnythingOfType("int")).
			Return(sellers.Seller{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, err := service.GetOne(1)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, expectedError, err.Err)

		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("Return the updated seller when successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		requestData := map[string]interface{}{
			"cid":          1,
			"company_name": "Mercado Solto",
			"address":      "Av. Fake das Dores",
			"telephone":    "12345",
			"locality_id":  "12345",
		}

		expectedSeller := sellers.Seller{
			Id:          1,
			Cid:         1,
			CompanyName: "Mercado Solto",
			Address:     "Av. Fake das Dores",
			Telephone:   "12345",
			LocalityId:  "12345",
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(fakeSellers[0], nil)
		mockedRepository.On("GetAll").Return([]sellers.Seller{}, nil)
		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(expectedSeller, nil)

		service := sellers.NewService(mockedRepository, mockedLocality)
		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)
		assert.Equal(t, expectedSeller, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Return null when seller id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)
		expectedError := errors.New("seller with id 1 not found")
		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(sellers.Seller{}, expectedError).Once()

		mockedRepository.On("GetAll").
			Return([]sellers.Seller{}, nil).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(sellers.Seller{}, nil).Once()

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("Test error case if seller CID already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		requestData := map[string]interface{}{
			"cid": 2.0,
		}

		expectedError := errors.New("cid already exists")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, nil).Once()
		mockedRepository.On("GetAll").Return(fakeSellers, nil).Once()
		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(sellers.Seller{}, expectedError).Once()

		service := sellers.NewService(mockedRepository, mockedLocality)

		_, err := service.Update(fakeSellers[0].Id, requestData)
		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, err.Code, http.StatusConflict)
	})
}
