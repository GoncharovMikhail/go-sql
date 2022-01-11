package sql

import (
	"github.com/testcontainers/testcontainers-go"
)

type pgContainer struct {
	testcontainers.Container
	URI string
}

/*
func getPgContainer(ctx context.Context) (*pgContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp", "5432/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080"),
		Name:         dbUsername,
		AutoRemove:   true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://root@%s:%s", hostIP, mappedPort.Port())

	return &pgContainer{Container: container, URI: uri}, nil
}

func initDB(ctx context.Context, db *sql.DB) error {
	// Actual SQL for initializing the database should probably live elsewhere
	const query = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user"
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT         NOT NULL
);

CREATE TABLE IF NOT EXISTS restore_data
(
    user_id          UUID         NOT NULL UNIQUE,
    CONSTRAINT fk__restore_data__user__one_to_one
        FOREIGN KEY (user_id)
            REFERENCES "user" (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    email            VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS authority
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_authority
(
    user_id      UUID NOT NULL UNIQUE,
    CONSTRAINT fk__user_authority__user__many_to_many
        FOREIGN KEY (user_id)
            REFERENCES "user" (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    authority_id UUID NOT NULL UNIQUE,
    CONSTRAINT fk__user_authority__authority__many_to_many
        FOREIGN KEY (authority_id)
            REFERENCES authority (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);`
	_, err := db.ExecContext(ctx, query)

	return err
}

func TestIntegrationDBInsertSelect(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	cdbContainer, err := getPgContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func(cdbContainer *pgContainer, ctx context.Context) {
		err := cdbContainer.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(cdbContainer, ctx)

	db, err := sql.Open("pgx", cdbContainer.URI)
	if err != nil {
		t.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	err = initDB(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	repository := &postgresUserRepository{db}
	save, err := repository.SaveInTx(ctx, &entity.UserEntity{
		Username:          "",
		Password:          "",
		RestoreDataEntity: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Assert(t, save != nil)
}

*/
