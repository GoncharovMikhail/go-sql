package main

import (
	"fmt"
	"github.com/jackc/pgx/v4"
)

func main() {
	config, err := pgx.ParseConfig("postgresql://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	config.User = "postgres"
	config.Password = "postgres"
	fmt.Println(config.User)
}
