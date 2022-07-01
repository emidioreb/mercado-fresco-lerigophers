package sections_test

import (
	"errors"

	"net/http"

	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

var inputSections = []sections.Section{
	{
		Id:                 1,
		SectionNumber:      10,
		CurrentTemperature: 25,
		MinimumTemperature: 0,
		CurrentCapacity:    130,
		MininumCapacity:    50,
		MaximumCapacity:    999,
		WarehouseId:        55,
		ProductTypeId:      70},
	{
		Id:                 2,
		SectionNumber:      11,
		CurrentTemperature: 26,
		MinimumTemperature: 1,
		CurrentCapacity:    131,
		MininumCapacity:    51,
		MaximumCapacity:    1000,
		WarehouseId:        56,
		ProductTypeId:      71},
	{
		Id:                 3,
		SectionNumber:      12,
		CurrentTemperature: 27,
		MinimumTemperature: 2,
		CurrentCapacity:    132,
		MininumCapacity:    52,
		MaximumCapacity:    1001,
		WarehouseId:        57,
		ProductTypeId:      72,
	},
}

var (
	errNotFound      = errors.New("section with id 1 not found")
	errAlreadyExists = errors.New("section number already exists")
)

func TestServiceCreate(t *testing.T) {

	t.Run("Test if create successfully", func(t *testing.T) {

		mockedRepository := new(mocks.Repository)

		input := inputSections[0]

		mockedRepository.On("GetAll").Return([]sections.Section{}, nil)

		mockedRepository.On("Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := sections.NewService(mockedRepository)

		result, err := service.Create(input.SectionNumber, input.CurrentTemperature, input.MinimumTemperature, input.CurrentCapacity, input.MininumCapacity, input.MaximumCapacity, input.WarehouseId, input.ProductTypeId)
		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)

		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test error case if section Section Number already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := inputSections[0]

		expectedError := errAlreadyExists
		listSections := []sections.Section{}
		listSections = append(listSections, input)

		mockedRepository.On("GetAll").Return(listSections, nil)

		mockedRepository.On("Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(sections.Section{}, expectedError)

		service := sections.NewService(mockedRepository)

		_, err := service.Create(input.SectionNumber, input.CurrentTemperature, input.MinimumTemperature, input.CurrentCapacity, input.MininumCapacity, input.MaximumCapacity, input.WarehouseId, input.ProductTypeId)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err.Error(), expectedError.Error())

		assert.Equal(t, err.Code, http.StatusConflict)
	})

}

func TestServiceDelete(t *testing.T) {

	t.Run("Verify the successfully case if the section is deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(nil)

		service := sections.NewService(mockedRepository)

		result := service.Delete(1)

		assert.Nil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNoContent)

		mockedRepository.AssertExpectations(t)

	})

	t.Run("Verify the error case if section do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errNotFound

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(expectedError)

		service := sections.NewService(mockedRepository)

		result := service.Delete(1)

		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNotFound)

		assert.Equal(t, result.Err, expectedError)

		mockedRepository.AssertExpectations(t)
	})

}

func TestServiceGetAll(t *testing.T) {

	t.Run("Test if an array of sections is returned when GetAll", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		input := inputSections

		mockedRepository.On("GetAll").Return(input, nil)

		service := sections.NewService(mockedRepository)

		result, err := service.GetAll()

		assert.Nil(t, err.Err)

		assert.Len(t, result, 3)

		assert.Equal(t, input[1].SectionNumber, result[1].SectionNumber)

		mockedRepository.AssertExpectations(t)
	})

}

func TestServiceGetOne(t *testing.T) {

	t.Run("Test if section is returned based on valid id", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		input := inputSections[0]

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(input, nil)

		service := sections.NewService(mockedRepository)

		result, err := service.GetOne(1)

		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)

		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test if requested section do not exists based on invalid id", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errNotFound

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(sections.Section{}, expectedError)

		service := sections.NewService(mockedRepository)

		_, err := service.GetOne(1)

		assert.NotNil(t, err.Err)

		assert.Equal(t, http.StatusNotFound, err.Code)

		assert.Equal(t, expectedError, err.Err)

		mockedRepository.AssertExpectations(t)

	})

}

func TestServiceUpdate(t *testing.T) {

	t.Run("Return the updated section when successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"section_number":      10,
			"current_temperature": 25,
			"minimum_temperature": 0,
			"current_capacity":    130,
			"minimum_capacity":    50,
			"maximum_capacity":    999,
			"warehouse_id":        55,
			"product_type_id":     70,
		}

		expectedSection := inputSections[0]

		input := sections.Section{
			Id:                 1,
			SectionNumber:      11,
			CurrentTemperature: 26,
			MinimumTemperature: 1,
			CurrentCapacity:    131,
			MininumCapacity:    51,
			MaximumCapacity:    1001,
			WarehouseId:        56,
			ProductTypeId:      71,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input, nil)

		mockedRepository.On("GetAll").
			Return([]sections.Section{}, nil)

		mockedRepository.On("Update",

			mock.AnythingOfType("int"),

			mock.Anything,
		).Return(expectedSection, nil)

		service := sections.NewService(mockedRepository)

		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)

		assert.Equal(t, expectedSection, result)

		mockedRepository.AssertExpectations(t)
	})

	t.Run("Return null when section id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errNotFound

		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(sections.Section{}, expectedError).Once()

		mockedRepository.On("GetAll").
			Return([]sections.Section{}, nil).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(sections.Section{}, nil).Once()

		service := sections.NewService(mockedRepository)

		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)

		assert.Equal(t, http.StatusNotFound, err.Code)

		assert.Equal(t, expectedError, err.Err)
	})

	t.Run("Return error when section_number already exists and id doesn't match", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		requestData := map[string]interface{}{
			"section_number": 10.0,
		}
		expectedError := errAlreadyExists
		input := inputSections

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input[1], nil).Once()
		mockedRepository.On("GetAll").
			Return(input, nil).Once()
		mockedRepository.On("Update", mock.AnythingOfType("int"), mock.Anything).Return(sections.Section{}, nil).Once()

		service := sections.NewService(mockedRepository)
		_, err := service.Update(2, requestData)
		assert.NotNil(t, err.Err)
		assert.Equal(t, expectedError.Error(), err.Err.Error())
		assert.Equal(t, http.StatusConflict, err.Code)
	})
}
