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
		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.Create(0, "", "", "", "")
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

		sellersRepo := NewMariaDbRepository(db)
		_, err = sellersRepo.Create(0, "", "", "", "")

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

	t.Run("Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSeller)).
			WillReturnError(errors.New("any error"))

		sellersRepo := NewMariaDbRepository(db)
		err = sellersRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, errDeleteSeller, err)
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

	query := `UPDATE sellers
	SET 
		company_name = ?,
		address = ?,
		telephone = ?,
		locality_id = ?,
		cid = ?
	WHERE id = ?`

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

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
			NewRows([]string{
				"id",
				"cid",
				"company_name",
				"address",
				"telephone",
				"locality_id"}).
			AddRow(1, 1, "Mercado Libre", "", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSeller)).
			WithArgs(1).WillReturnRows(newRow)

		sellersRepo := NewMariaDbRepository(db)

		seller, err := sellersRepo.Update(1, requestData)
		assert.Nil(t, err)

		assert.Equal(t, "Mercado Libre", seller.CompanyName)
	})

	t.Run("Exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(errors.New("any error"))

		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedSeller, err)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		newResult := sqlmock.NewResult(0, 0)
		newErrorResult := sqlmock.NewErrorResult(errors.New("any error"))
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnResult(newResult).
			WillReturnResult(newErrorResult)

		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedSeller, err)
	})

	t.Run("Return updated seller case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		newResult := sqlmock.NewResult(0, 0)
		newErrorResult := sqlmock.NewErrorResult(errors.New("any error"))
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnResult(newResult)

		mock.ExpectExec(regexp.QuoteMeta(queryGetOneSeller)).
			WillReturnResult(newErrorResult)

		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedSeller, err)
	})
}

func TestDBFindByCID(t *testing.T) {
	t.Run("Already exists case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"cid",
		}).
			AddRow(1, 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryFindByCID)).WillReturnRows(rows)

		sellersRepo := NewMariaDbRepository(db)
		id, err := sellersRepo.FindByCID(1)
		assert.Error(t, err)
		assert.Greater(t, id, 0)
		assert.Equal(t, "cid already exists", err.Error())
	})

	t.Run("No rows case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryFindByCID)).
			WillReturnError(sql.ErrNoRows)

		sellersRepo := NewMariaDbRepository(db)
		id, err := sellersRepo.FindByCID(1)
		assert.NoError(t, err)
		assert.Equal(t, 0, id)
	})

	t.Run("Another erro to verify cid", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryFindByCID)).
			WillReturnError(errors.New("any error"))

		sellersRepo := NewMariaDbRepository(db)
		id, err := sellersRepo.FindByCID(1)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, "failed to verify if cid already exists", err.Error())
	})
}
