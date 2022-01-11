package util

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"log"
)

func MustCloseDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Panic(err)
	}
}

func MustParseConfig(dbURL string) pgx.ConnConfig {
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		log.Panic(err)
	}
	return *config
}

func MustBeginTx(ctx context.Context, db *sql.DB, options *sql.TxOptions) *sql.Tx {
	tx, err := db.BeginTx(
		ctx,
		options,
	)
	if err != nil {
		log.Panic(err)
	}
	return tx
}
