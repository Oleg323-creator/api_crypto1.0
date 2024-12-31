package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

type TxData struct {
	Hash     string
	FromAddr string
	ToAddr   string
	Value    string
	Currency string
}

func (r *Repository) SaveTxDataToDB(data TxData) error {
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
