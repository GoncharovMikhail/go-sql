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
	repository = NewPostgresUserDataRepository()
}

func TestPostgresUserRepository_SaveInTx(t *testing.T) {
	save, errorz := saveInTx()
	saveInTxAsserts(t, save, errorz)
}

func TestPostgresUserRepository_FindOneByUsernameInTx(t *testing.T) {
	save, errorz := saveInTx()
	saveInTxAsserts(t, save, errorz)
	result, ok, errorz := repository.FindOneByUsernameInTx(test.CTX, entityToSave.Username, test.GetTX(test.DB))
	findOneByUsernameInTxAsserts(t, ok, result, errorz)
}

func saveInTx() (*entity.UserDataEntity, errors.Errors) {
	var e error
	uuidUsername, e = uuid.NewV1()
	if e != nil {
		log.Panic(e)
	}
	entityToSave = &entity.UserDataEntity{
		Username: uuidUsername.String(),
		Password: uuidUsername.String(),
	}
	return repository.SaveInTx(
		test.CTX,
		entityToSave,
		test.GetTX(test.DB),
	)
}

func saveInTxAsserts(t *testing.T, save *entity.UserDataEntity, errorz errors.Errors) {
	assert.Assert(t, errorz == nil)
	assert.Assert(t, save != nil)
	assert.Assert(t, &save.Id != nil)
}

func findOneByUsernameInTxAsserts(t *testing.T, ok bool, result *entity.UserDataEntity, errorz errors.Errors) {
	assert.Assert(t, result != nil)
	assert.Assert(t, ok == true)
	assert.Assert(t, errorz == nil)
	//assert.Assert(t, result.Username == desiredUsername)
}
