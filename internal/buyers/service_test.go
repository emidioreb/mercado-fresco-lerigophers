package buyers_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {
	t.Run("should return create buyer", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := buyers.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "José",
			LastName:     "Silva",
		}
		mockedRepository.On("GetAll").Return([]buyers.Buyer{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(input, nil)

		service := buyers.NewService(mockedRepository)

		result, err := service.Create(input.CardNumberId, input.FirstName, input.LastName)

		assert.Nil(t, err.Err)
		assert.Equal(t, result, input)

		mockedRepository.AssertExpectations(t)

	})

	t.Run("if CardNumberId exist should return erro", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		//
		expectedError := errors.New("CardNumberId already exists")

		input := buyers.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "José",
			LastName:     "Silva",
		}

		listBuyers := []buyers.Buyer{}
		listBuyers = append(listBuyers, input)

		mockedRepository.On("GetAll").Return(listBuyers, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(buyers.Buyer{}, expectedError)

		service := buyers.NewService(mockedRepository)
		_, err := service.Create(input.CardNumberId, input.FirstName, input.LastName)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, err.Code, http.StatusConflict)
	})

}

func TestServiceGetAll(t *testing.T) {
	t.Run("should return buyers list", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := []buyers.Buyer{
			{
				Id:           1,
				CardNumberId: "12345",
				FirstName:    "José",
				LastName:     "Silva",
			},
			{
				Id:           2,
				CardNumberId: "54321",
				FirstName:    "Maria",
				LastName:     "Pereira",
			},
		}

		mockedRepository.On("GetAll").Return(input, nil)

		service := buyers.NewService(mockedRepository)

		result, err := service.GetAll()

		assert.Nil(t, err.Err)
		assert.Len(t, result, 2)
		assert.Equal(t, result[1].CardNumberId, input[1].CardNumberId)

		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceGetOne(t *testing.T) {
	t.Run("should return one buyers", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := buyers.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "José",
			LastName:     "Silva",
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(input, nil)

		service := buyers.NewService(mockedRepository)

		result, err := service.GetOne(1)

		assert.Nil(t, err.Err)
		assert.Equal(t, result, input)

		mockedRepository.AssertExpectations(t)

	})

	t.Run("Verify the error case id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("Buyer with id 1 not found")

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, expectedError)

		service := buyers.NewService(mockedRepository)
		_, err := service.GetOne(1)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Code, http.StatusNotFound)
		assert.Equal(t, expectedError, err.Err)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("Verify if the buyer was deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(nil)
		service := buyers.NewService(mockedRepository)

		result := service.Delete(1)

		assert.Equal(t, result.Code, http.StatusNoContent)

		mockedRepository.AssertExpectations(t)

	})

	t.Run("Verify the error case if buyer do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("Buyer with id 1 not found")

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(expectedError)

		service := buyers.NewService(mockedRepository)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNotFound)
		assert.Equal(t, result.Err, expectedError)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("Return the updated buyer when successfully updated", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"card_number_id": "12345",
			"first_name":     "Fulano",
			"last_name":      "Beltrano",
		}

		expectedBuyer := buyers.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "Fulano",
			LastName:     "Beltrano",
		}

		input := buyers.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "José",
			LastName:     "Silva",
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input, nil)

		mockedRepository.On("GetAll").
			Return([]buyers.Buyer{}, nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(expectedBuyer, nil)

		service := buyers.NewService(mockedRepository)
		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)
		assert.Equal(t, expectedBuyer, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Return null when buyer id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("buyer not found")
		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(buyers.Buyer{}, expectedError).Once()

		mockedRepository.On("GetAll").
			Return([]buyers.Buyer{}, expectedError).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(buyers.Buyer{}, nil).Once()

		service := buyers.NewService(mockedRepository)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Code, http.StatusNotFound)
		assert.Equal(t, expectedError, err.Err)
	})

	t.Run("return error when card_number_id already exists and id doesn't match", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		requestData := map[string]interface{}{
			"card_number_id": "12345",
		}
		expectedError := errors.New("buyer number already exists")
		input := []buyers.Buyer{
			{
				Id:           1,
				CardNumberId: "12345",
				FirstName:    "José",
				LastName:     "Silva",
			},
			{
				Id:           2,
				CardNumberId: "12344",
				FirstName:    "José",
				LastName:     "Silva",
			},
		}
		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input[1], nil).Once()
		mockedRepository.On("GetAll").
			Return(input, nil).Once()
		mockedRepository.On("Update", mock.AnythingOfType("int"), mock.Anything).Return(buyers.Buyer{}, nil).Once()
		service := buyers.NewService(mockedRepository)
		_, err := service.Update(2, requestData)
		assert.NotNil(t, err.Err)
		assert.Equal(t, expectedError.Error(), err.Err.Error())
		assert.Equal(t, http.StatusConflict, err.Code)
	})
}
