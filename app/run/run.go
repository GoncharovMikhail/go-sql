package run

import (
	builtInSql "database/sql"
)

func Run() error {
	return nil
}

func closeDb(pg *builtInSql.DB) {
	err := pg.Close()
	if err != nil {
		panic(err)
	}
}
