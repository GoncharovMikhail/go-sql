package sql

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/stdlib"
	"gotest.tools/assert"
	"log"
	"testing"
)

const (
	dbURL      = `postgresql://localhost:5432/postgres`
	dbUsername = `postgres`
	dbPassword = `postgres`
)

//eagerInit
var (
	ctx = context.Background()
)

//lateInit
var (
	entityToSave *entity.UserDataEntity
	uuidUsername uuid.UUID
	repository   user.SQLUserRepository
	tx           *sql.Tx
)

func init() {
	var err error
	uuidUsername, err = uuid.NewV1()
	if err != nil {
		log.Panic(err)
	}
	entityToSave = &entity.UserDataEntity{
		Username: uuidUsername.String(),
		Password: uuidUsername.String(),
	}
	config := util.MustParseConfig(dbURL)
	config.User = dbUsername
	config.Password = dbPassword
	db := stdlib.OpenDB(config)
	defer util.MustCloseDb(db)
	tx = util.MustBeginTx(
		ctx,
		db,
		&sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		},
	)
	repository = NewPostgresUserRepository()
}

func TestPostgresUserRepository_Save(t *testing.T) {
	save, errors := repository.SaveInTx(
		ctx,
		entityToSave,
		tx,
	)
	assert.Assert(t, errors == nil)
	assert.Assert(t, save != nil)
	assert.Assert(t, &save.Id != nil)
}

func TestPostgresUserRepository_FindOneByUsernameInTx(t *testing.T) {
	TestPostgresUserRepository_Save(t)
	config := util.MustParseConfig(dbURL)
	config.User = dbUsername
	config.Password = dbPassword
	db := stdlib.OpenDB(config)
	defer util.MustCloseDb(db)
	tx = util.MustBeginTx(
		ctx,
		db,
		&sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		},
	)
	result, ok, errors := repository.FindOneByUsernameInTx(ctx, entityToSave.Username, tx)
	assert.Assert(t, errors == nil)
	assert.Assert(t, ok == true)
	assert.Assert(t, result != nil)
}
