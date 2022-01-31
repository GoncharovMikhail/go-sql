package test

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/jackc/pgx/v4"
	"log"
)

const (
	PGURL      = `postgresql://localhost:5432/postgres`
	PGUsername = `postgres`
	PGPassword = `postgres`
)

var (
	CTX = context.Background()

	connConfig      = util.MustParseConfig(PGURL)
	connConfigToUse = getConnConfigToUse(connConfig)
)

//todo переделать)))
func getConnConfigToUse(config *pgx.ConnConfig) *pgx.ConnConfig {
	config.User = PGUsername
	config.Password = PGPassword
	return config
}

func GetTX(db *sql.DB) *sql.Tx {
	tx, err := db.BeginTx(CTX, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		log.Panic(err)
	}
	return tx
}
