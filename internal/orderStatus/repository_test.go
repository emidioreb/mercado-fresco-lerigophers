package order_status

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBGetOneOrderStatus(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM order_status WHERE id = ?")).WillReturnRows(rows)

		orderStatusRepo := NewMariaDbRepository(db)
		err = orderStatusRepo.GetOne(1)
		assert.NoError(t, err)
	})
}
