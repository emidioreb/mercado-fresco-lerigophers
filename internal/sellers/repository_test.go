package sellers

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateSeller(t *testing.T) {

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateSeller)).
			WithArgs(
				123,
				"Mercado Libre",
				"Nações Unidas",
				"999",
				"123",
			).WillReturnResult(sqlmock.NewResult(1, 1))

		sellersRepo := NewMariaDbRepository(db)

		seller, err := sellersRepo.Create(
			123,
			"Mercado Libre",
			"Nações Unidas",
			"999",
			"123",
		)
		assert.Nil(t, err)

		assert.Equal(t, "Mercado Libre", seller.CompanyName)
	})

	t.Run("Insert to exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryCreateSeller)).WillReturnError(errors.New("internal db error"))
		localitiesRepo := NewMariaDbRepository(db)

		_, err = localitiesRepo.Create(0, "", "", "", "")
		assert.Error(t, err)
	})

	t.Run("Insert to return lastId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("ocurred an error to create seller"))
		mock.ExpectExec(regexp.QuoteMeta(queryCreateSeller)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		localitiesRepo := NewMariaDbRepository(db)
		_, err = localitiesRepo.Create(0, "", "", "", "")

		assert.Error(t, err)
		assert.Equal(t, "ocurred an error to create seller", err.Error())
	})

}

func TestDBGetOneSeller(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"cid",
			"company_name",
			"address",
			"telephone",
			"locality_id",
		}).
			AddRow(1, 1, "Mercado Libre", "Nações Unidas", "999", "1")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSeller)).WillReturnRows(rows)

		sellersRepo := NewMariaDbRepository(db)
		seller, err := sellersRepo.GetOne(1)
		assert.NoError(t, err)

		assert.Equal(t, "Mercado Libre", seller.CompanyName)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSeller)).WillReturnError(sql.ErrNoRows)

		sellersRepo := NewMariaDbRepository(db)
		_, err = sellersRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "seller with id 1 not found", err.Error())
	})

	t.Run("DB Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSeller)).
			WillReturnError(errors.New("unexpected error to get seller"))

		sellersRepo := NewMariaDbRepository(db)
		_, err = sellersRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to get seller", err.Error())
	})
}

func TestDBGetAllSellers(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"cid",
			"company_name",
			"address",
			"telephone",
			"locality_id",
		}).
			AddRow(1, 1, "Mercado Libre", "Nações Unidas", "123", "1").
			AddRow(2, 2, "Mercado Solto", "Nações Desunidas", "456", "1").
			AddRow(3, 3, "Mercado Freedom", "Nações Unificadas", "789", "1")
		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllSellers)).WillReturnRows(rows)

		sellersRepo := NewMariaDbRepository(db)

		sellerReports, err := sellersRepo.GetAll()
		assert.NoError(t, err)

		assert.Len(t, sellerReports, 3)
		assert.Equal(t, "Mercado Libre", sellerReports[0].CompanyName)
		assert.Equal(t, "Mercado Solto", sellerReports[1].CompanyName)
		assert.Equal(t, "Mercado Freedom", sellerReports[2].CompanyName)
	})

	t.Run("Wrong type error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"cid",
			"company_name",
			"address",
			"telephone",
			"locality_id",
		}).AddRow(1, "aaa", "Mercado Libre", "Nações Unidas", "123", "1")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllSellers)).WillReturnRows(rows)

		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.GetAll()
		assert.Error(t, err)
		assert.Equal(t, "couldn't get sellers", err.Error())
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllSellers)).
			WillReturnError(errors.New("couldn't get sellers"))

		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.GetAll()
		assert.Error(t, err)
	})
}

func TestDBDeleteSeller(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSeller)).WithArgs(int64(1)).
			WillReturnResult(affectedRowsResult)

		sellersRepo := NewMariaDbRepository(db)
		err = sellersRepo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewErrorResult(sql.ErrNoRows)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSeller)).
			WillReturnResult(affectedRowsResult)

		sellersRepo := NewMariaDbRepository(db)
		err = sellersRepo.Delete(1)
		assert.NotNil(t, err)
		assert.Equal(t, "seller with id 1 not found", err.Error())
	})

	t.Run("db Exec case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSeller)).
			WillReturnError(errors.New("any error"))

		sellersRepo := NewMariaDbRepository(db)
		err = sellersRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "any error", err.Error())
	})

	t.Run("DB Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		resultErr := sqlmock.NewErrorResult(errors.New("any error"))
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSeller)).
			WillReturnResult(affectedRowsResult).
			WillReturnResult(resultErr)

		sellersRepo := NewMariaDbRepository(db)
		err = sellersRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to delete seller", err.Error())
	})
}

func TestDBUpdateSeller(t *testing.T) {
	requestData := map[string]interface{}{
		"cid":          1.0,
		"address":      "Nações Separadas",
		"telephone":    "123",
		"locality_id":  "1",
		"company_name": "Mercado Free",
	}

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		query := `UPDATE sellers
		SET 
			company_name = ?,
			address = ?,
			telephone = ?,
			locality_id = ?,
			cid = ?
		WHERE id = ?`

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				"Mercado Free",
				"Nações Separadas",
				"123",
				"1",
				1,
				1, // Seller ID to update
			).WillReturnResult(sqlmock.NewResult(1, 1))

		newRow := mock.
			NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(1, 1, "Mercado Libre", "", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSeller)).
			WithArgs(1).WillReturnRows(newRow)

		sellersRepo := NewMariaDbRepository(db)

		seller, err := sellersRepo.Update(1, requestData)
		assert.Nil(t, err)

		assert.Equal(t, "Mercado Libre", seller.CompanyName)
	})

	t.Run("Not found case", func(t *testing.T) {})

	t.Run("DB Error case", func(t *testing.T) {})
}

func TestDBFindByCID(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {})

	t.Run("Not found case", func(t *testing.T) {})
}
