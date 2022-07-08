package sellers

import (
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
