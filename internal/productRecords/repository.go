package product_records

import (
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	CreateProductRecord(LastUpdateDate string, PurchasePrice float64, SalePrice float64, ProductId int) (ProductRecords, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) CreateProductRecord(LastUpdateDate string, PurchasePrice float64, SalePrice float64, ProductId int) (ProductRecords, error) {
	layout := "2006-01-02"
	lastUpdateDate, err := time.Parse(layout, LastUpdateDate)

	if err != nil {
		return ProductRecords{}, errors.New("invalid date input")
	}

	result, err := mariaDb.db.Exec(queryCreateProductRecord, lastUpdateDate, PurchasePrice, SalePrice, ProductId)
	if err != nil {
		return ProductRecords{}, errors.New("couldn't create a product_record")
	}

	newProductRecord := ProductRecords{
		LastUpdateDate: LastUpdateDate,
		PurchasePrice:  PurchasePrice,
		SalePrice:      SalePrice,
		ProductId:      ProductId,
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return ProductRecords{}, errors.New("ocurred an error to create product record")
	}

	newProductRecord.Id = int(lastId)

	return newProductRecord, nil
}
