package buyers

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedBuyer         = errors.New("ocurred an error while updating the buyer")
	errCreateBuyer          = errors.New("ocurred an error to create buyer")
	errGetBuyers            = errors.New("couldn't get buyers")
	errGetOneBuyer          = errors.New("unexpected error to get buyer")
	errDeleteBuyer          = errors.New("unexpected error to delete buyer")
	errReportPurchaseOrders = errors.New("error to report purchase_orders by buyer")
)

type Repository interface {
	Create(cardNumberId, firstName, lastName string) (Buyer, error)
	GetOne(id int) (Buyer, error)
	GetAll() ([]Buyer, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Buyer, error)
	GetReportPurchaseOrders(BuyerId int) ([]ReportPurchaseOrders, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(cardNumberId, firstName, lastName string) (Buyer, error) {
	newBuyer := Buyer{
		CardNumberId: cardNumberId,
		FirstName:    firstName,
		LastName:     lastName,
	}

	result, err := mariaDb.db.Exec(
		QueryCreateBuyer,
		cardNumberId,
		firstName,
		lastName,
	)

	if err != nil {
		return Buyer{}, errCreateBuyer
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Buyer{}, errCreateBuyer
	}

	newBuyer.Id = int(lastId)

	return newBuyer, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Buyer, error) {

	currentBuyer := Buyer{}

	row := mariaDb.db.QueryRow(QueryGetOneBuyer, id)
	err := row.Scan(
		&currentBuyer.Id,
		&currentBuyer.CardNumberId,
		&currentBuyer.FirstName,
		&currentBuyer.LastName,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Buyer{}, fmt.Errorf("buyer with id %d not found", id)
	}

	if err != nil {
		return Buyer{}, errGetOneBuyer
	}

	return currentBuyer, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Buyer, error) {
	buyers := []Buyer{}

	rows, err := mariaDb.db.Query(QueryGetAllBuyer)
	if err != nil {
		return []Buyer{}, errGetBuyers
	}

	for rows.Next() {
		var currentBuyer Buyer
		if err := rows.Scan(
			&currentBuyer.Id,
			&currentBuyer.CardNumberId,
			&currentBuyer.FirstName,
			&currentBuyer.LastName,
		); err != nil {
			return []Buyer{}, errGetBuyers
		}
		buyers = append(buyers, currentBuyer)
	}
	return buyers, nil
}
func (mariaDb mariaDbRepository) Delete(id int) error {
	result, err := mariaDb.db.Exec(QueryDeleteBuyer, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 {
		return fmt.Errorf("buyer with id %d not found", id)
	}

	if err != nil {
		return errDeleteBuyer
	}

	return nil
}

func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Buyer, error) {
	finalQuery, valuesToUse := QueryUpdateBuyer(requestData, id)

	result, err := mariaDb.db.Exec(finalQuery, valuesToUse...)
	if err != nil {
		return Buyer{}, errUpdatedBuyer
	}

	_, err = result.RowsAffected()
	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return Buyer{}, errUpdatedBuyer
	}

	currentbuyer, err := mariaDb.GetOne(id)
	if err != nil {
		return Buyer{}, errUpdatedBuyer
	}

	return currentbuyer, nil
}

func (mariaDb mariaDbRepository) GetReportPurchaseOrders(BuyerId int) ([]ReportPurchaseOrders, error) {
	reports := []ReportPurchaseOrders{}

	var (
		rows *sql.Rows
		err  error
	)

	if BuyerId != 0 {
		rows, err = mariaDb.db.Query(QueryGetReportOne, BuyerId)
	} else {
		rows, err = mariaDb.db.Query(QueryGetReportAll)
	}

	if err != nil {
		return []ReportPurchaseOrders{}, errReportPurchaseOrders
	}

	for rows.Next() {
		var currentReport ReportPurchaseOrders
		if err := rows.Scan(
			&currentReport.BuyerId,
			&currentReport.CardNumberId,
			&currentReport.FirstName,
			&currentReport.LastName,
			&currentReport.PurchaseOrdersCount,
		); err != nil {
			return []ReportPurchaseOrders{}, errReportPurchaseOrders
		}
		reports = append(reports, currentReport)
	}

	return reports, nil
}
