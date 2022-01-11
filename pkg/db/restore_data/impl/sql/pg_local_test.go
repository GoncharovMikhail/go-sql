package sql

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/restore_data"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	pgUuidType "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"testing"
)

//eagerInit
var (
	ctx = context.Background()
)

//lateInit
var (
	entityToSave *entity.RestoreDataEntity
	userUuid     pgUuidType.UUID
	repository   restore_data.SQLRestoreDataRepository
	tx           *sql.Tx
)

func init() {
	var err error
	uuiD, err := uuid.NewV1()
	if err != nil {
		log.Panic(err)
	}
	userUuid = pgUuidType.UUID{
		UUID: uuiD,
	}
	entityToSave = &entity.RestoreDataEntity{
		UserId: userUuid,
		Email:  uuiD.String(),
	}
	config := util.MustParseConfig(test.PGURL)
	config.User = test.PGUsername
	config.Password = test.PGPassword
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
	repository = NewPostgresRestoreDataRepository()
}

func TestPostgresRestoreDataRepository_SaveInTx(t *testing.T) {

}
