package inboundorders

import (
	"database/sql"
	"errors"
)

type Repository interface {
	CreateInboundOrders(orderNumber, orderDate string, employeeId, productBatchId, warehouseId int) (InboundOrder, error)
	GetReportInboundOrders(employeeId string) ([]ReportInboundOrder, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) CreateInboundOrders(orderNumber, orderDate string, employeeId, productBatchId, warehouseId int) (InboundOrder, error) {

	newInbound := InboundOrder{
		OrderNumber:    orderNumber,
		OrderDate:      orderDate,
		EmployeeId:     employeeId,
		ProductBatchId: productBatchId,
		WarehouseId:    warehouseId,
	}

	result, err := mariaDb.db.Exec(
		QueryCreate,
		orderNumber,
		orderDate,
		employeeId,
		productBatchId,
		warehouseId,
	)
	if err != nil {
		return InboundOrder{}, errors.New("couldn't create a inbound order")
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return InboundOrder{}, errors.New("couldn't load the inbound order created")
	}

	newInbound.Id = int(lastId)

	return newInbound, nil
}

func (mariaDb mariaDbRepository) GetReportInboundOrders(employeeId string) ([]ReportInboundOrder, error) {
	reports := []ReportInboundOrder{}

	var (
		rows *sql.Rows
		err  error
	)

	if employeeId != "" {
		rows, err = mariaDb.db.Query(QueryReportGetOne, employeeId)
	} else {
		rows, err = mariaDb.db.Query(QueryReportGetAll)
	}

	if err != nil {
		return []ReportInboundOrder{}, errors.New("error to report inbound_orders by employee")
	}

	for rows.Next() {
		var currentReport ReportInboundOrder
		if err := rows.Scan(
			&currentReport.Id,
			&currentReport.CardNumberId,
			&currentReport.FirstName,
			&currentReport.LastName,
			&currentReport.WarehouseId,
			&currentReport.InboundOrdersCount,
		); err != nil {
			return []ReportInboundOrder{}, errors.New("error to report inbound_orders by employee")
		}
		reports = append(reports, currentReport)
	}

	return reports, nil
}
