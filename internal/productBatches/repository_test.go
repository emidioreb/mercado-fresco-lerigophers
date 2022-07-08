package product_batches_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	"github.com/stretchr/testify/assert"
)

var mockProductBatch = product_batches.ProductBatches{
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
}

func TestCreate(t *testing.T) {

	query := `INSERT INTO product_batches (batch_number, current_quatity, current_temperature, initial_quantity, manufacturing_hour, minimum_temperature, product_id, section_id, due_date, manufacturing_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				mockProductBatch.BatchNumber,
				mockProductBatch.CurrentQuantity,
				mockProductBatch.CurrentTemperature,
				mockProductBatch.InitialQuantity,
				mockProductBatch.ManufacturingHour,
				mockProductBatch.MinimumTemperature,
				mockProductBatch.ProductId,
				mockProductBatch.SectionId,
				mockProductBatch.DueDate,
				mockProductBatch.ManufacturingDate,
			).WillReturnResult(sqlmock.NewResult(1, 1)) // last id, // rows affected

		productBatchRepo := product_batches.NewMariaDbRepository(db)

		pb, err := productBatchRepo.CreateProductBatch(
			mockProductBatch.BatchNumber,
			mockProductBatch.CurrentQuantity,
			mockProductBatch.CurrentTemperature,
			mockProductBatch.InitialQuantity,
			mockProductBatch.ManufacturingHour,
			mockProductBatch.MinimumTemperature,
			mockProductBatch.ProductId,
			mockProductBatch.SectionId,
			mockProductBatch.DueDate,
			mockProductBatch.ManufacturingDate)
		assert.NoError(t, err)

		expectedCurrentTemperature := 2

		assert.Equal(t, expectedCurrentTemperature, pb.CurrentTemperature)
	})

	t.Run("failed to create", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(0, 0, 0, 0, 0, 0, 0, 0).
			WillReturnResult(sqlmock.NewResult(1, 1)) // last id, // rows affected

		productBatchRepo := product_batches.NewMariaDbRepository(db)
		_, err = productBatchRepo.CreateProductBatch(
			mockProductBatch.BatchNumber,
			mockProductBatch.CurrentQuantity,
			mockProductBatch.CurrentTemperature,
			mockProductBatch.InitialQuantity,
			mockProductBatch.ManufacturingHour,
			mockProductBatch.MinimumTemperature,
			mockProductBatch.ProductId,
			mockProductBatch.SectionId,
			mockProductBatch.DueDate,
			mockProductBatch.ManufacturingDate)

		assert.Error(t, err)
	})
}

func TestGetOne(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"batch_number",
			"current_quantity",
			"current_temperature",
			"initial_quantity",
			"manufacturing_hour",
			"minimum_temperature",
			"product_id",
			"section_id",
			"due_date",
			"manufacturing_date",
		}).AddRow(1, 1, 1, 1, 1, 1, 1, 1, 1, duedate, manufacturingdate)

		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetOneProductBatch)).WillReturnRows(rows)

		productBatchRepo := product_batches.NewMariaDbRepository(db)

		pb, err := productBatchRepo.GetOne(1)
		assert.Nil(t, err)

		assert.Equal(t, 1, pb.Id)
		assert.Equal(t, 1, pb.CurrentQuantity)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetOneProductBatch)).WillReturnError(sql.ErrNoRows)

		productBatchRepo := product_batches.NewMariaDbRepository(db)

		_, err = productBatchRepo.GetOne(1)
		assert.NotNil(t, err)

		assert.Equal(t, "product_batch with batch_number 1 not found", err.Error())
	})

	t.Run("Another error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetOneProductBatch)).WillReturnError(errors.New(""))

		productBatchRepo := product_batches.NewMariaDbRepository(db)
		_, err = productBatchRepo.GetOne(1)
		assert.NotNil(t, err)
		assert.Equal(t, "error to find product_batch", err.Error())
	})
}

func TestDBGetReportSellers(t *testing.T) {
	t.Run("Get all reports", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"section_id",
			"section_number",
			"products_count",
		}).
			AddRow(1, 1, 1).
			AddRow(2, 2, 2).
			AddRow(3, 3, 3)

		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetReportAll)).WillReturnRows(rows)

		productBatchRepo := product_batches.NewMariaDbRepository(db)

		sectionReports, err := productBatchRepo.GetReportSection(0)
		assert.Nil(t, err)

		assert.Len(t, sectionReports, 3)
		assert.Equal(t, sectionReports[0].SectionId, 1)
		assert.Equal(t, sectionReports[1].SectionId, 2)
		assert.Equal(t, sectionReports[2].SectionId, 3)
	})

	t.Run("Get one report", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"section_id",
			"section_number",
			"products_count",
		}).
			AddRow(1, 1, 1)

		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetReportOne)).WillReturnRows(rows)

		productBatchRepo := product_batches.NewMariaDbRepository(db)

		sectionReports, err := productBatchRepo.GetReportSection(1)
		assert.Nil(t, err)

		assert.Len(t, sectionReports, 1)
		assert.Equal(t, 1, sectionReports[0].SectionId)
	})

	t.Run("Error to get report - case query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetReportOne)).WillReturnError(errors.New(""))

		sectionReports := product_batches.NewMariaDbRepository(db)

		_, err = sectionReports.GetReportSection(1)
		assert.NotNil(t, err)
		assert.Equal(t, "error to report sections by product_batches", err.Error())
	})

	t.Run("Error to get report - case scan", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"sellers_count",
		}).
			AddRow(1, "s", 1)

		mock.ExpectQuery(regexp.QuoteMeta(product_batches.QueryGetReportOne)).WithArgs(1).WillReturnRows(rows)

		ProductBatches := product_batches.NewMariaDbRepository(db)

		_, err = ProductBatches.GetReportSection(1)
		assert.NotNil(t, err)

		assert.Equal(t, "error to report sections by product_batches", err.Error())
	})
}
