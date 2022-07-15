package producttypes

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBGetOneProductType(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM product_type WHERE id = ?")).WillReturnRows(rows)

		productTypeRepo := NewMariaDbRepository(db)
		err = productTypeRepo.GetOne(1)
		assert.NoError(t, err)
	})
}
