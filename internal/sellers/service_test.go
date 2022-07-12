package sellers_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
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
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(0, nil).Once()
		mockedLocality.On("GetOne", fakeSellers[0].LocalityId).Return(localities.Locality{}, nil).Once()

		mockedRepository.On("Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(fakeSellers[0], nil).Once()

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

	t.Run("Test error case if seller CID already exists on create", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		expectedError := errors.New("cid already exists")
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(1, expectedError).Once()

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, resp := service.Create(
			fakeSellers[0].Cid,
			fakeSellers[0].CompanyName,
			fakeSellers[0].Address,
			fakeSellers[0].Telephone,
			fakeSellers[0].LocalityId)

		assert.Error(t, resp.Err)
		assert.Equal(t, expectedError.Error(), resp.Err.Error())
		mockedRepository.AssertExpectations(t)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("Test internal server error when verify CID", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		expectedError := errors.New("some error")
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(0, expectedError).Once()

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, resp := service.Create(
			fakeSellers[0].Cid,
			fakeSellers[0].CompanyName,
			fakeSellers[0].Address,
			fakeSellers[0].Telephone,
			fakeSellers[0].LocalityId)

		assert.Error(t, resp.Err)
		mockedRepository.AssertExpectations(t)
		assert.Equal(t, expectedError.Error(), resp.Err.Error())
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Test conflict if locality_id do not exist", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		expectedError := errors.New("some error")
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(1, nil).Once()
		mockedLocality.On("GetOne", mock.AnythingOfType("string")).Return(localities.Locality{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, resp := service.Create(
			fakeSellers[0].Cid,
			fakeSellers[0].CompanyName,
			fakeSellers[0].Address,
			fakeSellers[0].Telephone,
			fakeSellers[0].LocalityId)

		assert.Error(t, resp.Err)
		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test internal server error when create seller", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		expectedError := errors.New("some error")
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(1, nil).Once()
		mockedLocality.On("GetOne", mock.AnythingOfType("string")).Return(localities.Locality{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(sellers.Seller{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, resp := service.Create(
			fakeSellers[0].Cid,
			fakeSellers[0].CompanyName,
			fakeSellers[0].Address,
			fakeSellers[0].Telephone,
			fakeSellers[0].LocalityId)

		assert.Error(t, resp.Err)
		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("Verify the successfully case if the seller is deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, nil)
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

		mockedRepository.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).
			Return(
				sellers.Seller{},
				expectedError,
			)

		service := sellers.NewService(mockedRepository, mockedLocality)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, http.StatusNotFound, result.Code)
		assert.Equal(t, expectedError, result.Err)
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

	t.Run("Test internal server error when get all", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		expectedError := errors.New("any error")
		mockedRepository.On("GetAll").Return([]sellers.Seller{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)

		_, err := service.GetAll()
		assert.Error(t, err.Err)
		assert.Equal(t, expectedError, err.Err)
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
			"cid":          1.0,
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
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(0, nil)
		mockedLocality.On("GetOne", mock.AnythingOfType("string")).Return(localities.Locality{}, nil)
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

		service := sellers.NewService(mockedRepository, mockedLocality)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("Test error case if seller CID already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		requestData := map[string]interface{}{"cid": 2.0}

		expectedError := errors.New("cid already exists")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, nil)
		mockedRepository.On("FindByCID", mock.AnythingOfType("int")).Return(1, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)

		_, resp := service.Update(2, requestData)
		assert.Error(t, resp.Err)

		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("Update seller when locality_id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		requestData := map[string]interface{}{"locality_id": 2.0}

		expectedError := errors.New("cid already exists")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, nil)
		mockedLocality.On("GetOne", mock.AnythingOfType("string")).Return(localities.Locality{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)

		_, resp := service.Update(1, requestData)
		assert.Error(t, resp.Err)

		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("Internal server error when updating seller", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedLocality := new(mockLocalityRepository.Repository)

		requestData := map[string]interface{}{"address": "FAKE ADDRESS"}

		expectedError := errors.New("cid already exists")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(sellers.Seller{}, nil)
		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(sellers.Seller{}, expectedError)

		service := sellers.NewService(mockedRepository, mockedLocality)

		_, resp := service.Update(1, requestData)
		assert.Error(t, resp.Err)

		assert.Equal(t, expectedError, resp.Err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
