package db

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func GetDb(connConfig *pgx.ConnConfig) (*sql.DB, error) {
	db := stdlib.OpenDB(*connConfig)
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
