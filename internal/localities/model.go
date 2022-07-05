package localities

type Locality struct {
	Id           string `json:"id"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}

type ReportSellers struct {
	LocalityId   string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	SellersCount int    `json:"sellers_count"`
}
