package inboundorders_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	inboundorders "github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundOrders"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateLocality(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(inboundorders.QueryCreate)).
			WithArgs(
				"123",
				"123",
				123,
				123,
				123,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		inboundOrderRepo := inboundorders.NewMariaDbRepository(db)

		inboundOrder, err := inboundOrderRepo.CreateInboundOrders(
			"123",
			"123",
			123,
			123,
			123,
		)
		assert.Nil(t, err)

		assert.Equal(t, "123", inboundOrder.OrderNumber)
	})

	t.Run("Error exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(inboundorders.QueryCreate)).WillReturnError(errors.New("any"))
		inboundRepo := inboundorders.NewMariaDbRepository(db)

		inbound, err := inboundRepo.CreateInboundOrders("", "", 1, 1, 1)
		assert.NotNil(t, err)
		assert.Equal(t, "", inbound.OrderNumber)
	})

	t.Run("Error last insert id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("couldn't load the inbound order created"))
		mock.ExpectExec(regexp.QuoteMeta(inboundorders.QueryCreate)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		inboundRepo := inboundorders.NewMariaDbRepository(db)
		_, errLastInsert := inboundRepo.CreateInboundOrders("", "", 0, 0, 0)

		assert.Error(t, errLastInsert)
		assert.Equal(t, "couldn't load the inbound order created", errLastInsert.Error())
	})
}

func TestDBGetReportSellers(t *testing.T) {
	t.Run("Get all reports", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
			"inbound_orders_count",
		}).
			AddRow(1, "765", "Iu", "ri", 1, 10).
			AddRow(2, "7656", "Iuzin", "rizin", 2, 100).
			AddRow(3, "76567", "Iuzão", "rizão", 3, 1000)

		mock.ExpectQuery(regexp.QuoteMeta(inboundorders.QueryReportGetAll)).WillReturnRows(rows)

		inboundRepo := inboundorders.NewMariaDbRepository(db)

		inboundReports, err := inboundRepo.GetReportInboundOrders("")
		assert.Nil(t, err)

		assert.Len(t, inboundReports, 3)
		assert.Equal(t, inboundReports[0].CardNumberId, "765")
		assert.Equal(t, inboundReports[1].CardNumberId, "7656")
		assert.Equal(t, inboundReports[2].CardNumberId, "76567")
	})

	t.Run("Get one report", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
			"inbound_orders_count",
		}).
			AddRow(1, "765", "Iu", "ri", 1, 10)

		mock.ExpectQuery(regexp.QuoteMeta(inboundorders.QueryReportGetOne)).WillReturnRows(rows)

		inboundRepo := inboundorders.NewMariaDbRepository(db)

		inboundReports, err := inboundRepo.GetReportInboundOrders("1")
		assert.Nil(t, err)

		assert.Len(t, inboundReports, 1)
		assert.Equal(t, inboundReports[0].WarehouseId, 1)
	})

	t.Run("Error to get report - case query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(
			regexp.QuoteMeta(inboundorders.QueryReportGetOne),
		).WillReturnError(errors.New(""))

		inboundRepo := inboundorders.NewMariaDbRepository(db)

		_, err = inboundRepo.GetReportInboundOrders("1")
		assert.NotNil(t, err)
		assert.Equal(t, "error to report inbound_orders by employee", err.Error())
	})

	t.Run("Error to get report - case scan", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
			"inbound_orders_count",
		}).
			AddRow("", "", "", "", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(inboundorders.QueryReportGetOne)).WillReturnRows(rows)

		inboundRepo := inboundorders.NewMariaDbRepository(db)

		_, err = inboundRepo.GetReportInboundOrders("1")
		assert.Error(t, err)
	})
}
