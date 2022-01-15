package sql

import (
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	"gotest.tools/assert"
	"log"
	"testing"
)

var (
	entityToSave *entity.UserDataEntity
	uuidUsername uuid.UUID
	repository   user.SQLUserRepository
)

func init() {
	var e error
	uuidUsername, e = uuid.NewV1()
	if e != nil {
		log.Panic(e)
	}
	entityToSave = &entity.UserDataEntity{
		Username: uuidUsername.String(),
		Password: uuidUsername.String(),
	}
	repository = NewPostgresUserRepository()
}

func TestPostgresUserRepository_Save(t *testing.T) {
	save, errorz := saveInTx()
	saveAsserts(t, save, errorz)
}

func TestPostgresUserRepository_FindOneByUsernameInTx(t *testing.T) {
	save, errorz := saveInTx()
	saveAsserts(t, save, errorz)
	result, ok, errorz := repository.FindOneByUsernameInTx(test.CTX, entityToSave.Username, test.TX)
	assert.Assert(t, errorz == nil)
	assert.Assert(t, ok == true)
	assert.Assert(t, result != nil)
}

func saveInTx() (*entity.UserDataEntity, errors.Errors) {
	return repository.SaveInTx(
		test.CTX,
		entityToSave,
		test.TX,
	)
}

func saveAsserts(t *testing.T, save *entity.UserDataEntity, errorz errors.Errors) {
	assert.Assert(t, errorz == nil)
	assert.Assert(t, save != nil)
	assert.Assert(t, &save.Id != nil)
}
