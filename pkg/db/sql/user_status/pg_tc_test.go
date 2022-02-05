package user_status

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/tc"
	userData "github.com/GoncharovMikhail/go-sql/pkg/db/sql/user"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"gotest.tools/assert"
	"log"
	"testing"
)

var (
	randomUuid uuid.UUID
	container  testcontainers.Container
	db         *sql.DB
	ude        *entity.UserDataEntity
)

func init() {
	randomUuid = tc.InitUuid()
	container = tc.InitContainer()
	db = tc.InitDb()

	ude = &entity.UserDataEntity{
		Username: randomUuid.String(),
		Password: randomUuid.String(),
	}
	ude = MustSaveOrUpdateUserDataEntity(ude, db)
}

var (
	NullBoolTrue  = sql.NullBool{Bool: true, Valid: true}
	NullBoolFalse = sql.NullBool{Bool: false, Valid: true}
)

func TestSaveInTx(t *testing.T) {
	use := &entity.UserStatusEntity{
		UserId:                ude.Id,
		AccountNonExpired:     NullBoolTrue,
		AccountNonLocked:      NullBoolTrue,
		CredentialsNonExpired: NullBoolFalse,
		Enabled:               NullBoolTrue,
	}
	use = MustSaveOrUpdateUserStatusEntity(use, db)

	assert.Assert(t, ude.Id == use.UserId)
	assert.Assert(t, NullBoolTrue == use.AccountNonExpired)
	assert.Assert(t, NullBoolTrue == use.AccountNonLocked)
	assert.Assert(t, NullBoolFalse == use.CredentialsNonExpired)
	assert.Assert(t, NullBoolTrue == use.Enabled)
}

func TestUpdateInTx(t *testing.T) {
	// Save
	use := &entity.UserStatusEntity{
		UserId:                ude.Id,
		AccountNonExpired:     NullBoolTrue,
		AccountNonLocked:      NullBoolTrue,
		CredentialsNonExpired: NullBoolFalse,
		Enabled:               NullBoolTrue,
	}
	use = MustSaveOrUpdateUserStatusEntity(use, db)
	// Update
	use.AccountNonExpired = NullBoolFalse

	assert.Assert(t, ude.Id == use.UserId)
	assert.Assert(t, NullBoolFalse == use.AccountNonExpired)
	assert.Assert(t, NullBoolTrue == use.AccountNonLocked)
	assert.Assert(t, NullBoolFalse == use.CredentialsNonExpired)
	assert.Assert(t, NullBoolTrue == use.Enabled)
}

func TestFindOneByUsernameInTx(t *testing.T) {
	// Save
	use := &entity.UserStatusEntity{
		UserId:                ude.Id,
		AccountNonExpired:     NullBoolTrue,
		AccountNonLocked:      NullBoolTrue,
		CredentialsNonExpired: NullBoolTrue,
		Enabled:               NullBoolTrue,
	}
	use = MustSaveOrUpdateUserStatusEntity(use, db)
	// FindOneByUsernameInTx
	txToFindByUsername := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	var present bool
	var errorz errors.Errors
	use, present, errorz = FindOneByUsernameInTx(test.CTX, randomUuid.String(), txToFindByUsername)
	if errorz != nil {
		errorz, txToFindByUsername = util.TxRollbackErrorHandle(errorz.Get(), txToFindByUsername)
		log.Panic(errorz)
	}
	if !present {
		log.Panicf("couldn't find just saved user status entity for username: %s", randomUuid.String())
	}
}

func MustSaveOrUpdateUserDataEntity(ude *entity.UserDataEntity, db *sql.DB) *entity.UserDataEntity {
	txToSaveUserDataEntity := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	var errorz errors.Errors
	ude, errorz, txToSaveUserDataEntity = userData.SaveOrUpdateInTx(test.CTX, ude, txToSaveUserDataEntity)
	if errorz != nil {
		if errorz, _ = util.TxRollbackErrorHandle(errorz.Get(), txToSaveUserDataEntity); errorz != nil {
			log.Panic(errorz)
		}
	}
	util.MustCommitTx(txToSaveUserDataEntity)
	return ude
}

func MustSaveOrUpdateUserStatusEntity(use *entity.UserStatusEntity, db *sql.DB) *entity.UserStatusEntity {
	txToSaveUserStatusEntity := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	var errorz errors.Errors
	use, errorz, txToSaveUserStatusEntity = SaveOrUpdateInTx(test.CTX, use, txToSaveUserStatusEntity)
	if errorz != nil {
		errorz, txToSaveUserStatusEntity = util.TxRollbackErrorHandle(errorz.Get(), txToSaveUserStatusEntity)
		log.Panic(errorz)
	}
	util.MustCommitTx(txToSaveUserStatusEntity)
	return use
}
