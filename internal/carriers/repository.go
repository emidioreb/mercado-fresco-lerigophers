package carriers

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	errCreateCarry = errors.New("ocurred an error to create seller")
)

type Repository interface {
	Create(cid, companyName, address, telephone, localityId string) (Carry, error)
	GetOne(cid string) (Carry, error)
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) GetOne(cid string) (Carry, error) {
	query := `SELECT * FROM mercado_fresco.carriers WHERE cid=?`

	currentCarry := Carry{}

	row := mariaDb.db.QueryRow(query, cid)

	err := row.Scan(
		&currentCarry.Id,
		&currentCarry.Cid,
		&currentCarry.CompanyName,
		&currentCarry.Address,
		&currentCarry.Telephone,
		&currentCarry.LocalityId,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Carry{}, fmt.Errorf("Carry with cid %s not found", cid)
	}

	if err != nil {
		return Carry{}, errors.New("error to find Carry")
	}

	return currentCarry, nil
}

func (mariaDb mariaDbRepository) Create(cid, companyName, address, telephone, localityId string) (Carry, error) {
	insert := `INSERT INTO carriers (cid, company_name, address, telephone,locality_id) VALUES (?, ?, ?, ?,?)`

	newCarry := Carry{
		Cid:         cid,
		CompanyName: companyName,
		Address:     address,
		Telephone:   telephone,
		LocalityId:  localityId,
	}

	result, err := mariaDb.db.Exec(
		insert,
		cid,
		companyName,
		address,
		telephone,
		localityId,
	)

	if err != nil {
		return Carry{}, errCreateCarry
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return Carry{}, errCreateCarry
	}

	newCarry.Id = int(lastId)

	return newCarry, nil
}