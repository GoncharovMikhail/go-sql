package util

import (
	"log"
	"os"
	"testing"
)

func TestListAllFilesMatchingPatternsAllOverOsFromSpecifiedDir(t *testing.T) {
	pwd, _ := os.Getwd()
	files, errors := ListAllFilesMatchingPatternsAllOverOsFromSpecifiedDir(
		pwd,
		func(info os.FileInfo) bool { return !info.IsDir() },
		Conjunction,
		".*/resources.*", "\\.sql",
	)
	if errors != nil {
		panic(errors)
	}
	for _, file := range files {
		log.Print(file)
	}
}
