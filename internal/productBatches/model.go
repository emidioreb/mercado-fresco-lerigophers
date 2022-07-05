package product_batches

// ProductBatches
type ProductBatches struct {
	Id                 int    `json:"id"`
	BatchNumber        int    `json:"batch_number"`
	CurrentQuantity    int    `json:"current_quantity"`
	CurrentTemperature int    `json:"current_temperature"`
	InitialQuantity    int    `json:"initial_quantity"`
	ManufacturingHour  int    `json:"manufacturing_hour"`
	MinimumTemperature int    `json:"minumum_temperature"`
	ProductId          int    `json:"product_id"`
	SectionId          int    `json:"section_id"`
	DueDate            string `json:"due_date"`
	ManufacturingDate  string `json:"manufacturing_date"`
}

type ProductsQuantity struct {
	SectionId     string `json:"section_id"`
	SectionNumber string `json:"section_number"`
	ProductsCount int    `json:"products_count"`
}
