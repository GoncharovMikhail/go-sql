package main

import (
	"github.com/GoncharovMikhail/go-sql/app/run"
	"log"
)

func main() {
	err := run.Run()
	if err != nil {
		log.Fatalf("Err :<%s>", err)
	}
}
