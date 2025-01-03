package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

type DataToSave struct {
	PrivateKey string
	Address    string
	Currency   string
	Nonce      string
}

func (r *Repository) SaveNewAddr(data DataToSave) error {
	queryBuilder := squirrel.Insert("addresses").
		Columns("private_key", "address", "Currency", "Nonce").
		Values(data.PrivateKey, data.Address, data.Currency, data.Nonce)

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	log.Println("Tx data saved to DB")

	return nil
}

func (r *Repository) GetAllAddr() ([]string, []string, error) {
	queryBuilder := squirrel.Select("address", "currency").
		From("addresses")

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return nil, nil, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	defer rows.Close()

	var addresses []string
	var currencyes []string

	for rows.Next() {
		var addr string
		var curr string
		if err = rows.Scan(&addr, &curr); err != nil {
			log.Fatal(err)
		}

		addresses = append(addresses, addr)
		currencyes = append(currencyes, curr)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println(addresses)
	return addresses, currencyes, nil
}

func (r *Repository) GetPrivateKey(addr string) (string, string, error) {
	queryBuilder := squirrel.Select("private_key", "nonce").
		From("addresses").
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
