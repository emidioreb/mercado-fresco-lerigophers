package employees

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreate(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreate)).
			WithArgs(
				"365",
				"Iuri",
				"Caliman",
				1,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		employeesRepo := NewMariaDbRepository(db)

		employee, err := employeesRepo.Create(
			"365",
			"Iuri",
			"Caliman",
			1,
		)
		assert.Nil(t, err)

		assert.Equal(t, Employee{1, "365", "Iuri", "Caliman", 1}, employee)
	})

	t.Run("Insert to exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryCreate)).WillReturnError(errors.New("ocurred an error to create employee"))
		employeesRepo := NewMariaDbRepository(db)

		_, err = employeesRepo.Create("", "", "", 0)
		assert.Error(t, err)
	})

	t.Run("Insert to return lastId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("ocurred an error to create employee"))
		mock.ExpectExec(regexp.QuoteMeta(queryCreate)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		employeesRepo := NewMariaDbRepository(db)
		_, err = employeesRepo.Create("", "", "", 0)

		assert.Error(t, err)
		assert.Equal(t, "ocurred an error to create employee", err.Error())
	})
}

func TestRepositoryGetOne(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
		}).
			AddRow(1, "465", "Iuri", "Caliman", 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOne)).WillReturnRows(rows)

		employeesRepo := NewMariaDbRepository(db)
		employee, err := employeesRepo.GetOne(1)
		assert.NoError(t, err)

		assert.Equal(t, Employee{1, "465", "Iuri", "Caliman", 1}, employee)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOne)).WillReturnError(sql.ErrNoRows)

		employeesRepo := NewMariaDbRepository(db)
		_, err = employeesRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "employee with id 1 not found", err.Error())
	})

	t.Run("DB Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOne)).
			WillReturnError(errors.New("unexpected error to get employee"))

		employeesRepo := NewMariaDbRepository(db)
		_, err = employeesRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to get employee", err.Error())
	})
}

func TestRepositoryGetAll(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
		}).
			AddRow(1, "465", "Iuri", "Caliman", 1).
			AddRow(2, "4656", "Iurizin", "Calimanzin", 1).
			AddRow(3, "46564", "Iurizão", "Calimanzão", 1)
		mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)

		employeesRepo := NewMariaDbRepository(db)

		employeesReports, err := employeesRepo.GetAll()
		assert.NoError(t, err)

		assert.Len(t, employeesReports, 3)
		assert.Equal(t, Employee{1, "465", "Iuri", "Caliman", 1}, employeesReports[0])
		assert.Equal(t, Employee{2, "4656", "Iurizin", "Calimanzin", 1}, employeesReports[1])
		assert.Equal(t, Employee{3, "46564", "Iurizão", "Calimanzão", 1}, employeesReports[2])
	})

	t.Run("Wrong type error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"warehouse_id",
		}).
			AddRow(1, true, "Iuri", true, "aaa")

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)

		employeesRepo := NewMariaDbRepository(db)

		_, err = employeesRepo.GetAll()
		assert.Error(t, err)
		assert.Equal(t, "couldn't get employees", err.Error())
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).
			WillReturnError(errors.New("couldn't get employees"))

		employeesRepo := NewMariaDbRepository(db)

		_, err = employeesRepo.GetAll()
		assert.Error(t, err)
	})
}

func TestRepositoryDelete(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(int64(1)).
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
		mock.ExpectExec(regexp.QuoteMeta(queryDelete)).
			WillReturnResult(affectedRowsResult)

		employeesRepo := NewMariaDbRepository(db)
		err = employeesRepo.Delete(1)
		assert.NotNil(t, err)
		assert.Equal(t, "employee with id 1 not found", err.Error())
	})

	t.Run("db Exec case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDelete)).
			WillReturnError(errors.New("unexpected error to delete employee"))

		employeesRepo := NewMariaDbRepository(db)
		err = employeesRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to delete employee", err.Error())
	})

	t.Run("DB Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		resultErr := sqlmock.NewErrorResult(errors.New("unexpected error to delete employee"))
		mock.ExpectExec(regexp.QuoteMeta(queryDelete)).
			WillReturnResult(affectedRowsResult).
			WillReturnResult(resultErr)

		sellersRepo := NewMariaDbRepository(db)
		err = sellersRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to delete employee", err.Error())
	})
}

