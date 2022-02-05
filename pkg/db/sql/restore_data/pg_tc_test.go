package restore_data

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/tc"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/gofrs/uuid"
	"github.com/testcontainers/testcontainers-go"
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
	txToSaveRestoreDataEntity := util.MustBeginTx(test.CTX, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	print(txToSaveRestoreDataEntity)
	//SaveOrUpdateInTx()
}
