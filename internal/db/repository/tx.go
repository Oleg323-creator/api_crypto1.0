package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
	"strings"
)

type Params struct {
	Hash     string `from:"hash" json:"hash"`
	FromAddr string `form:"from_addr" json:"from_addr"`
	ToAddr   string `form:"to_addr" json:"to_addr"`
	Value    string `form:"value" json:"value"`
	Currency string `from:"currency" json:"currency"`
	Page     int    `form:"page" json:"page"`
	Limit    int    `form:"limit" json:"limit"`
	ID       int    `form:"id" json:"id"`
	Order    string `form:"order" json:"order"`
	OrderDir string `form:"order_dir" json:"order_dir"`
}

func (r *Repository) SaveTxDataToDB(data Params) error {
	queryBuilder := squirrel.Insert("transactions").
		Columns("hash", "from_addr", "to_addr", "value", "currency").
		Values(data.Hash, data.FromAddr, data.ToAddr, data.Value, data.Currency).
		Suffix("ON CONFLICT (hash) DO NOTHING")

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

type ResponseParams struct {
	Hash     string `from:"hash" json:"hash"`
	FromAddr string `form:"from_addr" json:"from_addr"`
	ToAddr   string `form:"to_addr" json:"to_addr"`
	Value    string `form:"value" json:"value"`
	Currency string `from:"currency" json:"currency"`
	ID       int    `form:"id" json:"id"`
}

func (r *Repository) GetTxFromDB(params Params) ([]ResponseParams, error) {

	offset := (params.Page - 1) * params.Limit

	queryBuilder := squirrel.Select("hash", "from_addr", "to_addr", "value", "currency", "id").
		From("transactions").
		Limit(uint64(params.Limit)).
		Offset(uint64(offset))

	if params.FromAddr != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"from_addr": "%" + params.FromAddr + "%"})
	}

	if params.ToAddr != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"to_addr": "%" + params.ToAddr + "%"})
	}

	if params.Order != "" {
		orderDirection := strings.ToUpper(params.OrderDir)

		if orderDirection != "ASC" && orderDirection != "DESC" {
			orderDirection = "ASC"
		}

		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("%s %s", params.Order, orderDirection))
	} else {
		queryBuilder = queryBuilder.OrderBy("id ASC")
	}

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	var rates []ResponseParams

	// GETTING RESULTS
	for rows.Next() {
		var rate ResponseParams
		if err = rows.Scan(&rate.Hash, &rate.FromAddr, &rate.ToAddr, &rate.Value, &rate.Currency, &rate.ID); err != nil {
			log.Fatal(err)
		}

		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rates, nil
}