func TestRepositoryGetOneByCardNumber(t *testing.T) {
	t.Run("Case of calling on Create", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
		}).
			AddRow(0)

		mock.ExpectQuery(regexp.QuoteMeta(queryByCardNumberCreate)).WillReturnRows(rows)

		employeesRepo := NewMariaDbRepository(db)
		errGetOne := employeesRepo.GetOneByCardNumber(0, "465")
		assert.NoError(t, errGetOne)
	})

	t.Run("Case of calling on Update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
		}).
			AddRow(0)

		mock.ExpectQuery(regexp.QuoteMeta(queryByCardNumberCreate)).WillReturnRows(rows)

		employeesRepo := NewMariaDbRepository(db)
		errGetOne := employeesRepo.GetOneByCardNumber(1, "465")
		assert.NoError(t, errGetOne)
	})

	t.Run("Case of error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryByCardNumberCreate)).WillReturnError(errors.New("unexpected error to get employee"))

		employeesRepo := NewMariaDbRepository(db)
		errGetOne := employeesRepo.GetOneByCardNumber(1, "465")
		assert.Error(t, errGetOne)
		assert.Equal(t, errors.New("unexpected error to get employee"), errGetOne)
	})

	t.Run("Case of card_number_id already exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
		}).
			AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(queryByCardNumberCreate)).WillReturnRows(rows)

		employeesRepo := NewMariaDbRepository(db)
		errGetOne := employeesRepo.GetOneByCardNumber(1, "465")
		assert.Error(t, errGetOne)
		assert.Equal(t, errors.New("card_number_id already exists"), errGetOne)
	})
}

func TestRepositoryUpdate(t *testing.T) {
	requestData := map[string]interface{}{
		"card_number_id": "456",
		"first_name":     "Iuri",
		"last_name":      "Caliman",
		"warehouse_id":   1.0,
	}

	query := `UPDATE employees
	SET 
		card_number_id = ?,
		first_name = ?,
		last_name = ?,
		warehouse_id = ?
	WHERE id = ?`

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				"456",
				"Iuri",
				"Caliman",
				1,
				1, // Employee ID to update
			).WillReturnResult(sqlmock.NewResult(1, 1))

		newRow := mock.
			NewRows([]string{
				"id",
				"card_number_id",
				"first_name",
				"last_name",
				"warehouse_id"}).
			AddRow(1, "456", "Iuri", "Caliman", 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOne)).
			WithArgs(1).WillReturnRows(newRow)

		employeesRepo := NewMariaDbRepository(db)

		employee, err := employeesRepo.Update(1, requestData)
		assert.Nil(t, err)

		assert.Equal(t, Employee{1, "456", "Iuri", "Caliman", 1}, employee)
	})

	t.Run("Exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(errors.New("any error"))

		employeesRepo := NewMariaDbRepository(db)

		_, err = employeesRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedEmployee, err)
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

		employeesRepo := NewMariaDbRepository(db)

		_, err = employeesRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedEmployee, err)
	})

	t.Run("Return updated employee case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		newResult := sqlmock.NewResult(0, 0)
		newErrorResult := sqlmock.NewErrorResult(errors.New("any error"))
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnResult(newResult)

		mock.ExpectExec(regexp.QuoteMeta(queryGetOne)).
			WillReturnResult(newErrorResult)

		sellersRepo := NewMariaDbRepository(db)

		_, err = sellersRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedEmployee, err)
	})
}
