package product_batches

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository interface {
	CreateProductBatch(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId int, DueDate, ManufacturingDate time.Time) (ProductBatches, error)
	GetReportSection(SectionId int) ([]ProductsQuantity, error)
	GetOne(BatchNumber int) (ProductBatches, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) GetOne(BatchNumber int) (ProductBatches, error) {
	currentProductBatch := ProductBatches{}

	row := mariaDb.db.QueryRow(QueryGetOneProductBatch, BatchNumber)
	err := row.Scan(
		&currentProductBatch.Id,
		&currentProductBatch.BatchNumber,
		&currentProductBatch.CurrentQuantity,
		&currentProductBatch.CurrentTemperature,
		&currentProductBatch.InitialQuantity,
		&currentProductBatch.ManufacturingHour,
		&currentProductBatch.MinimumTemperature,
		&currentProductBatch.ProductId,
		&currentProductBatch.SectionId,
		&currentProductBatch.DueDate,
		&currentProductBatch.ManufacturingDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return ProductBatches{}, fmt.Errorf("product_batch with batch_number %d not found", BatchNumber)
	}

	if err != nil {
		return ProductBatches{}, errors.New("error to find product_batch")
	}

	return currentProductBatch, nil
}

func (mariaDb mariaDbRepository) CreateProductBatch(BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId int, DueDate, ManufacturingDate time.Time) (ProductBatches, error) {
	_, err := mariaDb.db.Exec(QueryCreateProductBatch, BatchNumber, CurrentQuantity, CurrentTemperature, InitialQuantity, ManufacturingHour, MinimumTemperature, ProductId, SectionId, DueDate, ManufacturingDate)
	if err != nil {
		return ProductBatches{}, errors.New("couldn't create a product_batch")
	}

	newProductBatch := ProductBatches{
		BatchNumber:        BatchNumber,
		CurrentQuantity:    CurrentQuantity,
		CurrentTemperature: CurrentTemperature,
		InitialQuantity:    InitialQuantity,
		ManufacturingHour:  ManufacturingHour,
		MinimumTemperature: MinimumTemperature,
		ProductId:          ProductId,
		SectionId:          SectionId,
		DueDate:            DueDate,
		ManufacturingDate:  ManufacturingDate,
	}

	return newProductBatch, nil
}

func (mariaDb mariaDbRepository) GetReportSection(SectionId int) ([]ProductsQuantity, error) {
	reports := []ProductsQuantity{}

	var (
		rows *sql.Rows
		err  error
	)

	if SectionId != 0 {
		rows, err = mariaDb.db.Query(QueryGetReportOne, SectionId)
	} else {
		rows, err = mariaDb.db.Query(QueryGetReportAll)
	}

	if err != nil {
		return []ProductsQuantity{}, errors.New("error to report sections by product_batches")
	}

	for rows.Next() {
		var currentReport ProductsQuantity
		if err := rows.Scan(
			&currentReport.SectionId,
			&currentReport.SectionNumber,
			&currentReport.ProductsCount,
		); err != nil {
			return []ProductsQuantity{}, errors.New("error to report sections by product_batches")
		}
		reports = append(reports, currentReport)
	}

	return reports, nil
}
