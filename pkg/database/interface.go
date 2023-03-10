package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (db *PgClient) Execute(sql string) error {
	_, err := db.pool.Exec(context.Background(), sql)
	if err != nil {
		return err
	}
	return nil
}

func (db *PgClient) QueryRow(sql string) pgx.Row {
	return db.pool.QueryRow(context.Background(), sql)
}

func (db *PgClient) QueryRows(sql string) (pgx.Rows, error) {
	rows, err := db.pool.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
