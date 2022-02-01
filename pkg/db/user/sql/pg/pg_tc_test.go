package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	version  = "14.1-alpine"
	postgres = "postgres"
	port     = "5432"
)

var (
	container testcontainers.Container
	db        *sql.DB
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	initDbFiles, errors := util.ListAllFilesMatchingPatternsAllOverOsFromSpecifiedDir(
		pwd,
		func(info os.FileInfo) bool { return !info.IsDir() },
		util.Conjunction,
		".*/resources.*", "\\.sql",
	)
	if errors != nil {
		panic(errors)
	}
	//todo
	initDbDir := filepath.Dir(initDbFiles[0])
	request := testcontainers.ContainerRequest{
		Image:      postgres + ":" + version,
		Entrypoint: nil,
		Env: map[string]string{
			"POSTGRES_DB":       postgres,
			"PGUSER":            postgres,
			"POSTGRES_USER":     postgres,
			"POSTGRES_PASSWORD": postgres,
			//"PGDATA":            postgres,
		},
		ExposedPorts: []string{port},
		//todo ЭТО ПРОСТО ПИЗДЕЦ БЛЯТЬ В ЧЕМ РАЗНИЧА МЕЖДУ BindMounts И VolumeMounts
		BindMounts: map[string]string{
			"/docker-entrypoint-initdb.d": initDbDir,
		},
		Name: postgres,
		//User:       postgres,
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
		//AutoRemove: true,
	}
	container, err = testcontainers.GenericContainer(
		test.CTX,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		},
	)

	if err != nil {
		log.Panicln(err)
	}
	ports, _ := container.Ports(test.CTX)
	for porT := range ports {
		log.Println(porT)
	}
	host, _ := container.Host(test.CTX)
	log.Println(host)
	// DB
	defer mustTerminate(test.CTX, container)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, postgres, postgres)

	// todo сдаюсь, pgx это пиздец)))
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
}

func mustTerminate(ctx context.Context, container testcontainers.Container) {
	err := container.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

func TestSaveInTx(t *testing.T) {
}
