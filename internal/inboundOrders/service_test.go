package inboundorders_test

import (
	"errors"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	employeeRepository "github.com/emidioreb/mercado-fresco-lerigophers/internal/employees/mocks"
	inboundorders "github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundorders/mocks"
	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	productBatchesRepository "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	warehouseRepository "github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fakeInbounds = []inboundorders.InboundOrder{
	{
		Id:             1,
		OrderNumber:    "43",
		OrderDate:      "2006-01-02",
		EmployeeId:     1,
		ProductBatchId: 1,
		WarehouseId:    1,
	},
	{
		Id:             1,
		OrderNumber:    "434",
		OrderDate:      "2006-01-02",
		EmployeeId:     2,
		ProductBatchId: 1,
		WarehouseId:    2,
	},
}

var fakeReports = []inboundorders.ReportInboundOrder{
	{
		Id:                 1,
		CardNumberId:       "456",
		FirstName:          "Iuri",
		LastName:           "Oi",
		WarehouseId:        1,
		InboundOrdersCount: 1,
	},
	{
		Id:                 2,
		CardNumberId:       "4565",
		FirstName:          "Iurizin",
		LastName:           "Tchau",
		WarehouseId:        2,
		InboundOrdersCount: 1,
	},
}

func TestServiceCreate(t *testing.T) {
	t.Run("Test if create successfully", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, nil)

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, nil)

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, nil)

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(fakeInbounds[0], nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		result, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.Nil(t, err.Err)

		assert.Equal(t, fakeInbounds[0], result)
	})

	t.Run("Test if employee dont exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, errors.New("employee with id 1 not found"))

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, nil)

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, nil)

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("employee with id 1 not found"))
	})

	t.Run("Test if employee error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, errors.New("unexpected error to get employee"))

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, nil)

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, nil)

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("unexpected error to get employee"))
	})

	t.Run("Test if warehouse dont exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, nil)

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, errors.New("warehouse with id 1 not found"))

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, nil)

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("warehouse with id 1 not found"))
	})

	t.Run("Test if warehouse error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, nil)

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, errors.New("error"))

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, nil)

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("error"))
	})

	t.Run("Test if productbatches exists", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, nil)

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, nil)

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, errors.New("product_batch with batch_number 1 not found"))

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("product_batch with batch_number 1 not found"))
	})

	t.Run("Test if productbatches error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, nil)

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, nil)

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, errors.New("error"))

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("error"))
	})

	t.Run("Test if create error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		employeeRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(employees.Employee{}, nil)

		warehouseRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(warehouses.Warehouse{}, nil)

		productBatcheRepo.On(
			"GetOne",
			mock.AnythingOfType("int"),
		).Return(product_batches.ProductBatches{}, nil)

		mockedRepository.On(
			"CreateInboundOrders",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(inboundorders.InboundOrder{}, errors.New("error"))

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.CreateInboundOrders(
			fakeInbounds[0].OrderNumber,
			fakeInbounds[0].OrderDate,
			fakeInbounds[0].EmployeeId,
			fakeInbounds[0].ProductBatchId,
			fakeInbounds[0].WarehouseId,
		)

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("error"))
	})
}

func TestServiceGet(t *testing.T) {
	t.Run("Test if getreport success", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		mockedRepository.On(
			"GetReportInboundOrders",
			mock.AnythingOfType("string"),
		).Return(fakeReports, nil)

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		result, err := service.GetReportInboundOrders("1")

		assert.Nil(t, err.Err)

		assert.Equal(t, fakeReports, result)
	})

	t.Run("Test if getreport error", func(t *testing.T) {
		mockedRepository := new(mocks.Repository)
		employeeRepo := new(employeeRepository.Repository)
		warehouseRepo := new(warehouseRepository.Repository)
		productBatcheRepo := new(productBatchesRepository.Repository)

		mockedRepository.On(
			"GetReportInboundOrders",
			mock.AnythingOfType("string"),
		).Return([]inboundorders.ReportInboundOrder{}, errors.New("error"))

		service := inboundorders.NewService(mockedRepository, warehouseRepo, employeeRepo, productBatcheRepo)
		_, err := service.GetReportInboundOrders("1")

		assert.NotNil(t, err.Err)

		assert.Equal(t, err.Err, errors.New("error"))
	})
}
