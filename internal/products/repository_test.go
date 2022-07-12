package products

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateProduct(t *testing.T) {

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateProduct)).
			WithArgs(
				"ABX0001",
				"abacaxi",
				1.0,
				2.0,
				3.0,
				4.0,
				5.0,
				6.0,
				7.0,
				4,
				1,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		productsRepo := NewMariaDbRepository(db)

		product, err := productsRepo.Create(
			"ABX0001",
			"abacaxi",
			1.0,
			2.0,
			3.0,
			4.0,
			5.0,
			6.0,
			7.0,
			4,
			1,
		)
		assert.Nil(t, err)
		assert.Equal(t, "ABX0001", product.ProductCode)
	})

	t.Run("Insert to exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryCreateProduct)).WillReturnError(errors.New("internal db error"))
		productsRepo := NewMariaDbRepository(db)

		_, err = productsRepo.Create("0", "", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0)
		assert.Error(t, err)
	})

	t.Run("Insert to return lastId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("ocurred an error to create product"))
		mock.ExpectExec(regexp.QuoteMeta(queryCreateProduct)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		productsRepo := NewMariaDbRepository(db)
		_, err = productsRepo.Create("0", "", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0)

		assert.Error(t, err)
		assert.Equal(t, "ocurred an error to create product", err.Error())
	})
}

func TestDBGetOneProduct(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"product_code",
			"description",
			"width",
			"height",
			"length",
			"net_weight",
			"expiration_rate",
			"recommended_freezing_temperature",
			"freezing_rate",
			"product_type_id",
			"seller_id",
		}).AddRow(1, "ABX0001", "abacaxi", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 4, 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneProduct)).WillReturnRows(rows)

		productRepo := NewMariaDbRepository(db)
		product, err := productRepo.GetOne(1)
		assert.NoError(t, err)

		assert.Equal(t, "ABX0001", product.ProductCode)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneProduct)).
			WillReturnError(errors.New("unexpected error to get product"))

		productsRepo := NewMariaDbRepository(db)
		_, err = productsRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to get product", err.Error())
	})

	t.Run("DB Error Case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneProduct)).
			WillReturnError(errors.New("unexpected error to get product"))

		productsRepo := NewMariaDbRepository(db)
		_, err = productsRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to get product", err.Error())
	})
}

func TestDBgetAllProducts(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"product_code",
			"description",
			"width",
			"height",
			"length",
			"net_weight",
			"expiration_rate",
			"recommended_freezing_temperature",
			"freezing_rate",
			"product_type_id",
			"seller_id",
		}).
			AddRow(1, "ABX0001", "abacaxi", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 4, 1).
			AddRow(2, "ABX0002", "Melao", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 4, 1).
			AddRow(3, "ABX0003", "Goiaba", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 4, 1)
		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllProducts)).WillReturnRows(rows)

		productsRepo := NewMariaDbRepository(db)

		productsReports, err := productsRepo.GetAll()
		assert.NoError(t, err)

		assert.Len(t, productsReports, 3)
		assert.Equal(t, "ABX0001", productsReports[0].ProductCode)
		assert.Equal(t, "ABX0002", productsReports[1].ProductCode)
		assert.Equal(t, "ABX0003", productsReports[2].ProductCode)
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllProducts)).
			WillReturnError(errors.New("couldn't get products"))

		productsRepo := NewMariaDbRepository(db)

		_, err = productsRepo.GetAll()
		assert.Error(t, err)
	})
}

func TestDBDeleteProduct(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteProduct)).WithArgs(int64(1)).
			WillReturnResult(affectedRowsResult)

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewErrorResult(sql.ErrNoRows)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteProduct)).
			WillReturnResult(affectedRowsResult)

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.NotNil(t, err)
		assert.Equal(t, "product with id 1 not found", err.Error())
	})

	t.Run("db Exec case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteProduct)).
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

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteProduct)).
			WillReturnError(errors.New("unexpected error to delete product"))

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to delete product", err.Error())
	})
}

func TestDBUpdateProduct(t *testing.T) {
	requestData := map[string]interface{}{
		"product_code":                     "ABX0001",
		"description":                      "abacaxi",
		"width":                            1.0,
		"height":                           2.0,
		"length":                           3.0,
		"net_weight":                       4.0,
		"expiration_rate":                  5.0,
		"recommended_freezing_temperature": 6.0,
		"freezing_rate":                    7.0,
		"product_type_id":                  4,
		"seller_id":                        1,
	}

	query := `UPDATE products
		SET
			product_code =						?,
			description =						?,
			width =								?,
			height =							?,
			length =							?,
			net_weight =						?,
			expiration_rate =					?,
			recommended_freezing_temperature =	?,
			freezing_rate =						?,
			product_type_id =					?,
			seller_id =							?
		WHERE id = ?`

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				"ABX0001",
				"abacaxi",
				1.0,
				2.0,
				3.0,
				4.0,
				5.0,
				6.0,
				7.0,
				4,
				1,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		newRow := mock.
			NewRows([]string{
				"id",
				"product_code",
				"description",
				"width",
				"height",
				"length",
				"net_weight",
				"expiration_rate",
				"recommended_freezing_temperature",
				"freezing_rate",
				"product_type_id",
				"seller_id"}).
			AddRow(1, "ABX0001", "abacaxi", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 4, 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneProduct)).
			WithArgs(1).WillReturnRows(newRow)

		productsRepo := NewMariaDbRepository(db)

		products, err := productsRepo.Update(1, requestData)
		assert.Nil(t, err)

		assert.Equal(t, "ABX0001", products.ProductCode)
	})

}
