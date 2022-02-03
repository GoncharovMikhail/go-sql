//todo почему вместе не пробегают? По отдельности все ок
package pg

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/tc"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
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
	savedEntity, tx := saveInTx(uuid.Nil, uuidStringToUse, uuidStringToUse)
	util.MustCommitTx(tx)

	assert.Assert(t, savedEntity.Username == uuidStringToUse)
	assert.Assert(t, savedEntity.Password == uuidStringToUse)
	assert.Assert(t, savedEntity.Id.UUID.String() != "")
}

// todo пофиксить
func TestUpdateInTx(t *testing.T) {
	// Save
	randomUuidStringValue := randomUuid.String()
	savedEntity, txSaved := saveInTx(uuid.Nil, randomUuidStringValue, randomUuidStringValue)
	util.MustCommitTx(txSaved)
	// Update
	randomUuid = tc.InitUuid()
	randomUuidStringValue = randomUuid.String()
	updatedEntity, txUpdated := saveInTx(randomUuid, randomUuidStringValue, randomUuidStringValue)
	util.MustCommitTx(txUpdated)

	assert.Assert(t, savedEntity != nil)
	assert.Assert(t, randomUuid == updatedEntity.Id.UUID)
	assert.Assert(t, randomUuidStringValue == updatedEntity.Username)
	assert.Assert(t, randomUuidStringValue == updatedEntity.Password)
}

func TestFindOneByUsernameInTx(t *testing.T) {
	uuidStringToUse := randomUuid.String()
	savedUserEntity, txToSaveUser := saveInTx(uuid.Nil, uuidStringToUse, uuidStringToUse)
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

func saveInTx(id uuid.UUID, username, password string) (*entity.UserDataEntity, *sql.Tx) {
	tx := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})

	ude := &entity.UserDataEntity{
		Username: username,
		Password: password,
	}
	if id != uuid.Nil {
		ude.Id = pgUuidType.UUID{
			UUID: id,
		}
	}
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
