package products

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errUpdatedProduct = errors.New("ocurred an error while updating the product")
	errCreateProduct  = errors.New("ocurred an error to create product")
	errGetProducts    = errors.New("couldn`t get products")
	errGetOneProduct  = errors.New("unexpected error to get product")
	errDeleteProduct  = errors.New("unexpected error to delete product")
)

type Repository interface {
	Create(productCode, description string, width, height, length, netWeight, expirationRate, recommendedFreezingTemperaturechan,
		freezingRate float64, productTypeId, sellerId int) (Product, error)
	GetOne(id int) (Product, error)
	GetAll() ([]Product, error)
	Delete(id int) error
	Update(id int, requestData map[string]interface{}) (Product, error)
	GetReportProduct(ProductId int) ([]ProductRecords, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) Create(productCode, description string, width, height,
	length, netWeight, expirationRate, recommendedFreezingTemperature,
	freezingRate float64, productTypeId, sellerId int) (Product, error) {

	newProduct := Product{
		ProductCode:                    productCode,
		Description:                    description,
		Width:                          width,
		Height:                         height,
		Length:                         length,
		NetWeight:                      netWeight,
		ExpirationRate:                 expirationRate,
		RecommendedFreezingTemperature: recommendedFreezingTemperature,
		FreezingRate:                   freezingRate,
		ProductTypeId:                  productTypeId,
		SellerId:                       sellerId,
	}

	result, err := mariaDb.db.Exec(
		queryCreateProduct,
		productCode,
		description,
		width,
		height,
		length,
		netWeight,
		expirationRate,
		recommendedFreezingTemperature,
		freezingRate,
		productTypeId,
		sellerId,
	)

	if err != nil {
		return Product{}, errCreateProduct
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Product{}, errCreateProduct
	}

	newProduct.Id = int(lastId)

	return newProduct, nil
}

func (mariaDb mariaDbRepository) GetOne(id int) (Product, error) {
	currentProduct := Product{}

	row := mariaDb.db.QueryRow(queryGetOneProduct, id)
	err := row.Scan(&currentProduct.Id,
		&currentProduct.ProductCode,
		&currentProduct.Description,
		&currentProduct.Width,
		&currentProduct.Height,
		&currentProduct.Length,
		&currentProduct.NetWeight,
		&currentProduct.ExpirationRate,
		&currentProduct.RecommendedFreezingTemperature,
		&currentProduct.FreezingRate,
		&currentProduct.ProductTypeId,
		&currentProduct.SellerId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, fmt.Errorf("product with id %d not found", id)
	}

	if err != nil {
		return Product{}, errGetOneProduct
	}

	return currentProduct, nil
}

func (mariaDb mariaDbRepository) GetAll() ([]Product, error) {
	products := []Product{}

	rows, err := mariaDb.db.Query(queryGetAllProducts)
	if err != nil {
		return []Product{}, errGetProducts
	}

	for rows.Next() {
		var currentProduct Product
		err := rows.Scan(
			&currentProduct.Id,
			&currentProduct.ProductCode,
			&currentProduct.Description,
			&currentProduct.Width,
			&currentProduct.Height,
			&currentProduct.Length,
			&currentProduct.NetWeight,
			&currentProduct.ExpirationRate,
			&currentProduct.RecommendedFreezingTemperature,
			&currentProduct.FreezingRate,
			&currentProduct.ProductTypeId,
			&currentProduct.SellerId,
		)
		if err != nil {
			return []Product{}, errGetProducts
		}
		products = append(products, currentProduct)
	}
	return products, nil
}

func (mariaDb mariaDbRepository) Delete(id int) error {
	result, err := mariaDb.db.Exec(queryDeleteProduct, id)

	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	if err != nil {
		return errDeleteProduct
	}

	return nil
}

func (mariaDb mariaDbRepository) Update(id int, requestData map[string]interface{}) (Product, error) {
	finalQuery, valuesToUse := queryUpdateProduct(requestData, id)

	result, err := mariaDb.db.Exec(finalQuery, valuesToUse...)
	if err != nil {
		return Product{}, errUpdatedProduct
	}

	_, err = result.RowsAffected()
	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return Product{}, errUpdatedProduct
	}

	currentProduct, err := mariaDb.GetOne(id)
	if err != nil {
		return Product{}, errUpdatedProduct
	}

	return currentProduct, nil
}

func (mariaDb mariaDbRepository) GetReportProduct(ProductId int) ([]ProductRecords, error) {
	reports := []ProductRecords{}

	var (
		rows *sql.Rows
		err  error
	)

	if ProductId != 0 {
		rows, err = mariaDb.db.Query(queryGetReportOne, ProductId)
	} else {
		rows, err = mariaDb.db.Query(queryGetReportAll)
	}

	if err != nil {
		return []ProductRecords{}, errors.New("error to report products by product_id")
	}

	for rows.Next() {
		var currentReport ProductRecords
		if err := rows.Scan(
			&currentReport.ProductId,
			&currentReport.Description,
			&currentReport.RecordsCount,
		); err != nil {
			return []ProductRecords{}, errors.New("error to report products by product_id")
		}
		reports = append(reports, currentReport)
	}

	return reports, nil
}
