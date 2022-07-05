package carriers

type Carry struct {
	Id          int    `json:"id"`
	Cid         string   `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  string `json:"locality_id"`
}
