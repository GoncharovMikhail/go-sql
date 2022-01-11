package main

import (
	"database/sql"
	"fmt"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	sql2 "github.com/GoncharovMikhail/go-sql/pkg/db/user/impl/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func main() {
	var r user.SQLUserRepository = sql2.NewPostgresUserRepository(nil, nil, nil)
	print(r)
	config, err := pgx.ParseConfig("postgresql://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	config.User = "postgres"
	config.Password = "postgres"
	db := stdlib.OpenDB(*config)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Print("WP")
}
