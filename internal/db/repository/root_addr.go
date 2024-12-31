package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

func (r *Repository) SaveRootAddrToDB(data DataToSave) error {
	queryBuilder := squirrel.Insert("root_address").
		Columns("private_key", "address").
		Values(data.PrivateKey, data.Address)

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return execErr
	}
	log.Println("Root address saved to db")

	return nil
}

func (r *Repository) GetRootAddr() (string, error) {
	queryBuilder := squirrel.Select("address").
		From("root_address").
		Where(squirrel.Eq{"id": "1"})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return "", fmt.Errorf("There is no root address: %v", err)
	}

	var addr string

	err = r.DB.QueryRow(query, args...).Scan(&addr)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return addr, nil
}
