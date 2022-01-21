package sql

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/errors"
	user "github.com/GoncharovMikhail/go-sql/pkg/db/user/impl/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	pgUuidType "github.com/jackc/pgtype/ext/gofrs-uuid"
	"gotest.tools/assert"
	"log"
	"testing"
)

//todo REFACTOR
//eagerInit
var (
	ur           = user.NewPostgresUserDataRepository()
	repository   = NewPostgresRestoreDataRepository()
	entityToSave = &entity.RestoreDataEntity{}
)

//lateInit
var (
	userUuid pgUuidType.UUID
)

func setUserDataProps() {
	var err error
	uuidUuid, err := uuid.NewV1()
	if err != nil {
		log.Panic(err)
	}
	userUuid = pgUuidType.UUID{
		UUID: uuidUuid,
	}
	entityToSave = &entity.RestoreDataEntity{
		UserId: userUuid,
		Email:  uuidUuid.String(),
	}
}

func TestPostgresRestoreDataRepository_SaveInTx_Fails(t *testing.T) {
	result, errorz := saveInTx_Fails()
	saveInTxAsserts_Fails(t, result, errorz)
}

func TestPostgresRestoreDataRepository_SaveInTx_NullPhoneNumber_Ok(t *testing.T) {
	userDataEntity, _ := saveMockUser_Ok()
	userId := userDataEntity.Id
	entityToSave = &entity.RestoreDataEntity{
		UserId: userId,
		Email:  userId.UUID.String(),
	}
	entityToSave.Email = userId.UUID.String()
	result, errorz := repository.SaveInTx(
		test.CTX,
		entityToSave,
		test.GetTX(test.DB),
	)
	postSave_NonNullResult(t, result, errorz)
}

func TestPostgresRestoreDataRepository_SaveInTx_NonNullPhoneNumber_Ok(t *testing.T) {
	userDataEntity, errorz := saveMockUser_Ok()
	assert.Assert(t, errorz == nil)
	userId := userDataEntity.Id
	entityToSave = &entity.RestoreDataEntity{
		UserId: userId,
		Email:  userId.UUID.String(),
		PhoneNumber: sql.NullString{
			String: userId.UUID.String(),
			Valid:  true,
		},
	}
	entityToSave.Email = userId.UUID.String()
	result, errorz := repository.SaveInTx(
		test.CTX,
		entityToSave,
		test.GetTX(test.DB),
	)
	postSave_NonNullResult(t, result, errorz)
	value, err := result.PhoneNumber.Value()
	assert.NilError(t, err)
	assert.Assert(t, value != nil)
}

func TestPostgresRestoreDataRepository_FindOneByUsernameInTx_Ok(t *testing.T) {
	userDataEntity, errorz := saveMockUser_Ok()
	assert.Assert(t, errorz == nil)
	userId := userDataEntity.Id
	entityToSave = &entity.RestoreDataEntity{
		UserId: userId,
		Email:  userId.UUID.String(),
		PhoneNumber: sql.NullString{
			String: userId.UUID.String(),
			Valid:  true,
		},
	}
	entityToSave.Email = userId.UUID.String()
	result, errorz := repository.SaveInTx(
		test.CTX,
		entityToSave,
		test.GetTX(test.DB),
	)
	postSave_NonNullResult(t, result, errorz)
	result, ok, err := repository.FindOneByUsernameInTx(test.CTX, userDataEntity.Username, test.GetTX(test.DB))
	assert.NilError(t, err)
	assert.Assert(t, ok == true)
	assert.Assert(t, result != nil)
	assert.Assert(t, result.UserId == userDataEntity.Id)
}

func saveInTx_Fails() (*entity.RestoreDataEntity, errors.Errors) {
	setUserDataProps()
	return repository.SaveInTx(
		test.CTX,
		entityToSave,
		test.GetTX(test.DB),
	)
}

func saveInTxAsserts_Fails(t *testing.T, result *entity.RestoreDataEntity, errorz errors.Errors) {
	assert.Assert(t, result == nil)
	assert.Assert(t, errorz != nil)
}

func saveMockUser_Ok() (*entity.UserDataEntity, errors.Errors) {
	uuidUuid, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}
	uuidString := uuidUuid.String()
	return ur.SaveInTx(
		test.CTX,
		&entity.UserDataEntity{
			Username: uuidString,
			Password: uuidString,
		},
		test.GetTX(test.DB),
	)
}

func postSave_NonNullResult(t *testing.T, result *entity.RestoreDataEntity, errorz errors.Errors) {
	assert.Assert(t, errorz == nil)
	assert.Assert(t, result != nil)
}
