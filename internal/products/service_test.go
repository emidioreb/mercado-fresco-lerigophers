package products_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {
	t.Run("Test if create sucessfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := products.Product{
			Id:                             1,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}

		mockedRepository.On("GetAll").Return([]products.Product{}, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := products.NewService(mockedRepository)

		result, err := service.Create(input.ProductCode, input.Description,
			input.Width, input.Height, input.Length, input.NetWeight,
			input.ExpirationRate, input.RecommendedFreezingTemperature,
			input.FreezingRate, input.ProductTypeId)
		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test error case if product_code already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := []products.Product{{
			Id:                             1,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}, {
			Id:                             2,
			ProductCode:                    "FK0004",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}}

		expectedError := errors.New("Product_code already exists")

		mockedRepository.On("GetAll").Return(input, nil)
		mockedRepository.On("Create",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(products.Product{}, expectedError)

		service := products.NewService(mockedRepository)

		_, err := service.Create(input[0].ProductCode, input[0].Description,
			input[0].Width, input[0].Height, input[0].Length, input[0].NetWeight,
			input[0].ExpirationRate, input[0].RecommendedFreezingTemperature,
			input[0].FreezingRate, input[0].ProductTypeId)

		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Err.Error(), expectedError.Error())
		assert.Equal(t, err.Code, http.StatusConflict)
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("Verify the sucessfully case if the product is deleted", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("Delete", mock.AnythingOfType("int")).Return(nil)

		service := products.NewService(mockedRepository)
		result := service.Delete(1)
		assert.Nil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNoContent)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Verify the error case if product do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("product with id 1 not found")

		mockedRepository.On("Delete",
			mock.AnythingOfType("int")).Return(expectedError)

		service := products.NewService(mockedRepository)
		result := service.Delete(1)
		assert.NotNil(t, result.Err)

		assert.Equal(t, result.Code, http.StatusNotFound)
		assert.Equal(t, result.Err, expectedError)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceGetAll(t *testing.T) {
	t.Run("Test if an array of products is returned when GetAll", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := []products.Product{
			{
				Id:                             1,
				ProductCode:                    "FK0003",
				Description:                    "Fake Product",
				Width:                          23,
				Height:                         62,
				Length:                         101,
				NetWeight:                      27,
				ExpirationRate:                 88,
				RecommendedFreezingTemperature: 17,
				FreezingRate:                   23,
				ProductTypeId:                  7,
			}, {
				Id:                             2,
				ProductCode:                    "FK0004",
				Description:                    "Fake Product",
				Width:                          23,
				Height:                         62,
				Length:                         101,
				NetWeight:                      27,
				ExpirationRate:                 88,
				RecommendedFreezingTemperature: 17,
				FreezingRate:                   23,
				ProductTypeId:                  7,
			}, {
				Id:                             3,
				ProductCode:                    "FK0005",
				Description:                    "Fake Product",
				Width:                          23,
				Height:                         62,
				Length:                         101,
				NetWeight:                      27,
				ExpirationRate:                 88,
				RecommendedFreezingTemperature: 17,
				FreezingRate:                   23,
				ProductTypeId:                  7,
			},
		}

		mockedRepository.On("GetAll").Return(input, nil)

		service := products.NewService(mockedRepository)

		result, err := service.GetAll()
		assert.Nil(t, err.Err)

		assert.Len(t, result, 3)
		assert.Equal(t, input[1].ProductCode, result[1].ProductCode)
		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceGetOne(t *testing.T) {
	t.Run("Test if product is returned based on valid id", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		input := products.Product{
			Id:                             1,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(input, nil)

		service := products.NewService(mockedRepository)

		result, err := service.GetOne(1)
		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Test if requested product do not exists based on ivalid id", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("product with id 1 not found")

		mockedRepository.On("GetOne",
			mock.AnythingOfType("int")).Return(products.Product{}, expectedError)

		service := products.NewService(mockedRepository)
		_, err := service.GetOne(1)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, expectedError, err.Err)

		mockedRepository.AssertExpectations(t)
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("Return the updated product when sucessfully updated", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"ProductCode":                    "FK0003",
			"Description":                    "Fake Product",
			"Width":                          23,
			"Height":                         62,
			"Length":                         101,
			"NetWeight":                      27,
			"ExpirationRate":                 88,
			"RecommendedFreezingTemperature": 17,
			"FreezingRate":                   23,
			"ProductTypeId":                  7,
		}

		expectedProduct := products.Product{
			Id:                             2,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}

		input := products.Product{
			Id:                             2,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(input, nil)

		mockedRepository.On("GetAll").Return([]products.Product{}, nil)

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(expectedProduct, nil)

		service := products.NewService(mockedRepository)
		result, err := service.Update(1, requestData)

		assert.Nil(t, err.Err)
		assert.Equal(t, expectedProduct, result)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("Return null when product id do not exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		expectedError := errors.New("product with id 1 not found")
		requestData := map[string]interface{}{}

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).
			Return(products.Product{}, expectedError).Once()

		mockedRepository.On("GetAll").
			Return([]products.Product{}, expectedError).Once()

		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(products.Product{}, nil).Once()

		service := products.NewService(mockedRepository)
		_, err := service.Update(1, requestData)

		assert.NotNil(t, err.Err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Equal(t, err.Err, expectedError)
	})

	t.Run("Test error case product_code already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		requestData := map[string]interface{}{
			"product_code": "FK0003",
		}

		input := []products.Product{{
			Id:                             1,
			ProductCode:                    "FK0003",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		}, {
			Id:                             2,
			ProductCode:                    "FK0004",
			Description:                    "Fake Product",
			Width:                          23,
			Height:                         62,
			Length:                         101,
			NetWeight:                      27,
			ExpirationRate:                 88,
			RecommendedFreezingTemperature: 17,
			FreezingRate:                   23,
			ProductTypeId:                  7,
		},
		}

		expectedError := errors.New("product_code already exists")

		mockedRepository.On("GetOne",
			mock.AnythingOfType("int")).Return(products.Product{}, nil).Once()
		mockedRepository.On("GetAll").Return(input, nil).Once()
		mockedRepository.On("Update",
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(products.Product{}, expectedError).Once()

		service := products.NewService(mockedRepository)

		_, err := service.Update(input[1].Id, requestData)
		assert.NotNil(t, err.Err)

		assert.Equal(t, expectedError.Error(), err.Err.Error())
		assert.Equal(t, http.StatusConflict, err.Code)
	})
}
