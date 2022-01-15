package db

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func GetDb(connConfig *pgx.ConnConfig) (*sql.DB, errors.Errors) {
	db := stdlib.OpenDB(*connConfig)
	if errorz := pingDb(db); errorz != nil {
		return nil, errorz
	}
	return db, nil
}

func pingDb(db *sql.DB) errors.Errors {
	err := db.Ping()
	if err != nil {
		return errors.NewErrors(
			errors.BuildSimpleErrMsg("err", err),
			err,
			nil,
		)
	}
	return nil
}
