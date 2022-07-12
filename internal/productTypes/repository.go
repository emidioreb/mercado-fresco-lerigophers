package producttypes

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	GetOne(id int) error
}

type mariaDbRepository struct {
	db *sql.DB
}

func NewMariaDbRepository(db *sql.DB) Repository {
	return &mariaDbRepository{
		db: db,
	}
}

func (mariaDb mariaDbRepository) GetOne(id int) error {
	queryGetOne := "SELECT id FROM product_type WHERE id = ?"
	var idSelected int

	row := mariaDb.db.QueryRow(queryGetOne, id)
	err := row.Scan(
		&idSelected,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("product_type with id %d not found", id)
	}

	if err != nil {
		return errors.New("unexpected error to verify product_type")
	}

	return nil
}
