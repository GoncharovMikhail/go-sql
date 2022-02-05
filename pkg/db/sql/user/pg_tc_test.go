//todo почему вместе не пробегают? По отдельности все ок
package user

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/tc"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	pgUuidType "github.com/jackc/pgtype/ext/gofrs-uuid"
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
)

func init() {
	randomUuid = tc.InitUuid()
	container = tc.InitContainer()
	db = tc.InitDb()
}

func TestSaveInTx(t *testing.T) {
	uuidStringToUse := randomUuid.String()
	savedEntity, tx := mustSaveOrUpdateInTx(uuid.Nil, uuidStringToUse, uuidStringToUse)
	util.MustCommitTx(tx)

	assert.Assert(t, savedEntity.Username == uuidStringToUse)
	assert.Assert(t, savedEntity.Password == uuidStringToUse)
	assert.Assert(t, savedEntity.Id.UUID.String() != "")
}

func TestUpdateInTx(t *testing.T) {
	// Save
	randomUuidStringValue := randomUuid.String()
	savedEntity, txSaved := mustSaveOrUpdateInTx(uuid.Nil, randomUuidStringValue, randomUuidStringValue)
	util.MustCommitTx(txSaved)
	// Update
	randomUuidStringValue = savedEntity.Id.UUID.String()
	updatedEntity, txUpdated := mustSaveOrUpdateInTx(savedEntity.Id.UUID, randomUuidStringValue, randomUuidStringValue)
	util.MustCommitTx(txUpdated)

	assert.Assert(t, savedEntity != nil)
	assert.Assert(t, randomUuidStringValue == updatedEntity.Id.UUID.String())
	assert.Assert(t, randomUuidStringValue == updatedEntity.Username)
	assert.Assert(t, randomUuidStringValue == updatedEntity.Password)
}

func TestFindOneByUsernameInTx(t *testing.T) {
	uuidStringToUse := randomUuid.String()
	savedUserEntity, txToSaveUser := mustSaveOrUpdateInTx(uuid.Nil, uuidStringToUse, uuidStringToUse)
	util.MustCommitTx(txToSaveUser)
	txToFindUser := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	foundUserEntity, exists, errorz, tx := FindOneByUsernameInTx(test.CTX, uuidStringToUse, txToFindUser)
	if errorz != nil {
		log.Panicf("err: %s", errorz)
	}
	if !exists {
		log.Panicf("couldn't find just saved user with username: %s", uuidStringToUse)
	}
	util.MustCommitTx(tx)

	assert.Assert(t, foundUserEntity.Username == savedUserEntity.Username)
	assert.Assert(t, foundUserEntity.Password == savedUserEntity.Password)
	assert.Assert(t, foundUserEntity.Id == savedUserEntity.Id)
}

func mustSaveOrUpdateInTx(id uuid.UUID, username, password string) (*entity.UserDataEntity, *sql.Tx) {
	ude := &entity.UserDataEntity{
		Username: username,
		Password: password,
	}
	if id != uuid.Nil {
		ude.Id = pgUuidType.UUID{
			UUID: id,
		}
	}

	tx := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	saved, errorz, tx := SaveOrUpdateInTx(
		test.CTX,
		ude,
		tx,
	)
	if errorz != nil {
		log.Panic(errorz)
	}
	return saved, tx
}
