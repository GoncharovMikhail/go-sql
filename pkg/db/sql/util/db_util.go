package util

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/jackc/pgx/v4"
	"log"
)

func MustCloseDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Panic(err)
	}
}

func MustParseConfig(dbURL string) *pgx.ConnConfig {
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		log.Panic(err)
	}
	return config
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

func MustCommitTx(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func MustPing(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Panic(err)
	}
}

func TxRollbackErrorHandle(err error, tx *sql.Tx) (errors.Errors, *sql.Tx) {
	if err == nil {
		return nil,
			tx
	}
	errTxRollback := tx.Rollback()
	if errTxRollback != nil {
		return errors.NewErrors(
				errors.BuildSimpleErrMsg("err", err),
				err,
				errors.NewErrors(
					errors.BuildSimpleErrMsg("errTxRollback", errTxRollback),
					errTxRollback,
					nil,
				),
			),
			tx
	}
	return errors.NewErrors(
			errors.BuildSimpleErrMsg("err", err),
			err,
			nil,
		),
		tx

}
