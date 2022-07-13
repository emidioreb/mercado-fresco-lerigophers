package sections

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateSection(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryCreateSection)).
			WithArgs(
				1,
				10,
				2,
				100,
				50,
				500,
				1,
				1,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		sectionsRepo := NewMariaDbRepository(db)

		section, err := sectionsRepo.Create(
			1,
			10,
			2,
			100,
			50,
			500,
			1,
			1,
		)
		assert.Nil(t, err)

		assert.Equal(t, 10, section.CurrentTemperature)
	})

	t.Run("Insert to exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryCreateSection)).WillReturnError(errors.New("internal db error"))
		sectionsRepo := NewMariaDbRepository(db)

		_, err = sectionsRepo.Create(0, 0, 0, 0, 0, 0, 0, 0)
		assert.Error(t, err)
	})

	t.Run("Insert to return lastId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlDriverResultErr := sqlmock.NewErrorResult(errors.New("ocurred an error to create section"))
		mock.ExpectExec(regexp.QuoteMeta(queryCreateSection)).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlDriverResultErr)

		sectionsRepo := NewMariaDbRepository(db)
		_, err = sectionsRepo.Create(0, 0, 0, 0, 0, 0, 0, 0)

		assert.Error(t, err)
		assert.Equal(t, "ocurred an error to create section", err.Error())
	})
}

func TestDBGetOneSection(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"sectionNumber",
			"currentTemperature",
			"minimumTemperature",
			"currentCapacity",
			"mininumCapacity",
			"maximumCapacity",
			"warehouseId",
			"productTypeId",
		}).
			AddRow(1, 1, 10, 2, 100, 50, 500, 1, 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSection)).WillReturnRows(rows)

		sectionsRepo := NewMariaDbRepository(db)
		section, err := sectionsRepo.GetOne(1)
		assert.NoError(t, err)

		assert.Equal(t, 10, section.CurrentTemperature)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSection)).WillReturnError(sql.ErrNoRows)

		sectionsRepo := NewMariaDbRepository(db)
		_, err = sectionsRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "section with id 1 not found", err.Error())
	})

	t.Run("DB Error case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSection)).
			WillReturnError(errors.New("unexpected error to get section"))

		sectionsRepo := NewMariaDbRepository(db)
		_, err = sectionsRepo.GetOne(1)
		assert.Error(t, err)
		assert.Equal(t, "unexpected error to get section", err.Error())
	})
}

func TestDBGetAllSections(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"sectionNumber",
			"currentTemperature",
			"minimumTemperature",
			"currentCapacity",
			"mininumCapacity",
			"maximumCapacity",
			"warehouseId",
			"productTypeId",
		}).
			AddRow(1, 1, 10, 2, 100, 50, 500, 1, 1).
			AddRow(2, 2, 20, 3, 110, 60, 600, 2, 2).
			AddRow(3, 3, 30, 4, 120, 70, 700, 3, 3)
		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllSections)).WillReturnRows(rows)

		sectionsRepo := NewMariaDbRepository(db)

		sectionReports, err := sectionsRepo.GetAll()
		assert.NoError(t, err)

		assert.Len(t, sectionReports, 3)
		assert.Equal(t, 10, sectionReports[0].CurrentTemperature)
		assert.Equal(t, 20, sectionReports[1].CurrentTemperature)
		assert.Equal(t, 30, sectionReports[2].CurrentTemperature)
	})

	t.Run("Wrong type error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"sectionNumber",
			"currentTemperature",
			"minimumTemperature",
			"currentCapacity",
			"mininumCapacity",
			"maximumCapacity",
			"warehouseId",
			"productTypeId",
		}).AddRow(1, "abc", 10, 2, 100, 50, 500, 1, 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllSections)).WillReturnRows(rows)

		sectionsRepo := NewMariaDbRepository(db)

		_, err = sectionsRepo.GetAll()
		assert.Error(t, err)
		assert.Equal(t, "couldn't get sections", err.Error())
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(queryGetAllSections)).
			WillReturnError(errors.New("couldn't get sections"))

		sectionsRepo := NewMariaDbRepository(db)

		_, err = sectionsRepo.GetAll()
		assert.Error(t, err)
	})
}

func TestDBDeleteSection(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSection)).WithArgs(int64(1)).
			WillReturnResult(affectedRowsResult)

		sectionsRepo := NewMariaDbRepository(db)
		err = sectionsRepo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Error case - Exec query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSection)).
			WillReturnError(errors.New("any error"))

		sectionsRepo := NewMariaDbRepository(db)
		err = sectionsRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, errDeleteSection, err)
	})
}

func TestDBUpdateSection(t *testing.T) {
	requestData := map[string]interface{}{
		"section_number":      1.0,
		"current_temperature": 10.0,
		"minimum_temperature": 5.0,
		"current_capacity":    10.0,
		"minimum_capacity":    5.0,
		"maximum_capacity":    100.0,
		"warehouse_id":        1.0,
		"product_type_id":     1.0,
	}

	query := `	UPDATE sections SET 
						section_number = ?,
						current_temperature = ?,
						minimum_temperature = ?,
						current_capacity = ?,
						minimum_capacity = ?,
						maximum_capacity = ?,
						warehouse_id = ?,
						product_type_id = ?
					WHERE id = ?
			 `
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				1,
				10,
				5,
				10,
				5,
				100,
				1,
				1,
				1, // Section ID to update
			).WillReturnResult(sqlmock.NewResult(1, 1))

		newRow := mock.
			NewRows([]string{
				"id",
				"section_number",
				"current_temperature",
				"minimum_temperature",
				"current_capacity",
				"minimum_capacity",
				"maximum_capacity",
				"warehouse_id",
				"product_type_id",
			}).
			AddRow(1, 1, 15, 1, 1, 1, 1, 1, 1)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetOneSection)).
			WithArgs(1).WillReturnRows(newRow)

		sectionsRepo := NewMariaDbRepository(db)

		section, err := sectionsRepo.Update(1, requestData)
		assert.Nil(t, err)
		assert.Equal(t, 15, section.CurrentTemperature)
	})

	t.Run("Exec error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(errors.New("any error"))

		sectionsRepo := NewMariaDbRepository(db)

		_, err = sectionsRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errUpdatedSection, err)
	})
}
