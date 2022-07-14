package warehouses

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateWarehouse(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateWarehouse)).
			WithArgs(
				"32dsa1",
				"Rua das pedras",
				"8213312123",
				1,
				2,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		warehouseRepo := NewMariaDbRepository(db)

		warehouse, err := warehouseRepo.Create(
			"32dsa1",
			"Rua das pedras",
			"8213312123",
			1,
			2,
		)
		assert.Nil(t, err)

		assert.Equal(t, 1, warehouse.MinimumCapacity)
	})

	t.Run("Insert to exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryCreateWarehouse)).WillReturnError(errors.New("internal db error"))
		warehouseRepo := NewMariaDbRepository(db)

		_, err = warehouseRepo.Create("", "", "", 0, 0)
		assert.Error(t, err)
	})

	t.Run("Insert to return lastId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("ocurred an error to create warehouse"))
		mock.ExpectExec(regexp.QuoteMeta(queryCreateWarehouse)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		warehouseRepo := NewMariaDbRepository(db)
		_, err = warehouseRepo.Create("", "", "", 0, 0)
		assert.Error(t, err)
		assert.Equal(t, "ocurred an error to create warehouse", err.Error())
	})
}

func TestDBGetOneWarehouse(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"warehouse_code",
			"adress",
			"telephone",
			"minimum_temperature",
			"minimum_capacity",
		}).AddRow(1, "1", "3413412dasd", "2132312312", 2, 2)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneWarehouse)).WillReturnRows(rows)

		warehouseRepo := NewMariaDbRepository(db)
		warehouse, err := warehouseRepo.GetOne(1)
		assert.NoError(t, err)

		assert.Equal(t, 2, warehouse.MinimumTemperature)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneWarehouse)).WillReturnError(sql.ErrNoRows)

		warehouseRepo := NewMariaDbRepository(db)
		_, err = warehouseRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "warehouse with id 1 not found", err.Error())
	})

	t.Run("DB Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneWarehouse)).
			WillReturnError(errors.New("unexpected error to get warehouse"))

		warehouseRepo := NewMariaDbRepository(db)
		_, err = warehouseRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to get warehouse", err.Error())
	})
}

func TestDBGetAllWarehouses(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"warehouse_code",
			"address",
			"telephone",
			"minimum_temperature",
			"minimum_capacity",
		}).
			AddRow(1, "aa", "aa", "13123", 1, 10).
			AddRow(2, "bb", "bb", "123213", 1, 10).
			AddRow(3, "cc", "cc", "1123123", 1, 10)
		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllWarehouses)).WillReturnRows(rows)

		warehousesRepo := NewMariaDbRepository(db)

		warehouseReports, err := warehousesRepo.GetAll()
		assert.NoError(t, err)

		assert.Len(t, warehouseReports, 3)
		assert.Equal(t, 1, warehouseReports[0].MinimumTemperature)
		assert.Equal(t, 1, warehouseReports[1].MinimumTemperature)
		assert.Equal(t, 1, warehouseReports[2].MinimumTemperature)
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllWarehouses)).
			WillReturnError(errors.New("couldn't get warehouses"))

		warehouseRepo := NewMariaDbRepository(db)

		_, err = warehouseRepo.GetAll()
		assert.Error(t, err)
	})
}

func TestDBDeleteSection(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteWarehouse)).WithArgs(int64(1)).
			WillReturnResult(affectedRowsResult)

		sectionsRepo := NewMariaDbRepository(db)
		err = sectionsRepo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewErrorResult(sql.ErrNoRows)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteWarehouse)).
			WillReturnResult(affectedRowsResult)

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.NotNil(t, err)
		assert.Equal(t, "warehouse with id 1 not found", err.Error())
	})

	t.Run("db Exec case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteWarehouse)).
			WillReturnError(errors.New("any error"))

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "any error", err.Error())
	})

	t.Run("DB Error Case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteWarehouse)).
			WillReturnError(errors.New("unexpected error to delete warehouse"))

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to delete warehouse", err.Error())
	})
}

func TestDBUpdateWarehouse(t *testing.T) {
	requestData := map[string]interface{}{
		"warehouse_code":      "212",
		"address":             "rua do bobo",
		"telephone":           "0",
		"minimum_temperature": 30,
		"minimum_capacity":    10.0,
	}

	query := `	UPDATE warehouses SET 
						warehouse_code = ?,
						address = ?,
						telephone = ?,
						minimum_temperature = ?,
						minimum_capacity = ?
					WHERE id = ?
			 `

	t.Run("Exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(errors.New("any error"))

		sectionsRepo := NewMariaDbRepository(db)

		_, err = sectionsRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedWarehouse, err)
	})
}
