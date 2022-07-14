package buyers

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	mockBuyers := &Buyer{
		CardNumberId: "402324",
		FirstName:    "Fulano",
		LastName:     "Beltrano",
	}

	t.Run("success create_buyer_repository ", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(QueryCreateBuyer)).
			WithArgs(
				mockBuyers.CardNumberId,
				mockBuyers.FirstName,
				mockBuyers.LastName,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		carriersRepo := NewMariaDbRepository(db)

		carryCreate, err := carriersRepo.Create(mockBuyers.CardNumberId, mockBuyers.FirstName, mockBuyers.LastName)

		assert.NoError(t, err)

		expectedFirstName := "Fulano"

		assert.Equal(t, expectedFirstName, carryCreate.FirstName)

	})

	t.Run("failed create_buyer_repository", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(QueryCreateBuyer)).
			WithArgs(0, 0, 0, 0, 0).WillReturnResult(sqlmock.NewResult(1, 1))

		carriersRepo := NewMariaDbRepository(db)

		_, err = carriersRepo.Create(mockBuyers.CardNumberId, mockBuyers.FirstName, mockBuyers.LastName)

		assert.Error(t, err)
	})
}

func TestGetOne(t *testing.T) {
	mockBuyers := &Buyer{
		Id:           1,
		CardNumberId: "402324",
		FirstName:    "Fulano",
		LastName:     "Beltrano",
	}

	t.Run("success getOne_carry_repository", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
		}).
			AddRow(mockBuyers.Id, mockBuyers.CardNumberId, mockBuyers.FirstName, mockBuyers.LastName)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetOneBuyer)).WillReturnRows(rows)

		carriersRepo := NewMariaDbRepository(db)

		expectedFirstName := "Fulano"

		carryGetOne, err := carriersRepo.GetOne(1)
		assert.NoError(t, err)
		assert.NotNil(t, carryGetOne)
		assert.Equal(t, expectedFirstName, carryGetOne.FirstName)

	})

	t.Run("failed getOne_carry_repository", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
		}).
			AddRow("", "", "", "")

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetOneBuyer)).WillReturnRows(rows)

		carriersRepo := NewMariaDbRepository(db)

		_, err = carriersRepo.GetOne(1)

		assert.NotNil(t, err)
	})
}

func TestGetAllBuyers(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
		}).
			AddRow(1, "402324", "Fulano", "Beltrano").
			AddRow(2, "402325", "José", "Francisco").
			AddRow(3, "402326", "João", "Emídio")
		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAllBuyer)).WillReturnRows(rows)

		buyersRepo := NewMariaDbRepository(db)

		buyerGetAll, err := buyersRepo.GetAll()
		assert.NoError(t, err)

		assert.Len(t, buyerGetAll, 3)
		assert.Equal(t, "Fulano", buyerGetAll[0].FirstName)
		assert.Equal(t, "Francisco", buyerGetAll[1].LastName)
		assert.Equal(t, "402326", buyerGetAll[2].CardNumberId)
	})

	t.Run("DB Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetAllBuyer)).
			WillReturnError(errors.New("couldn't get buyers"))

		buyersRepo := NewMariaDbRepository(db)

		_, err = buyersRepo.GetAll()
		assert.Error(t, err)
	})
}

func TestDeleteBuyers(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewResult(0, 1)
		mock.ExpectExec(regexp.QuoteMeta(QueryDeleteBuyer)).WithArgs(int64(1)).
			WillReturnResult(affectedRowsResult)

		buyersRepo := NewMariaDbRepository(db)
		err = buyersRepo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Not found case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		affectedRowsResult := sqlmock.NewErrorResult(sql.ErrNoRows)
		mock.ExpectExec(regexp.QuoteMeta(QueryDeleteBuyer)).
			WillReturnResult(affectedRowsResult)

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.NotNil(t, err)
		assert.Equal(t, "buyer with id 1 not found", err.Error())
	})

	t.Run("db Exec case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(QueryDeleteBuyer)).
			WillReturnError(errors.New("any error"))

		productsRepo := NewMariaDbRepository(db)
		err = productsRepo.Delete(1)
		assert.Error(t, err)
		assert.Equal(t, "any error", err.Error())
	})
}

func TestUpdateBuyers(t *testing.T) {
	requestData := map[string]interface{}{
		"card_number_id": "402324",
		"first_name":     "Fulano",
		"last_name":      "Beltrano",
	}

	query := `UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?`

	t.Run("Success case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(
				"402324",
				"Fulano",
				"Beltrano",
				1,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		newRow := mock.
			NewRows([]string{
				"id",
				"card_number_id",
				"first_name",
				"last_name"}).
			AddRow(1, "402325", "João", "Emídio")

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetOneBuyer)).
			WithArgs(1).WillReturnRows(newRow)

		buyersRepo := NewMariaDbRepository(db)

		buyer, err := buyersRepo.Update(1, requestData)
		assert.Nil(t, err)

		assert.Equal(t, "João", buyer.FirstName)
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
		assert.Equal(t, errors.New("ocurred an error while updating the buyer"), err)
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

		buyersRepo := NewMariaDbRepository(db)

		_, err = buyersRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errors.New("ocurred an error while updating the buyer"), err)
	})

	t.Run("Return updated seller case", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		newResult := sqlmock.NewResult(0, 0)
		newErrorResult := sqlmock.NewErrorResult(errors.New("any error"))
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnResult(newResult)

		mock.ExpectExec(regexp.QuoteMeta(QueryGetOneBuyer)).
			WillReturnResult(newErrorResult)

		buyersRepo := NewMariaDbRepository(db)

		_, err = buyersRepo.Update(1, requestData)
		assert.Error(t, err)
		assert.Equal(t, errors.New("ocurred an error while updating the buyer"), err)
	})

}

func TestGetReportPurchaseOrders(t *testing.T) {
	t.Run("Get all reports purchase orders", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id",
			"card_number_id",
			"first_name",
			"last_name",
			"purchase_orders_count",
		}).
		AddRow(1, "402324", "Fulano", "Beltrano",123).
		AddRow(2, "402325", "José", "Francisco",124).
		AddRow(3, "402326", "João", "Emídio",125)

		mock.ExpectQuery(regexp.QuoteMeta(QueryGetReportAll)).WillReturnRows(rows)

		buyerRepo := NewMariaDbRepository(db)

		purchaseOrders, err := buyerRepo.GetReportPurchaseOrders(0)
		assert.Nil(t, err)

		assert.Len(t, purchaseOrders, 3)
		assert.Equal(t, purchaseOrders[0].FirstName, "Fulano")
		assert.Equal(t, purchaseOrders[1].FirstName, "José")
		assert.Equal(t, purchaseOrders[2].FirstName, "João")
	})
}
