package carriers

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	// "github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	mockCarriers := &Carry{
		Cid:         "CID#25vA",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		LocalityId:  "456",
	}

	query := `INSERT INTO carriers (cid, company_name, address, telephone,locality_id) VALUES (?, ?, ?, ?,?)`

	t.Run("success create_carry_repository ", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				mockCarriers.Cid,
				mockCarriers.CompanyName,
				mockCarriers.Address,
				mockCarriers.Telephone,
				mockCarriers.LocalityId,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		carriersRepo := NewMariaDbRepository(db)

		carryCreate, err := carriersRepo.Create(mockCarriers.Cid, mockCarriers.CompanyName, mockCarriers.Address, mockCarriers.Telephone, mockCarriers.LocalityId)

		assert.NoError(t, err)

		expectedCompanyName := "some name"

		assert.Equal(t, expectedCompanyName, carryCreate.CompanyName)

	})

	t.Run("failed create_carry_repository", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(0, 0, 0, 0, 0).WillReturnResult(sqlmock.NewResult(1, 1))

		carriersRepo := NewMariaDbRepository(db)

		_, err = carriersRepo.Create(mockCarriers.Cid, mockCarriers.CompanyName, mockCarriers.Address, mockCarriers.Telephone, mockCarriers.LocalityId)

		assert.Error(t, err)
	})
}

func TestGetOne(t *testing.T) {
	mockCarriers := &Carry{
		Id: 1,
		Cid:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		LocalityId:  "456",
	}

	t.Run("success getOne_carry_repository", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()


		query := `SELECT * FROM carriers WHERE cid=?`
		rows := sqlmock.NewRows([]string{
			"id",
			"cid",
			"company_name",
			"adress",
			"telephone",
			"locality_id",
		}).
			AddRow(mockCarriers.Id, mockCarriers.Cid, mockCarriers.CompanyName, mockCarriers.Address, mockCarriers.Telephone, mockCarriers.LocalityId)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		carriersRepo := NewMariaDbRepository(db)

		expectedCompanyName := "some name"

		carryGetOne, err := carriersRepo.GetOne("CID#1")
		assert.NoError(t, err)
		assert.NotNil(t, carryGetOne)
		assert.Equal(t,expectedCompanyName, carryGetOne.CompanyName)

	})

	t.Run("failed getOne_carry_repository", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()



		query := `SELECT * FROM carriers WHERE cid=?`
		rows := sqlmock.NewRows([]string{
			"id",
			"cid",
			"company_name",
			"adress",
			"telephone",
			"locality_id",
		}).
			AddRow("", "", "", "", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		carriersRepo := NewMariaDbRepository(db)

		_, err = carriersRepo.GetOne("CID#1")

		assert.NotNil(t, err)
	})
}
