package main

import (
	"log"
	"sql/app/run"
)

func main() {
	err := run.Run()
	if err != nil {
		log.Fatalf("Err :<%s>", err)
	}
}
