package product_batches_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var duedate, _ = time.Parse("2006-02-01", "2021-04-01")
var manufacturingdate, _ = time.Parse("2006-02-01", "2021-04-01")

var fakeProductBatches = []product_batches.ProductBatches{{
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            duedate,
	ManufacturingDate:  manufacturingdate,
}, {
	BatchNumber:        1,
	CurrentQuantity:    10,
	CurrentTemperature: 2,
	InitialQuantity:    500,
	ManufacturingHour:  10,
	MinimumTemperature: 890,
	ProductId:          23,
	SectionId:          56,
	DueDate:            duedate,
	ManufacturingDate:  manufacturingdate,
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
