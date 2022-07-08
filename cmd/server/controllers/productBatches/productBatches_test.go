package controllers_test

import (
	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
)

type ObjectResponseArr struct {
	Data []product_batches.ProductBatches
}

type ObjectResponse struct {
	Data product_batches.ProductBatches
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}
