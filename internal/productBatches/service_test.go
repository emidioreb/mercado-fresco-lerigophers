package product_batches_test

import (
	"errors"
	"net/http"
	"testing"

	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fakeProductBatches = []product_batches.ProductBatches{{
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            date,
	ManufacturingDate:  date,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            date,
	ManufacturingDate:  date,
}}

func TestServiceCreate(t *testing.T) {
	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(product_batches.ProductBatches{}, errors.New(""))
		mockedRepository.On("CreateProductBatch",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time")).Return(fakeProductBatches[0], nil)

		service := product_batches.NewService(mockedRepository)

		result, err := service.CreateProductBatch(
			fakeProductBatches[0].BatchNumber,
			fakeProductBatches[0].CurrentQuantity,
			fakeProductBatches[0].CurrentTemperature,
			fakeProductBatches[0].InitialQuantity,
			fakeProductBatches[0].ManufacturingHour,
			fakeProductBatches[0].MinimumTemperature,
			fakeProductBatches[0].ProductId,
			fakeProductBatches[0].SectionId,
			fakeProductBatches[0].DueDate,
			fakeProductBatches[0].ManufacturingDate)
		assert.Nil(t, err.Err)

		assert.Equal(t, fakeProductBatches[0], result)
	})

	t.Run("product_batch already exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(product_batches.ProductBatches{}, nil)

		service := product_batches.NewService(mockedRepository)

		_, err := service.CreateProductBatch(
			fakeProductBatches[0].BatchNumber,
			fakeProductBatches[0].CurrentQuantity,
			fakeProductBatches[0].CurrentTemperature,
			fakeProductBatches[0].InitialQuantity,
			fakeProductBatches[0].ManufacturingHour,
			fakeProductBatches[0].MinimumTemperature,
			fakeProductBatches[0].ProductId,
			fakeProductBatches[0].SectionId,
			fakeProductBatches[0].DueDate,
			fakeProductBatches[0].ManufacturingDate)

		assert.NotNil(t, err)
		assert.Equal(t, err.Code, http.StatusConflict)
	})

	t.Run("ProductBatchResult should return error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)

		mockedRepository.On("GetOne", mock.AnythingOfType("int")).Return(product_batches.ProductBatches{}, errors.New(""))
		mockedRepository.On("CreateProductBatch",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time")).Return(product_batches.ProductBatches{}, errors.New("couldn't create a product_batch"))

		service := product_batches.NewService(mockedRepository)
		_, err := service.CreateProductBatch(
			fakeProductBatches[0].BatchNumber,
			fakeProductBatches[0].CurrentQuantity,
			fakeProductBatches[0].CurrentTemperature,
			fakeProductBatches[0].InitialQuantity,
			fakeProductBatches[0].ManufacturingHour,
			fakeProductBatches[0].MinimumTemperature,
			fakeProductBatches[0].ProductId,
			fakeProductBatches[0].SectionId,
			fakeProductBatches[0].DueDate,
			fakeProductBatches[0].ManufacturingDate)
		assert.NotNil(t, err.Err)
		assert.Equal(t, err.Code, http.StatusInternalServerError)
	})
}

func TestServiceGetReport(t *testing.T) {
	t.Run("get report - success case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetReportSection", mock.AnythingOfType("int")).Return([]product_batches.ProductsQuantity{}, nil)
		service := product_batches.NewService(mockedRepository)

		result, err := service.GetReportSection(0)
		assert.NoError(t, err.Err)

		assert.Equal(t, result, []product_batches.ProductsQuantity{})
		assert.Equal(t, http.StatusOK, err.Code)
	})

	t.Run("get report - error case", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		mockedRepository.On("GetReportSection", mock.AnythingOfType("int")).Return([]product_batches.ProductsQuantity{}, errors.New("error to report sections by product_batches"))
		service := product_batches.NewService(mockedRepository)

		result, err := service.GetReportSection(0)
		assert.NotNil(t, err.Err)

		assert.Equal(t, result, []product_batches.ProductsQuantity{})
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "error to report sections by product_batches", err.Err.Error())
	})
}
