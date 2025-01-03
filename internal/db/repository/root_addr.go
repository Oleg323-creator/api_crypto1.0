package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

func (r *Repository) SaveRootData(data DataToSave) error {
	queryBuilder := squirrel.Insert("root_address").
		Columns("private_key", "address", "nonce").
		Values(data.PrivateKey, data.Address, data.Nonce)

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

func (r *Repository) GetRootPrivateKey(addr string) (string, string, error) {
	queryBuilder := squirrel.Select("private_key", "nonce").
		From("root_address").
		Where(squirrel.Eq{"address": addr})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return "", "", fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return "", "", fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	defer rows.Close()

	var key string
	var nonce string
	if rows.Next() {
		if err = rows.Scan(&key, &nonce); err != nil {
			return "", "", fmt.Errorf("failed to scan result: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return key, nonce, nil
}
