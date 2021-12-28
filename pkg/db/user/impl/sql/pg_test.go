package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"gotest.tools/assert"
	"testing"
)

var (
	dbURL      = `postgresql://localhost:5432/postgres`
	dbUsername = `postgres`
	dbPassword = `102030AaBb`

	uuidUsername, _ = uuid.NewV1()
	entityToSave    = &user.UserEntity{
		Username: uuidUsername.String(),
		Password: "pwd",
		RestoreData: &user.RestoreData{
			Email: uuidUsername.String(),
			PhoneNumber: sql.NullString{
				String: "phone_number",
				Valid:  true,
			},
		},
	}
)

//lateInit
var (
	repository user.UserRepository
)

func init() {
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		panic(err)
	}
	config.User = dbUsername
	config.Password = dbPassword
	openDB := stdlib.OpenDB(*config)
	repository = &PostgresUserRepository{
		openDB,
	}
}

func TestPostgresUserRepository_Save_NoRestoreData(t *testing.T) {
	uuidUsername, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}
	save, err := repository.Save(
		context.Background(),
		&user.UserEntity{
			Username:    uuidUsername.String(),
			Password:    "pwd",
			RestoreData: nil,
		},
	)
	assert.NilError(t, err)
	assert.Assert(t, save != nil)
	assert.Assert(t, &save.Id != nil)
}

func TestPostgresUserRepository_Save_EmailInRestoreData(t *testing.T) {
	uuidUsername, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}
	entityToSave := &user.UserEntity{
		Username: uuidUsername.String(),
		Password: "pwd",
		RestoreData: &user.RestoreData{
			Email: uuidUsername.String(),
		},
	}
	save, err := repository.Save(
		context.Background(),
		entityToSave,
	)
	assert.NilError(t, err)
	assert.Assert(t, save != nil)
	assert.Assert(t, &entityToSave.Id != nil)
	assert.Assert(t, fmt.Sprintf("%v", entityToSave.Id) == fmt.Sprintf("%v", entityToSave.RestoreData.UserId))
}

func TestPostgresUserRepository_Save_FullyQualifiedRestoreData(t *testing.T) {
	uuidUsername, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}
	save, err := repository.Save(
		context.Background(),
		&user.UserEntity{
			Username: uuidUsername.String(),
			Password: "pwd",
			RestoreData: &user.RestoreData{
				Email: uuidUsername.String(),
				PhoneNumber: sql.NullString{
					String: "phone_number",
					Valid:  true,
				},
			},
		},
	)
	assert.NilError(t, err)
	assert.Assert(t, save != nil)
	assert.Assert(t, &entityToSave.Id != nil)
	assert.Assert(t, fmt.Sprintf("%v", entityToSave.Id) == fmt.Sprintf("%v", entityToSave.RestoreData.UserId))
}
