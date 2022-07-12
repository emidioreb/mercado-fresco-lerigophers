package localities

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateLocality(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateLocality)).
			WithArgs(
				"123",
				"123",
				"123",
				"123",
			).WillReturnResult(sqlmock.NewResult(1, 1))

		localitiesRepo := NewMariaDbRepository(db)

		locality, err := localitiesRepo.CreateLocality(
			"123",
			"123",
			"123",
			"123",
		)
		assert.Nil(t, err)

		assert.Equal(t, "123", locality.LocalityName)
	})

	t.Run("Already exists case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		rows := sqlmock.NewRows([]string{
			"id",
			"locality_name",
			"province_name",
			"country_name",
		}).
			AddRow(
				"123",
				"123",
				"123",
				"123",
			)
		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneLocality)).WillReturnRows(rows)
		localitiesRepo := NewMariaDbRepository(db)

		locality, err := localitiesRepo.GetOne("123")
		assert.NoError(t, err)
		assert.Equal(t, "123", locality.Id)
	})

	t.Run("Already exists case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateLocality)).WillReturnError(errors.New(""))
		localitiesRepo := NewMariaDbRepository(db)

		emptyLocality, err := localitiesRepo.CreateLocality("123", "Presidente Dutra", "MA", "BR")
		assert.NotNil(t, err)
		assert.Equal(t, "", emptyLocality.CountryName)
	})

}

func TestDBGetOneLocality(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"locality_name",
			"province_name",
			"country_name",
		}).AddRow("123",
			"Presidente Dutra",
			"MA",
			"BR",
		)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneLocality)).WillReturnRows(rows)

		localitiesRepo := NewMariaDbRepository(db)
		locality, err := localitiesRepo.GetOne("123")
		assert.Nil(t, err)

		assert.Equal(t, "123", locality.Id)
		assert.Equal(t, "Presidente Dutra", locality.LocalityName)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneLocality)).WillReturnError(sql.ErrNoRows)

		localitiesRepo := NewMariaDbRepository(db)
		_, err = localitiesRepo.GetOne("123")
		assert.NotNil(t, err)
		assert.Equal(t, "locality with id 123 not found", err.Error())

	})
	t.Run("Another error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneLocality)).WillReturnError(errors.New(""))

		localitiesRepo := NewMariaDbRepository(db)
		_, err = localitiesRepo.GetOne("123")
		assert.NotNil(t, err)
		assert.Equal(t, "error to find locality", err.Error())
	})
}

func TestDBGetReportSellers(t *testing.T) {
	t.Run("Get all reports", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"sellers_count",
		}).
			AddRow("123", "Pres. Dutra", 200).
			AddRow("456", "Osasco", 150).
			AddRow("789", "São Luís", 250)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportAll)).WillReturnRows(rows)

		localityRepo := NewMariaDbRepository(db)

		sellerReports, err := localityRepo.GetReportSellers("")
		assert.Nil(t, err)

		assert.Len(t, sellerReports, 3)
		assert.Equal(t, sellerReports[0].LocalityId, "123")
		assert.Equal(t, sellerReports[1].LocalityId, "456")
		assert.Equal(t, sellerReports[2].LocalityId, "789")
	})

	t.Run("Get one report", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"sellers_count",
		}).
			AddRow("456", "Osasco", 150)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportOne)).WillReturnRows(rows)

		localityRepo := NewMariaDbRepository(db)

		sellerReports, err := localityRepo.GetReportSellers("456")
		assert.Nil(t, err)

		assert.Len(t, sellerReports, 1)
		assert.Equal(t, sellerReports[0].LocalityName, "Osasco")
	})

	t.Run("Error to get report - case query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(
			regexp.QuoteMeta(queryGetReportOne),
		).WillReturnError(errors.New(""))

		localityRepo := NewMariaDbRepository(db)

		_, err = localityRepo.GetReportSellers("456")
		assert.NotNil(t, err)
		assert.Equal(t, "error to report sellers by locality", err.Error())
	})

	t.Run("Error to get report - case scan", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"sellers_count",
		}).AddRow("123", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportOne)).WillReturnRows(rows)

		localityRepo := NewMariaDbRepository(db)

		_, err = localityRepo.GetReportSellers("123")
		assert.Error(t, err)
	})
}

func TestDBGetReportCarriers(t *testing.T) {
	t.Run("Get all reports", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"carriers_count",
		}).
			AddRow("123", "Pres. Dutra", 200).
			AddRow("456", "Osasco", 150).
			AddRow("789", "São Luís", 250)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportCarriersAll)).
			WillReturnRows(rows)

		localityRepo := NewMariaDbRepository(db)

		sellerReports, err := localityRepo.GetReportCarriers("")
		assert.Nil(t, err)

		assert.Len(t, sellerReports, 3)
		assert.Equal(t, sellerReports[0].LocalityId, "123")
		assert.Equal(t, sellerReports[1].LocalityId, "456")
		assert.Equal(t, sellerReports[2].LocalityId, "789")
	})

	t.Run("Get one report", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"carriers_count",
		}).
			AddRow("123", "Pres. Dutra", 200)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportCarriersOne)).
			WillReturnRows(rows)

		localityRepo := NewMariaDbRepository(db)

		sellerReports, err := localityRepo.GetReportCarriers("123")
		assert.Nil(t, err)

		assert.Len(t, sellerReports, 1)
		assert.Equal(t, sellerReports[0].LocalityId, "123")
	})

	t.Run("Error to exec query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportCarriersOne)).
			WillReturnError(errors.New("any error"))

		localityRepo := NewMariaDbRepository(db)

		_, err = localityRepo.GetReportCarriers("123")
		assert.Error(t, err)
		assert.Equal(t, "error to report carriers by locality", err.Error())
	})

	t.Run("Get all reports", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"locality_id",
			"locality_name",
			"carriers_count",
		}).
			AddRow("", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetReportCarriersAll)).
			WillReturnRows(rows)

		localityRepo := NewMariaDbRepository(db)

		carriersReport, err := localityRepo.GetReportCarriers("")
		assert.Error(t, err)
		assert.Equal(t, "error to report carriers by locality", err.Error())

		assert.Equal(t, []ReportCarriers{}, carriersReport)
	})
}
