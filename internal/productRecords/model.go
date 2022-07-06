package product_records

type ProductRecords struct {
	Id             int     `json:"id"`
	LastUpdateDate string  `json:"last_update_date"`
	PurchasePrice  float64 `json:"purchase_price"`
	SalePrice      float64 `json:"sale_price"`
	ProductId      int     `json:"product_id"`
}

type RecordsQuantity struct {
	ProductId    int    `json:"product_id"`
	ProductCode  string `json:"product_code"`
	RecordsCount int    `json:"records_count"`
}
