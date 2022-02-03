package pg

import (
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/db/tc"
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

func TestSaveOrUpdateInTx(t *testing.T) {

}
