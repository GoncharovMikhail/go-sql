package tc

import (
	"database/sql"
	"fmt"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/docker/go-connections/nat"
	"github.com/gofrs/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	version    = "14.1-alpine"
	postgres   = "postgres"
	pgUsername = postgres
	pgPassword = "password"
	pgPort     = "5432"
)

var (
	randomUuid uuid.UUID
	container  testcontainers.Container
	db         *sql.DB
)

func InitUuid() uuid.UUID {
	var err error
	randomUuid, err = uuid.NewV1()
	if err != nil {
		log.Panicf("couln't generate random uuid. err: %s", err)
	}
	return randomUuid
}

func getInitDbScriptsDir() string {
	goModDir, err := util.GetGoModDir()
	if err != nil {
		log.Panic(err)
	}
	initDbFiles, errors := util.ListAllFilesMatchingPatternsAllOverOsFromSpecifiedDir(
		goModDir,
		func(info os.FileInfo) bool { return !info.IsDir() },
		util.Conjunction,
		".*/resources/migrations.*", "up.sql",
	)
	if errors != nil {
		log.Panic(errors)
	}
	//todo
	return filepath.Dir(initDbFiles[0])
}

func InitContainer() testcontainers.Container {
	initDbScriptsDir := getInitDbScriptsDir()
	var err error
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        postgres + ":" + version,
			ExposedPorts: []string{pgPort},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env: map[string]string{
				"POSTGRES_DB":       postgres,
				"POSTGRES_USER":     pgUsername,
				"POSTGRES_PASSWORD": pgPassword,
			},
			BindMounts: map[string]string{
				"/docker-entrypoint-initdb.d": initDbScriptsDir,
			},
			WaitingFor: wait.ForSQL(pgPort, postgres, func(port nat.Port) string {
				return fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", postgres, pgUsername, pgPassword, port.Port(), postgres)
			}).Timeout(time.Second * 5),
			AutoRemove: true,
		},
		Started: true,
	}
	container, err = testcontainers.GenericContainer(test.CTX, req)
	if err != nil {
		log.Panicf("failed to start container: %s", err)
	}
	return container
}

func getMappedPort() nat.Port {
	if container == nil {
		log.Panicf("container is nil")
	}
	mappedPort, err := container.MappedPort(test.CTX, pgPort)
	if err != nil {
		log.Panicf("failed to get container external pgPort: %s", err)
	}
	return mappedPort
}

func InitDb() *sql.DB {
	mappedPort := getMappedPort()
	var err error
	url := fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", postgres, pgUsername, pgPassword, mappedPort.Port(), postgres)
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Panicf("failed to oped db: %s", err)
	}
	return db
}
