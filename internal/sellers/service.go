package sellers

type Service interface {
	Create(cid int, companyName, address, telephone string) (Seller, error)
	GetOne(id int) (Seller, error)
	GetAll() ([]Seller, error)
	Delete(id int) error
	Update(id, cid int, companyName, address, telephone string) (Seller, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Create(cid int, companyName, address, telephone string) (Seller, error) {
	seller, err := s.repository.Create(cid, companyName, address, telephone)

	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}

func (s service) GetOne(id int) (Seller, error) {
	seller, err := s.repository.GetOne(id)

	if err != nil {
		return Seller{}, err
	}
	return seller, nil
}

func (s service) GetAll() ([]Seller, error) {
	sellers, err := s.repository.GetAll()

	if err != nil {
		return []Seller{}, err
	}
	return sellers, nil
}
func (s service) Delete(id int) error {
	err := s.repository.Delete(id)

	if err != nil {
		return err
	}
	return nil
}
func (s service) Update(id, cid int, companyName, address, telephone string) (Seller, error) {
	seller, err := s.repository.Update(id, cid, companyName, address, telephone)

	if err != nil {
		return Seller{}, err
	}

	return seller, nil
}
