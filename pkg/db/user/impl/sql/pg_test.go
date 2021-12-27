package sql

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strings"
	"testing"
)

type pgContainer struct {
	testcontainers.Container
}

var (
	packageName = "db"

	//lateInit
	pgC *pgContainer
)

// Create the Postgres TestContainer
func init() {

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootDir := strings.Replace(workingDir, packageName, "", 1)
	mountFrom := fmt.Sprintf("%sresources/migrations/001__schema.up.sql", rootDir)
	mountTo := "/docker-entrypoint-initdb.d/init.sql"
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:11.6-alpine",
		ExposedPorts: []string{"5432/tcp"},
		BindMounts:   map[string]string{mountFrom: mountTo},
		Env: map[string]string{
			"POSTGRES_DB": os.Getenv("DBNAME"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	tcC, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		panic(err)
	}
	pgC = &pgContainer{tcC}
	print(pgC)
}

func TestMain(m *testing.M) {
	defer func(container pgContainer, ctx context.Context) {
		_ = container.Terminate(ctx)
	}(*pgC, context.Background())
	// Work out the path to the 'scripts' directory and set mount strings
	packageName := "database"
	workingDir, _ := os.Getwd()
	rootDir := strings.Replace(workingDir, packageName, "", 1)
	mountFrom := fmt.Sprintf("%s/scripts/init.sql", rootDir)
	mountTo := "/docker-entrypoint-initdb.d/init.sql"
	// Create the Postgres TestContainer
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:11.6-alpine",
		ExposedPorts: []string{"5432/tcp"},
		BindMounts:   map[string]string{mountFrom: mountTo},
		Env: map[string]string{
			"POSTGRES_DB": os.Getenv("DBNAME"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		// Panic and fail since there isn't much we can do if the container doesn't start
		panic(err)
	}
	defer func(postgresC testcontainers.Container, ctx context.Context) {
		err := postgresC.Terminate(ctx)
		if err != nil {
			panic(err)
		}
	}(postgresC, ctx)
	// Get the port mapped to 5432 and set as ENV
	p, _ := postgresC.MappedPort(ctx, "5432")
	_ = os.Setenv("DBPORT", p.Port())
	exitVal := m.Run()
	os.Exit(exitVal)

}
