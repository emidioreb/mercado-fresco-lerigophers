package carriers_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {
	input := carriers.Carry{
		Cid:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		LocalityId:  "456",
	}

	t.Run("should return create carriers", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetOne",
			mock.AnythingOfType("string")).Return(carriers.Carry{}, errors.New(""))

		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(input, nil)

		service := carriers.NewService(mockedRepository)

		result, err := service.Create(input.Cid, input.CompanyName, input.Address, input.Telephone, input.LocalityId)

		assert.Nil(t, err.Err)
		assert.Equal(t, result, input)

		mockedRepository.AssertExpectations(t)

	})

	t.Run("CID already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetOne",
			mock.AnythingOfType("string")).Return(carriers.Carry{}, nil)

		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(input, errors.New("CID already exists"))

		service := carriers.NewService(mockedRepository)

		_, resp := service.Create(input.Cid, input.CompanyName, input.Address, input.Telephone, input.LocalityId)

		assert.NotNil(t, resp.Err)
		assert.Equal(t, resp.Code, http.StatusConflict)

	})

	t.Run("CarryResult should return error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetOne",
			mock.AnythingOfType("string")).Return(carriers.Carry{}, errors.New(""))
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).Return(carriers.Carry{}, errors.New("StatusInternalServerError"))

		service := carriers.NewService(mockedRepository)
		_, resp := service.Create(input.Cid, input.CompanyName, input.Address, input.Telephone, input.LocalityId)

		assert.NotNil(t, resp.Err)
		assert.Equal(t, resp.Code, http.StatusInternalServerError)
	})

}
