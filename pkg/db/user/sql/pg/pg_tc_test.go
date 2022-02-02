package pg

import (
	"database/sql"
	"fmt"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	version  = "14.1-alpine"
	postgres = "postgres"
	username = postgres
	password = "password"
	port     = "5432"
)

var (
	container testcontainers.Container
	db        *sql.DB
)

func init() {
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
	initDbDir := filepath.Dir(initDbFiles[0])

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        postgres + ":" + version,
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},

			Env: map[string]string{
				"POSTGRES_DB":       postgres,
				"POSTGRES_USER":     username,
				"POSTGRES_PASSWORD": password,
			},
			BindMounts: map[string]string{
				"/docker-entrypoint-initdb.d": initDbDir,
			},
			WaitingFor: wait.ForSQL(port, postgres, func(port nat.Port) string {
				return fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", postgres, username, password, port.Port(), postgres)
			}).Timeout(time.Second * 5),
		},
		Started: true,
	}
	container, err = testcontainers.GenericContainer(test.CTX, req)
	if err != nil {
		log.Panicf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(test.CTX, port)
	if err != nil {
		log.Panicf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort)

	url := fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", postgres, username, password, mappedPort.Port(), postgres)
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Panicf("failed to establish database connection: %s", err)
	}
}

func TestSaveInTx(t *testing.T) {
	const username = "username"
	const password = "password"
	tx := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	saved, errors := SaveInTx(
		test.CTX,
		&entity.UserDataEntity{
			Username: username,
			Password: password,
		},
		tx,
	)
	if errors != nil {
		log.Panic(errors)
	}
	assert.Assert(t, saved.Username == username)
	assert.Assert(t, saved.Password == password)
	assert.Assert(t, saved.Id.UUID.String() != "")
}
