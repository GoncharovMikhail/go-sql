package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	version  = "14.1-alpine"
	postgres = "postgres"
	password = "password"
	port     = "5432"
)

func CreateTestContainer(ctx context.Context, dbname string) (testcontainers.Container, *sql.DB, error) {
	goModDir, err := util.GetGoModDir()
	if err != nil {
		return nil, nil, err
	}
	initDbFiles, errors := util.ListAllFilesMatchingPatternsAllOverOsFromSpecifiedDir(
		goModDir,
		func(info os.FileInfo) bool { return !info.IsDir() },
		util.Conjunction,
		".*/resources/migrations.*", "up.sql",
	)
	if errors != nil {
		panic(errors)
	}
	//todo
	initDbDir := filepath.Dir(initDbFiles[0])
	var port = "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", port.Port(), dbname)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},

			Env: map[string]string{
				"POSTGRES_DB":       dbname,
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "password",
			},
			BindMounts: map[string]string{
				"/docker-entrypoint-initdb.d": initDbDir,
			},
			WaitingFor: wait.ForSQL(nat.Port(port), "postgres", dbURL).Timeout(time.Second * 5),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort)

	url := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbname)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}

func TestSaveInTx(t *testing.T) {
	_, db, err := CreateTestContainer(test.CTX, postgres)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	tx := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	inTx, errors := SaveInTx(test.CTX, &entity.UserDataEntity{
		Username: "user",
		Password: "password",
	}, tx)
	if errors != nil {
		panic(errors)
	}
	log.Println(inTx)
}
