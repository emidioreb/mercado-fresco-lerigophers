package warehouses

type Warehouse struct {
	Id                 int    `json:"id"`
	WarehouseCode      string `json:"warehouse_code"`
	Address            string `json:"adress"`
	Telephone          string `json:"telephone"`
	MinimumCapacity    int    `json:"minimum_capacity"`
	MinimumTemperature int    `json:"minimum_temperature"`
}
