package localities_test

import (
	"errors"
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
}
