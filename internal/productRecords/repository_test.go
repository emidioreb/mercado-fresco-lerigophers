package product_records

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductRecord(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		layout := "2006-01-02"
		lastUpdateDate, err := time.Parse(layout, "2022-02-07")

		mock.ExpectExec(regexp.QuoteMeta(queryCreateProductRecord)).
			WithArgs(
				lastUpdateDate,
				3.0,
				4.0,
				4,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		productRepo := NewMariaDbRepository(db)

		productsRecords, err := productRepo.CreateProductRecord(
			"2022-02-07",
			3.0,
			4.0,
			4,
		)

		assert.NoError(t, err)
		assert.Equal(t, "2022-02-07", productsRecords.LastUpdateDate)
	})

	t.Run("Invalid date input", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		layout := "2006-01-02"
		lastUpdateDate, err := time.Parse(layout, "2022-30-07")

		mock.ExpectExec(regexp.QuoteMeta(queryCreateProductRecord)).
			WithArgs(
				lastUpdateDate,
				3.0,
				4.0,
				4,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		productRepo := NewMariaDbRepository(db)

		_, err = productRepo.CreateProductRecord(
			"2022-30-07",
			3.0,
			4.0,
			4,
		)

		assert.Error(t, err)
	})

	t.Run("Couldn`t create a product record", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateProductRecord)).
			WithArgs(
				"",
				"",
				"",
				"",
			).WillReturnResult(sqlmock.NewResult(1, 1))

		productRepo := NewMariaDbRepository(db)

		_, err = productRepo.CreateProductRecord(
			"2022-02-07",
			3.0,
			4.0,
			0,
		)

		assert.Error(t, err)
	})

	t.Run("Last insert id error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("ocurred an error to create product record"))
		mock.ExpectExec(regexp.QuoteMeta(queryCreateProductRecord)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		productRepo := NewMariaDbRepository(db)
		_, err = productRepo.CreateProductRecord("2022-02-07", 0.0, 0.0, 0)

		assert.Error(t, err)
		assert.Equal(t, "ocurred an error to create product record", err.Error())
	})
}

func TestDBGetOneProductRecord(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"last_update_date",
			"purchase_price",
			"sale_price",
			"product_id",
		}).AddRow(1, "2022-02-07", 3.0, 4.0, 4)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneProductRecord)).WillReturnRows(rows)

		productRepo := NewMariaDbRepository(db)
		err = productRepo.GetOne(1)
		assert.NoError(t, err)
	})
}
