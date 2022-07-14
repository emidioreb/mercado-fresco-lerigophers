package product_records_test

import (
	"errors"
	"testing"

	product_records "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords"
	mockedProductRecords "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	mockedProducts "github.com/emidioreb/mercado-fresco-lerigophers/internal/products/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceCreate(t *testing.T) {
	t.Run("Test if create sucessfully", func(t *testing.T) {
		mockedRepository := new(mockedProductRecords.Repository)
		mockedProdutsRepository := new(mockedProducts.Repository)

		input := product_records.ProductRecords{
			LastUpdateDate: "2022-02-07",
			PurchasePrice:  3.0,
			SalePrice:      4.0,
			ProductId:      1,
		}

		mockedProdutsRepository.On("GetOne", mock.AnythingOfType("int")).Return(products.Product{}, nil)
		mockedRepository.On("CreateProductRecord",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(input, nil)

		service := product_records.NewService(mockedRepository, mockedProdutsRepository)

		result, err := service.CreateProductRecord(input.LastUpdateDate,
			input.PurchasePrice, input.SalePrice, input.ProductId)

		assert.Nil(t, err.Err)

		assert.Equal(t, input, result)
		mockedProdutsRepository.AssertExpectations(t)
	})

	t.Run("Invalid product_id", func(t *testing.T) {
		mockedRepository := new(mockedProductRecords.Repository)
		mockedProdutsRepository := new(mockedProducts.Repository)

		input := product_records.ProductRecords{
			LastUpdateDate: "2022-02-07",
			PurchasePrice:  3.0,
			SalePrice:      4.0,
			ProductId:      1,
		}

		mockedProdutsRepository.On("GetOne", mock.AnythingOfType("int")).Return(products.Product{}, nil)
		mockedRepository.On("CreateProductRecord",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("float64"),
			mock.AnythingOfType("int"),
		).Return(input, errors.New("product_id do not exists"))

		service := product_records.NewService(mockedRepository, mockedProdutsRepository)

		_, err := service.CreateProductRecord(input.LastUpdateDate,
			input.PurchasePrice, input.SalePrice, 2)

		assert.NotNil(t, err.Err)
	})
}
