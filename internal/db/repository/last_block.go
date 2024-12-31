package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

func (r *Repository) SaveLastBlockToDB(block int64) error {

	queryBuilder := squirrel.Update("blocks").
		Set("block_number", block).
		Where(squirrel.Eq{"id": 1})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	log.Println("Block number saved to db:", block)

	return nil
}

func (r *Repository) GetLastBlockFromDB() (int64, error) {
	queryBuilder := squirrel.Select("block_number").
		From("blocks").
		Where(squirrel.Eq{"id": 1})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return 0, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	defer rows.Close()

	var lastBlock int64
	if rows.Next() {
		if err = rows.Scan(&lastBlock); err != nil {
			return 0, fmt.Errorf("failed to scan result: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return lastBlock, nil
}
